/**
 * @project Accounts
 * @filename defs.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/18 10:22
 * @version 1.0
 * @description
 * 变量、常亮、协议定义
 */

package base

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog"
	"sync"
)

// 配置文件路径
var ConfFile string

// 配置
var GConf = &AccountConf{}

// 空数据定义
var EmptyData = map[string]interface{}{}

// AccountBaseDb 基础库
var AccountBaseDb *sql.DB

// AccountMasterDbMap 主账号master db map
var AccountMasterDbMap sync.Map

// 主账号从库 db map
var AccountSlaveDbMap sync.Map

// 项目账号Master DB map
var GameUserMasterDbMap sync.Map

// 项目账号Slave DB map
var GameUserSlaveDbMap sync.Map

// Redis对象
var RedisClient *redis.Client

// 内存中存储信息
var MemoryStoreInfo sync.Map

// 游戏配置信息
var MemoryGameConfig sync.Map

// 记录登录注册数据的日志
var LoginDataLog zerolog.Logger
var RegisterDataLog zerolog.Logger

// 输出到控制台+日志文件
var MultipleLog zerolog.Logger

// 账号类型
const (
	//没值时传入的默认字符串
	DefaultNoValue = "-1"

	//AppId 类型, sdk, 客户端（白名单），服务器（登录校验）
	AppIdTypeSdk    = 1
	AppIdTypeClient = 2
	AppIdTypeServer = 3

	//密码盐长度
	PasswordSaltLength = 16

	//主账号库数量及表数量
	MainAccountDbNumber    = 3
	MainAccountTableNumber = 100

	// MainAccountHashTableNumber 绑定关系表数量，如游客绑定手机号、email
	MainAccountHashTableNumber = 100

	//项目用户库数量、表数量、第三方账号数量、注销账号表等相关表数量
	GameUserDbNumber          = 2
	GameUserTableNumber       = 100
	GameUserDeleteTableNumber = 5

	//账号类型
	AccountEmail    int = 1
	AccountMobile   int = 2
	AccountUsername int = 3
	AccountGuest    int = 4
	AccountThird    int = 5

	//生成 access token类型
	TokenTypeAccess int = 1
	//生成 refresh token类型
	TokenTypeRefresh int = 2

	//项目账号状态
	AccountNormal   = 1 //正常状态
	AccountDisabled = 2 //被禁用
	AccountDeleting = 3 //账号注销中

	//账号注销申请状态
	ApplyStatusPending        = 1 //等待注销中（或冷静期中）
	ApplyStatusSuccess        = 2 //注销成功
	ApplyStatusDeleted        = 3 //删除成功
	ApplyStatusRecover        = 4 //申请恢复
	ApplyStatusRecoverSuccess = 5 //恢复成功

	//白名单状态
	NotInWhiteList = 0 //不在白名单内
	InWhiteList    = 1 //在白名单内

	AppleValidateCodeUrl = "https://appleid.apple.com/auth/token"

	//密码盐字符串
	PasswordSaltChar = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_{}:<>?"

	// 第三方账号定义，防止调用方不统一问题，如微信传入weixin,后来改成wx, 造成账号查不到问题
	ThirdFacebook  = 1001
	ThirdGoogle    = 1002
	ThirdTwitter   = 1003
	ThirdYoutube   = 1004
	ThirdInstagram = 1005
	ThirdApple     = 1006
	ThirdWhatsapp  = 1007
	ThirdSkype     = 1008
	ThirdLinkedln  = 1009
	ThirdLine      = 1010
	ThirdVK        = 1011
	ThirdWeixin    = 1012
	ThirdReddit    = 1013
	ThirdWeibo     = 1014
	ThirdQQ        = 1015
	ThirdTiktok    = 1016

	//表名称 - base库基本信息表
	MailTplTable          = "mail_tpl"
	SmsTplTable           = "sms_tpl"
	UserDeleteConfigTable = "user_delete_config"
	WhiteUserListTable    = "white_user_list"
	HolidayTable          = "holiday"

	GameUserSlaveDb  = 1
	GameUserMasterDb = 2

	//用户uid自动生成Key
	RedisUidAutoIncrementKey = "_account_auto_increment_uid"

	//验证码类型  1: 注册, 2: 忘记密码, 3: 账号绑定, 4: 账号解绑, 5: 登录
	CodeTypeRegister       = 1
	CodeTypeForgetPassword = 2
	CodeTypeBindAccount    = 3
	CodeTypeUnBindAccount  = 4
	CodeTypeLogin          = 5
	//验证码格式，%d 为以上类型，%s 为账号
	CodeFormat = "_account_code_%d_%s"

	//header 内项目大区字符串
	HeaderGamePlatform = "game-platform"
	//日期时间格式
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"

	//图形验证码
	CaptchaExpire = 60 //图形验证码有效期
	CaptchaFormat = "_account_captcha_%s"

	//访问限制
	LimitIpLocKey      = "_account_limit_lock_ip_" //ip 锁key
	LimitIpKey         = "_account_limit_ip_"      //累加的ip key
	LimitRegisterIpKey = "_account_limit_ril_"     //注册ip锁key
	LimitLoginIpKey    = "_account_limit_li_"      //登录ip锁key
	LimitCodeIpKey     = "_account_limit_vcil_"    //验证码ip锁key
)

