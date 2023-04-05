/**
 * @project main.go
 * @filename users_test.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) 2022/7/25 kangyun@outlook.com
 * @datetime 2023/03/25 11:02
 * @version 1.0
 * @description
 *
 */

package controllers

import (
	"accounts/base"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	AppId      = 1000000001
	SecretKey  = "d1657859c5a0f6131ac22702f708fe2e"
	GameId     = 16
	PlatformId = 1
)

func TestRegister(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	//1: email, 2:手机号, 3: 用户名, 4: 游客, 5: 第三方

	pList := []base.RegisterFields{
		//Email + 密码
		{
			Account:    "kangyun@outlook.com",
			Code:       "-1",
			Password:   base.Md5Sum([]byte("ILoveU")),
			DeviceType: 2,
			Type:       1,
			Lang:       "zh-cn",
			ChannelId:  -1,
			DataExt:    "{}",
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
		/*
			//email+验证码
			{
				Account:    "a123456@outlook.com",
				Code:       "123456",
				Password:   "-1",
				DeviceType: 2,
				Type:       1,
				Lang:       "zh-cn",
				ChannelId:  -1,
				DataExt:    "{}",
				CommonFields: base.CommonFields{
					GameId:     GameId,
					PlatformId: PlatformId,
					AppId:      AppId,
				},
			},*/
		//手机号+验证码，密码不为默认值则写入
		/*		{
				Account:    "13838384388",
				Code:       "123688",
				Password:   password,
				DeviceType: 2,
				Type:       2,
				Lang:       "zh-cn",
				ChannelId:  -1,
				DataExt:    "{}",
				CommonFields: base.CommonFields{
					GameId:     GameId,
					PlatformId: PlatformId,
					AppId:      AppId,
				},
			},*/
		//手机号+密码登录
		/*		{
				Account:    "13838384388",
				Code:       "-1",
				Password:   password,
				DeviceType: 2,
				Type:       2,
				Lang:       "zh-cn",
				ChannelId:  -1,
				DataExt:    "{}",
				CommonFields: base.CommonFields{
					GameId:     GameId,
					PlatformId: PlatformId,
					AppId:      AppId,
				},
			},*/
		//用户名
		/*		{
				Account:    "kangkang",
				Code:       "-1",
				Password:   password,
				DeviceType: 2,
				Type:       3,
				Lang:       "zh-cn",
				ChannelId:  -1,
				DataExt:    "{}",
				CommonFields: base.CommonFields{
					GameId:     GameId,
					PlatformId: PlatformId,
					AppId:      AppId,
				},
			},*/
		//游客
		/*		{
				Account:    "bbb4b91d18df35d45cd4e832606a3ea8", //base.Md5Sum([]byte(base.GetUnixMilliString())),
				Code:       "-1",
				Password:   "-1",
				DeviceType: 2,
				Type:       4,
				Lang:       "zh-cn",
				ChannelId:  -1,
				DataExt:    "{}",
				CommonFields: base.CommonFields{
					GameId:     GameId,
					PlatformId: PlatformId,
					AppId:      AppId,
				},
			},*/
		//第三方
		/*		{
				Account:    "1001_" + base.Md5Sum([]byte("MasterWeixin")),
				Code:       "-1",
				Password:   "-1",
				DeviceType: 2,
				Type:       5,
				Lang:       "zh-cn",
				ChannelId:  -1,
				DataExt:    "{}",
				CommonFields: base.CommonFields{
					GameId:     GameId,
					PlatformId: PlatformId,
					AppId:      AppId,
				},
			},*/
	}

	for k, p := range pList {
		fmt.Printf("\nTest %d, value: %v\n", k, p)
		values := base.StructToString(&p)
		p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, SecretKey)))

		pJson, _ := json.Marshal(p)
		pString := string(pJson)
		fmt.Printf("\nrequset-data: %s\n", pString)

		req := httptest.NewRequest("POST", "/user/register", strings.NewReader(pString))
		req.Header.Set("Content-type", "application/json;charset=utf-8")
		Register(w, req)

		fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n")
		w.Body.Reset()
	}
}

