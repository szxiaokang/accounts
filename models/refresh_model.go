/**
 * @project Accounts
 * @filename refresh_model.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/10 11:22
 * @version 1.0
 * @description
 * 刷新相关数据库操作
 */

package models

import (
	"accounts/base"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"time"
)

type UserDeleteApply struct {
	Uid     int64
	MainUid int64
	Account string
	ExtInfo string
	Ext     string
	Type    int
}

func RefreshHoliday() {
	log.Info().Msg("RefreshHoliday to memory start")
	querySql := fmt.Sprintf("SELECT ymd FROM %s", base.HolidayTable)
	rows, err := base.AccountBaseDb.Query(querySql)
	if err != nil {
		log.Error().Msgf("query holiday table error: %s", err.Error())
		return
	}
	defer rows.Close()

	var holiday string
	for rows.Next() {
		err = rows.Scan(&holiday)
		if err != nil {
			log.Error().Msgf("query holiday table, scan error: %s", err.Error())
			continue
		}

		key := fmt.Sprintf("_account_holiday_%s", holiday)
		base.MemoryStoreInfo.Store(key, 1)
		log.Info().Msgf("%s load memory success!", key)
	}

	return
}

// RefreshGameConfig 定时读取game_config表数据,保存到内存中
// 超过有效期,则根据game_config表的apple_config数据,生成client_secret,并写入到game_config表的apple_client_id,apple_client_secret中
func RefreshGameConfig() {
	log.Info().Msg("RefreshGameConfig start")
	//1.查询game_config表数据
	querySql := fmt.Sprintf("SELECT game_id, platform_id, user_delete_wait_duration, apple_client_id, apple_client_secret, ext, apple_config, last_apple_update_time FROM %s", base.UserDeleteConfigTable)
	rows, err := base.AccountBaseDb.Query(querySql)
	if err != nil {
		log.Error().Msgf("query game config error: %s", err.Error())
		return
	}

	var gameConfig base.GameConfig
	var appleData string
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&gameConfig.GameId, &gameConfig.PlatformId, &gameConfig.UserDeleteWaitDuration, &gameConfig.AppleClientId, &gameConfig.AppleClientSecret, &gameConfig.Ext, &appleData, &gameConfig.LastAppleUpdateTime)
		if err != nil {
			log.Error().Msgf("scan game config error: %s", err.Error())
			continue
		}

		var appleConfig base.AppleConfig
		if appleData != "" {
			err = json.Unmarshal([]byte(appleData), &appleConfig)
			if err != nil {
				log.Error().Msgf("decode apple_config into struct failed, json unmarshall error: %s, apple data: %s", err.Error(), appleData)
				continue
			}
		}
		//2.开启协程
		//如果超过有效期一半,则生成新的client_secret,并写入到game_config表中
		//否则直接将game_config表数据写入内存中
		go saveGameConfigIntoMemory(appleConfig, gameConfig)
	}
}

func saveGameConfigIntoMemory(appleConfig base.AppleConfig, gameConfig base.GameConfig) {
	if gameConfig.GameId <= 0 || gameConfig.PlatformId <= 0 {
		log.Error().Msgf("saveGameConfigIntoMemory, invalid game_id or platform_id, value: %v", gameConfig)
		return
	}
	execTime := time.Now().Unix()
	//超过有效期一半
	if appleConfig.PrivateKey != "" && time.Second*time.Duration(execTime-gameConfig.LastAppleUpdateTime) > time.Hour*24*time.Duration(appleConfig.ValidityPeriod/2) {
		makeSecretAndUpdateConfig(appleConfig, &gameConfig, execTime)
	}
	key := fmt.Sprintf("_account_game_config_%d_%d", gameConfig.GameId, gameConfig.PlatformId)
	base.MemoryGameConfig.Store(key, gameConfig)
	log.Info().Interface("_account_game_config", gameConfig).Msgf("saveGameConfigIntoMemory, save game config success!!! key: %s", key)
	return
}

