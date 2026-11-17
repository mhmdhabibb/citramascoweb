package auth

import (
	"bytes"
	"citramascoweb-backend/config"
	"citramascoweb-backend/internal/modules/rbac"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	SignIn(req *SignInRequest) (*SignInResponse, error)
	GetProfile(userId string, token string) (*SignInResponse, error)
}

type authService struct {
	db         *gorm.DB
	httpClient *http.Client
}

func NewAuthService(db *gorm.DB) AuthService {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &authService{
		db: db,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   15 * time.Second,
		},
	}
}

func (s *authService) getExternalBaseURL() string {
	url := config.Config("EXTERNAL_AUTH_URL")
	if url == "" {
		url = "https://api-hr.lnginternasional.com"
	}
	return strings.TrimRight(url, "/")
}

func cleanStr(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, ".", "")
	return s
}

// Map department & position from HRIS to Citramasco role slug
func mapDepartmentAndPositionToRole(user *ExternalUser) string {
	roleClean := cleanStr(user.Role)
	usernameClean := cleanStr(user.Username)

	deptClean := ""
	posClean := ""
	if user.Employee != nil {
		if user.Employee.Department != nil {
			deptClean = cleanStr(user.Employee.Department.Name)
		}
		if user.Employee.Position != nil {
			posClean = cleanStr(user.Employee.Position.Name)
		}
	}

	// 1. System Administrator / Owner from HRIS
	if strings.Contains(roleClean, "admin") || strings.Contains(roleClean, "owner") || strings.Contains(usernameClean, "admin") || strings.Contains(posClean, "admin") {
		return "admin"
	}

	// 2. Executive / General Manager -> Hotel Manager
	if strings.Contains(posClean, "generalmanager") || strings.Contains(posClean, "hotelmanager") || strings.Contains(posClean, "director") || strings.Contains(posClean, "direktur") {
		return "manager"
	}

	// 3. Housekeeping / Room Attendant / Cleaning
	// Mendukung "House Keeping", "Staff House Keeping", "Housekeeper", "Room Attendant", dll
	if strings.Contains(posClean, "housekeep") || strings.Contains(deptClean, "housekeep") || strings.Contains(posClean, "roomattendant") || strings.Contains(posClean, "cleaner") || strings.Contains(posClean, "cleaning") || strings.Contains(deptClean, "cleaning") || strings.Contains(deptClean, "hk") {
		return "housekeeping"
	}

	// 4. Front Office / Front Desk / Receptionist
	if strings.Contains(posClean, "reception") || strings.Contains(posClean, "frontdesk") || strings.Contains(posClean, "frontoffice") || strings.Contains(posClean, "fo") || strings.Contains(posClean, "guestservice") || strings.Contains(deptClean, "front") || strings.Contains(deptClean, "reception") {
		return "reception"
	}

	// 5. Finance / Accounting / Cashier
	if strings.Contains(posClean, "finance") || strings.Contains(posClean, "account") || strings.Contains(posClean, "kasir") || strings.Contains(posClean, "cashier") || strings.Contains(deptClean, "finance") || strings.Contains(deptClean, "accounting") || strings.Contains(deptClean, "keuangan") {
		return "finance"
	}

	// 6. Inventory / Purchasing / Logistik / Gudang
	if strings.Contains(posClean, "purchas") || strings.Contains(posClean, "gudang") || strings.Contains(posClean, "inventory") || strings.Contains(posClean, "storekeeper") || strings.Contains(posClean, "warehouse") || strings.Contains(posClean, "logistik") || strings.Contains(deptClean, "purchas") || strings.Contains(deptClean, "gudang") || strings.Contains(deptClean, "inventory") {
		return "inventory"
	}

	// 7. General Managers / Supervisors
	if strings.Contains(roleClean, "manager") || strings.Contains(posClean, "manager") || strings.Contains(posClean, "supervisor") || strings.Contains(posClean, "spv") {
		return "manager"
	}

	// Default fallback
	return "reception"
}

