/**
 * @project Accounts
 * @filename users_model.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/21 12:02
 * @version 1.0
 * @description
 * 账号相关model
 */

package models

import (
	"accounts/base"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// AccountRegister 账号注册
func AccountRegister(user *base.RegisterFields, ip string) (*base.LoginReturnFields, *base.MyError) {
	//根据account获取uid,查询hash表, 如果存在，则登录
	mainUid := GetAccountUid(user.Account)
	if mainUid > 0 {
		err := base.LimitLogin(ip)
		if err != nil {
			return nil, err
		}

		loginFields := &base.LoginFields{
			Account:   user.Account,
			Code:      user.Code,
			Password:  user.Password,
			Type:      user.Type,
			ChannelId: user.ChannelId,
			DataExt:   user.DataExt,
			CommonFields: base.CommonFields{
				GameId:     user.GameId,
				PlatformId: user.PlatformId,
			},
		}
		ret, err := AccountLogin(loginFields, mainUid)
		if err.Code != base.LoginSuccess {
			base.LimitLoginIncr(ip, user.Account)
		}
		return ret, err
	}

	//不存在，则注册
	//注册需要插入三个表（account_hash, account, user（项目用户库））
	//---------------------------------------
	//此处采用的是多库独立事务+补偿的方式, 可能出现的情况：
	//第一个事务成功提交之后， 最后一个事务成功提交之前，如果出现服务器重启、数据库异常之类，可能会导致数据不一致
	//补偿是基于 hash表与账号表成功之后，项目用户表失败，补偿项目用户表的情况
	//
	//其他解决方案：
	// 1 XA方式 , Mysql 5.7.7 及以上，之前版本有bug，且需要开启XA支持，性能相对于单个事务降低 10倍
	// 2 分布式事务服务，如 DTM, Seata-go
	// https://github.com/dtm-labs/dtm
	// https://github.com/seata/seata-go

	//获取生成的账号uid
	accountUid, aErr := buildAccountUid()
	if aErr != nil {
		return nil, aErr
	}
	//根据账号uid生成项目uid
	gameUid := buildGameUid(accountUid, user.GameId, user.PlatformId)

	hashDbTable := base.GetHashDbTable(user.Account)
	dbTable := base.GetDbTable(accountUid, user.GameId, user.PlatformId)

	hashDbTx, err := hashDbTable.AccountMasterDb.Begin()
	if err != nil {
		return nil, &base.MyError{Code: base.GetHashDbTxError, Log: "get hash master db transaction error: " + err.Error()}
	}

	accountDbTx, err := dbTable.AccountMasterDb.Begin()
	if err != nil {
		return nil, &base.MyError{Code: base.GetAccountDbTxError, Log: "get account master db transaction error: " + err.Error()}
	}

	gameUserDbTx, err := dbTable.GameUserMasterDb.Begin()
	if err != nil {
		return nil, &base.MyError{Code: base.GetGameUserDbTxError, Log: "get game master db transaction error: " + err.Error()}
	}

	salt := ""
	password := ""
	//密码加盐处理
	//游客与第三方密码为空
	//手机号 + 密码或验证码
	//邮箱 + 密码或验证码
	//用户名 + 密码
	if user.Type == base.AccountUsername || (user.Type == base.AccountMobile && user.Password != base.DefaultNoValue) || (user.Type == base.AccountEmail && user.Password != base.DefaultNoValue) {
		salt = base.RandomString(base.PasswordSaltLength)
		password = base.Md5Sum([]byte(user.Password + salt))
	}

	currTime := base.GetTime()
	//插入语句
	insertHash := fmt.Sprintf("INSERT INTO %s(account, uid) VALUES(?,?)", hashDbTable.AccountHashTable)
	insertAccount := fmt.Sprintf("INSERT INTO %s (%s, uid, `password`, created_time, created_ip, last_login_time, type, device_type, salt, lang) VALUES(?,?,?,?,?,?,?,?,?,?)", dbTable.AccountTable, base.AccountType[user.Type])
	insertGameUser := fmt.Sprintf("INSERT INTO %s (account, uid, main_uid, created_time, type) VALUES (?,?,?,?,?)", dbTable.GameUserTable)

	_, err = hashDbTx.Exec(insertHash, user.Account, accountUid)
	if err != nil {
		return nil, &base.MyError{Code: base.TxExecInsertHashError, Log: "hash insert transaction exec error: " + err.Error()}
	}
	_, err = accountDbTx.Exec(insertAccount, user.Account, accountUid, password, currTime, ip, currTime, user.Type, user.DeviceType, salt, user.Lang)
	if err != nil {
		hashDbTx.Rollback()
		return nil, &base.MyError{Code: base.TxExecInsertAccountError, Log: "account insert transaction exec error: " + err.Error()}
	}

	_, err = gameUserDbTx.Exec(insertGameUser, user.Account, gameUid, accountUid, currTime, user.Type)
	if err != nil {
		hashDbTx.Rollback()
		accountDbTx.Rollback()
		return nil, &base.MyError{Code: base.TxExecInsertGameUserError, Log: "game user insert transaction exec error: " + err.Error()}
	}

	ret := &base.LoginReturnFields{
		Uid:       gameUid,
		Account:   user.Account,
		LoginTime: currTime,
		Tokens: base.LoginTokensFields{
			LoginToken:            "",
			ExpiresIn:             base.GConf.Base.LoginTokenExpires - 60,
			RefreshToken:          "",
			RefreshTokenExpiresIn: base.GConf.Base.RefreshTokenExpires,
		},
		Binds: base.BindsFields{
			Email:  "",
			Mobile: "",
			Thirds: []string{},
		},
		CardId: base.CardIdFields{
			IsRealName: 0,
			Adult:      0,
			PlayTime:   0,
		},
	}
	//获取生成token
	token, tokenErr := base.BuildLoginToken(user.GameId, user.PlatformId, user.ChannelId, ret, base.TokenTypeAccess, currTime+base.GConf.Base.LoginTokenExpires)
	if tokenErr != nil {
		hashDbTx.Rollback()
		accountDbTx.Rollback()
		gameUserDbTx.Rollback()
		return nil, tokenErr
	}
	ret.Tokens.LoginToken = token

	//获取refresh token
	refreshToken, tokenErr := base.BuildLoginToken(user.GameId, user.PlatformId, user.ChannelId, ret, base.TokenTypeRefresh, currTime+ret.Tokens.RefreshTokenExpiresIn)
	if tokenErr != nil {
		hashDbTx.Rollback()
		accountDbTx.Rollback()
		gameUserDbTx.Rollback()
		return nil, tokenErr
	}
	ret.Tokens.RefreshToken = refreshToken

	err = hashDbTx.Commit()
	if err != nil {
		accountDbTx.Rollback()
		gameUserDbTx.Rollback()
		return nil, &base.MyError{Code: base.TxHashDbCommitError, Log: "hashDbTx commit error: " + err.Error()}
	}
	err = accountDbTx.Commit()
	if err != nil {
		gameUserDbTx.Rollback()
		deleteSql := fmt.Sprintf("DELETE FROM %s WHERE account = ?", hashDbTable.AccountHashTable)
		hashDbTable.AccountMasterDb.Exec(deleteSql, user.Account)
	}
	gameUserDbTx.Commit()

	return ret, &base.MyError{Code: base.RegisterSuccess}
}

// 登录
func AccountLogin(user *base.LoginFields, mainUid int64) (*base.LoginReturnFields, *base.MyError) {
	dbTable := base.GetDbTable(mainUid, user.GameId, user.PlatformId)

	//当前时间戳
	currTime := base.GetTime()

	//查询主账号
	userInfo := &base.AccountQueryFields{}
	querySql := fmt.Sprintf("SELECT IFNULL(email, '') email, IFNULL(mobile, '') mobile, `password`, salt, `name`, card_id, type FROM %s WHERE uid = ?", dbTable.AccountTable)
	err := dbTable.AccountSlaveDb.QueryRow(querySql, mainUid).Scan(&userInfo.Email, &userInfo.Mobile, &userInfo.Password, &userInfo.Salt, &userInfo.Name, &userInfo.CardId, &userInfo.Type)
	if err != nil {
		return nil, &base.MyError{Code: base.AccountLoginException, Log: fmt.Sprintf("account %s, query uid: %d, exception: %s", user.Account, mainUid, err.Error())}
	}
	//是否实名认证绑定了身份证
	isRealName := 0
	if userInfo.Name != "" && userInfo.CardId != "" {
		isRealName = 1
	}

	//查询项目用户
	thirds := []string{}
	gameUserInfo := &base.GameUserQueryFields{}
	querySql = fmt.Sprintf("SELECT uid, main_uid, account, type, `status` FROM %s WHERE main_uid = ?", dbTable.GameUserTable)
	rows, err := dbTable.GameUserSlaveDb.Query(querySql, mainUid) //.QueryRow(querySql, mainUid).Scan(&gameUserInfo.Uid, &gameUserInfo.MainUid, &gameUserInfo.ThirdId, &gameUserInfo.Status)
	if err != nil {
		//不存在，则插入，不存在原因：
		//1 同一个账号，登录其他项目
		//2 多个单库事务因服务器异常造成失败，此处为补偿性插入
		if err == sql.ErrNoRows {
			gameUid := buildGameUid(mainUid, user.GameId, user.PlatformId)
			gameUserInsert := fmt.Sprintf("INSERT INTO %s (account, uid, main_uid, created_time, type) VALUES (?,?,?,?,?)", dbTable.GameUserTable)
			_, err = dbTable.GameUserMasterDb.Exec(gameUserInsert, user.Account, gameUid, mainUid, currTime, userInfo.Type)
			if err != nil {
				return nil, &base.MyError{Code: base.LoginInsertGameUserError, Log: "insert game users exec error: " + err.Error()}
			}
			gameUserInfo.MainUid = mainUid
			gameUserInfo.Uid = gameUid
		} else {
			return nil, &base.MyError{Code: base.LoginGameUserQueryError, Log: "query game users error: " + err.Error()}
		}
	} else {
		defer rows.Close()
		for rows.Next() {
			var (
				queryUid     int64
				queryMainUid int64
				queryAccount string
				queryType    int
				queryStatus  int
			)
			err = rows.Scan(&queryUid, &queryMainUid, &queryAccount, &queryType, &queryStatus)
			if err != nil {
				return nil, &base.MyError{Code: base.LoginInsertGameUserError, Log: "insert game users exec error: " + err.Error()}
			}
			err = rows.Scan(&gameUserInfo.Uid, &gameUserInfo.MainUid, &gameUserInfo.Account, &gameUserInfo.Type, &gameUserInfo.Status)
			if err != nil {
				return nil, &base.MyError{Code: base.QueryGameUserScanError, Log: "query game user scan error: " + err.Error()}
			}
			//当前登录的账号
			if gameUserInfo.Account == user.Account {
				gameUserInfo.MainUid = queryMainUid
				gameUserInfo.Uid = queryUid
				gameUserInfo.Type = queryType
				gameUserInfo.Account = queryAccount
				gameUserInfo.Status = queryStatus
			}
			//得到绑定的第三方账号
			if queryType == base.AccountThird && queryAccount != user.Account {
				thirds = append(thirds, queryAccount)
			}
		}
		//用户是否被禁用
		if gameUserInfo.Status == base.AccountDisabled {
			return nil, &base.MyError{Code: base.LoginAccountDisabled}
		}

		//账号注销中，返回错误码
		if gameUserInfo.Status == base.AccountDeleting {
			return nil, &base.MyError{Code: base.AccountIsBeingDeleted}
		}
	}

	saltPassword := base.Md5Sum([]byte(user.Password + userInfo.Salt))
	//用户名+密码方式，需要比对密码
	if user.Type == base.AccountUsername || user.Type == base.AccountMobile && user.Password != base.DefaultNoValue || user.Type == base.AccountEmail && user.Password != base.DefaultNoValue {
		if saltPassword != userInfo.Password {
			return nil, &base.MyError{Code: base.LoginPasswordError, Log: "password error"}
		}
	}
	adult, playTime := parseCardId(isRealName, userInfo.CardId)
	loginRet := &base.LoginReturnFields{
		Uid:       gameUserInfo.Uid,
		Account:   user.Account,
		LoginTime: currTime,
		Tokens: base.LoginTokensFields{
			LoginToken:            "",
			ExpiresIn:             base.GConf.Base.LoginTokenExpires - 60,
			RefreshToken:          "",
			RefreshTokenExpiresIn: base.GConf.Base.RefreshTokenExpires,
		},
		Binds: base.BindsFields{
			Email:  userInfo.Email,
			Mobile: userInfo.Mobile,
			Thirds: thirds,
		},
		CardId: base.CardIdFields{
			IsRealName: isRealName,
			Adult:      adult,
			PlayTime:   playTime,
		},
	}

	//获取生成token
	token, tokenErr := base.BuildLoginToken(user.GameId, user.PlatformId, user.ChannelId, loginRet, base.TokenTypeAccess, currTime+base.GConf.Base.LoginTokenExpires)
	if tokenErr != nil {
		return nil, tokenErr
	}
	loginRet.Tokens.LoginToken = token

	//获取生成token
	refreshToken, tokenErr := base.BuildLoginToken(user.GameId, user.PlatformId, user.ChannelId, loginRet, base.TokenTypeRefresh, currTime+loginRet.Tokens.RefreshTokenExpiresIn)
	if tokenErr != nil {
		return nil, tokenErr
	}
	loginRet.Tokens.RefreshToken = refreshToken

	return loginRet, &base.MyError{Code: base.LoginSuccess}
}

// 根据身份证id解析年龄信息，精确到天
// 返回是否成年、可玩时长（秒）
// 是否成年, 1 代表成年
// 如果未成年且已实名周五、六、 日和法定节假日的晚上20-21点，返回剩余的游戏时间；
// 如果成年，则无限制；
// 根据2021-08防沉迷规定：周五、六、 日和法定节假日的晚上20-21点，未成年人可以玩1小时，其余时间禁止
func parseCardId(isRealName int, id string) (int, int64) {
	//是否实名
	if isRealName == 0 {
		return 0, 0
	}
	//获取身份证的生日，与当前时间对比，得到年龄
	idYmd, err := strconv.Atoi(id[6:14])
	if err != nil {
		return 0, 0
	}
	now := time.Now()

	ymd, _ := strconv.Atoi(now.Format("20060102"))
	age := (ymd - idYmd) / 10000
	if age > 18 {
		return 1, 0
	}

	//是否星期五、六、日、法定节假日
	dayWeek := now.Weekday()
	hour := now.Hour()
	_, isHoliday := base.MemoryStoreInfo.Load("_account_holiday_" + now.Format(base.DateFormat))
	if hour >= 20 && hour < 21 && (dayWeek == 0 || dayWeek == 5 || dayWeek == 6 || isHoliday) {
		y, m, d := time.Now().Date()
		time21 := time.Date(y, m, d, 21, 0, 0, 0, time.Local)
		return 0, time21.Unix() - now.Unix()
	}

	return 0, 0
}

// 生成账号uid, 用Redis incr
func buildAccountUid() (int64, *base.MyError) {
	uid, err := base.RedisClient.Incr(base.RedisUidAutoIncrementKey).Result()
	if err != nil {
		return 0, &base.MyError{Code: base.BuildAccountUidError, Log: fmt.Sprintf("build account uid error: %s", err.Error())}
	}
	return uid, nil
}

// 生成项目uid
// 项目最大999999个
// 一个项目最大999个大区
// 每个项目的一个大区用户无限制，但在 9999999999 内容易区分
func buildGameUid(accountUid int64, gameId, platformId int) int64 {
	gId := strings.Trim(fmt.Sprintf("%6d", gameId), " ")
	pId := strings.Trim(fmt.Sprintf("%3d", platformId), " ")
	uidFormat := "%s%s0000000000"
	gameUid, _ := strconv.ParseInt(fmt.Sprintf(uidFormat, gId, pId), 10, 64)
	return gameUid + accountUid
}

// GetAccountUid 获取账号Uid
func GetAccountUid(account string) int64 {
	var uid int64
	dbTable := base.GetHashDbTable(account)
	userSql := fmt.Sprintf("SELECT uid FROM %s WHERE account = ?", dbTable.AccountHashTable)
	err := dbTable.AccountSlaveDb.QueryRow(userSql, account).Scan(&uid)
	if err != nil {
		return 0
	}
	return uid
}

// MailTplConfig 邮件模板配置
func MailTplConfig(verifyInfo *base.VerifyCodeFields) (*base.MailTpl, *base.MyError) {
	var config = &base.MailTpl{}
	mailSql := fmt.Sprintf("SELECT `type`, lang_id, title, content FROM %s WHERE `type` = ? AND lang_id = ?", base.MailTplTable)
	err := base.AccountBaseDb.QueryRow(mailSql, verifyInfo.CodeType, verifyInfo.LangId).Scan(&config.Type, &config.LangId, &config.Title, &config.Content)
	if err != nil {
		return nil, &base.MyError{Code: base.MailTplConfigQueryFailure, Log: fmt.Sprintf("query %s, error: %s", mailSql, err.Error())}
	}

	return config, nil
}

// SmsTplConfig 短信模板配置
func SmsTplConfig(verifyInfo *base.VerifyCodeFields) (*base.SmsTpl, *base.MyError) {
	var config = &base.SmsTpl{}
	tplSql := fmt.Sprintf("SELECT `type`, lang_id, sms_id, title FROM %s WHERE `type` = ? AND lang_id = ?", base.SmsTplTable)
	err := base.AccountBaseDb.QueryRow(tplSql, verifyInfo.CodeType, verifyInfo.LangId).Scan(&config.Type, &config.LangId, &config.SmsId, &config.Title)
	if err != nil {
		return nil, &base.MyError{Code: base.SmsTplConfigQueryFailure, Log: fmt.Sprintf("query %s, error: %s", tplSql, err.Error())}
	}

	return config, nil
}

// ForgetPassword 忘记密码-重置密码
func ForgetPassword(info *base.ForgetPasswordFields) *base.MyError {
	mainUid := GetAccountUid(info.Account)
	if mainUid == 0 {
		return &base.MyError{Code: base.ForgetPasswordAccountNotExists}
	}
	dbTable := base.GetDbTable(mainUid, info.GameId, info.PlatformId)
	currTime := base.GetTime()
	//密码加盐处理
	salt := base.RandomString(base.PasswordSaltLength)
	info.Password = base.Md5Sum([]byte(info.Password + salt))
	updateSql := fmt.Sprintf("UPDATE %s SET `password` = ?, salt = ?, updated_time = ? WHERE uid = ?", dbTable.AccountTable)

	_, dbErr := dbTable.AccountMasterDb.Exec(updateSql, info.Password, salt, currTime, mainUid)
	if dbErr != nil {
		return &base.MyError{Code: base.ForgetPasswordUpdateFailure, Log: fmt.Sprintf("update users exec error: %s", dbErr.Error())}
	}
	return nil
}

// 修改密码
func ChangePassword(info *base.ChangePasswordFields) *base.MyError {
	//先查询Uid
	var (
		password string
		salt     string
		gameUid  int64
		mainUid  int64
	)
	accountUid := GetAccountUid(info.Account)
	if accountUid == 0 {
		return &base.MyError{Code: base.ChangePasswordAccountNotExists}
	}

	dbTable := base.GetDbTable(accountUid, info.GameId, info.PlatformId)

	//先查询项目用户是否存在
	gameSql := fmt.Sprintf("SELECT uid, main_uid FROM %s WHERE account = ?", dbTable.GameUserTable)
	err := dbTable.GameUserSlaveDb.QueryRow(gameSql, info.Account).Scan(&gameUid, &mainUid)
	if err != nil {
		if err == sql.ErrNoRows {
			return &base.MyError{Code: base.ChangePasswordGameUserQueryNotExists, Log: fmt.Sprintf("game users query row empty, uid: %d", info.Uid)}
		}
		return &base.MyError{Code: base.ChangePasswordGameUserQueryError, Log: "query game users error: " + err.Error()}
	}
	if gameUid != info.Uid {
		return &base.MyError{Code: base.ChangePasswordUidNotMatch, Log: fmt.Sprintf("params uid :%d, by account query uid: %d", info.Uid, gameUid)}
	}

	accountSql := fmt.Sprintf("SELECT `password`, salt  FROM %s WHERE uid = ?", dbTable.AccountTable)
	err = dbTable.AccountSlaveDb.QueryRow(accountSql, accountUid).Scan(&password, &salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &base.MyError{Code: base.ChangePasswordGameUserQueryNotExists, Log: fmt.Sprintf("game users query row empty, uid: %d", info.Uid)}
		}
		return &base.MyError{Code: base.ChangePasswordGameUserQueryError, Log: "query game users error: " + err.Error()}
	}

	info.OldPassword = base.Md5Sum([]byte(info.OldPassword + salt))
	if password != info.OldPassword {
		return &base.MyError{Code: base.OldPasswordError, Log: "user password not equal to OldPassword"}
	}

	currTime := base.GetTime()
	//密码加盐处理
	newSalt := base.RandomString(base.PasswordSaltLength)
	info.NewPassword = base.Md5Sum([]byte(info.NewPassword + newSalt))
	updateSql := fmt.Sprintf("UPDATE %s SET `password` = ?, salt = ?, updated_time = ? WHERE uid = ?", dbTable.AccountTable)
	_, dbErr := dbTable.AccountMasterDb.Exec(updateSql, info.NewPassword, newSalt, currTime, accountUid)
	if dbErr != nil {
		return &base.MyError{Code: base.ChangePasswordUpdateFailure, Log: fmt.Sprintf("update users exec error: %s", dbErr.Error())}
	}
	return nil
}

