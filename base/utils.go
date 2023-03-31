/**
 * @project Accounts
 * @filename utils.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/19 15:02
 * @version 1.0
 * @description
 * 常用方法
 */

package base

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"hash/crc32"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 失败输出
func ResponseFail(w http.ResponseWriter, myError *MyError, logger zerolog.Logger) {
	resp := &AccountResponse{}
	resp.Code = myError.Code
	if myError.Error() != "" {
		resp.Msg = myError.Error()
	} else {
		resp.Msg = ErrorMsg[resp.Code]
	}

	resp.Data = EmptyData
	if myError.Data != nil {
		resp.Data = myError.Data
	}
	body, err := json.Marshal(resp)
	if err != nil {
		logger.Error().Msgf("response json error %s", err.Error())
	}

	startTimeStr := w.Header().Get("StartTime")
	w.Header().Del("StartTime")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(body)

	if err != nil {
		logger.Error().Msgf("response write error %s", err.Error())
	}

	//日志打印耗时
	startTime, _ := strconv.ParseInt(startTimeStr, 10, 64)
	var timeTotal int64
	if startTime != 0 {
		timeTotal = time.Now().UnixNano()/1e6 - startTime
	}

	logger.Error().Int("error_code", myError.Code).Int64("time_total", timeTotal).Interface("return_info", resp).Msg(myError.Log)
}

// 成功输出
func ResponseOK(w http.ResponseWriter, returnData interface{}, logger zerolog.Logger) {
	resp := &AccountResponse{}
	resp.Code = Success
	resp.Msg = ErrorMsg[Success]
	resp.Data = returnData
	body, err := json.Marshal(resp)
	if err != nil {
		logger.Error().Msgf("response json error %s", err.Error())
	}

	startTimeStr := w.Header().Get("StartTime")
	w.Header().Del("StartTime")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(body)

	if err != nil {
		logger.Error().Msgf("response write error %s", err.Error())
	}

	//日志打印耗时
	startTime, _ := strconv.ParseInt(startTimeStr, 10, 64)
	var timeTotal int64
	if startTime != 0 {
		timeTotal = time.Now().UnixNano()/1e6 - startTime
	}

	logger.Info().Int("error_code", Success).Int64("time_total", timeTotal).Interface("return_info", resp).Msg("")
}

// 日志添加信息
func (h RequestHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Interface("req_body", h.RequestBody).Str("req_ip", h.IP).Str("hgp", h.HeaderGamePlatform)

	if h.GameId != 0 {
		e.Int("game_id", h.GameId)
	}
	if h.Uid != 0 {
		e.Int64("uid", h.Uid)
	}
}

// RequestHandler 请求处理: 验证请求方式; 解析数据; 字段校验
func RequestHandler(r *http.Request, data interface{}) *MyError {
	if r.Method != "POST" {
		return &MyError{Code: NotPostRequest}
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &MyError{Code: RequestDataIncorrect, Log: fmt.Sprintf("io.ReadAll error %s: %s", err.Error(), body)}
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return &MyError{Code: RequestDataParserError, Log: fmt.Sprintf("json.Unmarshal error %s: %s", err.Error(), body)}
	}

	//校验每个字段
	validate := validator.New()
	err = validate.Struct(data)
	if err != nil {
		return &MyError{Code: RequestDataValidatorFail, Log: fmt.Sprintf("validate.Struct error %s: %s", err.Error(), body)}
	}

	//ip 限制访问策略
	limitErr := LimitIp(GetRealAddr(r).String(), r.URL.Path)
	if limitErr != nil {
		return limitErr
	}

	return nil
}

// 连接数据库
func connMysql(config MysqlConfig) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&charset=utf8", config.User, config.Pass, config.Host, config.Port, config.DbName, GConf.MysqlTimeout.MysqlTimeout, GConf.MysqlTimeout.MysqlReadTimeout, GConf.MysqlTimeout.MysqlWriteTimeout)

	db, MysqlError := sql.Open("mysql", dsn)
	if MysqlError != nil {
		MultipleLog.Fatal().Msgf("dsn: %s Open fail, error: %s", dsn, MysqlError.Error())
		return nil
	}

	db.SetMaxOpenConns(config.MaxConn)
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)

	MysqlError = db.Ping()
	if MysqlError != nil {
		MultipleLog.Fatal().Msgf("dsn: %s Ping fail, error: %s", dsn, MysqlError.Error())
		return nil
	}
	return db
}

