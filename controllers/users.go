/**
 * @project Accounts
 * @filename users.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/19 10:02
 * @version 1.0
 * @description
 * 登录、注册等
 */

package controllers

import (
	"accounts/base"
	"accounts/models"
	"fmt"
	"github.com/rs/zerolog/hlog"
	"net/http"
	"regexp"
	"strings"
)

// Register 注册、登录
// 当账号不存在时则注册，适合于游客、第三方账号
// 邮箱和手机号支持通过验证码注册、登录
func Register(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	userLog := hlog.FromRequest(req)
	ip := base.GetRealAddr(req).String()
	logHook := base.RequestHook{IP: ip}
	data := &base.RegisterFields{}
	err := base.RequestHandler(req, data)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(logHook))
		return
	}
	logHook.RequestBody = data
	logHook.GameId = data.GameId
	logHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(logHook))
		return
	}

	//检查账号格式
	data.Account, err = base.CheckUserAccountFormat(data.Account, data.Type)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(logHook))
		return
	}

	//邮箱或手机号注册时，验证码和密码不允许都是空;
	if (data.Type == base.AccountEmail || data.Type == base.AccountMobile) && data.Code == base.DefaultNoValue && data.Password == base.DefaultNoValue {
		base.ResponseFail(resp, &base.MyError{Code: base.RegisterCodeAndPasswordEmpty}, userLog.Hook(logHook))
		return
	}

	//检查ip是否达到限制数量
	err = base.LimitRegister(ip)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(logHook))
		return
	}

	//检查验证码
	var codeKey string
	if (data.Type == base.AccountMobile || data.Type == base.AccountEmail) && data.Code != base.DefaultNoValue {
		codeKey = fmt.Sprintf(base.CodeFormat, base.CodeTypeRegister, data.Account)
		err = models.CheckVerifyCode(codeKey, data.Code)
		if err != nil {
			base.ResponseFail(resp, err, userLog.Hook(logHook))
			return
		}
	}

	ret, err := models.AccountRegister(data, ip)
	if err.Code != base.RegisterSuccess && err.Code != base.LoginSuccess {
		base.ResponseFail(resp, err, userLog.Hook(logHook))
		return
	}

	dataLogId := "1314520"
	dataLog := base.LoginDataLog
	//注册成功后累加此ip的数量
	if err.Code == base.RegisterSuccess {
		dataLog = base.RegisterDataLog
		dataLogId = "1314521"
		err = base.LimitRegisterIncr(ip)
		if err != nil {
			userLog.Err(err).Msg("limit register incr error")
		}
	}

	//删除已使用验证码
	if (data.Type == base.AccountMobile || data.Type == base.AccountEmail) && data.Code != base.DefaultNoValue {
		models.DeleteVerifyCode(codeKey)
	}

	//数据写入
	base.DataExtLog(ret.Uid, data.DataExt, dataLogId, ip, dataLog)

	//成功返回
	base.ResponseOK(resp, ret, userLog.Hook(logHook))
	return
}

// Login 登录
func Login(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	userLog := hlog.FromRequest(req)
	ip := base.GetRealAddr(req).String()
	logHook := base.RequestHook{IP: ip}
	data := &base.LoginFields{}
	err := base.RequestHandler(req, data)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(logHook))
		return
	}
	logHook.RequestBody = data
	logHook.GameId = data.GameId
	logHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(logHook))
		return
	}

	//检查账号格式
	data.Account, err = base.CheckUserAccountFormat(data.Account, data.Type)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(logHook))
		return
	}

	//邮箱或手机号登录时，验证码和密码不允许都是空;
	if (data.Type == base.AccountEmail || data.Type == base.AccountMobile) && data.Code == base.DefaultNoValue && data.Password == base.DefaultNoValue {
		base.ResponseFail(resp, &base.MyError{Code: base.LoginCodeAndPasswordEmpty}, userLog.Hook(logHook))
		return
	}

	err = base.LimitLogin(ip)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(logHook))
		return
	}

	//检查验证码
	var codeKey string
	if (data.Type == base.AccountMobile || data.Type == base.AccountEmail) && data.Code != base.DefaultNoValue {
		codeKey = fmt.Sprintf(base.CodeFormat, base.CodeTypeLogin, data.Account)
		err = models.CheckVerifyCode(codeKey, data.Code)
		if err != nil {
			base.ResponseFail(resp, err, userLog.Hook(logHook))
			return
		}
	}

	accountUid := models.GetAccountUid(data.Account)
	if accountUid == 0 {
		err = base.LimitLoginIncr(ip, data.Account)
		if err != nil {
			userLog.Err(err).Msg("limit login incr error")
		}
		//此处的错误码为账号或密码错误，而不是账号不存在，防止利用登录探测账号是否存在
		base.ResponseFail(resp, &base.MyError{Code: base.LoginUserOrPasswordError}, userLog.Hook(logHook))
		return
	}

	ret, err := models.AccountLogin(data, accountUid)
	if err.Code != base.LoginSuccess {
		limitErr := base.LimitLoginIncr(ip, data.Account)
		if limitErr != nil {
			userLog.Err(limitErr).Msg("limit login incr error")
		}
		base.ResponseFail(resp, err, userLog.Hook(logHook))
		return
	}

	//删除已使用验证码
	if (data.Type == base.AccountMobile || data.Type == base.AccountEmail) && data.Code != base.DefaultNoValue {
		models.DeleteVerifyCode(codeKey)
	}

	//数据写入
	base.DataExtLog(ret.Uid, data.DataExt, "1314520", ip, base.LoginDataLog)

	//成功返回
	base.ResponseOK(resp, ret, userLog.Hook(logHook))
	return
}

