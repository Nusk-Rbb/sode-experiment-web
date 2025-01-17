package util

import (
	"fmt"
	"io/ioutil"
	"net/smtp"

	"backend/models"

	"gopkg.in/yaml.v2"
)

// ---------------------------------------
// SMTPサーバ情報取得
// ---------------------------------------
func SmtpServerConfig() (models.SmtpServerConfig, error) {

	var smtpServerInfo []byte
	var err error

	smtpServerInfo, err = ioutil.ReadFile("config/smtpserver.yaml")
	if err != nil {
		return models.SmtpServerConfig{}, err
	}

	var smtpConfig models.SmtpServerConfig
	err = yaml.Unmarshal([]byte(smtpServerInfo), &smtpConfig)
	if err != nil {
		return models.SmtpServerConfig{}, err
	}

	return smtpConfig, nil
}

// ---------------------------------------
// メール送信
// param1. 送信するメールアドレス  string
// param2. 送信するメールの件名   string
// param3. 送信するメールの本文   string
// return: error
// ---------------------------------------

func SmtpSendMail(mailAddress string, mailSubject string, mailBody string) error {

	var smtpConn models.SmtpServerConfig
	var err error

	smtpConn, err = SmtpServerConfig()
	if err != nil {
		return err
	}

	// メールの内容を定義
	toMailAddress := []string{mailAddress}

	mailMessage := []byte("To: " + mailAddress + "\r\n" +
		"Subject: " + mailSubject + "\r\n" +
		"\r\n" +
		mailBody + "\r\n")

	fmt.Println("SMTPサーバ接続開始")
	// SMTPサーバ接続
	auth := smtp.PlainAuth("", smtpConn.AuthAddress, smtpConn.AuthPassword, smtpConn.SmtpServer)

	// メール送信
	err = smtp.SendMail(smtpConn.SmtpServer+":"+smtpConn.SmtpPort, auth, smtpConn.AuthAddress, toMailAddress, mailMessage)
	if err != nil {
		return err
	}

	return nil
}