func TestLogin(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	password := base.Md5Sum([]byte("i-love-u"))
	//1: email, 2:手机号, 3: 用户名, 4: 游客, 5: 第三方

	pList := []base.LoginFields{
		//email + 密码
		{
			Account:   "kangyun@outlook.com",
			Code:      "-1",
			Password:  password,
			Type:      1,
			ChannelId: -1,
			DataExt:   "{}",
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
		//email+验证码
		{
			Account:   "a123456@outlook.com",
			Code:      "123456",
			Password:  "-1",
			Type:      1,
			ChannelId: -1,
			DataExt:   "{}",
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
		//手机号+验证码
		{
			Account:   "13838384388",
			Code:      "123688",
			Password:  password,
			Type:      2,
			ChannelId: -1,
			DataExt:   "{}",
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
		//手机号+密码登录
		{
			Account:   "13838384388",
			Code:      "-1",
			Password:  password,
			Type:      2,
			ChannelId: -1,
			DataExt:   "{}",
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
		//用户名
		{
			Account:   "kangkang",
			Code:      "-1",
			Password:  password,
			Type:      3,
			ChannelId: -1,
			DataExt:   "{}",
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
		//游客
		{
			Account:   "bbb4b91d18df35d45cd4e832606a3ea8", //base.Md5Sum([]byte(base.GetUnixMilliString())),
			Code:      "-1",
			Password:  "-1",
			Type:      4,
			ChannelId: -1,
			DataExt:   "{}",
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
		//第三方
		{
			Account:   "1001_" + base.Md5Sum([]byte("MasterWeixin")),
			Code:      "-1",
			Password:  "-1",
			Type:      5,
			ChannelId: -1,
			DataExt:   "{}",
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
	}

	for k, p := range pList {

		values := base.StructToString(&p)
		p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, SecretKey)))

		pJson, _ := json.Marshal(p)
		pString := string(pJson)
		fmt.Printf("\nTest: %d, requset-data: %s\n", k, pString)

		req := httptest.NewRequest("POST", "/user/login", strings.NewReader(pString))
		req.Header.Set("Content-type", "application/json;charset=utf-8")
		Login(w, req)

		fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n\n")
		w.Body.Reset()
	}
}

func TestLoginAuth(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	p := &base.LoginAuthFields{
		Uid:   1610000100014,
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTQsIkxvZ2luVGltZSI6MTY4MDE1NjkwMiwiVG9rZW5UeXBlIjoxLCJleHAiOjE2ODAxNjQxMDIsImlzcyI6ImFjY291bnRfc2VydmVyIn0.Ui9nLyx4XPVn-4xxcl42-vUqCfqzFJCvzbNCfAKfEXY",
		CommonFields: base.CommonFields{
			GameId:     GameId,
			PlatformId: PlatformId,
			AppId:      1000000003,
		},
	}
	values := base.StructToString(p)
	p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, "506a19ddb043ff757cc8a4389ade351a")))

	pJson, _ := json.Marshal(p)
	pString := string(pJson)
	fmt.Printf("data: %s\n", pString)
	req := httptest.NewRequest("POST", "/user/loginAuth", strings.NewReader(pString))
	req.Header.Set("Content-type", "application/json;charset=utf-8")
	LoginAuth(w, req)

	fmt.Println(string(w.Body.Bytes()))
}

func TestSendVerifyCode(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()

	//email
	/*	p := &base.VerifyCodeFields{
		Account:  "273030282@qq.com",
		CodeType: 1, //验证码类型  1: 注册, 2: 忘记密码, 3: 账号绑定, 4: 账号解绑, 5: 登录
		LangId:   "zh-CN",
		Type:     1, //1 email, 2 短信
		CommonFields: base.CommonFields{
			GameId:     16,
			PlatformId: 1,
			AppId:      1000000001,
		},
	}*/

	p := &base.VerifyCodeFields{
		Account:  "13838384388",
		CodeType: 1, //验证码类型  1: 注册, 2: 忘记密码, 3: 账号绑定, 4: 账号解绑, 5: 登录
		LangId:   "zh-CN",
		Type:     2, //1 email, 2 短信
		CommonFields: base.CommonFields{
			GameId:     16,
			PlatformId: 1,
			AppId:      1000000001,
		},
	}

	values := base.StructToString(p)
	p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, SecretKey)))
	pJson, _ := json.Marshal(p)
	pString := string(pJson)
	req := httptest.NewRequest("POST", "/user/sendVerifyCode", strings.NewReader(pString))
	req.Header.Set("Content-type", "application/json;charset=utf-8")
	SendVerifyCode(w, req)

	fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n")
}