// 保存验证码
func SetVerifyCode(codeKey string, codeRand string) *base.MyError {
	expireTime := time.Duration(base.GConf.Base.CodeExpires*60) * time.Second
	//写入
	err := base.RedisClient.Set(codeKey, codeRand, expireTime).Err()
	if err != nil {
		return &base.MyError{Code: base.VerifyCodeInsertError, Log: fmt.Sprintf("set code %s exec error: %s", codeKey, err.Error())}
	}
	return nil
}

// 校验验证码
func CheckVerifyCode(codeKey string, codeValue string) *base.MyError {
	val := base.RedisClient.Get(codeKey).Val()
	if codeValue == val {
		return nil
	}
	if val != "" {
		return &base.MyError{Code: base.VerifyCodeNotExists}
	}

	return &base.MyError{Code: base.VerifyCodeError}
}

// 删除已使用验证码
func DeleteVerifyCode(codeKey string) *base.MyError {
	err := base.RedisClient.Del(codeKey).Err()
	if err != nil {
		return &base.MyError{Code: base.DeleteVerifyCodeError, Log: fmt.Sprintf("delete verify code exec error: %s", err.Error())}
	}
	return nil
}

// 绑定账号
func BindAccount(bindInfo *base.BindAccountFields) *base.MyError {
	accountUid := GetAccountUid(bindInfo.Account)
	if accountUid == 0 {
		return &base.MyError{Code: base.BindAccountNotExists, Log: fmt.Sprintf("account not exists: %s", bindInfo.Account)}
	}
	//检查要绑定的账号是否存在
	bindAccountUid := GetAccountUid(bindInfo.BindAccount)
	if bindAccountUid != 0 {
		return &base.MyError{Code: base.BindAccountAlreadyExists, Log: fmt.Sprintf("bind account already exists: %s", bindInfo.BindAccount)}
	}

	dbTable := base.GetDbTable(accountUid, bindInfo.GameId, bindInfo.PlatformId)

	//查询邮箱或手机号是否已绑定, 不是空则不能绑定
	if bindInfo.Type != base.AccountThird {
		accountInfo := &base.AccountQueryFields{}
		querySql := fmt.Sprintf("SELECT IFNULL(email, '') email, IFNULL(mobile, '') mobile FROM %s WHERE uid = ?", dbTable.AccountTable)
		err := dbTable.AccountSlaveDb.QueryRow(querySql, accountUid).Scan(&accountInfo.Email, &accountInfo.Mobile)
		if err != nil {
			return &base.MyError{Code: base.BindQueryAccountInfoError, Log: fmt.Sprintf("bind account already exists: %s", bindInfo.BindAccount)}
		}
		if bindInfo.Type == base.AccountEmail && accountInfo.Email != "" {
			return &base.MyError{Code: base.BindEmailAlreadyExists}
		}
		if bindInfo.Type == base.AccountMobile && accountInfo.Mobile != "" {
			return &base.MyError{Code: base.BindMobileAlreadyExists}
		}
	}

	hashDbTable := base.GetHashDbTable(bindInfo.BindAccount)
	hashDbTx, err := hashDbTable.AccountMasterDb.Begin()
	if err != nil {
		return &base.MyError{Code: base.BindGetHashTxError, Log: fmt.Sprintf("get hash db transaction error: %s", err.Error())}
	}
	dbTx, err := dbTable.AccountMasterDb.Begin()
	if err != nil {
		return &base.MyError{Code: base.BindGetAccountDbTxError, Log: fmt.Sprintf("get account db transaction error: %s", err.Error())}
	}

	gameTx, err := dbTable.GameUserMasterDb.Begin()
	if err != nil {
		return &base.MyError{Code: base.BindGetGameUserDbTxError, Log: fmt.Sprintf("get game user db transaction error: %s", err.Error())}
	}

	//如果是第三方，插入 项目用户表、account hash 表
	//如果是email或mobile， 更新account表、插入account hash表
	insertHash := fmt.Sprintf("INSERT INTO %s(account, uid) VALUES(?,?)", hashDbTable.AccountHashTable)
	_, err = hashDbTx.Exec(insertHash, bindInfo.BindAccount, accountUid)
	if err != nil {
		return &base.MyError{Code: base.BindHashTxExecInsertError, Log: "bind, hash insert transaction exec error: " + err.Error()}
	}
	currTime := base.GetTime()
	if bindInfo.Type == base.AccountThird {
		insertGameUser := fmt.Sprintf("INSERT INTO %s (account, uid, main_uid, created_time, type) VALUES (?,?,?,?,?)", dbTable.GameUserTable)
		_, err = gameTx.Exec(insertGameUser, bindInfo.BindAccount, bindInfo.Uid, accountUid, currTime, bindInfo.Type)
		if err != nil {
			hashDbTx.Rollback()
			return &base.MyError{Code: base.BindGameTxExecInsertError, Log: "bind, game user insert transaction exec error: " + err.Error()}
		}
	}
	if bindInfo.Type == base.AccountEmail || bindInfo.Type == base.AccountMobile {
		password := ""
		salt := ""
		if bindInfo.Password != base.DefaultNoValue {
			salt = base.RandomString(base.PasswordSaltLength)
			password = base.Md5Sum([]byte(bindInfo.Password + salt))
		}
		accountUpdate := fmt.Sprintf("UPDATE %s SET %s = ?, updated_time = ?, `password` = ?, salt = ? WHERE uid = ?", dbTable.AccountTable, base.AccountType[bindInfo.Type])
		_, err = dbTx.Exec(accountUpdate, bindInfo.BindAccount, currTime, password, salt, accountUid)
		if err != nil {
			hashDbTx.Rollback()
			return &base.MyError{Code: base.BindAccountTxExecUpdateError, Log: "bind, account update transaction exec error: " + err.Error()}
		}
	}
	err = hashDbTx.Commit()
	if err != nil {
		dbTx.Rollback()
		gameTx.Rollback()
		deleteSql := fmt.Sprintf("DELETE FROM %s WHERE account = ?", hashDbTable.AccountHashTable)
		hashDbTable.AccountMasterDb.Exec(deleteSql, bindInfo.BindAccount)
	}
	dbTx.Commit()
	gameTx.Commit()

	return nil
}

