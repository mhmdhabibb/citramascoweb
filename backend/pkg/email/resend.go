package email

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client    *resend.Client
	fromEmail string
	domain    string
}

func NewEmailService() *EmailService {
	apiKey := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("RESEND_FROM_EMAIL")
	if from == "" {
		from = "CM Living <reservations@lngdemo.my.id>"
	}

	domain := os.Getenv("APP_DOMAIN")
	if domain == "" {
		// Extract domain from fromEmail (e.g., "CM Living <reservations@customdomain.com>" -> "customdomain.com")
		if idx := strings.Index(from, "@"); idx != -1 {
			domain = strings.Trim(from[idx+1:], " >\"'")
		} else {
			domain = "lngdemo.my.id"
		}
	}

	if apiKey == "" {
		log.Println("[EmailService] WARNING: RESEND_API_KEY is not set. Emails will be skipped in dev mode.")
	}

	return &EmailService{
		client:    resend.NewClient(apiKey),
		fromEmail: from,
		domain:    domain,
	}
}

// SendBookingConfirmation sends an elegant booking confirmation e-receipt to the guest
func (s *EmailService) SendBookingConfirmation(
	toEmail string,
	guestName string,
	roomName string,
	roomNumber string,
	checkinDate string,
	checkoutDate string,
	bookingCode string,
	totalPayment float64,
	status string,
) error {
	if s.client == nil || os.Getenv("RESEND_API_KEY") == "" {
		log.Printf("[EmailService] (Dev Mode) Skipped sending booking confirmation to %s (#%s)", toEmail, bookingCode)
		return nil
	}

	if strings.TrimSpace(toEmail) == "" {
		return nil
	}

	roomDisplay := roomName
	if roomNumber != "" {
		roomDisplay = fmt.Sprintf("%s (Kamar No. %s)", roomName, roomNumber)
	}

	statusDisplay := strings.ToUpper(status)
	if status == "approve" || status == "approved" {
		statusDisplay = "CONFIRMED & APPROVED"
	} else if status == "pending" {
		statusDisplay = "PENDING CONFIRMATION"
	}

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Konfirmasi Reservasi CM Living</title>
  <style>
    body { font-family: 'Plus Jakarta Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #faf8f5; color: #1e1711; margin: 0; padding: 24px; }
    .container { max-width: 580px; margin: 0 auto; background: #ffffff; border: 1px solid #faebc6; border-radius: 20px; overflow: hidden; box-shadow: 0 12px 36px rgba(180, 140, 60, 0.08); }
    .header { background: #1a1612; padding: 36px 24px; text-align: center; color: #ffffff; border-bottom: 3px solid #faebc6; }
    .logo-monogram { font-family: Georgia, serif; font-size: 32px; font-weight: 700; color: #faebc6; letter-spacing: 0.15em; margin: 0; }
    .subtitle { font-size: 11px; text-transform: uppercase; letter-spacing: 0.2em; color: rgba(250, 235, 198, 0.85); margin-top: 6px; }
    .content { padding: 36px 28px; }
    .greeting { font-size: 17px; font-weight: 700; color: #1a1612; margin-top: 0; margin-bottom: 12px; }
    .intro-text { font-size: 14px; color: #6b5c47; line-height: 1.6; margin-bottom: 24px; }
    .card-box { background: #fdfbf7; border: 1px solid #faebc6; border-radius: 16px; padding: 22px; margin-bottom: 24px; }
    .card-title { font-size: 13px; font-weight: 800; text-transform: uppercase; letter-spacing: 0.08em; color: #8c6716; margin-top: 0; margin-bottom: 16px; border-bottom: 1px solid #faebc6; padding-bottom: 8px; }
    .detail-row { display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 13.5px; }
    .detail-label { color: #8c7a62; }
    .detail-value { font-weight: 700; color: #1a1612; text-align: right; }
    .total-row { border-top: 2px dashed #faebc6; padding-top: 14px; margin-top: 14px; display: flex; justify-content: space-between; font-size: 16px; font-weight: 800; color: #1a1612; }
    .total-price { color: #8c6716; }
    .badge { display: inline-block; padding: 4px 12px; border-radius: 100px; font-size: 11px; font-weight: 800; background: #faebc6; color: #7a5713; border: 1px solid rgba(140, 103, 22, 0.2); }
    .instruction-box { background: #faf8f5; border-left: 4px solid #faebc6; padding: 14px 18px; border-radius: 0 12px 12px 0; margin-bottom: 24px; font-size: 13px; color: #6b5c47; line-height: 1.5; }
    .footer { background: #faf8f5; padding: 24px; text-align: center; font-size: 12px; color: #9c8b74; border-top: 1px solid #faebc6; }
    .footer a { color: #8c6716; text-decoration: none; font-weight: 600; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1 class="logo-monogram">CM LIVING</h1>
      <div class="subtitle">Luxury Hotel & Suites</div>
    </div>
    <div class="content">
      <h2 class="greeting">Halo, %s</h2>
      <p class="intro-text">
        Terima kasih telah memilih <strong>CM Living</strong> sebagai akomodasi Anda. Reservasi Anda telah berhasil tercatat di sistem kami dengan status <span class="badge">%s</span>.
      </p>

      <div class="card-box">
        <div class="card-title">Ringkasan Reservasi</div>
        
        <div class="detail-row">
          <span class="detail-label">Kode Booking</span>
          <span class="detail-value">#%s</span>
        </div>
        
        <div class="detail-row">
          <span class="detail-label">Tipe Kamar</span>
          <span class="detail-value">%s</span>
        </div>
        
        <div class="detail-row">
          <span class="detail-label">Jadwal Check-In</span>
          <span class="detail-value">%s (14:00 WIB)</span>
        </div>
        
        <div class="detail-row">
          <span class="detail-label">Jadwal Check-Out</span>
          <span class="detail-value">%s (12:00 WIB)</span>
        </div>

        <div class="total-row">
          <span>Total Pembayaran</span>
          <span class="total-price">Rp %s</span>
        </div>
      </div>

      <div class="instruction-box">
        <strong>Informasi Kedatangan:</strong><br>
        Silakan tunjukkan email konfirmasi atau Kode Booking ini kepada resepsionis kami saat proses Check-In di lobi CM Living.
      </div>

      <p class="intro-text" style="margin-bottom: 0;">
        Jika Anda membutuhkan layanan khusus atau perubahan reservasi, silakan hubungi tim kami melalui email atau resepsionis.
      </p>
    </div>
    <div class="footer">
      &copy; 2026 <strong>CM Living</strong> &bull; <a href="https://%s">%s</a><br>
      Kompleks CM Living, Indonesia
    </div>
  </div>
</body>
</html>
`,
		guestName,
		statusDisplay,
		bookingCode,
		roomDisplay,
		checkinDate,
		checkoutDate,
		formatRupiah(totalPayment),
		s.domain,
		s.domain,
	)

	params := &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("Konfirmasi Reservasi CM Living - Booking #%s", bookingCode),
		Html:    htmlContent,
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Error sending email to %s: %v", toEmail, err)
		return err
	}

	log.Printf("[EmailService] Successfully sent email to %s (Resend ID: %s)", toEmail, sent.Id)
	return nil
}

// SendReservationStatusNotification sends notification when status changes (e.g. Approved / Cancelled / Checked-In)
func (s *EmailService) SendReservationStatusNotification(
	toEmail string,
	guestName string,
	bookingCode string,
	statusTitle string,
	statusMessage string,
) error {
	if s.client == nil || os.Getenv("RESEND_API_KEY") == "" || strings.TrimSpace(toEmail) == "" {
		return nil
	}

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: 'Plus Jakarta Sans', sans-serif; background-color: #faf8f5; color: #1e1711; margin: 0; padding: 24px; }
    .container { max-width: 560px; margin: 0 auto; background: #ffffff; border: 1px solid #faebc6; border-radius: 20px; overflow: hidden; box-shadow: 0 12px 36px rgba(180, 140, 60, 0.08); }
    .header { background: #1a1612; padding: 30px 24px; text-align: center; color: #ffffff; border-bottom: 3px solid #faebc6; }
    .logo-monogram { font-family: Georgia, serif; font-size: 28px; font-weight: 700; color: #faebc6; letter-spacing: 0.15em; margin: 0; }
    .content { padding: 32px 26px; }
    .title { font-size: 18px; font-weight: 700; color: #1a1612; margin-top: 0; }
    .message-box { background: #fdfbf7; border: 1px solid #faebc6; border-radius: 14px; padding: 18px; margin: 20px 0; font-size: 14px; line-height: 1.6; color: #5c4b38; }
    .footer { background: #faf8f5; padding: 20px; text-align: center; font-size: 12px; color: #9c8b74; border-top: 1px solid #faebc6; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1 class="logo-monogram">CM LIVING</h1>
    </div>
    <div class="content">
      <h2 class="title">%s</h2>
      <p>Halo <strong>%s</strong> (Kode Booking: #%s),</p>
      <div class="message-box">
        %s
      </div>
      <p style="font-size: 13px; color: #8c7a62;">Terima kasih atas kepercayaan Anda kepada CM Living.</p>
    </div>
    <div class="footer">
      &copy; 2026 CM Living &bull; %s
    </div>
  </div>
</body>
</html>
`, statusTitle, guestName, bookingCode, statusMessage, s.domain)

	params := &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("%s - CM Living (#%s)", statusTitle, bookingCode),
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	return err
}

func formatRupiah(amount float64) string {
	val := fmt.Sprintf("%.0f", amount)
	n := len(val)
	if n <= 3 {
		return val
	}
	var res []string
	rem := n % 3
	if rem > 0 {
		res = append(res, val[:rem])
	}
	for i := rem; i < n; i += 3 {
		res = append(res, val[i:i+3])
	}
	return strings.Join(res, ".")
}