// ForgetPassword 忘记密码-重置密码
// 仅支持注册账号或是被绑定的email或手机号
func ForgetPassword(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)

	data := &base.ForgetPasswordFields{}
	err := base.RequestHandler(req, data)

	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查账号格式
	data.Account, err = base.CheckUserAccountFormat(data.Account, data.Type)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查验证码
	codeKey := fmt.Sprintf(base.CodeFormat, base.CodeTypeForgetPassword, data.Account)
	err = models.CheckVerifyCode(codeKey, data.Code)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	err = models.ForgetPassword(data)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//删除已使用验证码
	err = models.DeleteVerifyCode(codeKey)
	base.ResponseOK(resp, base.EmptyData, userLog.Hook(requestHook))
	return
}

// ChangePassword 修改密码
// 仅用于 username+password, email+password, mobile+password
// 不支持 email+验证码、 mobile+验证码方式注册的修改密码，因为这些方式密码是空的，无法匹配到旧密码
// 不支持 游客、第三方
func ChangePassword(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)

	data := &base.ChangePasswordFields{}
	err := base.RequestHandler(req, data)
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.Uid = data.Uid
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查登录token内的uid与当前绑定的uid是否一致
	err = base.LoginTokenCheck(data.LoginToken, data.Uid)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	err = models.ChangePassword(data)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	base.ResponseOK(resp, base.EmptyData, userLog.Hook(requestHook))
	return
}

// SendVerifyCode 发送验证码, account 只能为邮箱或手机号
// lang_id: zh-TW || zh-CN || en
func SendVerifyCode(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	ip := base.GetRealAddr(req).String()
	requestHook := base.RequestHook{IP: ip}
	userLog := hlog.FromRequest(req)
	data := &base.VerifyCodeFields{}
	err := base.RequestHandler(req, data)
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	if _, ok := base.CodeTypeMap[data.CodeType]; !ok {
		base.ResponseFail(resp, &base.MyError{Code: base.VerifyCodeTypeUnknown}, userLog.Hook(requestHook))
		return
	}

	//检查账号格式
	data.Account, err = base.CheckUserAccountFormat(data.Account, data.Type)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	// 发送验证码前检查是否达到限制数量
	err = base.LimitVerifyCode(ip, data.Account)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	codeRand := base.GetRandom(6)
	codeKey := fmt.Sprintf(base.CodeFormat, data.CodeType, data.Account)
	//账号类型为 邮件方式
	if data.Type == base.AccountEmail {
		tplConfig, err := models.MailTplConfig(data)
		if err != nil {
			base.ResponseFail(resp, err, userLog.Hook(requestHook))
			return
		}
		err = models.SetVerifyCode(codeKey, codeRand)
		if err != nil {
			base.ResponseFail(resp, err, userLog.Hook(requestHook))
			return
		}

		tplConfig.Content = strings.Replace(tplConfig.Content, "{CODE}", codeRand, -1)
		fmt.Printf("tplContent: %+v\n", tplConfig)
		go base.SendMail(&base.GConf.MailConfig, tplConfig, data.Account, userLog) //单元测试时，协程去掉
		err = base.LimitVerifyCodeIncr(ip, data.Account)
		if err != nil {
			userLog.Err(err).Msg("limit verify code incr error")
		}
	}
	//短信方式
	if data.Type == base.AccountMobile {
		smsConfig, err := models.SmsTplConfig(data)
		if err != nil {
			base.ResponseFail(resp, err, userLog.Hook(requestHook))
			return
		}
		if smsConfig == nil {
			base.ResponseFail(resp, &base.MyError{Code: base.SmsTplConfigEmpty}, userLog.Hook(requestHook))
			return
		}

		err = models.SetVerifyCode(codeKey, codeRand)
		if err != nil {
			base.ResponseFail(resp, err, userLog.Hook(requestHook))
			return
		}

		code := fmt.Sprintf("{\"code\":\"%s\"}", codeRand)
		go base.AlibabaSmsSend(data.Account, smsConfig.SmsId, smsConfig.Title, code, userLog) //单元测试时，协程去掉
		if err != nil {
			userLog.Err(err).Msg("limit verify code incr error")
		}
	}

	base.ResponseOK(resp, base.EmptyData, userLog.Hook(requestHook))
	return
}