func TestBindAccount(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTMsIkxvZ2luVGltZSI6MTY4MDE4MDQ2MywiVG9rZW5UeXBlIjoxLCJleHAiOjE2ODAxODc2NjMsImlzcyI6ImFjY291bnRfc2VydmVyIn0.YT6w3YR9zTAyJNaF65jilf_jQPpAE1Uj8QrYe81cqUs"
	pList := []base.BindAccountFields{
		//游客绑第三方
		/*		{
					Account:     "bbb4b91d18df35d45cd4e832606a3ea8", //base.Md5Sum([]byte(base.GetUnixMilliString())),
					Uid:         1610000100013,
					Code:        "-1",
					Password:    "-1",
					BindAccount: "1001_f1d560b9e73fca3e04ed114c23775c86",
					Type:        5,
					LoginToken:  token,
					CommonFields: base.CommonFields{
						GameId:     GameId,
						PlatformId: PlatformId,
						AppId:      AppId,
					},
				},
				//游客绑email
				{
					Account:     "bbb4b91d18df35d45cd4e832606a3ea8", //base.Md5Sum([]byte(base.GetUnixMilliString())),
					Uid:         1610000100013,
					Code:        "168168",
					Password:    base.Md5Sum([]byte("1314520")),
					BindAccount: "ab123456@qq.com",
					Type:        1,
					LoginToken:  token,
					CommonFields: base.CommonFields{
						GameId:     GameId,
						PlatformId: PlatformId,
						AppId:      AppId,
					},
				},*/

		//游客绑手机号
		{
			Account:     "bbb4b91d18df35d45cd4e832606a3ea8", //base.Md5Sum([]byte(base.GetUnixMilliString())),
			Uid:         1610000100013,
			Code:        "168168",
			Password:    base.Md5Sum([]byte("1314520")),
			BindAccount: "13824802596",
			Type:        2,
			LoginToken:  token,
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
	}

	for k, p := range pList {
		values := base.StructToString(&p)
		p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, SecretKey)))
		pJson, _ := json.Marshal(p)
		pString := string(pJson)
		fmt.Printf("\nTest: %d, requset-data: %s\n", k+1, pString)

		req := httptest.NewRequest("POST", "/user/bindAccount", strings.NewReader(pString))
		req.Header.Set("Content-type", "application/json;charset=utf-8")
		BindAccount(w, req)

		fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n\n")
		w.Body.Reset()
	}

}

