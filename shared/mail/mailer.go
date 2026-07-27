package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"strconv"
)

type Config struct {
	Host      string
	Port      string
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

type Mailer struct {
	config Config
}

func NewMailer(config Config) *Mailer {
	return &Mailer{
		config: config,
	}
}

type Message struct {
	To      []string
	Subject string
	HTML    string
	ReplyTo string
}

func (m *Mailer) Send(message Message) error {
	if len(message.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	addr := net.JoinHostPort(
		m.config.Host,
		m.config.Port,
	)

	port, err := strconv.Atoi(m.config.Port)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	client, err := smtp.NewClient(conn, m.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}

	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName: m.config.Host,
			MinVersion: tls.VersionTLS12,
		}

		if err := client.StartTLS(
			tlsConfig,
		); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	auth := smtp.PlainAuth(
		"",
		m.config.Username,
		m.config.Password,
		m.config.Host,
	)

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	if err := client.Mail(m.config.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range message.To {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set reciepnt %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open SMTP data writer: %w", err)
	}

	var buffer bytes.Buffer

	fromName := mime.QEncoding.Encode(
		"UTF-8",
		m.config.FromEmail,
	)

	buffer.WriteString(
		fmt.Sprintf(
			"From: %s <%s>\r\n",
			fromName,
			m.config.FromEmail,
		),
	)

	buffer.WriteString(
		fmt.Sprintf(
			"To: %s\r\n",
			message.To[0],
		),
	)

	buffer.WriteString(
		fmt.Sprintf(
			"Subject: %s\r\n",
			mime.QEncoding.Encode(
				"UTF-8",
				message.Subject,
			),
		),
	)

	if message.ReplyTo != "" {
		buffer.WriteString(
			fmt.Sprintf(
				"Reply-To: %s\r\n",
				message.ReplyTo,
			),
		)
	}

	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	buffer.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	buffer.WriteString("\r\n")
	buffer.WriteString(message.HTML)

	if _, err := writer.Write(
		buffer.Bytes(),
	); err != nil {
		_ = writer.Close()

		return fmt.Errorf("failed to write email bode: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close SMTP data writer: %w", err)
	}

	if err := client.Quit(); err != nil {
		return fmt.Errorf("failed to quit SMTP session: %w", err)
	}

	_ = port

	return nil
}