// 验证码类型
var CodeTypeMap = map[int]string{
	CodeTypeRegister:       "Register",
	CodeTypeForgetPassword: "ForgetPassword",
	CodeTypeBindAccount:    "BindAccount",
	CodeTypeUnBindAccount:  "UnBindAccount",
	CodeTypeLogin:          "Login",
}

// 除第三方的账号数据库字段
var AccountType = map[int]string{
	AccountEmail:    "email",
	AccountMobile:   "mobile",
	AccountUsername: "username",
	AccountGuest:    "guest",
	AccountThird:    "third",
}

// 支持的第三方账号id
var ThirdIds = map[int]string{
	ThirdFacebook:  "Facebook",
	ThirdGoogle:    "Google",
	ThirdTwitter:   "Twitter",
	ThirdYoutube:   "Youtube",
	ThirdInstagram: "Instagram",
	ThirdApple:     "Apple",
	ThirdWhatsapp:  "Whatsapp",
	ThirdSkype:     "Skype",
	ThirdLinkedln:  "Linkedln",
	ThirdLine:      "Line",
	ThirdVK:        "VK",
	ThirdWeixin:    "Weixin",
	ThirdReddit:    "Reddit",
	ThirdWeibo:     "Weibo",
	ThirdQQ:        "QQ",
	ThirdTiktok:    "Tiktok",
}

// 日志钩子结构，添加日志信息
type RequestHook struct {
	RequestBody        interface{}
	IP                 string
	GameId             int
	Uid                int64
	HeaderGamePlatform string
}

// 以下Etcd的配置内容为 json

// 配置信息
type AccountConf struct {
	Server                  ServerConf
	RedisConfig             RedisConf
	MysqlAccountBase        MysqlConfig
	MysqlAccountMasterList  map[string]MysqlConfig
	MysqlAccountSlaveList   map[string]MysqlConfig
	MysqlGameUserMasterList map[string]MysqlConfig
	MysqlGameUserSlaveList  map[string]MysqlConfig
	AliSmsConfig            AlibabaSms
	Base                    Base
	HttpTimeout             HttpTimeout
	RefreshTime             RefreshTime
	MysqlTimeout            MysqlTimeout
	MailConfig              MailConfig
	RequestLimitRule        ReqLimitRule
}

// 请求限制规则
type ReqLimitRule struct {
	Enabled               bool     `validate:"required"`
	WhiteList             []string `validate:"required"`
	WhiteListMap          map[string]int
	LoginAuthWhiteList    []string `validate:"required"`
	LoginAuthWhiteListMap map[string]int
	Ip                    []int `validate:"required"`
	Login                 []int `validate:"required"`
	VerifyCode            []int `validate:"required"`
	Register              []int `validate:"required"`
}

// RedisConf Redis配置结构体
type RedisConf struct {
	Addr        string `validate:"required"`
	Pass        string `validate:"required"`
	Db          int    `validate:"required"`
	PoolSize    int    `validate:"required"`
	MinIdle     int    `validate:"required"`
	MaxLifetime int    `validate:"required"`
}

// 服务基本信息
type ServerConf struct {
	LogRoot      string `validate:"required"`
	LogLevel     int8   `validate:"required"`
	DataLogsPath string `validate:"required"`
	Name         string `validate:"required"` //服务
	ID           int64  `validate:"required"` //服务id
	Host         string `validate:"required"` //监听host:port
}