func makeSecretAndUpdateConfig(appleConfig base.AppleConfig, gameConfig *base.GameConfig, execTime int64) {
	log.Info().Msgf("makeSecretAndUpdateConfig, start, execTime: %d, lastUpdateTime: %d", execTime, gameConfig.LastAppleUpdateTime)
	//生成jwt格式的client secret
	claims := &jwt.StandardClaims{
		Issuer:    appleConfig.TeamId,
		IssuedAt:  execTime,
		ExpiresAt: execTime + appleConfig.ValidityPeriod*86400, // ValidityPeriod days, 最大180
		Audience:  "https://appleid.apple.com",
		Subject:   appleConfig.ClientId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["alg"] = "ES256"
	token.Header["kid"] = appleConfig.KeyId
	block, _ := pem.Decode([]byte(appleConfig.PrivateKey))
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Error().Msgf("makeSecretAndUpdateConfig, make jwt client secret error, empty block after decoding")
		return
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Error().Msgf("makeSecretAndUpdateConfig, make jwt client secret error, parse error: %s", err.Error())
		return
	}
	//直接使用appleConfig.PrivateKey去sign,在es256算法下,会报错key is of invalid type
	clientSecret, err := token.SignedString(privateKey)
	if err != nil {
		log.Error().Msgf("makeSecretAndUpdateConfig, make jwt client secret error, sign error: %s", err.Error())
		return
	}
	//将apple_client_id, apple_client_secret,last_apple_update_time更新写入game_config表
	_, err = base.AccountBaseDb.Exec(fmt.Sprintf("UPDATE %s SET apple_client_id = ?, apple_client_secret = ?, last_apple_update_time = ? WHERE game_id = %d AND platform_id = %d", base.UserDeleteConfigTable, gameConfig.GameId, gameConfig.PlatformId), appleConfig.ClientId, clientSecret, execTime)
	if err != nil {
		log.Error().Msgf("makeSecretAndUpdateConfig, exec error, error: %s", err.Error())
		return
	}
	log.Info().Msgf("makeSecretAndUpdateConfig, refresh apple client id and client secret into table success, game_id: %d, platform: %d, apple_client_id: %s, apple_client_secret: %s", gameConfig.GameId, gameConfig.PlatformId, appleConfig.ClientId, clientSecret)
	gameConfig.AppleClientId = appleConfig.ClientId
	gameConfig.AppleClientSecret = clientSecret
	return
}

// RefreshUserDelete 定时处理项目账号注销和恢复
func RefreshUserDelete() {
	dbSlaveMap := base.GetAllGameUserDb(base.GameUserSlaveDb)
	dbMasterMap := base.GetAllGameUserDb(base.GameUserMasterDb)
	base.MemoryGameConfig.Range(func(key, value interface{}) bool {
		gameConfig, ok := value.(base.GameConfig)
		if !ok {
			log.Error().Msgf("RefreshUserDelete, game config assert failed, key: %s, value: %v", key, value)
			return true
		}
		if gameConfig.GameId <= 0 || gameConfig.PlatformId <= 0 {
			log.Error().Msgf("RefreshUserDelete, game config invalid value, key: %s, value: %v", key, value)
			return true
		}
		l := len(base.GConf.MysqlGameUserSlaveList)
		for i := 1; i <= l; i++ {
			slaveDbKey := fmt.Sprintf("slave_db_%d_%d_%d", gameConfig.GameId, gameConfig.PlatformId, i)
			masterDbKey := fmt.Sprintf("master_db_%d_%d_%d", gameConfig.GameId, gameConfig.PlatformId, i)
			slaveDb := dbSlaveMap[slaveDbKey]
			masterDb := dbMasterMap[masterDbKey]
			for k := 1; k <= base.GameUserDeleteTableNumber; k++ {
				userDeleteHandle(masterDb, slaveDb, k, &gameConfig)
				//处理需要恢复的账号
				userRecoverHandle(masterDb, slaveDb, k, &gameConfig)
			}
		}
		return true
	})
}

