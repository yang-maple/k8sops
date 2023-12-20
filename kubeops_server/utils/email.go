package utils

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"kubeops/config"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// Emails 发送邮件
func Emails(emails string, content string, subject string) (err error) {
	server := mail.NewSMTPClient()
	//SMTP Server
	server.Host = config.HostSmtp
	server.Port = config.Port
	server.Username = config.Username
	server.Password = config.Password

	server.KeepAlive = false

	server.ConnectTimeout = 10 * time.Second

	server.SendTimeout = 10 * time.Second

	//SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal(err)
	}

	//New email simple html with inline and CC
	email := mail.NewMSG()

	email.SetFrom(config.Address).
		AddTo(emails).
		//AddCc("otherto@example.com").
		SetSubject(subject)

	email.SetBody(mail.TextHTML, content)
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// FormatEmailBody 格式化邮件内容
func FormatEmailBody(path string, data interface{}) string {
	builder := &strings.Builder{}
	funcs := map[string]interface{}{
		"dateformat": func(t *time.Time) string {
			if t == nil {
				return ""
			}
			return t.Format("2006-01-02 15:04:05")
		},
	}
	tpl := template.Must(template.New("tpl").Funcs(funcs).ParseFiles(path))
	err := tpl.ExecuteTemplate(builder, filepath.Base(path), data)
	if err != nil {
		return err.Error()
	}
	return builder.String()
}