// 根据account 获取其已绑定的信息
func GetAccountBindsInfo(account string, gameId, platformId int) (*base.UserInfoRespFields, *base.MyError) {
	accountUid := GetAccountUid(account)
	if accountUid == 0 {
		return nil, &base.MyError{Code: base.GetAlreadyBindInfoNotFound}
	}
	dbTable := base.GetDbTable(accountUid, gameId, platformId)
	//查询主表
	var (
		email  string
		mobile string
		cardId string
		name   string
	)
	thirds := []string{}

	querySql := fmt.Sprintf("SELECT IFNULL(email, '') email, IFNULL(mobile, '') mobile, card_id, `name` FROM %s WHERE uid = ?", dbTable.AccountTable)
	err := dbTable.AccountSlaveDb.QueryRow(querySql, accountUid).Scan(&email, &mobile, &cardId, &name)
	if err != nil {
		return nil, &base.MyError{Code: base.GetAlreadyBindInfoNotFoundByUid, Log: fmt.Sprintf("query account uid: %d error: %s", accountUid, err.Error())}
	}

	//查找项目用户表中第三方账号
	querySql = fmt.Sprintf("SELECT account, type, `status` FROM %s WHERE main_uid = ?", dbTable.GameUserTable)
	rows, err := dbTable.GameUserSlaveDb.Query(querySql, accountUid)
	if err != nil {
		return nil, &base.MyError{Code: base.GetAlreadyBindThirdError, Log: "get bind_third_info error:" + err.Error()}
	}
	defer rows.Close()
	for rows.Next() {
		var (
			tmpAccount string
			tmpType    int
			tmpStatus  int
		)
		err = rows.Scan(&tmpAccount, &tmpType, &tmpStatus)
		if err != nil {
			return nil, &base.MyError{Code: base.QueryAlreadyBindThirdScanError, Log: "query game user by main_uid scan error: " + err.Error()}
		}
		//用户是否被禁用
		if tmpStatus == base.AccountDisabled {
			return nil, &base.MyError{Code: base.LoginAccountDisabled}
		}

		//账号注销中，返回错误码
		if tmpStatus == base.AccountDeleting {
			return nil, &base.MyError{Code: base.AccountIsBeingDeleted}
		}
		//得到绑定的第三方账号, 除当前账号外
		if tmpType == base.AccountThird && tmpAccount != account {
			thirds = append(thirds, tmpAccount)
		}
	}
	isRealName := 1
	if cardId == "" && name == "" {
		isRealName = 0
	}
	adult, playTime := parseCardId(isRealName, cardId)
	ret := &base.UserInfoRespFields{
		Binds: base.BindsFields{
			Email:  email,
			Mobile: mobile,
			Thirds: thirds,
		},
		Cards: base.CardIdFields{
			IsRealName: isRealName,
			Adult:      adult,
			PlayTime:   playTime,
		},
	}

	return ret, nil
}