func userDeleteHandle(masterDb *sql.DB, slaveDb *sql.DB, tableIndex int, gameConfig *base.GameConfig) bool {
	//处理需要注销的账号
	execTime := base.GetTime()
	applySql := fmt.Sprintf("SELECT uid, main_uid, ext_info, account FROM user_delete_apply_%d WHERE `status` = %d AND execute_delete_time < %d", tableIndex, base.ApplyStatusPending, execTime)
	rows, err := slaveDb.Query(applySql)
	if err != nil {
		log.Error().Msgf("RefreshUserDelete, query apply error, error: %s", err.Error())
		return false
	}
	defer rows.Close()
	var applyInfo UserDeleteApply
	for rows.Next() {
		err = rows.Scan(&applyInfo.Uid, &applyInfo.MainUid, &applyInfo.ExtInfo, &applyInfo.Account)
		if err != nil {
			log.Error().Msgf("RefreshUserDelete, scan apply error, error: %s", err.Error())
			continue
		}

		go userDeleteTransaction(masterDb, &applyInfo, gameConfig, execTime, tableIndex)
	}
	return true
}

func userRecoverHandle(masterDb *sql.DB, slaveDb *sql.DB, tableIndex int, gameConfig *base.GameConfig) bool {
	recoverSql := fmt.Sprintf("SELECT main_uid, uid, account, type, ext FROM user_delete_apply_%d WHERE `status` = %d", tableIndex, base.ApplyStatusRecover)
	rows, err := slaveDb.Query(recoverSql)
	if err != nil {
		log.Error().Msgf("RefreshUserDelete, query delete log error, error: %s", err.Error())
		return true
	}
	defer rows.Close()
	var info UserDeleteApply
	execTime := base.GetTime()
	for rows.Next() {
		err = rows.Scan(&info.MainUid, &info.Uid, &info.Account, &info.Type, &info.Ext)
		if err != nil {
			log.Error().Msgf("RefreshUserDelete, scan delete error, error: %s", err.Error())
			continue
		}
		go userRecoverTransaction(masterDb, &info, gameConfig, execTime, tableIndex)
	}

	return true
}

// 协程处理账号注销
func userDeleteTransaction(masterDb *sql.DB, applyInfo *UserDeleteApply, gameConfig *base.GameConfig, execTime int64, tableIndex int) {
	userLog := log.With().Str("req_id", fmt.Sprintf("script_delete_%d_%d_%d", gameConfig.GameId, gameConfig.PlatformId, applyInfo.Uid)).Logger()
	userLog.Info().Int64("uid", applyInfo.Uid).Msgf("delete trans, start, applyInfo: %v, gameConfig: %v, execTime: %d", applyInfo, gameConfig, execTime)

	dbTx, _ := masterDb.Begin()
	//1.更新status(从1到2)
	//判断affected,多服务并发时,只有一个服务去处理
	updateSql := fmt.Sprintf("UPDATE user_delete_apply_%d SET `status` = %d WHERE uid = ? AND `status` = %d AND execute_delete_time < %d", tableIndex, base.ApplyStatusSuccess, base.ApplyStatusPending, execTime)
	res, err := dbTx.Exec(updateSql, applyInfo.Uid)
	if err != nil {
		userLog.Error().Int64("uid", applyInfo.Uid).Msgf("delete trans, update apply, exec error: %s", err.Error())
		dbTx.Rollback()
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		userLog.Error().Int64("uid", applyInfo.Uid).Msgf("delete trans, update apply, rows affected error: %s", err.Error())
		dbTx.Rollback()
		return
	}
	if int(affected) < 1 {
		userLog.Error().Int64("uid", applyInfo.Uid).Msgf("delete trans, update apply, affected: [%d]", affected)
		dbTx.Rollback()
		return
	}

	gameUserTable := base.GetGameUserDeleteTable(applyInfo.MainUid)

	//2.删除项目用户表
	deleteSql := fmt.Sprintf("DELETE FROM %s WHERE main_uid = ?", gameUserTable)
	res, err = dbTx.Exec(deleteSql, applyInfo.MainUid)
	if err != nil {
		userLog.Error().Int64("uid", applyInfo.Uid).Msgf("delete trans, delete user, exec error: %s", err.Error())
		dbTx.Rollback()
		return
	}

	//3.更新delete_apply表
	deleteSql = fmt.Sprintf("UPDATE user_delete_apply_%d SET `status` = ? WHERE main_uid = ? AND `status` = ? AND execute_delete_time < ?", tableIndex)
	res, err = dbTx.Exec(deleteSql, base.ApplyStatusDeleted, applyInfo.MainUid, base.ApplyStatusSuccess, execTime)
	if err != nil {
		userLog.Error().Int64("uid", applyInfo.Uid).Msgf("delete trans, delete user apply, exec error: %s", err.Error())
		dbTx.Rollback()
		return
	}

	dbTx.Commit()
	userLog.Info().Int64("uid", applyInfo.Uid).Msgf("delete trans, user [%d], transaction success!!!", applyInfo.Uid)
	//7.调用苹果revoke,判断结果,不参与事务提交
	revokeRes := appleAuthRevoke(applyInfo, gameConfig)
	if !revokeRes {
		userLog.Error().Int64("uid", applyInfo.Uid).Msgf("delete trans, call apple revoke token failed")
	}
	return
}

