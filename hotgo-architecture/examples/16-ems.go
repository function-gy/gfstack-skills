// ================================================================
// 示例: Email 邮件发送 (internal/library/ems/ems.go)
// 基于 SMTP 带 TLS + 超时机制的邮件发送
// ================================================================

package ems

import (
	"crypto/tls"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"

	"hotgo/internal/model"
	"hotgo/utility/validate"
)

// Send 发送邮件入口
func Send(config *model.EmailConfig, to string, subject string, body string) error {
	return sendToMail(config, to, subject, body, "html")
}

func sendToMail(config *model.EmailConfig, to, subject, body, mailType string) error {
	if config == nil {
		return gerror.New("邮件配置不能为空")
	}

	var (
		contentType string
		sendTo      = strings.Split(to, ";")
	)

	if len(sendTo) == 0 {
		return gerror.New("收件人不能为空")
	}

	for _, em := range sendTo {
		if !validate.IsEmail(em) {
			return gerror.Newf("邮件格式不正确，请检查：%v", em)
		}
	}

	if mailType == "html" {
		contentType = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + config.SendName + "<" + config.User + ">" +
		"\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)

	return sendMailWithTimeout(config, sendTo, msg)
}

// sendMailWithTimeout 带10秒超时的SMTP发送
func sendMailWithTimeout(config *model.EmailConfig, to []string, msg []byte) error {
	timeout := 10 * time.Second

	conn, err := net.DialTimeout("tcp", config.Addr, timeout)
	if err != nil {
		return gerror.Wrapf(err, "无法连接到SMTP服务器 %s", config.Addr)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(timeout))

	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		return gerror.Wrapf(err, "创建SMTP客户端失败，主机名：%s", config.Host)
	}
	defer client.Close()

	// STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: config.Host, InsecureSkipVerify: false}
		if err = client.StartTLS(tlsConfig); err != nil {
			return gerror.Wrap(err, "启动TLS加密失败")
		}
	}

	// 认证
	auth := smtp.PlainAuth("", config.User, config.Password, config.Host)
	if err = client.Auth(auth); err != nil {
		return gerror.Wrapf(err, "SMTP身份验证失败，用户名：%s", config.User)
	}

	// 发件人
	if err = client.Mail(config.User); err != nil {
		return gerror.Wrapf(err, "设置发件人失败：%s", config.User)
	}

	// 收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return gerror.Wrapf(err, "设置收件人失败：%s", addr)
		}
	}

	// 邮件内容
	w, err := client.Data()
	if err != nil {
		return gerror.Wrap(err, "准备发送邮件内容失败")
	}
	if _, err = w.Write(msg); err != nil {
		w.Close()
		return gerror.Wrap(err, "写入邮件内容失败")
	}
	w.Close()

	return client.Quit()
}
