package services

import (
	"net/smtp"
	"os"

	"starter/pkg/logger"

	"github.com/jordan-wright/email"
)

type EmailService struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewEmailService() *EmailService {
	return &EmailService{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

func (s *EmailService) SendEmail(to []string, subject, body string) error {
	log := logger.GetLogger()

	e := email.NewEmail()
	e.From = s.Username
	e.To = to
	e.Subject = subject
	e.Text = []byte(body)

	addr := s.Host + ":" + s.Port
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	err := e.Send(addr, auth)
	if err != nil {
		log.Error().Err(err).Str("to", to[0]).Msg("Failed to send email")
		return err
	}

	log.Info().Str("to", to[0]).Msg("Email sent successfully")
	return nil
}
