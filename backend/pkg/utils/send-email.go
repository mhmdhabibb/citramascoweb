package utils

import (
	"citramascoweb-backend/config"
	"fmt"
	"net/smtp"
)

// SendEmail sends an HTML email to the specified recipient using SMTP.
func SendEmail(toEmail string, subject string, htmlContent string) error {
	host := config.Config("MAIL_TRIP_HOST")
	port := config.Config("MAIL_TRIP_PORT")
	username := config.Config("MAIL_TRIP_USERNAME")
	password := config.Config("MAIL_TRIP_PASSWORD")

	if host == "" || port == "" || username == "" || password == "" {
		return fmt.Errorf("SMTP configuration is missing in environment variables")
	}

	from := "info@citramascohotel.com"
	addr := fmt.Sprintf("%s:%s", host, port)

	// Set up authentication information.
	auth := smtp.PlainAuth("", username, password, host)

	// Format HTML mail headers and body
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subjectHeader := fmt.Sprintf("Subject: %s\n", subject)
	fromHeader := fmt.Sprintf("From: %s\n", from)
	toHeader := fmt.Sprintf("To: %s\n", toEmail)
	msg := []byte(fromHeader + toHeader + subjectHeader + mime + htmlContent)

	err := smtp.SendMail(addr, auth, from, []string{toEmail}, msg)
	if err != nil {
		fmt.Printf("SMTP Error: failed to send email to %s: %v\n", toEmail, err)
		return err
	}

	return nil
}

// FormatDotSeparator formats an integer with dot thousand separators (e.g. 11000000 -> 11.000.000)
func FormatDotSeparator(n int) string {
	isNegative := false
	if n < 0 {
		isNegative = true
		n = -n
	}
	
	in := fmt.Sprintf("%d", n)
	numOfDigits := len(in)
	var result []byte
	if isNegative {
		result = append(result, '-')
	}
	for i := 0; i < numOfDigits; i++ {
		if i > 0 && (numOfDigits-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, in[i])
	}
	return string(result)
}