// 获取app信息
func getAppIdInfo(appId int64, appType int) (AppIdConfig, *MyError) {
	key := fmt.Sprintf("_account_app_id_%d", appId)
	value, _ := MemoryStoreInfo.Load(key)
	var appInfo AppIdConfig
	if value == nil {
		RefreshApiAppKey(false)
		value, _ = MemoryStoreInfo.Load(key)
		if value == nil {
			return appInfo, &MyError{Code: AppIdQueryError, Log: fmt.Sprintf("app id parser fail, app_id=%s, value:%v", key, value)}
		}
	}

	var ok bool
	appInfo, ok = value.(AppIdConfig)
	if ok == false {
		return appInfo, &MyError{Code: AppIdError, Log: fmt.Sprintf("app id parser fail, app_id=%s, value:%v", key, value)}
	}

	if appInfo.Type != appType {
		return appInfo, &MyError{Code: AppIdTypeError, Log: fmt.Sprintf("app id type error, app_id=%s, value:%v", key, appInfo)}
	}

	return appInfo, nil
}

// 刷新app_key
func RefreshApiAppKey(isFirst bool) {
	appIdMap, err := GetAllConfFiles(GConf.Base.AppIdConfPath)
	if err != nil {
		//如果时启动时加载, 有错误则停止
		if isFirst {
			MultipleLog.Fatal().Err(err).Msgf("load appid config error: %s", err.Error())
			return
		}
		log.Error().Err(err).Msgf("load appid config error: %s", err.Error())
		return
	}
	for appId, appInfo := range appIdMap {
		key := fmt.Sprintf("_account_app_id_%d", appId)
		MemoryStoreInfo.Store(key, appInfo)
		log.Info().Interface("_account_app_id_", appInfo).Msg(key)
	}
}

// 签名校验
func SignValidator(appId int64, sign string, gameId int, dataFields interface{}, appType int) *MyError {
	appInfo, err := getAppIdInfo(appId, appType)
	if err != nil {
		return err
	}

	if appInfo.GameId != gameId {
		return &MyError{Code: AppIdQueryError, Log: fmt.Sprintf("RAM app game_id %d, request game_id %d", appInfo.GameId, gameId)}
	}

	values := StructToString(dataFields)
	localSignStr := fmt.Sprintf("%s&%s", values, appInfo.SecretKey)
	localSign := Md5Sum([]byte(localSignStr))
	if sign != localSign {
		return &MyError{Code: SignError, Log: fmt.Sprintf("request sign %s, local sign %s, local sign string %s", sign, localSign, localSignStr)}
	}

	//sdk调用的接口，检查game_id是否存在
	if appType == AppIdTypeSdk {
		gameExist := false
		for _, id := range GConf.Base.EnabledGameList {
			if id == gameId {
				gameExist = true
				break
			}
		}
		if !gameExist {
			return &MyError{Code: GameIdNotExists, Log: "game_id not exist"}
		}
	}

	return nil
}

// 生成登录Token
func BuildLoginToken(gameId int, platformId int, channelId int, loginRet *LoginReturnFields, tokenType int, expireTime int64) (string, *MyError) {
	signingKey := []byte(GConf.Base.ServerKey)

	claims := CustomClaims{
		gameId,
		platformId,
		channelId,
		loginRet.Uid,
		loginRet.LoginTime,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "account_server",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(signingKey)
	if err != nil {
		return "", &MyError{Code: BuildTokenFailure, Log: "jwt build token error: " + err.Error()}
	}

	return token, nil
}

// 检查登录token内的uid与传入的uid是否一致
func LoginTokenCheck(loginToken string, uid int64) *MyError {
	//检查token解析，是否过期
	token, tokenErr := jwt.ParseWithClaims(loginToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GConf.Base.ServerKey), nil
	})
	if tokenErr != nil {
		return &MyError{Code: LoginTokenParseError, Log: fmt.Sprintf("login token uid: %d, auth error: %s", uid, tokenErr)}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.Uid != uid {
			return &MyError{Code: LoginTokenUidUnequal, Log: fmt.Sprintf("bind account LoginToken parse uid: %d, params uid: %d", claims.Uid, uid)}
		}
	}
	return nil
}

