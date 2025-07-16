package email

import (
	"bytes"
	"fmt"
	"html/template"
	"imageboard/config"
	"imageboard/database"
	"imageboard/models"
	"net/smtp"
	"regexp"
)

func extractEmailAddress(from string) string {
	re := regexp.MustCompile(`<([^>]+)>`)
	matches := re.FindStringSubmatch(from)
	if len(matches) == 2 {
		return matches[1]
	}
	return from
}

func SendMail(to, subject, body string) error {
	var auth smtp.Auth
	if config.SMTP.Username != "" {
		auth = smtp.PlainAuth("", config.SMTP.Username, config.SMTP.Password, config.SMTP.Host)
	} else {
		auth = nil
	}
	fromHeader := config.SMTP.From
	fromAddress := extractEmailAddress(config.SMTP.From)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		fromHeader, to, subject, body)
	addr := fmt.Sprintf("%s:%d", config.SMTP.Host, config.SMTP.Port)
	return smtp.SendMail(addr, auth, fromAddress, []string{to}, []byte(msg))
}

func SendVerificationEmail(user *models.User) error {
	token, err := database.GenerateEmailToken(int(user.ID), models.EmailTokenTypeVerification)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	tmpl, err := template.ParseFiles("templates/email/verification.html")
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}
	verificationLink := fmt.Sprintf("%s%s?token=%s", config.Server.AppBaseURL, config.URL_VERIFY_EMAIL, token.Token)
	data := struct {
		Username string
		Appname  string
		Link     string
	}{
		Username: user.Username,
		Appname:  config.Server.AppName,
		Link:     verificationLink,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	subject := fmt.Sprintf("Verify your email for %s", config.Server.AppName)
	return SendMail(user.Email, subject, body.String())
}

// func SendPasswordResetEmail(user *models.User) error {
// 	token, err := user.GenerateToken(database.DB, models.EmailTokenTypePasswordReset)
// 	if err != nil {
// 		return fmt.Errorf("failed to generate password reset token: %w", err)
// 	}

// 	tmpl, err := template.ParseFiles("templates/email/password_reset.html")
// 	if err != nil {
// 		return fmt.Errorf("failed to parse email template: %w", err)
// 	}
// 	resetLink := fmt.Sprintf("%s/account/reset-password?token=%s", config.Server.AppBaseURL, token.Token)
// 	data := struct {
// 		Username string
// 		Link     string
// 	}{
// 		Username: user.Username,
// 		Link:     resetLink,
// 	}

// 	var body bytes.Buffer
// 	if err := tmpl.Execute(&body, data); err != nil {
// 		return fmt.Errorf("failed to execute email template: %w", err)
// 	}

// 	subject := fmt.Sprintf("Password reset for %s", config.Server.AppName)
// 	return SendMail(user.Email, subject, body.String())
// }

// func SendEmailChangeConfirmation(user *models.User, newEmail string) error {
// 	token, err := user.GenerateToken(database.DB, models.EmailTokenTypeChangeEmail)
// 	if err != nil {
// 		return fmt.Errorf("failed to generate email change token: %w", err)
// 	}

// 	tmpl, err := template.ParseFiles("templates/email/email_change.html")
// 	if err != nil {
// 		return fmt.Errorf("failed to parse email template: %w", err)
// 	}
// 	confirmLink := fmt.Sprintf("%s/account/confirm-email-change?token=%s&email=%s", config.Server.AppBaseURL, token.Token, newEmail)
// 	data := struct {
// 		Username string
// 		NewEmail string
// 		Link     string
// 	}{
// 		Username: user.Username,
// 		NewEmail: newEmail,
// 		Link:     confirmLink,
// 	}

// 	var body bytes.Buffer
// 	if err := tmpl.Execute(&body, data); err != nil {
// 		return fmt.Errorf("failed to execute email template: %w", err)
// 	}

// 	subject := fmt.Sprintf("Confirm email change for %s", config.Server.AppName)
// 	return SendMail(newEmail, subject, body.String())
// }
