package notification

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ServiceAccountConfig struct {
	Type        string `json:"type"`
	ProjectId   string `json:"project_id"`
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	TokenUri    string `json:"token_uri"`
}

type FCMV1MessagePayload struct {
	Message FCMV1Message `json:"message"`
}

type FCMV1Message struct {
	Token        string                 `json:"token"`
	Notification *FCMV1Notification     `json:"notification,omitempty"`
	Data         map[string]string      `json:"data,omitempty"`
	Android      *FCMV1AndroidConfig    `json:"android,omitempty"`
	Webpush      *FCMV1WebpushConfig    `json:"webpush,omitempty"`
}

type FCMV1Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type FCMV1AndroidConfig struct {
	Priority     string                 `json:"priority"`
	Notification *FCMV1AndroidNotif     `json:"notification,omitempty"`
	Data         map[string]string      `json:"data,omitempty"`
}

type FCMV1AndroidNotif struct {
	Sound       string `json:"sound,omitempty"`
	ClickAction string `json:"click_action,omitempty"`
}

type FCMV1WebpushConfig struct {
	Notification *FCMV1WebpushNotif `json:"notification,omitempty"`
}

type FCMV1WebpushNotif struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Icon  string `json:"icon,omitempty"`
}

var (
	oauthTokenMu   sync.Mutex
	cachedToken    string
	tokenExpiresAt time.Time
	parsedKey      *rsa.PrivateKey
	saConfig       *ServiceAccountConfig
)

func loadServiceAccount() error {
	if saConfig != nil && parsedKey != nil {
		return nil
	}

	var data []byte
	var err error

	// 1. Cek jika credentials ditaruh langsung di .env (misal untuk Railway / Docker / Cloud)
	envJSON := os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON")
	if envJSON != "" {
		data = []byte(envJSON)
	} else {
		// 2. Cek file lokal
		paths := []string{
			os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH"),
			"service-account.json",
			"../service-account.json",
			"backend/service-account.json",
		}

		for _, p := range paths {
			if p == "" {
				continue
			}
			data, err = os.ReadFile(p)
			if err == nil {
				break
			}
		}
	}

	if len(data) == 0 {
		return fmt.Errorf("service-account.json atau FIREBASE_SERVICE_ACCOUNT_JSON tidak ditemukan di .env/file")
	}

	var sa ServiceAccountConfig
	if err := json.Unmarshal(data, &sa); err != nil {
		return fmt.Errorf("invalid service account json format: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(sa.PrivateKey))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	saConfig = &sa
	parsedKey = key
	return nil
}

func getGoogleAccessToken() (string, error) {
	oauthTokenMu.Lock()
	defer oauthTokenMu.Unlock()

	if cachedToken != "" && time.Now().Before(tokenExpiresAt) {
		return cachedToken, nil
	}

	if err := loadServiceAccount(); err != nil {
		return "", err
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   saConfig.ClientEmail,
		"scope": "https://www.googleapis.com/auth/firebase.messaging",
		"aud":   "https://oauth2.googleapis.com/token",
		"iat":   now.Unix(),
		"exp":   now.Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedJWT, err := token.SignedString(parsedKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign Google OAuth JWT: %w", err)
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  {signedJWT},
	})
	if err != nil {
		return "", fmt.Errorf("failed to request OAuth token: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OAuth token endpoint returned %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return "", err
	}

	cachedToken = res.AccessToken
	tokenExpiresAt = time.Now().Add(time.Duration(res.ExpiresIn-60) * time.Second)
	return cachedToken, nil
}

// SendFCMPush broadcasts a push notification via FCM HTTP v1 to given tokens (Android HP + Web)
func SendFCMPush(tokens []string, title, body string, data map[string]interface{}) error {
	if len(tokens) == 0 {
		return nil
	}

	accessToken, err := getGoogleAccessToken()
	if err != nil {
		log.Printf("[FCM][WARN] Gagal mendapatkan Google Access Token: %v. Cek service-account.json.", err)
		return err
	}

	// Ubah data map ke map[string]string untuk FCM v1
	stringData := make(map[string]string)
	for k, v := range data {
		stringData[k] = fmt.Sprintf("%v", v)
	}

	fcmURL := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", saConfig.ProjectId)
	client := &http.Client{Timeout: 10 * time.Second}

	for _, token := range tokens {
		go func(targetToken string) {
			payload := FCMV1MessagePayload{
				Message: FCMV1Message{
					Token: targetToken,
					Notification: &FCMV1Notification{
						Title: title,
						Body:  body,
					},
					Data: stringData,
					Android: &FCMV1AndroidConfig{
						Priority: "high",
						Notification: &FCMV1AndroidNotif{
							Sound:       "default",
							ClickAction: "FLUTTER_NOTIFICATION_CLICK",
						},
						Data: stringData,
					},
					Webpush: &FCMV1WebpushConfig{
						Notification: &FCMV1WebpushNotif{
							Title: title,
							Body:  body,
							Icon:  "/logo.png",
						},
					},
				},
			}

			jsonBytes, err := json.Marshal(payload)
			if err != nil {
				return
			}

			req, err := http.NewRequest("POST", fcmURL, bytes.NewBuffer(jsonBytes))
			if err != nil {
				return
			}
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("[FCM][ERROR] Gagal mengirim push notification ke token %s: %v", targetToken[:10], err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				log.Printf("[FCM][SUCCESS] FCM v1 Push terkirim sukses ke token %s... (Status: %s)", targetToken[:10], resp.Status)
			} else {
				respBody, _ := io.ReadAll(resp.Body)
				log.Printf("[FCM][INFO] FCM v1 Response (%s): %s", resp.Status, string(respBody))
			}
		}(token)
	}

	return nil
}