// 生成指定长度的随机数，不超过16位，验证码使用
func GetRandom(len int) string {
	if len > 16 {
		return ""
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	lenFormat := int64(math.Pow10(len))
	return fmt.Sprintf(fmt.Sprintf("%%0%dv", len), rnd.Int63n(lenFormat))
}

// 生成随机字符串，密码盐
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = PasswordSaltChar[rand.Intn(len(PasswordSaltChar))]
	}
	return string(b)
}

// 获取当前时间，秒
func GetTime() int64 {
	return time.Now().Unix()
}

// 检查用户账号格式
func CheckUserAccountFormat(account string, accountType int) (string, *MyError) {
	switch accountType {
	case AccountEmail: //检查邮箱格式
		reg := regexp.MustCompile(`^[\w-]+(\.[\w-]*)*@[\w-]+(\.[\w-]+)+$`)
		if !reg.MatchString(account) {
			return account, &MyError{Code: EmailFormatError}
		}
	case AccountMobile: //检查手机号格式
		//去掉+号和空格
		account = strings.Replace(account, "+", "", -1)
		//account = strings.Replace(account, " ", "", -1)

		reg := regexp.MustCompile(`^\d{5,}$`)
		if !reg.MatchString(account) {
			return account, &MyError{Code: PhoneNumFormatError}
		}
	case AccountUsername: //注册类型为用户名时，检查长度，检查格式
		if len(account) < 4 || len(account) > 60 {
			return account, &MyError{Code: UsernameLengthError}
		}

		reg := regexp.MustCompile(`^[\w\d-_.]{4,60}$`)
		if !reg.MatchString(account) {
			return account, &MyError{Code: UsernameFormatError}
		}
	case AccountGuest, AccountThird: //注册类型为游客和第三方，检查长度，检查格式
		if len(account) < 4 || len(account) > 128 {
			return account, &MyError{Code: GuestOrThirdLengthError}
		}

		reg := regexp.MustCompile(`^[\w\d-_.]{4,128}$`)
		if !reg.MatchString(account) {
			return account, &MyError{Code: GuestOrThirdFormatError}
		}
		if accountType == AccountThird {
			third := strings.Split(account, "_")
			if len(third) < 2 {
				return account, &MyError{Code: ThirdFormatError}
			}

			id, err := strconv.Atoi(third[0])
			if err != nil {
				return account, &MyError{Code: ThirdIdParseFailure}
			}
			if _, ok := ThirdIds[id]; !ok {
				return account, &MyError{Code: ThirdIdUnsupported}
			}

			if third[1] == "" {
				return account, &MyError{Code: ThirdUidEmpty, Log: "thirdUid empty"}
			}
		}
	}
	return account, nil
}

// form表单http请求
func HttpPostForm(postUrl string, postData map[string]string) ([]byte, *MyError) {
	postDataReal := url.Values{}
	for k, v := range postData {
		postDataReal.Set(k, v)
	}
	resp, err := http.PostForm(postUrl, postDataReal)
	if err != nil {
		return nil, &MyError{Code: PostFormRequestError, Log: fmt.Sprintf("make post form request error: %s", err.Error())}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &MyError{Code: PostFormReadError, Log: fmt.Sprintf("post form request read resp body error: %s", err.Error())}
	}

	return body, nil
}

func GetRealAddr(r *http.Request) net.IP {
	// parse X-ORIGINAL-FORWARDED-FOR header
	if xoff := strings.Trim(r.Header.Get("X-ORIGINAL-FORWARDED-FOR"), ","); len(xoff) > 0 {
		addrs := strings.Split(xoff, ",")
		if ip := net.ParseIP(addrs[0]); ip != nil {
			return ip
		}
	}

	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		if ip := net.ParseIP(addrs[0]); ip != nil {
			return ip
		}
	}
	// parse X-Real-Ip header
	if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			return ip
		}
	}

	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) >= 1 {
		return net.ParseIP(parts[0])
	}
	return nil
}

// GetUnixMilliString 获取当前毫秒数字符串
func GetUnixMilliString() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