// 基本配置
type Base struct {
	ServerKey           string `validate:"required"`
	LoginTokenExpires   int64  `validate:"required"`
	RefreshTokenExpires int64  `validate:"required"`
	CodeExpires         int64  `validate:"required"`
	CodeCheckMax        int    `validate:"required"`
	CallRestfulTimeout  int    `validate:"required"`
	Env                 string `validate:"required"`
	EnabledGameList     []int  `validate:"required"` //可以使用此系统的游戏id
	AppIdConfPath       string `validate:"required"`
	AutoIncrementUid    int64  `validate:"required"`
}

type HttpTimeout struct {
	ReadTimeout  int `validate:"required"`
	WriteTimeout int `validate:"required"`
	IdleTimeout  int `validate:"required"`
}

type RefreshTime struct {
	AppKeyRefreshTime     int `validate:"required"`
	GameConfigRefreshTime int `validate:"required"`
	UserDeleteRefreshTime int `validate:"required"`
	HolidayRefreshTime    int `validate:"required"`
}

type MysqlTimeout struct {
	MysqlTimeout      string `validate:"required"`
	MysqlReadTimeout  string `validate:"required"`
	MysqlWriteTimeout string `validate:"required"`
}

// 阿里短信配置信息
type AlibabaSms struct {
	RegionId  string `validate:"required"`
	AccessId  string `validate:"required"`
	SecretKey string `validate:"required"`
}

type MysqlConfig struct {
	Host        string `validate:"required"`
	User        string `validate:"required"`
	Port        int    `validate:"required"`
	Pass        string `validate:"required"`
	MaxConn     int    `validate:"required"`
	MaxIdle     int    `validate:"required"`
	MaxLifetime int    `validate:"required"`
	DbName      string `validate:"required"`
}

// 输出信息
type AccountResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 邮件配置信息
type MailConfig struct {
	Hostname string `validate:"required"`
	Port     int    `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
	Charset  string `validate:"required"`
}

// 邮件模板
type MailTpl struct {
	Type    string
	LangId  string
	Title   string
	Content string
}

// 短信模板模板
type SmsTpl struct {
	Type   string
	LangId string
	Title  string
	SmsId  string
}

// AppIdConfig SDK/服务器/客户端调用接口签名时所用的app_id, secret key等信息
type AppIdConfig struct {
	GameId     int    `json:"game_id" validate:"required"`
	AppId      int64  `json:"app_id" validate:"required"`
	SecretKey  string `json:"secret_key" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Type       int    `json:"type" validate:"required"` //1 为SDK, 2 为服务器， 3 为客户端
	Enabled    int    `json:"enabled" validate:"required"`
	PlatformId int64  `json:"platform_id" validate:"required"`
}

// 用户注册协议
type RegisterFields struct {
	Account    string `json:"account" validate:"required"` //第三方账号拼接：如facebook账号，拼接为 1000_abc123
	Code       string `json:"code" validate:"required"`    //手机号注册支持
	Password   string `json:"password" validate:"required"`
	DeviceType int    `json:"device_type" validate:"required"`
	Type       int    `json:"type" validate:"min=1,max=5"`
	Lang       string `json:"lang" validate:"required"`
	ChannelId  int    `json:"channel_id" validate:"required"`
	DataExt    string `json:"data_ext" validate:"required"`
	CommonFields
}

