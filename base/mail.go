/**
 * @project Accounts
 * @filename mail.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/10 10:22
 * @version 1.0
 * @description
 * 邮件发送
 */

package base

import (
	"github.com/rs/zerolog"
	"gopkg.in/gomail.v2"
)

// 发送邮件
func SendMail(mailConfig *MailConfig, mailTpl *MailTpl, toEmail string, logger *zerolog.Logger) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", mailConfig.Username)
	mail.SetHeader("To", toEmail)
	mail.SetHeader("Subject", mailTpl.Title)
	mail.SetBody("text/html", mailTpl.Content)
	ret := gomail.NewDialer(mailConfig.Hostname, mailConfig.Port, mailConfig.Username, mailConfig.Password)
	err := ret.DialAndSend(mail)

	if err != nil {
		logger.Error().Msgf("send %s to %s mail error: %s", mailTpl.Title, toEmail, err.Error())
		return
	}
	logger.Info().Msgf("send %s to %s mail success", mailTpl.Title, toEmail)
}
