package services

import (
	"fmt"
	"math/rand"
	"net/smtp"

	"github.com/pkg/errors"
)

type EmailsService interface {
	SendConfirmationEmail(to, token string) error
	GenerateRandomEmailToken() string
}

const (
	// The format of the message that will be sent according to RFC 822
	msgFormat = "From: %s\nTo: %s\nSubject: %s\n\n%s: %s?%s=%s&%s=%s"
	// The valid charset that can be used unencoded within URLs
	validCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~-_.!*()',"
)

type emailsService struct {
	addr                    string
	from                    string
	subject                 string
	msg                     string
	confirmationURL         string
	tokenQueryParameterName string
	emailQueryParameterName string
	tokenLength             int
	randGenerator           *rand.Rand
	auth                    smtp.Auth
}

func NewEmailsService(host, port, from, username, password, subject, genericMessage, confirmationURL, tokenQueryParameterName, emailQueryParameterName string, tokenLength int, randGenerator *rand.Rand) EmailsService {
	return &emailsService{
		addr:                    fmt.Sprintf("%s:%s", host, port),
		from:                    from,
		subject:                 subject,
		msg:                     genericMessage,
		confirmationURL:         confirmationURL,
		tokenQueryParameterName: tokenQueryParameterName,
		emailQueryParameterName: emailQueryParameterName,
		tokenLength:             tokenLength,
		randGenerator:           randGenerator,
		auth:                    smtp.PlainAuth("", username, password, host),
	}
}

func (es *emailsService) SendConfirmationEmail(to, token string) error {
	msg := fmt.Sprintf(msgFormat, es.from, to, es.subject, es.msg, es.confirmationURL, es.emailQueryParameterName, to, es.tokenQueryParameterName, token)
	err := smtp.SendMail(es.addr, es.auth, es.from, []string{to}, []byte(msg))

	return errors.Wrap(err, "could not send mail")
}

func (es *emailsService) GenerateRandomEmailToken() string {
	b := make([]byte, es.tokenLength)
	for i := range b {
		b[i] = validCharset[es.randGenerator.Intn(len(validCharset))]
	}

	return string(b)
}