// BindAccount 绑定账号
// 支持绑邮箱、手机号、第三方（可以多个）
func BindAccount(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)

	data := &base.BindAccountFields{}
	err := base.RequestHandler(req, data)
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.Uid = data.Uid
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查账号格式
	data.BindAccount, err = base.CheckUserAccountFormat(data.BindAccount, data.Type)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查登录token内的uid与当前绑定的uid是否一致
	err = base.LoginTokenCheck(data.LoginToken, data.Uid)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查验证码
	var codeKey string
	if data.Type != base.AccountThird {
		codeKey = fmt.Sprintf(base.CodeFormat, base.CodeTypeBindAccount, data.BindAccount)
		err = models.CheckVerifyCode(codeKey, data.Code)
		if err != nil {
			base.ResponseFail(resp, err, userLog.Hook(requestHook))
			return
		}
	}

	err = models.BindAccount(data)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//删除已使用验证码
	if data.Type != base.AccountThird {
		models.DeleteVerifyCode(codeKey)
	}

	ret, err := models.GetAccountBindsInfo(data.Account, data.GameId, data.PlatformId)
	if err != nil {
		userLog.Info().Err(err).Msg("get already bind error")
	}
	base.ResponseOK(resp, ret, userLog.Hook(requestHook))
	return
}

// UnBindAccount 解绑账号
// 可以解绑邮箱、手机号、第三方
// 不支持解绑注册时为邮箱、手机号的类型
// 不支持解绑当前账号
func UnBindAccount(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)

	data := &base.UnBindAccountFields{}
	err := base.RequestHandler(req, data)
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.Uid = data.Uid
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查账号格式
	data.UnBindAccount, err = base.CheckUserAccountFormat(data.UnBindAccount, data.Type)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	//不能解绑当前账号
	if data.Account == data.UnBindAccount {
		base.ResponseFail(resp, &base.MyError{Code: base.UnBindUnSupportCurrentAccount}, userLog.Hook(requestHook))
		return
	}

	//检查登录token内的uid与当前绑定的uid是否一致
	err = base.LoginTokenCheck(data.LoginToken, data.Uid)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查验证码
	var codeKey string
	if data.Type != base.AccountThird {
		codeKey = fmt.Sprintf(base.CodeFormat, base.CodeTypeUnBindAccount, data.UnBindAccount)
		err = models.CheckVerifyCode(codeKey, data.Code)
		if err != nil {
			base.ResponseFail(resp, err, userLog.Hook(requestHook))
			return
		}
	}

	binds, err := models.UnBindAccount(data)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//删除已使用验证码
	if data.Type != base.AccountThird {
		_ = models.DeleteVerifyCode(codeKey)
	}

	base.ResponseOK(resp, binds, userLog.Hook(requestHook))
	return
}

