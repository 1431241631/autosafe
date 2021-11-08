package core

import (
	"fmt"
	"mime"
	"strconv"

	"gopkg.in/gomail.v2"
)

// 发送邮件
func sendMail(mailTo string, subject string, body string, authCode string, host string, port_ int) {
	mailConn := map[string]string{
		"username": mailTo,
		"authCode": authCode,
		"host":     host,
		"port":     fmt.Sprint(port_),
	}
	port, _ := strconv.Atoi(mailConn["port"])
	m := gomail.NewMessage()
	m.SetHeader("From", mime.QEncoding.Encode("UTF-8", "Support")+"<"+mailConn["username"]+">")
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(mailConn["host"], port, mailConn["username"], mailConn["authCode"])
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println("To:", mailTo, "##", "Send Email Failed!Err:", err)
	} else {
		fmt.Println("To:", mailTo, "##", "Send Email Successfully!")
	}

}