func TestUnBindAccount(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTMsIkxvZ2luVGltZSI6MTY4MDE4ODcyNiwiVG9rZW5UeXBlIjoxLCJleHAiOjE2ODAxOTU5MjYsImlzcyI6ImFjY291bnRfc2VydmVyIn0.6w7VEq7zjGHCllX2l1YEnKqQEjUTgX78uleIxsO63Q8"
	pList := []base.UnBindAccountFields{
		//解绑绑手机号
		/*		{
				Uid:           1610000100013,
				Account:       "bbb4b91d18df35d45cd4e832606a3ea8", //base.Md5Sum([]byte(base.GetUnixMilliString())),
				Code:          "168168",
				UnBindAccount: "13824802596",
				Type:          2,
				LoginToken:    token,
				CommonFields: base.CommonFields{
					GameId:     GameId,
					PlatformId: PlatformId,
					AppId:      AppId,
				},
			},*/
		//第三方
		{
			Uid:           1610000100013,
			Account:       "bbb4b91d18df35d45cd4e832606a3ea8", //base.Md5Sum([]byte(base.GetUnixMilliString())),
			Code:          "-1",
			UnBindAccount: "1001_f1d560b9e73fca3e04ed114c23775c86",
			Type:          5,
			LoginToken:    token,
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
		//解绑当前账号
		{
			Uid:           1610000100014,
			Account:       "kangyun@outlook.com", //base.Md5Sum([]byte(base.GetUnixMilliString())),
			Code:          "-1",
			UnBindAccount: "kangyun@outlook.com",
			Type:          1,
			LoginToken:    token,
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
	}

	for k, p := range pList {
		values := base.StructToString(&p)
		p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, SecretKey)))
		pJson, _ := json.Marshal(p)
		pString := string(pJson)
		fmt.Printf("\nTest: %d, requset-data: %s\n", k+1, pString)

		req := httptest.NewRequest("POST", "/user/unBindAccount", strings.NewReader(pString))
		req.Header.Set("Content-type", "application/json;charset=utf-8")
		UnBindAccount(w, req)

		fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n\n")
		w.Body.Reset()
	}

}

func TestChangePassword(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTQsIkxvZ2luVGltZSI6MTY4MDIyNDY5OCwiVG9rZW5UeXBlIjoxLCJleHAiOjE2ODAyMzE4OTgsImlzcyI6ImFjY291bnRfc2VydmVyIn0.CWqLOMKlmfmrOvK5hwHCV4Wk7DryIaXOaoWaEMUmOpE"
	pList := []base.ChangePasswordFields{

		{
			Uid:         1610000100014,
			Account:     "kangyun@outlook.com", //base.Md5Sum([]byte(base.GetUnixMilliString())),
			LoginToken:  token,
			OldPassword: base.Md5Sum([]byte("i-love-u")),
			NewPassword: base.Md5Sum([]byte("i-love-u,2")),
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
	}

	for k, p := range pList {
		values := base.StructToString(&p)
		p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, SecretKey)))
		pJson, _ := json.Marshal(p)
		pString := string(pJson)
		fmt.Printf("\nTest: %d, requset-data: %s\n", k+1, pString)

		req := httptest.NewRequest("POST", "/user/changePassword", strings.NewReader(pString))
		req.Header.Set("Content-type", "application/json;charset=utf-8")
		ChangePassword(w, req)

		fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n\n")
		w.Body.Reset()
	}

}

func TestForgetPassword(t *testing.T) {
	baseInit()

	w := httptest.NewRecorder()
	pList := []base.ForgetPasswordFields{
		{
			Account:  "kangyun@outlook.com",
			Code:     "168168",
			Password: base.Md5Sum([]byte("ILoveU")),
			Type:     1, //email
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
	}

	for k, p := range pList {
		values := base.StructToString(&p)
		p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, SecretKey)))

		pJson, _ := json.Marshal(p)
		pString := string(pJson)
		fmt.Printf("\nTest: %d, requset-data: %s\n", k+1, pString)
		req := httptest.NewRequest("POST", "/user/forgetPassword", strings.NewReader(pString))
		req.Header.Set("Content-type", "application/json;charset=utf-8")
		ForgetPassword(w, req)

		fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n\n")
		w.Body.Reset()
	}

}

func TestWhiteList(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	pList := []base.WhiteListFields{
		{
			Uid: 1610000100014,
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      1000000004,
			},
		},
	}

	for k, p := range pList {
		values := base.StructToString(&p)
		p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, "a5235128807fcbd84b9dd437eacef521")))
		pJson, _ := json.Marshal(p)
		pString := string(pJson)
		fmt.Printf("\nTest: %d, requset-data: %s\n", k+1, pString)

		req := httptest.NewRequest("POST", "/user/whiteList", strings.NewReader(pString))
		req.Header.Set("Content-type", "application/json;charset=utf-8")
		WhiteList(w, req)

		fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n\n")
		w.Body.Reset()
	}
}

