package notications

import (
	"fmt"
	"net"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SimpleEmail struct {
	To      string
	Subject string
	Body    string
}

type EmailSender struct {
	config *SMTPConfig
}

func NewEmailSender(config *SMTPConfig) *EmailSender {
	return &EmailSender{config: config}
}

func (s *EmailSender) SendSimpleEmail(email *SimpleEmail) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Establish a connection to the SMTP server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create an SMTP client
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate if credentials are provided
	if s.config.Username != "" && s.config.Password != "" {
		auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
		}
	}

	// Set the sender and recipient
	if err := client.Mail(s.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(email.To); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Get the data writer from the SMTP client
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send email data: %w", err)
	}
	defer wc.Close()

	// Construct the email message and write it to the SMTP server
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.config.From, email.To, email.Subject, email.Body)
	if _, err := wc.Write([]byte(message)); err != nil {
		return fmt.Errorf("failed to write email message: %w", err)
	}

	return nil
}

func (s *EmailSender) SendLoginNotification(userEmail, userName string) error {
	subject := "Login Notification"
	body := fmt.Sprintf("Hello %s,\n\nYou have successfully logged in to your account.\n\nBest regards,\nYour Team", userName)

	email := &SimpleEmail{
		To:      userEmail,
		Subject: subject,
		Body:    body,
	}

	return s.SendSimpleEmail(email)
}
