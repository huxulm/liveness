// This file is auto-generated, don't edit it. Thanks.
package sms

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/huxulm/liveness/internal/provider"
	"github.com/sirupsen/logrus"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(key, secret string) (_result *dysmsapi20170525.Client, _err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
	}
	if len(*config.AccessKeyId) == 0 {
		config.AccessKeyId = &key
	}
	if len(*config.AccessKeySecret) == 0 {
		config.AccessKeySecret = &secret
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func (p *smsProvider) sendMsg(req *dysmsapi20170525.SendSmsRequest) (_err error) {

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err = p.client.SendSmsWithOptions(req, runtime)
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}

type smsProvider struct {
	client *dysmsapi20170525.Client
	config *SmsConfig
}

type SmsConfig struct {
	Type         string   `yaml:"type"`
	Provider     string   `yaml:"provider"`
	Key          string   `yaml:"key"`
	Secret       string   `yaml:"secret"`
	TemplateCode string   `yaml:"template"`
	SignName     string   `yaml:"sign_name"`
	Phones       []string `yaml:"phones"`
}

type MultiErr struct {
	errs []error
}

func (me *MultiErr) Add(err error) {
	me.errs = append(me.errs, err)
}

func (me *MultiErr) Error() string {
	a := []string{}
	for _, err := range me.errs {
		a = append(a, err.Error())
	}
	return strings.Join(a, ", ")
}

func (me *MultiErr) Empty() bool {
	return len(me.errs) == 0
}

func (p *smsProvider) textify(s string) string {
	t := strings.ReplaceAll(s, "https://", "")
	t = strings.ReplaceAll(t, ".", "[.]")
	t = strings.ReplaceAll(t, "/", "[/]")
	return t
}

func (p *smsProvider) Send(args ...string) error {
	var errs MultiErr
	for _, phone := range p.config.Phones {
		req := &dysmsapi20170525.SendSmsRequest{
			SignName:      tea.String(p.config.SignName),
			TemplateCode:  tea.String(p.config.TemplateCode),
			PhoneNumbers:  tea.String(phone),
			TemplateParam: tea.String(fmt.Sprintf(`{"name": "%s", "time": "%s"}`, p.textify(args[0]), args[1])),
		}
		if err := p.sendMsg(req); err != nil {
			logrus.Error(err)
			errs.Add(err)
		}
	}
	if errs.Empty() {
		return nil
	}
	return &errs
}

// NewSms create a sms provider
func NewSms(c *SmsConfig) (provider.Provider, error) {
	client, err := CreateClient(c.Key, c.Secret)
	if err != nil {
		return nil, err
	}
	return &smsProvider{client: client, config: c}, nil
}
