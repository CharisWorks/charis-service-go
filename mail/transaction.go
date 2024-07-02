package mail

import (
	"log"
	"net/smtp"

	"github.com/charisworks/charisworks-service-go/util"
)

func SendEmail(to string, subject string, body string) error {
	authmail := util.AUTH_EMAIL
	from := util.MAIL_FORM
	password := util.MAIL_AUTH_PASS
	smtp_server_addr := util.SMTP_SERVER_ADDR
	smtp_server := util.SMTP_SERVER

	// メールヘッダーの作成
	header := make(map[string]string)
	header["From"] = "CharisWorks本部" + " <" + from + ">"
	header["To"] = to
	header["Subject"] = subject
	log.Printf(`
**********************************************************************************************
Sending Email...
from: %s
to: %s
subject: %s
**********************************************************************************************
`,
		header["From"],
		header["To"],
		header["Subject"],
	)
	// メール本文の作成
	message := ""
	for key, value := range header {
		message += key + ": " + value + "\r\n"
	}
	message += "\r\n" + body
	err := smtp.SendMail(smtp_server_addr,
		smtp.PlainAuth("", authmail, password, smtp_server),
		from, []string{to}, []byte(message))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Email sent successfully!")
	}
	return nil
}
