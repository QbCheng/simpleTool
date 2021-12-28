package simpleEmail

import (
	"crypto/tls"
	"errors"
	"gopkg.in/gomail.v2"
)

// EmailMessage 返回消息对象
// from: 发件人
// subject: 标题
// contentType: 内容的类型 text/plain text/html
// attach: 附件
// to: 收件人
// cc: 抄送人
type EmailMessage struct {
	From        string
	To          []string
	Cc          []string
	Subject     string
	ContentType string
	Content     string
	Attach      string
}

var (
	ErrorToTooMany = errors.New(" Too many email recipients ")
	ErrorCcTooMany = errors.New(" Too many email cc ")
)

const (
	defaultToMaxNumber = 10 // 默认最大同时收件人为10.
	defaultCcMaxNumber = 10 // 默认最大同时抄收为10.
)

type Option func(log *EmailClient)

// WithToMaxNumber SMTP 服务一份邮件最大接收人数量. 需要根据具体 SMTP 服务 来进行设置
func WithToMaxNumber(toMaxNumber int) Option {
	return func(emailClient *EmailClient) {
		emailClient.toMaxNumber = toMaxNumber
	}
}

// WithCcMaxNumber SMTP 服务一份邮件最大抄手人数量. 需要根据具体 SMTP 服务 来进行设置
func WithCcMaxNumber(ccMaxNumber int) Option {
	return func(emailClient *EmailClient) {
		emailClient.ccMaxNumber = ccMaxNumber
	}
}

// WithSkipTls 表示跳过 tls 效验. 生产环境下不允许使用
func WithSkipTls(tlsConfig *tls.Config) Option {
	return func(emailClient *EmailClient) {
		emailClient.tls = tlsConfig
	}
}

// EmailClient 发送客户端
type EmailClient struct {
	host     string
	port     int
	username string
	password string

	tls         *tls.Config
	ccMaxNumber int // 邮件最大同时接收人
	toMaxNumber int // 邮件最大同时抄收人

}

// NewEmailClient 返回一个邮件客户端
// host smtp地址
// username 用户名
// password 密码(或 授权码)
// port 端口,  常见端口 465 ssl 或 587 非ssl
func NewEmailClient(host, username, password string, port int, options ...Option) *EmailClient {
	ret := &EmailClient{
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		ccMaxNumber: defaultCcMaxNumber,
		toMaxNumber: defaultToMaxNumber,
		tls:         &tls.Config{InsecureSkipVerify: true}, // 默认跳过 tls 效验.
	}

	for i := range options {
		options[i](ret)
	}

	return ret
}

// SendMessage 发送邮件
func (c *EmailClient) SendMessage(email EmailMessage) (bool, error) {

	if len(email.To) > c.toMaxNumber {
		return false, ErrorToTooMany
	}

	if len(email.Cc) > c.ccMaxNumber {
		return false, ErrorCcTooMany
	}

	e := gomail.NewDialer(c.host, c.port, c.username, c.password)
	dm := gomail.NewMessage()
	dm.SetHeader("From", email.From)
	dm.SetHeader("To", email.To...)

	if len(email.Cc) != 0 {
		dm.SetHeader("Cc", email.Cc...)
	}

	dm.SetHeader("Subject", email.Subject)
	dm.SetBody(email.ContentType, email.Content)

	if email.Attach != "" {
		dm.Attach(email.Attach)
	}

	if err := e.DialAndSend(dm); err != nil {
		return false, err
	}
	return true, nil
}
