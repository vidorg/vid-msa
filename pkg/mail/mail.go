package mail

import (
	"github.com/spf13/viper"

	"gopkg.in/gomail.v2"
)

// SendMail 发送验证码
func SendMail(code string, recipient string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("smtp.name"))
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "【vid-msa】验证码")
	m.AddAlternative("text/plain", "您好：\n\n"+code+"\n这是您vid-msa的验证码，有效时间30分钟。\n")
	d := gomail.NewDialer(viper.GetString("smtp.host"), viper.GetInt("smtp.port"),
		viper.GetString("smtp.username"),
		viper.GetString("smtp.password"))
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
