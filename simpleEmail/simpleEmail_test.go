package simpleEmail

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
使用qq邮箱发送
邮箱	POP3服务器（端口995）	SMTP服务器（端口465或587）
qq.com	pop.qq.com			smtp.qq.com
*/

const (
	testName   = "12345678@qq.com" // 邮箱
	testPasswd = "11111"           // 授权码
)

func TestEmailClient_SendMessage(t *testing.T) {
	client := NewEmailClient("smtp.qq.com", testName, testPasswd, 587)
	_, err := client.SendMessage(EmailMessage{
		From:        "422134976@qq.com",
		To:          []string{"chenbin940404@163.com"},
		Subject:     "测试邮件",
		ContentType: "text/plain",
		Content:     "测试邮件内容",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmailClient_SendMessageMultiToSucc(t *testing.T) {
	client := NewEmailClient(
		"smtp.qq.com",
		testName,
		testPasswd,
		587,
		WithCcMaxNumber(2),
		WithToMaxNumber(2),
	)
	_, err := client.SendMessage(EmailMessage{
		From:        "422134976@qq.com",
		To:          []string{"chenbin940404@163.com", "chenbin940404@163.com"},
		Subject:     "测试邮件",
		ContentType: "text/plain",
		Content:     "测试邮件内容",
	})
	if err != nil {
		t.Fatal()
	}
}

func TestEmailClient_SendMessageMultiTo(t *testing.T) {
	client := NewEmailClient(
		"smtp.qq.com",
		testName,
		testPasswd,
		587,
		WithCcMaxNumber(2),
		WithToMaxNumber(2),
	)
	_, err := client.SendMessage(EmailMessage{
		From:        "422134976@qq.com",
		To:          []string{"chenbin940404@163.com", "chenbin940404@163.com", "chenbin940404@163.com"},
		Subject:     "测试邮件",
		ContentType: "text/plain",
		Content:     "测试邮件内容",
	})
	assert.ErrorIs(t, err, ErrorToTooMany)
}

func TestEmailClient_SendMessageMultiCc(t *testing.T) {
	client := NewEmailClient(
		"smtp.qq.com",
		testName,
		testPasswd,
		587,
		WithCcMaxNumber(2),
		WithToMaxNumber(2),
	)
	_, err := client.SendMessage(EmailMessage{
		From:        "422134976@qq.com",
		To:          []string{"chenbin940404@163.com"},
		Cc:          []string{"chenbin940404@163.com", "chenbin940404@163.com", "chenbin940404@163.com"},
		Subject:     "测试邮件",
		ContentType: "text/plain",
		Content:     "测试邮件内容",
	})
	assert.ErrorIs(t, err, ErrorCcTooMany)
}