// LoginAuth 服务器登录校验
func LoginAuth(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)
	//白名单
	if base.GConf.RequestLimitRule.Enabled && len(base.GConf.RequestLimitRule.LoginAuthWhiteListMap) > 0 {
		domain := req.URL.Host
		if strings.Index(domain, ":") != -1 {
			domains := strings.Split(domain, ":")
			domain = domains[0]
		}
		if _, ok := base.GConf.RequestLimitRule.LoginAuthWhiteListMap[domain]; !ok {
			//fmt.Printf("request domain: %s, host:%s\n", domain, req.URL.Host)
			base.ResponseFail(resp, &base.MyError{Code: base.LoginAuthDomainNotWhiteList, Log: "request domain:" + domain}, userLog.Hook(requestHook))
			return
		}
	}
	data := &base.LoginAuthFields{}
	err := base.RequestHandler(req, data)
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.Uid = data.Uid
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeClient)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查登录token内的uid与当前绑定的uid是否一致
	err = base.LoginTokenCheck(data.Token, data.Uid)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	base.ResponseOK(resp, base.EmptyData, userLog.Hook(requestHook))
	return
}

// WhiteList 白名单校验，客户端直接使用
func WhiteList(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)

	data := &base.WhiteListFields{}
	err := base.RequestHandler(req, data)
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.Uid = data.Uid
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeServer)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	result := models.GetWhiteList(data, requestHook.IP)
	base.ResponseOK(resp, result, userLog.Hook(requestHook))
	return
}

// ApplyLogout 账号注销
func ApplyLogout(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)

	data := &base.LogoutAccountFields{}
	err := base.RequestHandler(req, data)
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.Uid = data.Uid
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查登录token内的uid与当前绑定的uid是否一致
	err = base.LoginTokenCheck(data.Token, data.Uid)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	err = models.AddDeleteApply(data, userLog)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	base.ResponseOK(resp, base.EmptyData, userLog.Hook(requestHook))
	return
}

// UndoLogout 撤销账号注销
// 由用户申请撤销、管理后台审批通过后，更改其状态
func UndoLogout(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)

	data := &base.UndoLogoutFields{}
	err := base.RequestHandler(req, data)
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.Uid = data.Uid
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查登录token内的uid与当前绑定的uid是否一致
	err = base.LoginTokenCheck(data.Token, data.Uid)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	err = models.UndoDeleteApply(data)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	base.ResponseOK(resp, base.EmptyData, userLog.Hook(requestHook))
	return
}

// GetUserInfo ...
func GetUserInfo(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)

	data := &base.UserInfoReqFields{}
	err := base.RequestHandler(req, data)
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.Uid = data.Uid
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查登录token内的uid与当前绑定的uid是否一致
	err = base.LoginTokenCheck(data.Token, data.Uid)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	ret, err := models.GetAccountBindsInfo(data.Account, data.GameId, data.PlatformId)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	base.ResponseOK(resp, ret, userLog.Hook(requestHook))
}

// RealNameAuth 实名认证
// 实名认证根据国家不同而认证接口不同，防沉迷也类似，政策各不相同
func RealNameAuth(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("StartTime", base.GetUnixMilliString())
	requestHook := base.RequestHook{IP: base.GetRealAddr(req).String()}
	userLog := hlog.FromRequest(req)

	data := &base.UserRealNameAuthReqFields{}
	err := base.RequestHandler(req, data)
	requestHook.RequestBody = data
	requestHook.GameId = data.GameId
	requestHook.Uid = data.Uid
	requestHook.HeaderGamePlatform = req.Header.Get(base.HeaderGamePlatform)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	userLog.Info().Interface("req_body", data).Msg("")

	err = base.SignValidator(data.AppId, data.Sign, data.GameId, data, base.AppIdTypeSdk)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}

	//检查登录token内的uid与当前绑定的uid是否一致
	err = base.LoginTokenCheck(data.Token, data.Uid)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	//检查姓名
	reg := regexp.MustCompile(`^[\x{4e00}-\x{9fa5}]{2,5}$`)
	if !reg.MatchString(data.Name) {
		base.ResponseFail(resp, &base.MyError{Code: base.RealNameNameError}, userLog.Hook(requestHook))
		return
	}
	//检查身份证
	if !base.CheckCardId(data.CardId) {
		base.ResponseFail(resp, &base.MyError{Code: base.RealNameCardIdError}, userLog.Hook(requestHook))
		return
	}
	err = models.RealNameAuth(data)
	if err != nil {
		base.ResponseFail(resp, err, userLog.Hook(requestHook))
		return
	}
	base.ResponseOK(resp, base.EmptyData, userLog.Hook(requestHook))
}
