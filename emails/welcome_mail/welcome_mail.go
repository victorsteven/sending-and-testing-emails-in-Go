package welcome_mail

import (
	"bytes"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"mail-sending/helpers"
	"net/smtp"
	"os"
	"strings"
)

type service struct{}

type Service interface {
	SendWelcomeMail(*helpers.WelcomeMail) (bool, error)
}

func NewService() *service {
	return &service{}
}

//The service struct implements the Service interface
var _ Service = &service{}

func (s *service) SendWelcomeMail(cred *helpers.WelcomeMail) (bool, error) {

	em := &helpers.WelcomeMail{
		Temp:    "welcome_mail.html",
		Name:    cred.Name,
		Email:   cred.Email,
		Subject: "Welcome Onboard",
	}

	//Just to demonstrate how you can use multiple mail services
	if strings.Contains(em.Email, "@yahoo") {
		_, err := s.sendEmailUsingGmail(em)
		if err != nil {
			return false, err
		}
	} else {
		_, err := s.sendEmailUsingSendGrid(em)
		if err != nil {
			return false, err
		}
	}

	return true, nil

}

func (s *service) sendEmailUsingGmail(m *helpers.WelcomeMail) (bool, error) {

	gmailUsername := os.Getenv("GMAIL_USERNAME")
	gmailPassword := os.Getenv("GMAIL_PASSWORD")
	gmailServer := os.Getenv("GMAIL_SERVER")
	gmailPort := os.Getenv("GMAIL_PORT")

	smtpData := &helpers.WelcomeMail{
		Name:  m.Name,
		Email: m.Email,
	}

	auth := smtp.PlainAuth("",
		gmailUsername,
		gmailPassword,
		gmailServer,
	)

	dir, _ := os.Getwd()
	file := fmt.Sprintf("%s/templates/%s", dir, m.Temp)

	t, err := template.ParseFiles(file)
	if err != nil {
		return false, err
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	var doc bytes.Buffer

	doc.Write([]byte(fmt.Sprintf("Subject:"+m.Subject+"\n%s\n\n", mime)))

	err = t.Execute(&doc, smtpData)
	if err != nil {
		return false, err
	}

	emails := []string{m.Email} //or

	err = smtp.SendMail(gmailServer+":"+gmailPort,
		auth,
		gmailUsername,
		emails,
		doc.Bytes())

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *service) sendEmailUsingSendGrid(m *helpers.WelcomeMail) (bool, error) {

	fromAdmin := os.Getenv("SENDGRID_FROM")
	apiKey := os.Getenv("SENDGRID_API_KEY")

	from := mail.NewEmail("Wonderful Company", fromAdmin)
	subject := m.Subject
	to := mail.NewEmail("Client", m.Email)

	dir, _ := os.Getwd()
	file := fmt.Sprintf("%s/templates/%s", dir, m.Temp)

	t, err := template.ParseFiles(file)
	if err != nil {
		return false, err
	}

	smtpData := &helpers.WelcomeMail{
		Name:  m.Name,
		Email: m.Email,
	}

	var doc bytes.Buffer

	err = t.Execute(&doc, smtpData)
	if err != nil {
		return false, err
	}

	message := mail.NewSingleEmail(from, subject, to, doc.String(), doc.String())
	client := sendgrid.NewSendClient(apiKey)
	_, err = client.Send(message)
	if err != nil {
		return false, err
	}

	return true, nil
}