func (s *authService) SignIn(req *SignInRequest) (*SignInResponse, error) {
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/api/auth/sign-in", s.getExternalBaseURL())
	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) CitramascoBackend/1.0")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth server: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth server response: %w", err)
	}

	trimmed := strings.TrimSpace(string(bodyBytes))
	if strings.HasPrefix(trimmed, "<") {
		return nil, fmt.Errorf("auth server returned HTML page (HTTP %d). Make sure EXTERNAL_AUTH_URL points to https://api-hr.lnginternasional.com", resp.StatusCode)
	}

	var authResp ExternalAuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, fmt.Errorf("invalid response from auth server: %w", err)
	}

	if resp.StatusCode != http.StatusOK || authResp.Data == nil {
		errMsg := authResp.Error
		if errMsg == "" {
			errMsg = authResp.Message
		}
		if errMsg == "" {
			errMsg = "Invalid username or password"
		}
		return nil, errors.New(errMsg)
	}

	user := authResp.Data.User
	citramascoRole := ""

	// Check local user_roles in Citramasco DB
	var userRole rbac.UserRole
	err = s.db.Preload("Role").Where("user_id = ?", user.Id).First(&userRole).Error

	mappedSlug := mapDepartmentAndPositionToRole(&user)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// Auto-map based on Department & Position from HRIS
		var targetRole rbac.Role
		if err := s.db.Where("slug = ?", mappedSlug).First(&targetRole).Error; err == nil {
			assignedBy := "auto-hr-mapping"
			newUserRole := rbac.UserRole{
				UserId:     user.Id,
				RoleId:     targetRole.Id,
				AssignedBy: &assignedBy,
				AssignedAt: time.Now(),
			}
			s.db.Create(&newUserRole)
			citramascoRole = targetRole.Slug
		}
	} else if err == nil && userRole.Role != nil {
		citramascoRole = userRole.Role.Slug
		// If mapped role is admin or manager but DB has a lower role, upgrade to admin/manager
		if (mappedSlug == "admin" || mappedSlug == "manager") && citramascoRole != mappedSlug {
			var upgradedRole rbac.Role
			if err := s.db.Where("slug = ?", mappedSlug).First(&upgradedRole).Error; err == nil {
				s.db.Model(&userRole).Update("role_id", upgradedRole.Id)
				citramascoRole = upgradedRole.Slug
			}
		}
	}

	deptName := ""
	posName := ""
	if user.Employee != nil {
		if user.Employee.Department != nil {
			deptName = user.Employee.Department.Name
		}
		if user.Employee.Position != nil {
			posName = user.Employee.Position.Name
		}
	}

	// Fetch active permissions for this user in Citramasco
	var permissions []string
	s.db.Table("permissions").
		Distinct("permissions.name").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.deleted_at IS NULL", user.Id).
		Pluck("permissions.name", &permissions)

	// Enforce Shift-Based Access Control (Admin and Manager are 24/7 bypassed)
	if err := s.checkShiftAccess(authResp.Data.AccessToken, &user, citramascoRole); err != nil {
		return nil, err
	}

	return &SignInResponse{
		AccessToken:    authResp.Data.AccessToken,
		RefreshToken:   authResp.Data.RefreshToken,
		User:           user,
		CitramascoRole: citramascoRole,
		Department:     deptName,
		Position:       posName,
		Permissions:    permissions,
	}, nil
}

