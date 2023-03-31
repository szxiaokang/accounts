/**
 * @project AP2
 * @filename captcha.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/3/29 12:02
 * @version 1.0
 * @description
 * 图形验证码， 用于拦截恶意请求
 */

package controllers

import (
	"accounts/base"
	"github.com/dchest/captcha"
	"github.com/rs/zerolog/hlog"
	"net/http"
)

// 验证码显示、刷新
func Image(resp http.ResponseWriter, req *http.Request) {
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)
	data := &base.CaptchaImageFields{}
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	err := base.RequestHandler(req, data)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	if data.Width == -1 {
		data.Width = captcha.StdWidth
	}
	if data.Height == -1 {
		data.Height = captcha.StdHeight
	}
	if data.Refresh == 1 {
		captcha.Reload(data.Id)
	}

	resp.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	resp.Header().Set("Content-Type", "image/png")
	captcha.WriteImage(resp, data.Id, data.Width, data.Height)
}

// 验证，通过则删除限制
func Verify(resp http.ResponseWriter, req *http.Request) {
	ip := base.GetRealAddr(req).String()
	requestHook := base.RequestHook{IP: ip}
	userLog := hlog.FromRequest(req)
	data := &base.CaptchaVerifyFields{}
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	err := base.RequestHandler(req, data)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	ok := captcha.VerifyString(data.CaptchaId, data.CaptchaCode)
	if !ok {
		base.ResponseFail(resp, &base.MyError{Code: base.RequestLimitCodeError}, userLog.Hook(requestHook))
		return
	}

	locKey := data.CaptchaType + ip
	ipKey := base.LimitIpKey + ip
	base.RedisClient.Del(locKey)
	base.RedisClient.Del(ipKey)
	base.ResponseOK(resp, base.EmptyData, userLog.Hook(requestHook))
}