// 解绑账号
func UnBindAccount(unBindInfo *base.UnBindAccountFields) (*base.UserInfoRespFields, *base.MyError) {
	accountUid := GetAccountUid(unBindInfo.Account)
	if accountUid == 0 {
		return nil, &base.MyError{Code: base.UnbindAccountNotExists, Log: fmt.Sprintf("account not exists: %s", unBindInfo.Account)}
	}
	//检查要绑定的账号是否存在
	unBindAccountUid := GetAccountUid(unBindInfo.UnBindAccount)
	if unBindAccountUid == 0 {
		return nil, &base.MyError{Code: base.BeUnBindAccountNotExists, Log: fmt.Sprintf("be unbind account not exists: %s", unBindInfo.UnBindAccount)}
	}

	//验证解绑的账号和当前账号是同一个账号
	if unBindAccountUid != accountUid {
		return nil, &base.MyError{Code: base.BindAccountAndBeUnBindNotMatch, Log: fmt.Sprintf("unbind account uid: %d, be unbind account uid :%d", accountUid, unBindAccountUid)}
	}

	dbTable := base.GetDbTable(unBindAccountUid, unBindInfo.GameId, unBindInfo.PlatformId)

	//1 查找解绑的账号是否存在且不能解绑注册时的类型
	//2 得到已绑定的信息
	var (
		email        string
		mobile       string
		accountType  int
		thirdAccount string
		gameUserType int
		name         string
		cardId       string
	)
	thirdAccounts := []string{}
	//检查解绑的账号是否存在
	querySql := fmt.Sprintf("SELECT IFNULL(email, '') email, IFNULL(mobile, '') mobile, IFNULL(third, '') third, type, card_id, `name` FROM %s WHERE uid = ?", dbTable.AccountTable)
	err := dbTable.AccountSlaveDb.QueryRow(querySql, accountUid).Scan(&email, &mobile, &thirdAccount, &accountType, &cardId, &name)
	if err != nil {
		return nil, &base.MyError{Code: base.BeUnBindAccountNotExists2, Log: fmt.Sprintf("be unbind account not exists error: %s", err.Error())}
	}

	//不支持解绑注册时的类型,第三方不支持解绑注册时的账号
	if (unBindInfo.Type == base.AccountThird && thirdAccount == unBindInfo.UnBindAccount) || accountType == unBindInfo.Type {
		return nil, &base.MyError{Code: base.UnBindUnSupportRegisterType}
	}

	if unBindInfo.Type == base.AccountEmail {
		if email == "" {
			return nil, &base.MyError{Code: base.EmailAlreadyUnbind}
		} else { //处理返回已绑定信息
			email = ""
		}
	}
	if unBindInfo.Type == base.AccountMobile {
		if mobile == "" {
			return nil, &base.MyError{Code: base.MobileAlreadyUnbind}
		} else {
			mobile = ""
		}
	}

	//第三方
	gameUserSql := fmt.Sprintf("SELECT account, type FROM %s WHERE main_uid = ?", dbTable.GameUserTable)
	rows, err := dbTable.GameUserSlaveDb.Query(gameUserSql, accountUid)
	if err != nil {
		return nil, &base.MyError{Code: base.UnBindQueryThirdError, Log: fmt.Sprintf("be unbind, query third info error: %s", err.Error())}
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&thirdAccount, &gameUserType)
		if err != nil {
			return nil, &base.MyError{Code: base.UnBindQueryThirdScanError, Log: fmt.Sprintf("be unbind, query third info scan error: %s", err.Error())}
		}
		//处理已绑定信息
		if thirdAccount == unBindInfo.UnBindAccount || gameUserType != base.AccountThird {
			continue
		}
		thirdAccounts = append(thirdAccounts, thirdAccount)
	}

	//解绑，要删除hash表记录、account、项目用户表
	hashDbTable := base.GetHashDbTable(unBindInfo.UnBindAccount)

	gameDbTx, err := dbTable.GameUserMasterDb.Begin()
	//如果是解绑第三方，则先删除第三方
	if unBindInfo.Type == base.AccountThird {
		deleteGameUser := fmt.Sprintf("DELETE FROM %s WHERE account = ?", dbTable.GameUserTable)
		if err != nil {
			return nil, &base.MyError{Code: base.UnBindGetGameUserDbTxError, Log: fmt.Sprintf("get game db transaction error: %s", err.Error())}
		}
		_, err = gameDbTx.Exec(deleteGameUser, unBindInfo.UnBindAccount)
		if err != nil {
			return nil, &base.MyError{Code: base.UnBindGameUserTxExecDeleteError, Log: fmt.Sprintf("unbind,delete game user transaction error: %s", err.Error())}
		}
	}

	hashDbTx, err := hashDbTable.AccountMasterDb.Begin()
	if err != nil {
		if unBindInfo.Type == base.AccountThird {
			gameDbTx.Rollback()
		}
		return nil, &base.MyError{Code: base.UnBindGetHashTxError, Log: fmt.Sprintf("get hash db transaction error: %s", err.Error())}
	}
	dbTx, err := dbTable.AccountMasterDb.Begin()
	if err != nil {
		if unBindInfo.Type == base.AccountThird {
			gameDbTx.Rollback()
		}
		return nil, &base.MyError{Code: base.UnBindGetAccountDbTxError, Log: fmt.Sprintf("get account db transaction error: %s", err.Error())}
	}

	deleteHashSql := fmt.Sprintf("DELETE FROM %s WHERE account = ?", hashDbTable.AccountHashTable)
	updateAccountSql := fmt.Sprintf("UPDATE %s SET %s = null, updated_time = ? WHERE uid = ?", dbTable.AccountTable, base.AccountType[unBindInfo.Type])
	_, err = hashDbTx.Exec(deleteHashSql, unBindInfo.UnBindAccount)
	if err != nil {
		if unBindInfo.Type == base.AccountThird {
			gameDbTx.Rollback()
		}
		return nil, &base.MyError{Code: base.UnBindHashTxExecUpdateError, Log: fmt.Sprintf("unbind, account hash update transaction error: %s", err.Error())}
	}

	_, err = dbTx.Exec(updateAccountSql, base.GetTime(), unBindAccountUid)
	if err != nil {
		if unBindInfo.Type == base.AccountThird {
			gameDbTx.Rollback()
		}
		hashDbTx.Rollback()
		return nil, &base.MyError{Code: base.UnBindAccountTxExecUpdateError, Log: fmt.Sprintf("unbind, account update transaction error: %s", err.Error())}
	}
	if unBindInfo.Type == base.AccountThird {
		gameDbTx.Commit()
	}
	hashDbTx.Commit()
	dbTx.Commit()

	isRealName := 0
	if cardId != "" && name != "" {
		isRealName = 1
	}
	adult, playTime := parseCardId(isRealName, cardId)
	binds := &base.UserInfoRespFields{
		Binds: base.BindsFields{
			Email:  email,
			Mobile: mobile,
			Thirds: thirdAccounts,
		},
		Cards: base.CardIdFields{
			IsRealName: isRealName,
			Adult:      adult,
			PlayTime:   playTime,
		},
	}

	return binds, nil
}

