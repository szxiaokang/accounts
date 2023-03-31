/**
 * @project Accounts
 * @filename constants.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) 2023/2/19
 * @datetime 2023/2/19 10:12
 * @version 1.0
 * @description
 * 错误码定义
 * 0 - 1212 通用
 * 其余开头
 * 23 注册
 * 33 登录
 * 43 忘记密码
 * 53 修改密码
 * 63 发送验证码
 * 73 绑定账号
 * 83 解绑
 * 93 服务器登录校验
 * 103 注销
 * 113 撤销注销
 * 123 实名认证
 */

package base

const (
	Success                              = 0     //成功
	Failure                              = 1     //失败
	NotPostRequest                       = 2     //不是post请求
	AppIdError                           = 101   //获取app_id错误,可能不存在
	AppIdQueryError                      = 102   //查询appid遇到错误, 可能不存在
	AppIdTypeError                       = 104   //app_id类型错误
	SignError                            = 105   //签名错误
	GameIdNotExists                      = 106   //不存在的game_id
	RequestDataParserError               = 107   //请求的数据解析错误
	RequestDataValidatorFail             = 108   //请求数据校验失败
	RequestDataIncorrect                 = 109   //请求的数据不正确
	PostFormRequestError                 = 110   //发送post请求失败
	PostFormReadError                    = 111   //post请求,读取数据失败
	LoginTokenParseError                 = 112   //登录token解析错误
	LoginTokenUidUnequal                 = 113   //登录token内的uid与传入的uid不相等
	RequestLimitRuleIpTrigger            = 114   //触发ip 访问限制规则
	RequestLimitIpConfigError            = 115   //ip访问限制规则配置错误
	RequestLimitLoginIpLock              = 116   //登录限制，ip已被锁定
	RequestLimitLoginConfigError         = 117   //登录限制配置错误
	RequestLimitVerifyCodeConfigError    = 118   //验证码发送频率限制配置错误
	RequestLimitVerifyCodeAccount        = 119   //验证码发送频率：60秒内仅允许发送一次
	RequestLimitVerifyCodeLockIp         = 120   //验证码发送频率：此ip超过了限定上线锁ip
	RequestLimitRegisterConfigError      = 121   //注册限制：配置错误
	RequestLimitRegisterLockIp           = 122   //注册限制：ip已被锁定
	RequestLimitCodeError                = 123   //图形验证码验证错误
	VerifyCodeNotExists                  = 1201  //验证码不存在
	VerifyCodeError                      = 1202  //验证码错误
	DeleteVerifyCodeError                = 1204  //删除验证码出错
	BuildTokenFailure                    = 1205  //生成token失败
	EmailFormatError                     = 1207  //邮箱格式错误
	PhoneNumFormatError                  = 1208  //手机号格式错误
	UsernameLengthError                  = 1209  //用户名长度错误
	UsernameFormatError                  = 1210  //用户名格式错误
	GuestOrThirdLengthError              = 1211  //游客或第三方ID长度错误
	GuestOrThirdFormatError              = 1212  //游客或第三方ID格式错误
	ThirdFormatError                     = 2308  //第三方账号格式错误
	ThirdIdParseFailure                  = 2309  //第三方id解析失败
	ThirdIdUnsupported                   = 2310  //不支持的第三方id
	ThirdUidEmpty                        = 2311  //第三方账号uid为空
	GetHashDbTxError                     = 2312  //获取hash master db 事务错误
	GetAccountDbTxError                  = 2313  //获取账号master db事务错误
	GetGameUserDbTxError                 = 2314  //获取项目用户master db事务错误
	TxExecInsertHashError                = 2315  //账号hash表插入执行错误
	TxExecInsertAccountError             = 2316  //账号表执行插入错误
	TxExecInsertGameUserError            = 2317  //项目用户执行插入错误
	TxHashDbCommitError                  = 2318  //hash库事务提交失败
	RegisterCodeAndPasswordEmpty         = 2319  //邮箱或手机号注册时，验证码和密码不允许都是空
	RegisterSuccess                      = 2320  //注册成功
	BuildAccountUidError                 = 2321  //生成账号uid错误
	LoginAccountDisabled                 = 3324  //账号被禁用
	LoginPasswordError                   = 3326  //登录密码错误
	LoginUserOrPasswordError             = 3327  //账号或密码错误
	LoginCodeAndPasswordEmpty            = 3328  //验证码和密码不能都是空值
	LoginGameUserQueryError              = 3332  //查询项目用户表出错
	AccountIsBeingDeleted                = 3333  //账号注销中
	AccountLoginException                = 3337  //主账号登录异常
	LoginInsertGameUserError             = 3338  //插入项目用户表时失败
	LoginSuccess                         = 3339  //登录成功
	QueryGameUserScanError               = 3340  //查询项目用户表Scan时出错
	ForgetPasswordUpdateFailure          = 4340  //忘记密码更新数据失败
	ForgetPasswordAccountNotExists       = 4341  //账号不存在
	OldPasswordError                     = 5360  //原密码错误
	ChangePasswordAccountNotExists       = 5361  //修改密码账号不存在
	ChangePasswordUpdateFailure          = 5362  //修改密码更新数据失败
	ChangePasswordGameUserQueryNotExists = 5363  //修改密码查询项目用户表不存在
	ChangePasswordGameUserQueryError     = 5364  //修改密码查询项目用户表错误
	ChangePasswordUidNotMatch            = 5365  //账号与uid不匹配
	MailTplConfigQueryFailure            = 6380  //邮件模板配置查询失败
	SmsTplConfigQueryFailure             = 6382  //短信配置查询失败
	SmsTplConfigEmpty                    = 6383  //短信配置为空
	SendCodeAccountNotExists             = 6384  //发送验证码账号不存在
	SendCodeAccountAlreadyExists         = 6385  //发送注册验证码账号已存在
	VerifyCodeInsertError                = 6387  //验证码写入失败
	VerifyCodeUpdateFailure              = 6388  //验证码更新失败
	VerifyCodeTypeUnknown                = 6389  //未知的短信类型
	BindAccountNotExists                 = 7301  //绑定账号，账号不存在
	GetAlreadyBindInfoNotFound           = 7335  //查询已绑定账号失败，找不到
	GetAlreadyBindThirdError             = 7336  //查询第三方绑定信息错误
	BindAccountAlreadyExists             = 7338  //要绑定的账号已经存在
	BindGetHashTxError                   = 7339  //获取hash库事务错误
	BindGetAccountDbTxError              = 7340  //获取账号、项目用户事务错误
	BindHashTxExecInsertError            = 7341  //hash 账号插入事务执行错误
	BindGameTxExecInsertError            = 7342  //项目表执行插入事务错误
	BindAccountTxExecUpdateError         = 7343  //账号表执行更新事务错误
	GetAlreadyBindInfoNotFoundByUid      = 7344  //查询uid查询错误，找不到
	QueryAlreadyBindThirdScanError       = 7345  //查询第三方已绑定信息Scan时错误
	BindGetGameUserDbTxError             = 7346  //获取项目事务失败
	BindQueryAccountInfoError            = 7347  //查询账号信息错误
	BindEmailAlreadyExists               = 7348  //Email已绑定
	BindMobileAlreadyExists              = 7349  //手机号不为空，已绑定
	UnbindAccountNotExists               = 8320  //解绑，账号不存在
	EmailAlreadyUnbind                   = 8323  //邮箱已经解绑
	MobileAlreadyUnbind                  = 8324  //手机号已经解绑
	BeUnBindAccountNotExists             = 8332  //被解绑账号不存在
	BindAccountAndBeUnBindNotMatch       = 8333  //被解绑的账号与当前账号不属于同一个账号下
	BeUnBindAccountNotExists2            = 8334  //被解绑账号不存在，查询主账号
	UnBindQueryThirdError                = 8335  //查询第三方账号时错误
	UnBindQueryThirdScanError            = 8336  //查询第三方账号Scan时错误
	UnBindUnSupportRegisterType          = 8337  //不支持解绑注册时的类型
	UnBindGetHashTxError                 = 8338  //获取hash库事务错误
	UnBindGetAccountDbTxError            = 8339  //获取账号用户事务错误
	UnBindGetGameUserDbTxError           = 8340  //获取项目用户事务错误
	UnBindGameUserTxExecDeleteError      = 8341  //项目用户删除事务错误
	UnBindHashTxExecUpdateError          = 8342  //hash账号更新事务错误
	UnBindAccountTxExecUpdateError       = 8343  //账号表执行更新事务错误
	UnBindUnSupportCurrentAccount        = 8344  //不支持解绑当前账号
	LoginAuthError                       = 9340  //登录校验失败
	LoginAuthTokenParseError             = 9341  //登录Token解析失败
	LoginTokenUidNotMatch                = 9342  //登录Token与uid不匹配
	LoginAuthDomainNotWhiteList          = 9343  //请求域名不是白名单内
	AddDeleteApplyError                  = 10365 //添加注销申请失败
	DeleteApplyUpdateUserStatusError     = 10366 //账号注销修改用户状态失败
	DeleteApplyAlreadyExists             = 10367 //添加注销申请已存在
	DeleteAccountNotExists               = 10368 //申请注销的账号不存在
	DeleteAccountQueryNoRows             = 10369 //账号不存在
	DeleteAccountQueryError              = 10370 //账号查询错误
	DeleteAccountAndUidNotMatch          = 10371 //账号uid不匹配
	UndoDeleteUpdateUserStatusError      = 11384 //撤销账号注销修改用户状态失败
	UndoDeleteApplyNotExists             = 11385 //撤销账号注销记录不存在
	UndoDeleteAccountNotExists           = 11386 //撤销账号不存在
	UndoDeleteAccountAndRecordNotMatch   = 11387 //撤销账号与记录中的不匹配
	RealNameNameError                    = 12300 //姓名错误
	RealNameCardIdError                  = 12301 //身份证号错误
	RealNameApiReturnError               = 12302 //调用身份证实名接口错误
	RealNameUpdateError                  = 12303 //更新错误
	RealNameGetAccountUidNotExists       = 12304 //根据account获取账号id不存在
	RealNameUpdateAccountError           = 12305 //更新账号信息错误
)

