/**
 * @project Accounts
 * @filename limit.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/3/28 10:15
 * @version 1.0
 * @description
 * 恶意访问限制
 * 当达到一段数量后，返回图形验证码id, 验证成功后可以继续访问
 */

package base

import (
	"fmt"
	"time"
)

// ip 限制
func LimitIp(ip, urlPath string) *MyError {
	if !GConf.RequestLimitRule.Enabled {
		return nil
	}

	if _, ok := GConf.RequestLimitRule.WhiteListMap[urlPath]; ok {
		return nil
	}

	lockKey := LimitIpLocKey + ip // li
	val, _ := RedisClient.Get(lockKey).Int()
	if val == 1 {
		ret := &LimitLockRetFields{
			CaptchaId:   BuildCaptchaId(),
			CaptchaType: LimitIpKey,
		}
		return &MyError{Code: RequestLimitRuleIpTrigger, Data: ret}
	}

	conf := GConf.RequestLimitRule.Ip
	if len(conf) != 3 {
		return &MyError{Code: RequestLimitIpConfigError}
	}

	expire := time.Duration(conf[0]) * time.Second
	limit := int64(conf[1])
	lockTime := time.Duration(conf[2]) * time.Second

	key := LimitIpKey + ip
	result := RedisClient.SetNX(key, 1, expire)
	if result.Val() {
		return nil
	}

	value := RedisClient.Incr(key).Val()
	if value >= limit {
		RedisClient.Set(lockKey, 1, lockTime)
	}

	return nil
}

// 是否超过限定
func LimitLogin(ip string) *MyError {
	if !GConf.RequestLimitRule.Enabled {
		return nil
	}

	key := LimitLoginIpKey + ip //li: login ip
	value, _ := RedisClient.Get(key).Int()
	if value == 1 {
		return &MyError{Code: RequestLimitLoginIpLock}
	}

	return nil
}

// 累加登录失败的账号
func LimitLoginIncr(ip, account string) *MyError {
	if !GConf.RequestLimitRule.Enabled {
		return nil
	}
	conf := GConf.RequestLimitRule.Login
	if len(conf) != 3 {
		return &MyError{Code: RequestLimitLoginConfigError}
	}
	allowTime := time.Duration(conf[0]) * time.Second
	limit := int64(conf[1])
	lockTime := time.Duration(conf[2]) * time.Second

	accountKey := fmt.Sprintf("_account_limit_l_%s", account) //l: login
	lockKey := LimitLoginIpKey + ip                           //li: login ip
	result := RedisClient.SetNX(accountKey, 1, allowTime)

	//如果是第一次 则直接返回正常
	if result.Val() {
		return nil
	}
	val := RedisClient.Incr(accountKey).Val()
	if val >= limit {
		RedisClient.Set(lockKey, 1, lockTime)
	}

	return nil
}

// 发送后累加数量
func LimitVerifyCodeIncr(ip, account string) *MyError {
	if !GConf.RequestLimitRule.Enabled {
		return nil
	}
	conf := GConf.RequestLimitRule.VerifyCode
	if len(conf) != 3 {
		return &MyError{Code: RequestLimitVerifyCodeConfigError}
	}
	allowTime := time.Duration(conf[0]) * time.Second
	allowNumber := int64(conf[2])
	accountKey := fmt.Sprintf("_account_limit_vc_%s", account) //vc: verify code
	ipKey := fmt.Sprintf("_account_limit_vci_%s", ip)          //vci: verify code ip

	//锁账号
	accountRes := RedisClient.SetNX(accountKey, 1, allowTime)
	ipRes := RedisClient.SetNX(ipKey, 1, allowTime)
	if !accountRes.Val() {
		RedisClient.Incr(accountKey)
	}
	if !ipRes.Val() {
		val := RedisClient.Incr(ipKey).Val()
		if val >= allowNumber {
			lockKey := LimitCodeIpKey + ip //vcil: verify code ip lock
			RedisClient.Set(lockKey, 1, allowTime)
		}
	}
	return nil
}

// 发送验证码前检查是否达到限制数量
func LimitVerifyCode(ip, account string) *MyError {
	if !GConf.RequestLimitRule.Enabled {
		return nil
	}
	conf := GConf.RequestLimitRule.VerifyCode
	if len(conf) != 3 {
		return &MyError{Code: RequestLimitVerifyCodeConfigError}
	}
	allowNumber := conf[1]

	accountKey := fmt.Sprintf("_account_limit_vc_%s", account) //vc: verify code
	accountLimitVal, _ := RedisClient.Get(accountKey).Int()
	if accountLimitVal >= allowNumber {
		return &MyError{Code: RequestLimitVerifyCodeAccount}
	}

	ipLockKey := LimitCodeIpKey + ip //vcil: verify code ip lock
	ipLimitVal, _ := RedisClient.Get(ipLockKey).Int()
	if ipLimitVal == 1 {
		return &MyError{Code: RequestLimitVerifyCodeLockIp}
	}

	return nil
}

// 注册前检查是否达到限制数量
func LimitRegister(ip string) *MyError {
	if !GConf.RequestLimitRule.Enabled {
		return nil
	}
	lockKey := LimitRegisterIpKey + ip //ril = limit register ip lock
	val, _ := RedisClient.Get(lockKey).Int()
	if val == 1 {
		return &MyError{Code: RequestLimitRegisterLockIp}
	}
	return nil
}

// 注册成功后累加数量
func LimitRegisterIncr(ip string) *MyError {
	if !GConf.RequestLimitRule.Enabled {
		return nil
	}
	conf := GConf.RequestLimitRule.Register
	if len(conf) != 3 {
		return &MyError{Code: RequestLimitRegisterConfigError}
	}
	limitTime := time.Duration(conf[0]) * time.Second
	limitNumber := int64(conf[1])
	lockTime := time.Duration(conf[2]) * time.Second

	key := fmt.Sprintf("_account_limit_ri_%s", ip) //lri = register ip
	lockKey := LimitRegisterIpKey + ip             //ril = register ip lock

	res := RedisClient.SetNX(key, 1, limitTime)
	if !res.Val() {
		val := RedisClient.Incr(key).Val()
		if val >= limitNumber {
			RedisClient.Set(lockKey, 1, lockTime)
		}
	}

	return nil
}