// 查询是否在白名单中
func GetWhiteList(whiteListInfo *base.WhiteListFields, ip string) *base.WhiteListResponse {
	//查询ip是否在白名单内
	var env string
	ipSql := fmt.Sprintf("SELECT env FROM %s WHERE game_id = ? AND platform_id = ? AND ip = ?", base.WhiteUserListTable)
	err := base.AccountBaseDb.QueryRow(ipSql, whiteListInfo.GameId, whiteListInfo.PlatformId, ip).Scan(&env)

	if err == nil {
		if env == "" {
			return &base.WhiteListResponse{Status: base.NotInWhiteList, Env: []string{}}
		}
		return &base.WhiteListResponse{Status: base.InWhiteList, Env: strings.Split(env, "|")}
	}

	//查询uid是否在白名单内
	userSql := fmt.Sprintf("SELECT env FROM %s WHERE game_id = ? AND platform_id = ? AND uid = ?", base.WhiteUserListTable)
	err = base.AccountBaseDb.QueryRow(userSql, whiteListInfo.GameId, whiteListInfo.PlatformId, whiteListInfo.Uid).Scan(&env)
	if err == nil {
		if env == "" {
			return &base.WhiteListResponse{Status: base.NotInWhiteList, Env: []string{}}
		}
		return &base.WhiteListResponse{Status: base.InWhiteList, Env: strings.Split(env, "|")}
	}

	return &base.WhiteListResponse{Status: base.NotInWhiteList, Env: []string{}}
}

