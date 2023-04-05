/**
 * @project Accounts
 * @filename routers.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) 2023/2/19
 * @datetime 2023/2/17 10:12
 * @version 1.0
 * @description
 * 路由
 */

package routers

import (
	"accounts/base"
	"accounts/controllers"
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// InitRouterService 初始化路由
func InitRouterService() {
	//日志格式化
	mid := base.Middleware{}
	mid = mid.Append(hlog.NewHandler(log.Logger)).Append(hlog.RequestHandler("request")).Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	http.Handle("/user/register", mid.Then(http.HandlerFunc(controllers.Register)))             //注册、登录
	http.Handle("/user/login", mid.Then(http.HandlerFunc(controllers.Login)))                   //登录, 相对于Register接口区别在于 在用户不存在的情况下，不会注册，上面接口适用于游客、第三方
	http.Handle("/user/forgetPassword", mid.Then(http.HandlerFunc(controllers.ForgetPassword))) //忘记密码
	http.Handle("/user/changePassword", mid.Then(http.HandlerFunc(controllers.ChangePassword))) //修改密码
	http.Handle("/user/sendSmsCode", mid.Then(http.HandlerFunc(controllers.SendVerifyCode)))    //发送验证码
	http.Handle("/user/bindAccount", mid.Then(http.HandlerFunc(controllers.BindAccount)))       //绑定账号
	http.Handle("/user/unBindAccount", mid.Then(http.HandlerFunc(controllers.UnBindAccount)))   //解绑
	http.Handle("/user/loginAuth", mid.Then(http.HandlerFunc(controllers.LoginAuth)))           //服务器登录校验
	http.Handle("/user/whiteList", mid.Then(http.HandlerFunc(controllers.WhiteList)))           //白名单校验
	http.Handle("/user/applyLogout", mid.Then(http.HandlerFunc(controllers.ApplyLogout)))       //账号注销
	http.Handle("/user/undoLogout", mid.Then(http.HandlerFunc(controllers.UndoLogout)))         //撤销账号注销
	http.Handle("/user/getUserInfo", mid.Then(http.HandlerFunc(controllers.GetUserInfo)))       //服务器登录校验
	http.Handle("/user/heartbeat", mid.Then(http.HandlerFunc(controllers.Heartbeat)))           //心跳监测
	http.Handle("/user/realNameAuth", mid.Then(http.HandlerFunc(controllers.RealNameAuth)))     //实名认证
	http.Handle("/captcha/image", mid.Then(http.HandlerFunc(controllers.Image)))                //图形验证码展示
	http.Handle("/captcha/verify", mid.Then(http.HandlerFunc(controllers.Verify)))              //图形验证码验证

	srv := &http.Server{
		Addr:         base.GConf.Server.Host,
		ReadTimeout:  time.Duration(base.GConf.HttpTimeout.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(base.GConf.HttpTimeout.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(base.GConf.HttpTimeout.IdleTimeout) * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		base.MultipleLog.Fatal().Msg(err.Error())
	}
}