func (s *authService) checkShiftAccess(token string, user *ExternalUser, roleSlug string) error {
	// 0. Optional environment bypass (useful for development / testing without active shift roster)
	if config.Config("BYPASS_SHIFT_CHECK") == "true" {
		return nil
	}

	// 1. Role Admin and Manager have 24/7 bypass access
	if roleSlug == "admin" || roleSlug == "manager" {
		return nil
	}

	// 2. Identify Employee ID
	empId := ""
	if user.EmployeeId != nil && *user.EmployeeId != "" {
		empId = *user.EmployeeId
	} else if user.Employee != nil && user.Employee.Id != "" {
		empId = user.Employee.Id
	}

	if empId == "" {
		return errors.New("akses ditolak: akun Anda belum terhubung dengan data karyawan HR untuk verifikasi shift")
	}

	// 3. Current time in WIB (UTC+7)
	loc := time.FixedZone("WIB", 7*3600)
	now := time.Now().In(loc)
	todayStr := now.Format("2006-01-02")

	// 4. Query Shift Roster from HRIS API
	endpoint := fmt.Sprintf("%s/api/shift-roster/?start_date=%s&end_date=%s", s.getExternalBaseURL(), todayStr, todayStr)
	httpReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("gagal membuat request validasi shift: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 CitramascoBackend/1.0")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("gagal menghubungkan ke layanan shift HRIS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("gagal memvalidasi jadwal shift dari HRIS")
	}

	var rosterResp ExternalRosterResponse
	if err := json.NewDecoder(resp.Body).Decode(&rosterResp); err != nil {
		return fmt.Errorf("format data shift tidak valid: %w", err)
	}

	var userRoster *ExternalShiftRoster
	for i := range rosterResp.Data.Data {
		if rosterResp.Data.Data[i].EmployeeId == empId {
			userRoster = &rosterResp.Data.Data[i]
			break
		}
	}

	if userRoster == nil {
		return errors.New("akses ditolak: Anda tidak memiliki jadwal shift kerja terdaftar untuk hari ini. Web terkunci")
	}

	if userRoster.IsDayOff {
		return errors.New("akses ditolak: Hari ini adalah jadwal Libur (Day Off) Anda. Web terkunci")
	}

	if userRoster.ShiftMaster == nil {
		return errors.New("akses ditolak: Detail shift kerja Anda belum ditentukan untuk hari ini. Web terkunci")
	}

	startParts := strings.Split(userRoster.ShiftMaster.StartTime, ":")
	endParts := strings.Split(userRoster.ShiftMaster.EndTime, ":")
	if len(startParts) >= 2 && len(endParts) >= 2 {
		startHour, _ := strconv.Atoi(startParts[0])
		startMin, _ := strconv.Atoi(startParts[1])
		endHour, _ := strconv.Atoi(endParts[0])
		endMin, _ := strconv.Atoi(endParts[1])

		shiftStart := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMin, 0, 0, loc)
		shiftEnd := time.Date(now.Year(), now.Month(), now.Day(), endHour, endMin, 0, 0, loc)

		if userRoster.ShiftMaster.IsOvernight || shiftEnd.Before(shiftStart) {
			if now.Hour() < endHour || (now.Hour() == endHour && now.Minute() <= endMin) {
				shiftStart = shiftStart.AddDate(0, 0, -1)
			} else {
				shiftEnd = shiftEnd.AddDate(0, 0, 1)
			}
		}

		// Berikan toleransi 30 menit sebelum shift (persiapan/handover) dan 30 menit setelah shift (laporan/closing)
		allowedStart := shiftStart.Add(-30 * time.Minute)
		allowedEnd := shiftEnd.Add(30 * time.Minute)

		if now.Before(allowedStart) || now.After(allowedEnd) {
			return fmt.Errorf("akses ditolak: Anda berada di luar jam shift kerja Anda (%s: %s - %s). Web otomatis terkunci", userRoster.ShiftMaster.Name, userRoster.ShiftMaster.StartTime, userRoster.ShiftMaster.EndTime)
		}
	}

	return nil
}

func (s *authService) GetProfile(userId string, token string) (*SignInResponse, error) {
	endpoint := fmt.Sprintf("%s/api/auth/me", s.getExternalBaseURL())
	httpReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err == nil {
		httpReq.Header.Set("Authorization", "Bearer "+token)
		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := s.httpClient.Do(httpReq)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var meResp struct {
				Data ExternalUser `json:"data"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&meResp); err == nil && meResp.Data.Id != "" {
				var roleSlug string
				s.db.Table("user_roles").
					Select("roles.slug").
					Joins("JOIN roles ON roles.id = user_roles.role_id").
					Where("user_roles.user_id = ? AND roles.deleted_at IS NULL", userId).
					Limit(1).
					Pluck("roles.slug", &roleSlug)

				if roleSlug == "" {
					roleSlug = mapDepartmentAndPositionToRole(&meResp.Data)
				}

				deptName := ""
				posName := ""
				if meResp.Data.Employee != nil {
					if meResp.Data.Employee.Department != nil {
						deptName = meResp.Data.Employee.Department.Name
					}
					if meResp.Data.Employee.Position != nil {
						posName = meResp.Data.Employee.Position.Name
					}
				}

				var permissions []string
				s.db.Table("permissions").
					Distinct("permissions.name").
					Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
					Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
					Where("user_roles.user_id = ? AND permissions.deleted_at IS NULL", userId).
					Pluck("permissions.name", &permissions)

				// Enforce Shift-Based Access Control
				if err := s.checkShiftAccess(token, &meResp.Data, roleSlug); err != nil {
					return nil, err
				}

				return &SignInResponse{
					AccessToken:    token,
					User:           meResp.Data,
					CitramascoRole: roleSlug,
					Department:     deptName,
					Position:       posName,
					Permissions:    permissions,
				}, nil
			}
		}
	}

	return nil, errors.New("failed to retrieve profile")
}