// 添加删除申请
// 往删除申请表里插入一条， 即使有多个项目uid，然后根据main_uid更新所有的uid状态
func AddDeleteApply(deleteInfo *base.LogoutAccountFields, userLog *zerolog.Logger) *base.MyError {
	accountUid := GetAccountUid(deleteInfo.Account)
	if accountUid == 0 {
		return &base.MyError{Code: base.DeleteAccountNotExists}
	}

	dbTable := base.GetDbTable(accountUid, deleteInfo.GameId, deleteInfo.PlatformId)
	//查询用户是否存在
	var (
		mainUid      int64
		gameUsertype int
	)
	querySql := fmt.Sprintf("SELECT main_uid, type FROM %s WHERE account = ?", dbTable.GameUserTable)
	err := dbTable.GameUserSlaveDb.QueryRow(querySql, deleteInfo.Account).Scan(&mainUid, &gameUsertype)
	if err != nil {
		if err == sql.ErrNoRows {
			return &base.MyError{Code: base.DeleteAccountQueryNoRows, Log: fmt.Sprintf("game users query row empty, uid: %d", deleteInfo.Uid)}
		} else {
			return &base.MyError{Code: base.DeleteAccountQueryError, Log: "query game users error: " + err.Error()}
		}
	}
	if mainUid != accountUid {
		return &base.MyError{Code: base.DeleteAccountAndUidNotMatch, Log: fmt.Sprintf("account uid: %d, by %d query account uid: %d", accountUid, deleteInfo.Uid, mainUid)}
	}

	//查看是否已申请过
	querySql = fmt.Sprintf("SELECT uid FROM %s WHERE uid = ?", dbTable.GameUserDeleteApplyTable)
	err = dbTable.GameUserSlaveDb.QueryRow(querySql, deleteInfo.Uid).Scan(&deleteInfo.Uid)
	if err == nil {
		return &base.MyError{Code: base.DeleteApplyAlreadyExists}
	}

	//内存中获取冷静期时长
	key := fmt.Sprintf("_account_game_config_%d_%d", deleteInfo.GameId, deleteInfo.PlatformId)
	gameConfig, _ := base.MemoryGameConfig.Load(key)
	gameConfigInfo, ok := gameConfig.(base.GameConfig)
	if !ok {
		gameConfigInfo.UserDeleteWaitDuration = 15
	}

	//如果有apple的信息，先获取token
	thirdAccount := &base.ThirdAccount{}
	deleteAccountExt := &base.DeleteAccountExt{}
	jsonErr := json.Unmarshal([]byte(deleteInfo.ThirdInfo), thirdAccount)
	if jsonErr == nil && thirdAccount.ThirdId == base.ThirdApple && thirdAccount.AuthorizationCode != "" && gameConfigInfo.AppleClientId != "" && gameConfigInfo.AppleClientSecret != "" {
		code, decodeErr := base64.StdEncoding.DecodeString(thirdAccount.AuthorizationCode)
		if decodeErr == nil {
			appleTokenParams := map[string]string{
				"client_id":     gameConfigInfo.AppleClientId,
				"client_secret": gameConfigInfo.AppleClientSecret,
				"code":          string(code),
				"grant_type":    "authorization_code",
			}
			result, err := base.HttpPostForm(base.AppleValidateCodeUrl, appleTokenParams)
			userLog.Info().Interface("apple_validate_params", appleTokenParams).RawJSON("apple_validate_return", result).Msg("apple validate code info")

			if err == nil {
				var appleValidateInfo struct {
					RefreshToken string `json:"refresh_token"`
				}
				jsonErr = json.Unmarshal(result, &appleValidateInfo)
				if jsonErr == nil && appleValidateInfo.RefreshToken != "" {
					deleteAccountExt.AppleRefreshToken = appleValidateInfo.RefreshToken
				}
			} else {
				userLog.Error().Str("apple_validate_error", err.Log).Msg("apple validate code error")
			}
		} else {
			userLog.Error().Str("authorization_code", thirdAccount.AuthorizationCode).Msg("authorization code base64 decode error: " + decodeErr.Error())
		}
	}

	currTime := base.GetTime()
	executeTime := currTime + int64(gameConfigInfo.UserDeleteWaitDuration*86400)
	deleteAccountExtString, _ := json.Marshal(deleteAccountExt)
	gameUserTx, _ := dbTable.GameUserMasterDb.Begin()

	//添加申请,
	applySql := "INSERT INTO %s (uid, main_uid, account, type, apply_time, execute_delete_time, ext_info, ext)"
	applySql += "SELECT uid, main_uid, account, type, '%d' apply_time, %d execute_delete_time, '%s' ext_info, ext FROM %s WHERE main_uid = ?"
	applySql = fmt.Sprintf(applySql, dbTable.GameUserDeleteApplyTable, currTime, executeTime, deleteAccountExtString, dbTable.GameUserTable)

	res, err := gameUserTx.Exec(applySql, mainUid)
	if err != nil {
		gameUserTx.Rollback()
		return &base.MyError{Code: base.AddDeleteApplyError, Log: "insert users_delete_apply error: " + err.Error()}
	}
	affected, _ := res.RowsAffected()
	userLog.Info().Msgf("account apply logout, main uid: %d, affected rows: %d", mainUid, affected)

	//修改用户状态
	gameUsersSql := fmt.Sprintf("UPDATE %s SET `status` = ? WHERE main_uid = ? AND `status` = ?", dbTable.GameUserTable)
	_, err = gameUserTx.Exec(gameUsersSql, base.AccountDeleting, mainUid, base.AccountNormal)
	if err != nil {
		gameUserTx.Rollback()
		return &base.MyError{Code: base.DeleteApplyUpdateUserStatusError, Log: "add delete apply update users status error: " + err.Error()}
	}

	gameUserTx.Commit()
	return nil
}

