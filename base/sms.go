/**
 * @project Accounts
 * @filename sms.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/19 10:22
 * @version 1.0
 * @description
 * 短信发送
 */

package base

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/rs/zerolog"
)

// 阿里短信发送
func AlibabaSmsSend(mobile string, smsTplId string, title string, code string, logger *zerolog.Logger) {
	client, err := dysmsapi.NewClientWithAccessKey(GConf.AliSmsConfig.RegionId, GConf.AliSmsConfig.AccessId, GConf.AliSmsConfig.SecretKey)
	if err != nil {
		logger.Error().Msgf("ali sms new client error: %s", err.Error())
		return
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = mobile
	request.SignName = title
	request.TemplateCode = smsTplId
	request.TemplateParam = code
	request.OutId = mobile
	response, err := client.SendSms(request)
	if err != nil {
		logger.Error().Msgf("ali sms %s send to %s error: %s", smsTplId, mobile, err.Error())
		return
	}
	logger.Info().Msgf("ali sms %s send to %s result: %#v", smsTplId, mobile, response)
}