func TestApplyLogout(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTQsIkxvZ2luVGltZSI6MTY4MDIzMzI0NCwiVG9rZW5UeXBlIjoxLCJleHAiOjE2ODAyNDA0NDQsImlzcyI6ImFjY291bnRfc2VydmVyIn0.-TZN_1lE6rk2dW_AsuW1dHbmX0x7Y8gv3qJDeHBb_Z0"
	pList := []base.LogoutAccountFields{
		{
			Uid:       1610000100014,
			Account:   "kangyun@outlook.com",
			Token:     token,
			ThirdInfo: "-1",
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
	}

	for k, p := range pList {
		values := base.StructToString(&p)
		p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, SecretKey)))
		pJson, _ := json.Marshal(p)
		pString := string(pJson)
		fmt.Printf("\nTest: %d, requset-data: %s\n", k+1, pString)

		req := httptest.NewRequest("POST", "/user/ApplyLogout", strings.NewReader(pString))
		req.Header.Set("Content-type", "application/json;charset=utf-8")
		ApplyLogout(w, req)

		fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n\n")
		w.Body.Reset()
	}
}

func TestUndoLogout(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTQsIkxvZ2luVGltZSI6MTY4MDIzMzI0NCwiVG9rZW5UeXBlIjoxLCJleHAiOjE2ODAyNDA0NDQsImlzcyI6ImFjY291bnRfc2VydmVyIn0.-TZN_1lE6rk2dW_AsuW1dHbmX0x7Y8gv3qJDeHBb_Z0"
	pList := []base.UndoLogoutFields{
		{
			Uid:     1610000100014,
			Account: "kangyun@outlook.com",
			Token:   token,
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
	}

	for k, p := range pList {
		values := base.StructToString(&p)
		p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, SecretKey)))
		pJson, _ := json.Marshal(p)
		pString := string(pJson)
		fmt.Printf("\nTest: %d, requset-data: %s\n", k+1, pString)

		req := httptest.NewRequest("POST", "/user/UndoLogout", strings.NewReader(pString))
		req.Header.Set("Content-type", "application/json;charset=utf-8")
		UndoLogout(w, req)

		fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n\n")
		w.Body.Reset()
	}
}

func TestGetUserInfo(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTQsIkxvZ2luVGltZSI6MTY4MDIzMzI0NCwiVG9rZW5UeXBlIjoxLCJleHAiOjE2ODAyNDA0NDQsImlzcyI6ImFjY291bnRfc2VydmVyIn0.-TZN_1lE6rk2dW_AsuW1dHbmX0x7Y8gv3qJDeHBb_Z0"
	pList := []base.UndoLogoutFields{
		{
			Uid:     1610000100014,
			Account: "kangyun@outlook.com",
			Token:   token,
			CommonFields: base.CommonFields{
				GameId:     GameId,
				PlatformId: PlatformId,
				AppId:      AppId,
			},
		},
	}

	for k, p := range pList {
		values := base.StructToString(&p)
		p.Sign = base.Md5Sum([]byte(fmt.Sprintf("%s&%s", values, SecretKey)))
		pJson, _ := json.Marshal(p)
		pString := string(pJson)
		fmt.Printf("\nTest: %d, requset-data: %s\n", k+1, pString)

		req := httptest.NewRequest("POST", "/user/getUserInfo", strings.NewReader(pString))
		req.Header.Set("Content-type", "application/json;charset=utf-8")
		GetUserInfo(w, req)

		fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n\n")
		w.Body.Reset()
	}
}

func TestDbQueryTest(t *testing.T) {
	baseInit()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/dbQuery", nil)
	req.Header.Set("Content-type", "application/json;charset=utf-8")
	DbQueryTest(w, req)
	fmt.Println("\nresult:" + string(w.Body.Bytes()) + "\n")
}

func baseInit() {
	base.InitAppService()
	//routers.InitRouterService()
}