// UnmarshalToml 解析toml配置到到out
func UnmarshalToml(file string, out interface{}) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return toml.Unmarshal(content, out)
}

// Md5Sum 获取md5值
func Md5Sum(data []byte) string {
	sign := md5.Sum(data)
	return fmt.Sprintf("%x", sign)
}

// 结构体转字符串
func StructToString(data interface{}) string {
	if data == nil {
		return ""
	}
	var values []string
	iVal := reflect.ValueOf(data).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		tag := typ.Field(i).Tag.Get("json")
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', -1, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', -1, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			if tag == "sign" {
				continue
			}
			v = f.String()
		default:
			if f.Kind() == reflect.Struct {
				values = append(values, StructToString(f.Addr().Interface()))
				continue
			}
		}
		if v != "" {
			values = append(values, v)
		}
	}

	return strings.Join(values, "&")
}

// 数据落地
func DataExtLog(accountId int64, dataExt, eventId, ip string, dataLog zerolog.Logger) {
	zone := time.Now().Format("-07")
	zoneNum, _ := strconv.Atoi(zone)

	dataLog.Log().Str("#account_id", strconv.FormatInt(accountId, 10)).
		Str("#time", time.Now().Format(DateTimeFormat)).
		Str("#ip", ip).
		Str("#event_name", eventId).
		Str("#type", "track").
		Int("#zone_offset", zoneNum).
		RawJSON("properties", json.RawMessage(dataExt)).
		Msg("")
}

// GetAllConfFiles 读取path目录下所有.json文件,解析到 configMap中去
func GetAllConfFiles(path string) (map[int64]AppIdConfig, error) {
	configMap := map[int64]AppIdConfig{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := file.Name()
		lowerName := strings.ToLower(name)
		if strings.Index(lowerName, ".json") == -1 || lowerName[len(lowerName)-5:] != ".json" {
			continue
		}

		filename := path + "/" + name
		appInfo := AppIdConfig{}
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Error().Err(err).Msgf("read name: %s error, continue", name)
			continue
		}
		err = json.Unmarshal(content, &appInfo)
		if err != nil {
			log.Error().Err(err).Msgf("name: %s parse error, continue", name)
			continue
		}
		log.Info().Msgf("load name: %s to configMap success", name)
		configMap[appInfo.AppId] = appInfo
	}
	if len(configMap) < 1 {
		return nil, errors.New("configMap length = 0")
	}

	return configMap, nil
}

// GetAccountDbHashId 主账号分库
func GetAccountDbHashId(account string) uint32 {
	return crc32Algo(account)%MainAccountDbNumber + 1
}

// GetAccountTableHashId 主账号分表
func GetAccountTableHashId(uid int64) uint32 {
	return crc32Algo(strconv.FormatInt(uid, 10))%MainAccountTableNumber + 1
}

// hash表
func GetAccountHashTableHashId(account string) uint32 {
	return crc32Algo(account)%MainAccountHashTableNumber + 1
}

// GetGameDbHashId 项目用户分库
func GetGameDbHashId(uid int64) uint32 {
	return crc32Algo(strconv.FormatInt(uid, 10))%GameUserDbNumber + 1
}

// GetGameTableHashId 项目用户分表
func GetGameTableHashId(uid int64) uint32 {
	return crc32Algo(strconv.FormatInt(uid, 10))%GameUserTableNumber + 1
}

// 项目用户申请删除及删除记录表hash
func GetGameUserDeleteTableHashId(uid int64) uint32 {
	return crc32Algo(strconv.FormatInt(uid, 10))%GameUserDeleteTableNumber + 1
}

// 统一hash算法, md5后crc32
func crc32Algo(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(Md5Sum([]byte(s))))
}

// 主账号master db
func GetAccountMasterDb(account string) *sql.DB {
	id := GetAccountDbHashId(account)
	dbKey := fmt.Sprintf("account_master_db_%d", id)
	db, ok := AccountMasterDbMap.Load(dbKey)
	if !ok {
		log.Info().Msgf("get account master db error by %s", account)
		os.Exit(1)
	}
	return db.(*sql.DB)
}

