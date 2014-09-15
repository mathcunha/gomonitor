package action

import (
	"bytes"
	"github.com/mathcunha/gomonitor/prop"
	"net/smtp"
	"text/template"
)

type SmtpTemplateData struct {
	From    string
	To      []string
	Subject string
	Body    string
}

const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}

Sincerely,

{{.From}}
`

func SimpleSendMail(from string, to []string, subject string, body string) error {
	var doc bytes.Buffer
	var err error

	context := SmtpTemplateData{from, to, subject, body}

	t := template.New("emailTemplate")
	t, err = t.Parse(emailTemplate)
	if err != nil {
		return err
	}

	err = t.Execute(&doc, context)
	if err != nil {
		return err
	}

	return simpleSendMail(prop.Property("smtp"), context.From, context.To, doc.Bytes())
}

func simpleSendMail(endpoint string, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(endpoint)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()

	if err != nil {
		return err
	}
	return c.Quit()
}