// 登录协议
type LoginFields struct {
	Account   string `json:"account" validate:"required"`
	Code      string `json:"code" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Type      int    `json:"type" validate:"min=1,max=5"`
	ChannelId int    `json:"channel_id" validate:"required"`
	DataExt   string `json:"data_ext" validate:"required"`
	CommonFields
}

// 第三方账号请求参数
type ThirdAccount struct {
	ThirdId           int    `json:"third_id"`  //如 1000, 代表 Facebook
	ThirdUid          string `json:"third_uid"` //第三方账号账号id
	ThirdUsername     string `json:"third_username"`
	ThirdEmail        string `json:"third_email"`
	AccessToken       string `json:"access_token"`
	AuthorizationCode string `json:"authorization_code"` //苹果AuthorizationCode
}

// 登录和注册成功返回字段
type LoginReturnFields struct {
	Uid       int64             `json:"uid"`
	Account   string            `json:"account"`
	LoginTime int64             `json:"time"`
	Tokens    LoginTokensFields `json:"tokens"`
	Binds     BindsFields       `json:"binds"`
	CardId    CardIdFields      `json:"card_id_info"`
}

// 身份信息,是否实名
type CardIdFields struct {
	IsRealName int   `json:"is_real_name"` //是否实名认证， 1 是
	Adult      int   `json:"adult"`        //是否成人， 1 是
	PlayTime   int64 `json:"play_time"`    //可以玩的时间，秒数， 如果是成人则为 -1
}

// 登录token及刷新token
type LoginTokensFields struct {
	LoginToken            string `json:"token"`
	ExpiresIn             int64  `json:"token_expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}

// 已绑定的信息
type BindsFields struct {
	Email  string   `json:"email"`
	Mobile string   `json:"mobile"`
	Thirds []string `json:"thirds"`
}

// 公用字段
type CommonFields struct {
	GameId     int    `json:"game_id" validate:"required"`
	PlatformId int    `json:"platform_id" validate:"required"`
	AppId      int64  `json:"app_id" validate:"required"`
	Sign       string `json:"sign" validate:"required"`
}

// 登录token数据信息结构
type CustomClaims struct {
	GameId     int
	PlatformId int
	ChannelId  int
	Uid        int64
	LoginTime  int64
	TokenType  int //1 为刷新token, 2 为访问token(创建订单)
	jwt.StandardClaims
}

// 登录校验返回信息
type LoginVerifyResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 主账号信息
type AccountQueryFields struct {
	Uid      int64
	Password string
	Salt     string
	Email    string
	Mobile   string
	Username string
	Guest    string
	Third    string
	CardId   string
	Name     string
	Type     int
}

// 项目用户信息
type GameUserQueryFields struct {
	Uid        int64
	MainUid    int64
	Account    string
	Status     int
	Ext        string
	Type       int
	CreateTime int
}

// 发送验证码协议
type VerifyCodeFields struct {
	Account  string `json:"account" validate:"required"`
	CodeType int    `json:"code_type" validate:"required"` //类型, 1: 注册, 2: 忘记密码, 3: 账号绑定, 4: 账号解绑, 5:登录
	LangId   string `json:"lang_id" validate:"required"`   //语言标识, 简体:zh-CN, 繁体: zh-TW, 英文: en-US
	Type     int    `json:"type" validate:"oneof=1 2"`     //1 email, 2 短信
	CommonFields
}

// 忘记密码-重置密码
type ForgetPasswordFields struct {
	Account  string `json:"account" validate:"required"`
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
	Type     int    `json:"type" validate:"oneof=1 2"`
	CommonFields
}

// 修改密码
type ChangePasswordFields struct {
	Uid         int64  `json:"uid" validate:"required"`
	Account     string `json:"account" validate:"required"`
	LoginToken  string `json:"token" validate:"required"` //登录token, 用于验证uid是否正确
	OldPassword string `json:"opassword" validate:"required"`
	NewPassword string `json:"npassword" validate:"required"`
	CommonFields
}

// 绑定账号
type BindAccountFields struct {
	Uid         int64  `json:"uid" validate:"required"`
	Code        string `json:"code" validate:"required"`         //验证码， 邮箱或手机号需要验证码，第三方 为-1
	Account     string `json:"account" validate:"required"`      //当前账号， 用于hash
	Password    string `json:"password" validate:"required"`     //邮箱或手机号绑定时，设置的密码，第三方可以为-1,不处理
	BindAccount string `json:"bind_account" validate:"required"` //要绑定的账号, 第三方的的信息也包含在内
	Type        int    `json:"type" validate:"oneof=1 2 5"`      //要绑定的账号类型，如游客绑定手机号，则传入手机号的类型 2
	LoginToken  string `json:"token" validate:"required"`        //登录token, 用于验证uid是否正确
	CommonFields
}

// 解绑账号
type UnBindAccountFields struct {
	Uid           int64  `json:"uid" validate:"required"`
	Account       string `json:"account" validate:"required"` //当前登录的账号
	Code          string `json:"code" validate:"required"`
	UnBindAccount string `json:"unbind_account" validate:"required"` //要解绑的账号
	Type          int    `json:"type" validate:"oneof=1 2 5"`        //要解绑账号的类型
	LoginToken    string `json:"token" validate:"required"`          //登录token, 用于验证uid是否正确
	CommonFields
}

// 服务器登录校验
type LoginAuthFields struct {
	Uid   int64  `json:"uid" validate:"required"`
	Token string `json:"token" validate:"required"`
	CommonFields
}

// 白名单参数
type WhiteListFields struct {
	Uid int64 `json:"uid" validate:"required"`
	CommonFields
}

// 白名单返回信息
type WhiteListResponse struct {
	Status int      `json:"status"`
	Env    []string `json:"env"`
}

// 账号注销参数
type LogoutAccountFields struct {
	Uid       int64  `json:"uid" validate:"required"`
	Account   string `json:"account" validate:"required"`
	Token     string `json:"token" validate:"required"`      //登录token
	ThirdInfo string `json:"third_info" validate:"required"` //如果需要撤销苹果授权，需要传此信息
	CommonFields
}

// 撤销账号注销参数
type UndoLogoutFields struct {
	Uid     int64  `json:"uid" validate:"required"`
	Account string `json:"account" validate:"required"`
	Token   string `json:"token" validate:"required"` //登录token
	CommonFields
}

// 账号注销Ext结构
type DeleteAccountExt struct {
	AppleRefreshToken string `json:"apple_refresh_token"`
}

// 游戏配置信息
type GameConfig struct {
	GameId                 int    //项目id
	PlatformId             int    //大区id
	UserDeleteWaitDuration int    //用户注销冷静期时长，单位天
	AppleClientId          string //苹果客户端id
	AppleClientSecret      string //苹果客户端秘钥
	Ext                    string
	LastAppleUpdateTime    int64 //最后一次更新苹果secret时间
}

// 游戏配置信息
type AppleConfig struct {
	TeamId         string `json:"team_id" validate:"required"`         //team id,来自SDK,用于组成jwt payload的iss(Issuer)
	ClientId       string `json:"client_id" validate:"required"`       //client id,来自SDK,用于组成jwt payload的sub(Subject)
	KeyId          string `json:"key_id" validate:"required"`          //p8文件名称,来自SDK,用于组成jwt header的kid
	PrivateKey     string `json:"private_key" validate:"required"`     //p8文件内容(不含头尾部分),来自SDK,用于jwt签名
	ValidityPeriod int64  `json:"validity_period" validate:"required"` //有效时长,用于组成jwt payload的exp(ExpiresAt)
	//Audience       string `json:"audience" validate:"required"`        //用户,用于组成jwt payload的aud(Audience)
}

// 用户信息查询结构
type UserInfoReqFields struct {
	Uid     int64  `json:"uid" validate:"required"`
	Account string `json:"account" validate:"required"`
	Token   string `json:"token" validate:"required"`
	CommonFields
}

// 用户查询返回结构
type UserInfoRespFields struct {
	Binds BindsFields  `json:"binds"`        //绑定的信息
	Cards CardIdFields `json:"card_id_info"` //身份证相关信息
}

// UserRealNameAuthReqFields 用户实名认证结构
type UserRealNameAuthReqFields struct {
	Uid     int64  `json:"uid" validate:"required"`     //项目uid
	Account string `json:"account" validate:"required"` //当前登录账号
	Token   string `json:"token" validate:"required"`   //登录token
	Name    string `json:"name" validate:"required"`    //姓名
	CardId  string `json:"card_id" validate:"required"` //身份证id
	CommonFields
}

// 图形验证码请求参数
type CaptchaImageFields struct {
	Id      string `json:"id"  validate:"required"`
	Width   int    `json:"width"  validate:"required"`  //不指定宽度传-1，使用默认宽度 240
	Height  int    `json:"height" validate:"required"`  //不指定高度传-1，使用默认高度 80
	Refresh int    `json:"refresh" validate:"required"` //是否刷新，1 是
	CommonFields
}

// 图形验证码验证参数
type CaptchaVerifyFields struct {
	CaptchaId   string `json:"captcha_id" validate:"required"`
	CaptchaType string `json:"type" validate:"required"`
	CaptchaCode string `json:"code" validate:"required"`
	CommonFields
}

// 触发限制访问时返回的数据格式
type LimitLockRetFields struct {
	CaptchaId   string `json:"captcha_id" validate:"required"`
	CaptchaType string `json:"type" validate:"required"`
}

// 根据账号 uid hash表、库结构
type DbTable struct {
	AccountMasterDb          *sql.DB
	AccountSlaveDb           *sql.DB
	AccountTable             string
	GameUserMasterDb         *sql.DB
	GameUserSlaveDb          *sql.DB
	GameUserTable            string
	GameUserDeleteApplyTable string
}

// 根据 account(email/mobile/guest/third/username) hash 其库、表
type HashDbTable struct {
	AccountMasterDb  *sql.DB
	AccountSlaveDb   *sql.DB
	AccountHashTable string
}