// 主账号slave db
func GetAccountSlaveDb(accountOrUid string) *sql.DB {
	id := GetAccountDbHashId(accountOrUid)
	dbKey := fmt.Sprintf("account_slave_db_%d", id)
	db, ok := AccountSlaveDbMap.Load(dbKey)
	if !ok {
		log.Info().Msgf("get account slave db error by %s", accountOrUid)
		os.Exit(1)
	}
	return db.(*sql.DB)
}

// 主账号表
func GetAccountTable(uid int64) string {
	return fmt.Sprintf("account_%d", GetAccountTableHashId(uid))
}

// 主账号表
func GetAccountHashTable(account string) string {
	return fmt.Sprintf("account_hash_%d", GetAccountHashTableHashId(account))
}

// GetGameUserMasterDb 项目master db
func GetGameUserMasterDb(uid int64, gameId, platformId int) *sql.DB {
	key := fmt.Sprintf("master_db_%d_%d_%d", gameId, platformId, GetGameDbHashId(uid))
	db, ok := GameUserMasterDbMap.Load(key)
	if !ok {
		log.Info().Msgf("get game user master db error by %d", uid)
		os.Exit(1)
	}
	return db.(*sql.DB)
}

// 项目slave db
func GetGameUserSlaveDb(uid int64, gameId, platformId int) *sql.DB {
	key := fmt.Sprintf("slave_db_%d_%d_%d", gameId, platformId, GetGameDbHashId(uid))
	db, ok := GameUserSlaveDbMap.Load(key)
	if !ok {
		log.Info().Msgf("get game user slave db error by %d", uid)
		os.Exit(1)
	}
	return db.(*sql.DB)
}

func GetAllGameUserDb(dbType int) map[string]*sql.DB {
	if dbType != GameUserSlaveDb && dbType != GameUserMasterDb {
		return nil
	}
	dbList := GConf.MysqlGameUserSlaveList
	gameUserDbMap := GameUserSlaveDbMap
	if dbType == GameUserMasterDb {
		dbList = GConf.MysqlGameUserMasterList
		gameUserDbMap = GameUserMasterDbMap
	}
	var dbMap = map[string]*sql.DB{}
	for key := range dbList {
		db, ok := gameUserDbMap.Load(key)
		if !ok {
			log.Error().Msgf("get game user slave db error by %s", key)
			continue
		}
		dbMap[key] = db.(*sql.DB)
	}

	return dbMap
}

// 项目用户表
func GetGameUserTable(uid int64) string {
	return fmt.Sprintf("user_%d", GetGameTableHashId(uid))
}

// 账号注销申请表
func GetGameUserDeleteTable(uid int64) string {
	return fmt.Sprintf("user_delete_apply_%d", GetGameUserDeleteTableHashId(uid))
}

// 根据uid获取相关的db, table
func GetDbTable(uid int64, gameId, platformId int) *DbTable {
	dbTable := &DbTable{
		AccountMasterDb:          GetAccountMasterDb(strconv.FormatInt(uid, 10)),
		AccountSlaveDb:           GetAccountSlaveDb(strconv.FormatInt(uid, 10)),
		AccountTable:             GetAccountTable(uid),
		GameUserTable:            GetGameUserTable(uid),
		GameUserDeleteApplyTable: GetGameUserDeleteTable(uid),
	}
	if gameId != -1 && platformId != -1 {
		dbTable.GameUserMasterDb = GetGameUserMasterDb(uid, gameId, platformId)
		dbTable.GameUserSlaveDb = GetGameUserSlaveDb(uid, gameId, platformId)
	}
	return dbTable
}

// 根据账号获取hash db, table
func GetHashDbTable(account string) *HashDbTable {
	hashDbTable := &HashDbTable{
		AccountMasterDb:  GetAccountMasterDb(account),
		AccountSlaveDb:   GetAccountSlaveDb(account),
		AccountHashTable: GetAccountHashTable(account),
	}

	return hashDbTable
}

// 身份证校验
func CheckCardId(cardId string) bool {
	if strings.HasSuffix(cardId, "x") {
		cardId = strings.ToUpper(cardId)
	}
	l := len(cardId)
	if l != 18 {
		return false
	}
	weight := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	validate := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
	sum := 0
	for i := 0; i < len(weight); i++ {
		sum += weight[i] * int(byte(cardId[i])-'0')
	}
	m := sum % 11
	return validate[m] == cardId[l-1]
}