// 调用苹果revoke接口
func appleAuthRevoke(applyInfo *UserDeleteApply, gameConfig *base.GameConfig) bool {
	userLog := log.With().Str("req_id", fmt.Sprintf("script_delete_%d_%d_%d", gameConfig.GameId, gameConfig.PlatformId, applyInfo.Uid)).Logger()
	userLog.Info().Int64("uid", applyInfo.Uid).Msgf("call apple, start")
	//1.如果ext_info解json失败,或目标字段值为空字符串
	//则视为普通(非苹果)用户注销,直接返回成功
	data := &base.DeleteAccountExt{}
	//json.Valid()
	if applyInfo.ExtInfo != "" {
		err := json.Unmarshal([]byte(applyInfo.ExtInfo), data)
		if err != nil {
			userLog.Info().Int64("uid", applyInfo.Uid).Msgf("call apple, user apply ext_info: %s, json unmarshall error: %s, return success", applyInfo.ExtInfo, err.Error())
			return true
		}
	}
	if data.AppleRefreshToken == "" {
		userLog.Info().Int64("uid", applyInfo.Uid).Msgf("call apple, refresh token empty, return success")
		return true
	}
	if gameConfig.AppleClientId == "" || gameConfig.AppleClientSecret == "" {
		userLog.Error().Int64("uid", applyInfo.Uid).Msgf("call apple, game config missing data, client id: %s, secret: %s", gameConfig.AppleClientId, gameConfig.AppleClientSecret)
		return false
	}
	//2.调用revoke接口
	postData := map[string]string{
		"client_id":       gameConfig.AppleClientId,
		"client_secret":   gameConfig.AppleClientSecret,
		"token_type_hint": "refresh_token",
		"token":           data.AppleRefreshToken,
	}
	res, err := base.HttpPostForm("https://appleid.apple.com/auth/revoke", postData)
	userLog.Info().Int64("uid", applyInfo.Uid).Interface("apple_revoke_params", postData).RawJSON("apple_revoke_return", res).Msg("apple auth revoke info")
	if err != nil {
		userLog.Error().Int64("uid", applyInfo.Uid).Msgf("call apple, post failed, error log: %s", err.Log)
		return false
	}
	//返回空字符串,代表成功
	if string(res) == "" {
		userLog.Info().Int64("uid", applyInfo.Uid).Msgf("call apple, success, http res empty string")
		return true
	}
	//其他情况,代表失败
	userLog.Error().Int64("uid", applyInfo.Uid).Msgf("call apple, res not empty string, res: %s", res)
	return false
}

