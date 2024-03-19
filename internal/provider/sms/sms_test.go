package sms

import (
	"testing"
)

func TestSms(t *testing.T) {
	sms, err := NewSms(&SmsConfig{
		Key:          "",           // SDK key， aliyun provided
		Secret:       "",           // SDK key secret, aliyun provided
		TemplateCode: "",           // SMS template code
		SignName:     "",           // SMS template signature
		Phones:       []string{""}, // phone numbers
	})
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"", ""}
	_ = sms.Send(args...)
}
