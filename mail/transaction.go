package mail

import (
	"encoding/json"
	"log"
	"net/smtp"
	"os"

	"github.com/charisworks/charisworks-service-go/util"
)

func SendEmail(to string, subject string, body string) error {
	data, err := os.ReadFile("./email_credentials.json")
	if err != nil {
		log.Fatalf("JSONファイルの読み込みに失敗しました：%v", err)
		return err
	}
	emailCredentials := make(map[string]string)
	err = json.Unmarshal(data, &emailCredentials)
	if err != nil {
		log.Fatalf("JSONデータの解析に失敗しました：%v", err)
		return err
	}
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

	// メール本文の作成
	message := ""
	for key, value := range header {
		message += key + ": " + value + "\r\n"
	}
	message += "\r\n" + body
	err = smtp.SendMail(smtp_server_addr,
		smtp.PlainAuth("", authmail, password, smtp_server),
		from, []string{to}, []byte(message))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Email sent successfully!")
	}
	return nil
}
