package utils

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"kubeops/config"
	"path/filepath"
	"strings"
	"time"
)

var EmailClient *mail.SMTPClient

// InitEmail 初始化邮箱服务
func InitEmail() {
	EmailServer := &mail.SMTPServer{
		Authentication: 0,
		Encryption:     0,
		Username:       config.Username,
		Password:       config.Password,
		Helo:           "",
		ConnectTimeout: 10 * time.Second,
		SendTimeout:    10 * time.Second,
		Host:           config.HostSmtp,
		Port:           config.Port,
		KeepAlive:      false,
		TLSConfig:      nil,
		CustomConn:     nil,
	}
	var err error
	EmailClient, err = EmailServer.Connect()
	if err != nil {
		Logger.Fatal("Failed to initialize the mailbox service" + err.Error())
	}
	Logger.Info("The mailbox service was initialized")
}

// Emails 发送邮件
func Emails(emails string, content string, subject string) (err error) {
	//server := mail.NewSMTPClient()
	////SMTP Server
	//server.Host = config.HostSmtp
	//server.Port = config.Port
	//server.Username = config.Username
	//server.Password = config.Password
	//
	//server.KeepAlive = false
	//
	//server.ConnectTimeout = 10 * time.Second
	//
	//server.SendTimeout = 10 * time.Second

	//SMTP client
	//smtpClient, err := EmailServer.Connect()

	//if err != nil {
	//	Logger.Fatal(err)
	//}

	//New email simple html with inline and CC
	email := mail.NewMSG()

	email.SetFrom(config.Address).
		AddTo(emails).
		//AddCc("otherto@example.com").
		SetSubject(subject)

	email.SetBody(mail.TextHTML, content)
	err = email.Send(EmailClient)
	if err != nil {
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