var ErrorMsg = map[int]string{
	Success:                              "success",
	Failure:                              "failure",
	NotPostRequest:                       "not a post request",
	AppIdError:                           "get app id error, may not exist",
	AppIdQueryError:                      "query appid encountered an error and may not exist",
	AppIdTypeError:                       "app id type error",
	SignError:                            "signature error",
	GameIdNotExists:                      "nonexistent game id",
	RequestDataParserError:               "data parsing error requested",
	RequestDataValidatorFail:             "request data check failed",
	RequestDataIncorrect:                 "the requested data is incorrect",
	PostFormRequestError:                 "send post request failed",
	PostFormReadError:                    "post request, read data failed",
	LoginTokenParseError:                 "logon token parsing error",
	LoginTokenUidUnequal:                 "login token uid is not equal to incoming uid",
	RequestLimitRuleIpTrigger:            "trigger ip access restriction rule",
	RequestLimitIpConfigError:            "ip access restriction rule configuration error",
	RequestLimitLoginIpLock:              "logon restriction, ip locked",
	RequestLimitLoginConfigError:         "logon restriction configuration error",
	RequestLimitVerifyCodeConfigError:    "authentication code send frequency limit configuration error",
	RequestLimitVerifyCodeAccount:        "authentication code send frequency: only one send is allowed in 60 seconds",
	RequestLimitVerifyCodeLockIp:         "authentication code send frequency: this ip exceeds the limit on online lock ip",
	RequestLimitRegisterConfigError:      "registration limit: configuration error",
	RequestLimitRegisterLockIp:           "registration limit: ip is locked",
	RequestLimitCodeError:                "graphic authentication error",
	VerifyCodeNotExists:                  "authentication code does not exist",
	VerifyCodeError:                      "authentication code error",
	DeleteVerifyCodeError:                "error deleting verification code",
	BuildTokenFailure:                    "generate token failed",
	EmailFormatError:                     "email format error",
	PhoneNumFormatError:                  "incorrectly formatted mobile number",
	UsernameLengthError:                  "user name length error",
	UsernameFormatError:                  "username format error",
	ThirdFormatError:                     "third-party format error",
	GuestOrThirdLengthError:              "guest or third-party id length error",
	GuestOrThirdFormatError:              "guest or third-party id format error",
	ThirdIdParseFailure:                  "registration - third party id resolution failed",
	ThirdIdUnsupported:                   "unsupported third party id",
	ThirdUidEmpty:                        "third-party account uid is empty",
	GetHashDbTxError:                     "get hash master db transaction error",
	GetAccountDbTxError:                  "get account master db transaction error",
	GetGameUserDbTxError:                 "get project user master db transaction error",
	TxExecInsertHashError:                "account hash table insert execution error",
	TxExecInsertAccountError:             "account table execution insertion error",
	TxExecInsertGameUserError:            "project user performs insert error",
	TxHashDbCommitError:                  "hash library transaction commit failed",
	RegisterCodeAndPasswordEmpty:         "the verification code and password cannot both be empty",
	RegisterSuccess:                      "successful registration",
	BuildAccountUidError:                 "generate account uid error",
	LoginAccountDisabled:                 "account is disabled",
	LoginPasswordError:                   "logon password error",
	LoginUserOrPasswordError:             "account or password error",
	LoginCodeAndPasswordEmpty:            "authentication code and password cannot both be null values",
	LoginGameUserQueryError:              "error querying project user table",
	AccountIsBeingDeleted:                "account logoff status",
	AccountLoginException:                "main account logon exception",
	LoginInsertGameUserError:             "insert project user table failed",
	LoginSuccess:                         "successful login",
	QueryGameUserScanError:               "error querying project user table scan",
	ForgetPasswordUpdateFailure:          "forget password update data failure",
	ForgetPasswordAccountNotExists:       "account does not exist",
	OldPasswordError:                     "original password error",
	ChangePasswordAccountNotExists:       "modified password, account does not exist",
	ChangePasswordUpdateFailure:          "failure to modify password update data",
	ChangePasswordGameUserQueryNotExists: "modify password, query item user table does not exist",
	ChangePasswordGameUserQueryError:     "modify password, query item user table error",
	ChangePasswordUidNotMatch:            "account does not match uid",
	MailTplConfigQueryFailure:            "mail template configuration query failed",
	SmsTplConfigQueryFailure:             "sms configuration query failed",
	SmsTplConfigEmpty:                    "sms configuration is empty",
	SendCodeAccountNotExists:             "send authentication number account does not exist",
	SendCodeAccountAlreadyExists:         "send registration verification number account already exists",
	VerifyCodeInsertError:                "authenticode write failure",
	VerifyCodeUpdateFailure:              "verify code update failed",
	VerifyCodeTypeUnknown:                "unknown sms type",
	BindAccountNotExists:                 "bind account, account does not exist",
	GetAlreadyBindInfoNotFound:           "query failed for bound account, could not be found",
	GetAlreadyBindThirdError:             "error in querying third-party binding information",
	BindAccountAlreadyExists:             "account to bind already exists",
	BindGetHashTxError:                   "get hash library transaction error",
	BindGetAccountDbTxError:              "get account, project user transaction error",
	BindHashTxExecInsertError:            "hash account insertion transaction execution error",
	BindGameTxExecInsertError:            "project table insert transaction error",
	BindAccountTxExecUpdateError:         "account table execution update transaction error",
	GetAlreadyBindInfoNotFoundByUid:      "query uid query error, could not be found",
	QueryAlreadyBindThirdScanError:       "error querying third party bound information scan",
	BindGetGameUserDbTxError:             "get game user transaction error",
	UnbindAccountNotExists:               "unbind, account does not exist",
	EmailAlreadyUnbind:                   "mailbox unbound",
	MobileAlreadyUnbind:                  "mobile number unbound",
	BeUnBindAccountNotExists:             "unbound account does not exist",
	BeUnBindAccountNotExists2:            "unbound account does not exist 2",
	UnBindQueryThirdError:                "error querying third-party accounts",
	UnBindQueryThirdScanError:            "error querying third-party account scan",
	UnBindUnSupportRegisterType:          "type when unbound registration is not supported",
	UnBindGetHashTxError:                 "get hash library transaction error",
	UnBindGetAccountDbTxError:            "get account user transaction error",
	UnBindGetGameUserDbTxError:           "get project user transaction error",
	UnBindGameUserTxExecDeleteError:      "project user delete transaction error",
	UnBindHashTxExecUpdateError:          "hash account update transaction error",
	UnBindAccountTxExecUpdateError:       "account table execution update transaction error",
	UnBindUnSupportCurrentAccount:        "unbinding the current account is not supported",
	LoginAuthError:                       "logon verification failed",
	LoginAuthTokenParseError:             "login token resolution failure",
	LoginTokenUidNotMatch:                "login token does not match uid",
	LoginAuthDomainNotWhiteList:          "request domain name is not on the whitelist",
	AddDeleteApplyError:                  "add logoff request failed",
	DeleteApplyUpdateUserStatusError:     "account logoff failed to modify user status",
	DeleteApplyAlreadyExists:             "add logoff request already exists",
	DeleteAccountNotExists:               "the account not exists",
	DeleteAccountQueryNoRows:             "account query blank line, account table",
	DeleteAccountQueryError:              "account query error",
	DeleteAccountAndUidNotMatch:          "account uid mismatch",
	UndoDeleteUpdateUserStatusError:      "undo account logoff modify user status failed",
	UndoDeleteApplyNotExists:             "revoke account logoff record does not exist",
	UndoDeleteAccountNotExists:           "cancel account does not exist",
	UndoDeleteAccountAndRecordNotMatch:   "undo account does not match the record",
	RealNameNameError:                    "name error",
	RealNameCardIdError:                  "id number error",
	RealNameApiReturnError:               "call authentication name interface error",
	RealNameUpdateError:                  "update error",
	RealNameGetAccountUidNotExists:       "getting account id from account does not exist",
	RealNameUpdateAccountError:           "update account information error",
}