// 协程处理账号恢复
func userRecoverTransaction(gameMasterDb *sql.DB, deleteInfo *UserDeleteApply, gameConfig *base.GameConfig, execTime int64, tableIndex int) {
	table := fmt.Sprintf("user_delete_apply_%d", tableIndex)
	userLog := log.With().Str("req_id", fmt.Sprintf("script_recover_%d_%d_%d", gameConfig.GameId, gameConfig.PlatformId, deleteInfo.MainUid)).Logger()
	userLog.Info().Int64("uid", deleteInfo.Uid).Msgf("recover trans, start, deleteInfo: %v, gameConfig: %v", deleteInfo, gameConfig)

	gameUserTx, _ := gameMasterDb.Begin()
	//1.更新status(从4到5)
	//判断affected,多服务并发时,只有一个服务去处理
	updateSql := fmt.Sprintf("UPDATE %s SET `status` = %d WHERE account = ? AND `status` = %d", table, base.ApplyStatusRecoverSuccess, base.ApplyStatusRecover)
	res, err := gameUserTx.Exec(updateSql, deleteInfo.Account)
	if err != nil {
		userLog.Error().Int64("uid", deleteInfo.Uid).Msgf("recover trans, update delete log, exec error: %s", err.Error())
		gameUserTx.Rollback()
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		userLog.Error().Int64("uid", deleteInfo.Uid).Msgf("recover trans, update delete log, rows affected error: %s", err.Error())
		gameUserTx.Rollback()
		return
	}
	if int(affected) != 1 {
		userLog.Error().Int64("uid", deleteInfo.Uid).Msgf("recover trans, update delete log, affected: [%d] not 1", affected)
		gameUserTx.Rollback()
		return
	}
	userTable := base.GetGameUserTable(deleteInfo.MainUid)
	//2.恢复users表数据
	insertSql := fmt.Sprintf("INSERT INTO %s (account, main_uid, uid, created_time, ext) SELECT account, main_uid, uid, %d created_time, ext FROM %s WHERE account = ?", userTable, execTime, table)
	res, err = gameUserTx.Exec(insertSql, deleteInfo.Account)
	if err != nil {
		userLog.Error().Int64("uid", deleteInfo.Uid).Msgf("recover trans, copy user delete log into user, exec error: %s", err.Error())
		gameUserTx.Rollback()
		return
	}
	affected, err = res.RowsAffected()
	if err != nil {
		userLog.Error().Int64("uid", deleteInfo.Uid).Msgf("recover trans, copy user delete log into user, rows affected error: %s", err.Error())
		gameUserTx.Rollback()
		return
	}
	if int(affected) != 1 {
		userLog.Error().Int64("uid", deleteInfo.Uid).Msgf("recover trans, copy user delete log into user, affected: [%d] not 1", affected)
		gameUserTx.Rollback()
		return
	}
	//3.删除数据
	deleteSql := fmt.Sprintf("DELETE FROM %s WHERE account = ?", table)
	res, err = gameUserTx.Exec(deleteSql, deleteInfo.Account)
	if err != nil {
		userLog.Error().Int64("uid", deleteInfo.Uid).Msgf("recover trans, delete user delete log, exec error: %s", err.Error())
		gameUserTx.Rollback()
		return
	}
	affected, err = res.RowsAffected()
	if err != nil {
		userLog.Error().Int64("uid", deleteInfo.Uid).Msgf("recover trans, delete user delete log, rows affected error: %s", err.Error())
		gameUserTx.Rollback()
		return
	}
	if int(affected) != 1 {
		userLog.Error().Int64("uid", deleteInfo.Uid).Msgf("recover trans, delete user delete log, affected: [%d] not 1", affected)
		gameUserTx.Rollback()
		return
	}
	gameUserTx.Commit()
	userLog.Info().Int64("uid", deleteInfo.Uid).Msgf("recover trans, user [%v], transaction success!!!", deleteInfo)
	return
}