// 撤销删除申请
func UndoDeleteApply(undoDeleteInfo *base.UndoLogoutFields) *base.MyError {
	accountUid := GetAccountUid(undoDeleteInfo.Account)
	if accountUid == 0 {
		return &base.MyError{Code: base.UndoDeleteAccountNotExists}
	}

	dbTable := base.GetDbTable(accountUid, undoDeleteInfo.GameId, undoDeleteInfo.PlatformId)
	//查询是否有申请记录
	querySql := fmt.Sprintf("SELECT account FROM %s WHERE account = ?", dbTable.GameUserDeleteApplyTable)
	var account string
	err := dbTable.GameUserSlaveDb.QueryRow(querySql, undoDeleteInfo.Account).Scan(&account)
	if err != nil {
		return &base.MyError{Code: base.UndoDeleteApplyNotExists, Log: "delete apply not exists, error: " + err.Error()}
	}
	if account != undoDeleteInfo.Account {
		return &base.MyError{Code: base.UndoDeleteAccountAndRecordNotMatch, Log: fmt.Sprintf("undo delete account: %s, db record account: %s", undoDeleteInfo.Account, account)}
	}
	//修改用户状态
	gameUsersSql := fmt.Sprintf("UPDATE %s SET `status` = ? WHERE main_uid = ? AND `status` = ?", dbTable.GameUserDeleteApplyTable)
	_, err = dbTable.GameUserMasterDb.Exec(gameUsersSql, base.ApplyStatusRecover, accountUid, base.AccountDeleting)
	if err != nil {
		return &base.MyError{Code: base.UndoDeleteUpdateUserStatusError, Log: "undo delete apply update users status error: " + err.Error()}
	}

	return nil
}

// 调用第三方实名认证，检测姓名与身份证是否正确
func RealNameAuth(userInfo *base.UserRealNameAuthReqFields) *base.MyError {
	// todo 请求接口
	//-----
	//-------
	//---------

	//如果请求接口OK, 则更新账号信息
	accountUid := GetAccountUid(userInfo.Account)
	if accountUid == 0 {
		return &base.MyError{Code: base.RealNameGetAccountUidNotExists}
	}

	dbTable := base.GetDbTable(accountUid, -1, -1)
	querySql := fmt.Sprintf("UPDATE %s SET `name` = ?, card_id = ?, updated_time = ? WHERE uid = ?", dbTable.AccountTable)
	_, err := dbTable.AccountMasterDb.Exec(querySql, userInfo.Name, userInfo.CardId, base.GetTime(), accountUid)
	if err != nil {
		return &base.MyError{Code: base.RealNameUpdateAccountError, Log: "update account name,card_id error: " + err.Error()}
	}

	return nil
}

func DbTest() (interface{}, error) {

	/*update := "UPDATE user_1 SET created_time = 1 WHERE account = '2730302821@qq.com'"
	res, err := base.AccountBaseDb.Exec(update)
	if err != nil {
		fmt.Printf("update user_1 error: %s", err.Error())
		return nil, &base.MyError{Code: base.Failure, Msg: "err: " + err.Error()}
	}
	affected, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("update user_1 affected error: %s", err.Error())
		return nil, &base.MyError{Code: base.Failure, Msg: "err: " + err.Error()}
	}

	fmt.Printf("update user_1 affected: %d\n", affected)
	return nil, &base.MyError{Code: base.Success}

	insert := "INSERT INTO user_1(uid, main_uid, third_id) VALUES(?, ?, ?)"
	uid := 2000000000
	mainUid := 100001
	thirdIdValue := ""
	thirdId := sql.NullString{
		String: thirdIdValue,
		Valid:  false,
	}

	res, err = base.AccountBaseDb.Exec(insert, uid, mainUid, thirdId)
	if err != nil {
		fmt.Printf("insert into user_1 error: %s", err.Error())
		return nil, &base.MyError{Code: base.Failure, Msg: "err: " + err.Error()}
	}
	lastId, err2 := res.LastInsertId()
	if err2 != nil {
		fmt.Printf("get LastInsertId error: %s", err2.Error())
		return nil, &base.MyError{Code: base.Failure, Msg: "err: " + err2.Error()}
	}
	fmt.Printf("insert into user_1 error: %d", lastId)
	return res, &base.MyError{Code: base.Success}*/

	//account := "273030282@qq.com"
	//dbTable := base.GetDbTable(1, 16, 1)
	//sql := fmt.Sprintf("SELECT COUNT(*) num FROM %s WHERE email = ?", dbTable.AccountTable)
	//querySql := fmt.Sprintf("SELECT IFNULL(email, '') email, IFNULL(mobile, '') mobile,card_id FROM %s WHERE uid = 2460000000001", dbTable.AccountTable)
	querySql := fmt.Sprintf("SELECT card_id, `name` FROM tpl_account_1 WHERE uid = 1234")
	fmt.Println(querySql)
	//count := 0
	var (
		/*email  string
		mobile string*/
		result []string
		cardId string
		name   string
	)

	err := base.AccountBaseDb.QueryRow(querySql).Scan(&cardId, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, &base.MyError{Code: base.Failure, Log: "err: " + err.Error(), Msg: "ErrNoRows err: " + err.Error()}
		}
		return result, &base.MyError{Code: base.Failure, Log: "err: " + err.Error(), Msg: "err: " + err.Error()}
	}
	result = append(result, cardId)
	result = append(result, name)

	return result, nil
}
