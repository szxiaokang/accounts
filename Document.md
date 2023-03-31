
## 接口文档
<hr>

### 1 注册、登录

##### 简要描述

- 用户注册接口
- 支持 email、手机号+验证码或密码方式
- 当账号存在时则登录，否则注册

##### 请求URL
- ` /user/register `

##### 请求方式
- POST application/json

##### 参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|account |是  |string |用户账号，Email/手机号/用户名/游客id/第三方id(类型编号_第三方账号id, 如微信id是 1234567， 则传入：1012_1234567， 类型编号见下面列表)   |
|code |是  |string |短信或邮件验证码，非验证码方式注册、登录传入-1     |
|password |是  |string |密码，md5后，没有传-1     |
|device_type |是  |int |设备类型，1:ios，2:android，3:pc/h5，4:其他     |
|type |是  |int |账号类型，1:email，2:手机号，3:用户名，4:游客，5:第三方     |
|lang |是  |string |语言简拼，符合i18n规范，如简体中文:zh_CN     |
|channel_id |是  |int |用户登录包的渠道id     |
|data_ext |是  |string json |数据埋点json字符串，可以传任意值，没有传空 "{}"     |
|game_id     |是  |int | 游戏ID    |
|platform_id     |是  |int | 大区ID    |
|app_id     |是  |int | 分配的APPID    |
|sign     |是  |string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)    |

##### 返回示例

``` 
  {
	"code": 0,
	"msg": "success",
	"data": {
		"uid": 1610000100014,
		"account": "kangyun@outlook.com",
		"time": 1680233244,
		"tokens": {
			"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTQsIkxvZ2luVGltZSI6MTY4MDIzMzI0NCwiVG9rZW5UeXBlIjoxLCJleHAiOjE2ODAyNDA0NDQsImlzcyI6ImFjY291bnRfc2VydmVyIn0.-TZN_1lE6rk2dW_AsuW1dHbmX0x7Y8gv3qJDeHBb_Z0",
			"token_expires_in": 7140,
			"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTQsIkxvZ2luVGltZSI6MTY4MDIzMzI0NCwiVG9rZW5UeXBlIjoyLCJleHAiOjE3MTE3NjkyNDQsImlzcyI6ImFjY291bnRfc2VydmVyIn0.IIeUoIgO4H-yt4hyFUHnK34kfYoTimxir17v-fDUTVY",
			"refresh_token_expires_in": 31536000
		},
		"binds": {
			"email": "kangyun@outlook.com",
			"mobile": "",
			"thirds": []
		},
		"card_id_info": {
			"is_real_name": 0,
			"adult": 0,
			"play_time": 0
		}
	}
}
```

##### 返回参数说明

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|uid |int   |用户UID  |
|account |string   |用户账号  |
|time |int   |时间戳  |
|tokens |json   |token信息  |
|tokens.token |string   |登录token  |
|tokens.token_expires_in |int   |登录token有效期，单位：秒  |
|tokens.refresh_token |string   |刷新token  |
|tokens.rerefsh_token_expires_in |int   |刷新token有效期，单位：秒  |
|binds |json   |绑定信息  |
|binds.email |string   |绑定的email  |
|binds.mobile |string   |绑定的手机号  |
|binds.thirds |array   |绑定的第三方  |
|card_id_info |json   |身份信息  |
|card_id_info.is_real_name |int   |是否实名，1 是  |
|card_id_info.adult |int   |是否成人，1 是  |
|card_id_info.play_time |int   |反沉迷，可玩时长，秒数  |

##### 错误码
见 错误码及常量
<hr>


### 2 用户登录

##### 简要描述

- 用户登录接口
- 支持email、手机号+验证码或密码

##### 请求URL
- ` /user/login `

##### 请求方式
- POST application/json

##### 参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|account |是  |string |用户账号，Email/手机号/用户名/游客id/第三方id   |
|code |是  |string |短信验证码，非手机号传-1     |
|password |是  |string |密码，md5后，没有传-1     |
|type |是  |int |账号类型，1:email，2:手机号，3:用户名，4:游客，5:第三方     |
|channel_id |是  |int |用户登录包的渠道id     |
|data_ext |是  |string |数据埋点json字符串，可以传任意值，没有传空 "{}"     |
|game_id     |是  |int | 游戏ID    |
|platform_id     |是  |int | 大区ID    |
|app_id     |是  |int | 分配的APPID    |
|sign     |是  |string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)    |

##### 返回示例

``` 
  {
    "code": 0,
    "msg": "success",
    "data": {
        "uid": 1610000100014,
        "account": "kangyun@outlook.com",
        "time": 1680233244,
        "tokens": {
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTQsIkxvZ2luVGltZSI6MTY4MDIzMzI0NCwiVG9rZW5UeXBlIjoxLCJleHAiOjE2ODAyNDA0NDQsImlzcyI6ImFjY291bnRfc2VydmVyIn0.-TZN_1lE6rk2dW_AsuW1dHbmX0x7Y8gv3qJDeHBb_Z0",
            "token_expires_in": 7140,
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJHYW1lSWQiOjE2LCJQbGF0Zm9ybUlkIjoxLCJDaGFubmVsSWQiOi0xLCJVaWQiOjE2MTAwMDAxMDAwMTQsIkxvZ2luVGltZSI6MTY4MDIzMzI0NCwiVG9rZW5UeXBlIjoyLCJleHAiOjE3MTE3NjkyNDQsImlzcyI6ImFjY291bnRfc2VydmVyIn0.IIeUoIgO4H-yt4hyFUHnK34kfYoTimxir17v-fDUTVY",
            "refresh_token_expires_in": 31536000
        },
        "binds": {
            "email": "kangyun@outlook.com",
            "mobile": "",
            "thirds": []
        },
        "card_id_info": {
            "is_real_name": 0,
            "adult": 0,
            "play_time": 0
        }
    }
}

```

##### 返回参数说明

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|uid |int   |用户UID  |
|account |string   |用户账号  |
|time |int   |时间戳  |
|tokens |json   |token信息  |
|tokens.token |string   |登录token  |
|tokens.token_expires_in |int   |登录token有效期，单位：秒  |
|tokens.refresh_token |string   |刷新token  |
|tokens.rerefsh_token_expires_in |int   |刷新token有效期，单位：秒  |
|binds |json   |绑定信息  |
|binds.email |string   |绑定的email  |
|binds.mobile |string   |绑定的手机号  |
|binds.thirds |array   |绑定的第三方  |
|card_id_info |json   |身份信息  |
|card_id_info.is_real_name |int   |是否实名，1 是  |
|card_id_info.adult |int   |是否成人，1 是  |
|card_id_info.play_time |int   |反沉迷，可玩时长，秒数  |

##### 错误码
见 错误码及常量
<hr>


### 3 绑定账号
##### 简要描述

- 绑定账号，可以绑定邮箱、手机号、第三方（可以多个）

##### 请求URL
- ` /user/bindAccount `

##### 请求方式
- POST application/json

##### 参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|uid |是  |int |用户UID     |
|code |是  |string |验证码，第三方传入-1     |
|account |是  |string |当前账号     |
|password |是  |string |密码，md5后，没有传-1     |
|bind_account |是  |string |要绑定的账号，Email/手机号/第三方id   |
|type |是  |int |账号类型，1:email，2:手机号，5:第三方     |
|token |是  |string |登录token     |
|game_id     |是  |int | 游戏ID    |
|platform_id     |是  |int | 大区ID    |
|app_id     |是  |int | 分配的APPID    |
|sign     |是  |string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)    |

##### 返回示例
binds, card_id_info 意思同登录、注册
``` 
  {
    "code": 0,
    "msg": "OK",
    "data": {
		"binds": {
			"email": "kangyun@outlook.com",
			"mobile": "",
			"thirds": []
		},
		"card_id_info": {
			"is_real_name": 0,
			"adult": 0,
			"play_time": 0
		}
    }
  }
```

##### 错误码
见 错误码及常量
<hr>


### 4 服务器登录校验
##### 简要描述

- 服务器登录校验接口

##### 请求URL
- ` /user/loginAuth `

##### 请求方式
- POST application/json

##### 参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|uid |是  |int |用户UID     |
|token |是  |string |登录token     |
|game_id     |是  |int | 游戏ID    |
|platform_id     |是  |int | 大区ID    |
|app_id     |是  |int | 分配的APPID    |
|sign     |是  |string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)    |

##### 返回示例

``` 
  {
    "code": 0,
    "msg": "OK",
    "data": {}
  }
```

##### 错误码
见 错误码及常量
<hr>





### 5 修改密码
##### 简要描述

- 修改密码接口

##### 请求URL
- ` /user/changePassword `

##### 请求方式
- POST application/json

##### 参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|uid |是  |int |用户UID     |
|account |是  |string |当前登录的账号     |
|token |是  |string |登录token     |
|opassword |是  |string |旧密码,md5后     |
|npassword |是  |string |新密码,md5后     |
|game_id     |是  |int | 游戏ID    |
|platform_id     |是  |int | 大区ID    |
|app_id     |是  |int | 分配的APPID    |
|sign     |是  |string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)    |

##### 返回示例

``` 
  {
    "code": 0,
    "msg": "OK",
    "data": {}
  }
```

##### 错误码
见 错误码及常量
<hr>




### 6 忘记密码-重置密码
##### 简要描述

- 忘记密码-重置密码接口
- 支持注册或是被绑定的email或手机号

##### 请求URL
- ` /user/forgetPassword `

##### 请求方式
- POST application/json

##### 参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|account |是  |string |用户账号，Email/手机号   |
|code |是  |string |短信或邮箱验证码     |
|password |是  |string |密码，md5后，没有传-1     |
|type |是  |int |账号类型，1:email，2:手机号     |
|game_id     |是  |int | 游戏ID    |
|platform_id     |是  |int | 大区ID    |
|app_id     |是  |int | 分配的APPID    |
|sign     |是  |string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)    |

##### 返回示例

``` 
  {
    "code": 0,
    "msg": "OK",
    "data": {}
  }
```
##### 错误码
见 错误码及常量
<hr>






### 7 发送验证码
##### 简要描述

- 发送验证码接口
- 一个账号60秒内仅允许发送一次

##### 请求URL
- ` /user/sendSmsCode `

##### 请求方式
- POST application/json

##### 参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|account |是  |string |用户账号，Email/手机号   |
|code_type |是  |string |类型, register: 注册, forgetPassword: 忘记密码, bindAccount: 账号绑定, unBindAccount: 账号解绑     |
|lang_id |是  |string |语言标识, 简体:zh-CN, 繁体: zh-TW, 英文: en     |
|type |是  |int |账号类型，1:email，2:手机号     |
|game_id     |是  |int | 游戏ID    |
|platform_id     |是  |int | 大区ID    |
|app_id     |是  |int | 分配的APPID    |
|sign     |是  |string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)    |

##### 返回示例

``` 
  {
    "code": 0,
    "msg": "OK",
    "data": {}
  }
```

##### 错误码
见 错误码及常量
<hr>





### 8 白名单校验
##### 简要描述
- 白名单校验接口，客户端直接使用

##### 请求URL
- ` /user/whiteList `

##### 请求方式
- POST application/json

##### 参数
|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|uid |是  |int |用户UID，与客户端ip为或关系，两者有一个满足即在白名单内     |
|game_id     |是  |int | 游戏ID    |
|platform_id     |是  |int | 大区ID    |
|app_id     |是  |int | 分配的APPID    |
|sign     |是  |string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)    |

**返回示例**

```
{
    "code": 0,
    "msg": "OK",
    "data": {
    	"status": 1
    	"env": ["dev","review","test","prerelease","produce"]
    }
}
```

**返回参数说明**

| **参数名** | **类型** | **说明**                                                     |
| ---------- | -------- | ------------------------------------------------------------ |
| status     | int      | 白名单状态，0：不在白名单内，1：在白名单内                   |
| env        | array      | 可进入的环境，为"dev","review","test","prerelease","produce"其中多个或一个 |

##### 错误码
见 错误码及常量
<hr>

### 9 解绑账号
##### 简要描述

- 解绑账号接口
- 仅支持绑定了 邮箱、手机号或是第三方的账号

##### 请求URL
- ` /user/unBindAccount `

##### 请求方式
- POST application/json

##### 参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|uid |是  |int |用户UID     |
|account |是  |string |用户账号   |
|code |是  |string |验证码，第三方传-1     |
|unbind_account |是  |string |用户账号，Email/手机号/第三方账号信息   |
|type |是  |int |解绑的账号类型，1:email，2:手机号，5：第三方     |
|token |是  |string |登录token     |
|game_id     |是  |int | 游戏ID    |
|platform_id     |是  |int | 大区ID    |
|app_id     |是  |int | 分配的APPID    |
|sign     |是  |string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)    |

##### 返回示例

``` 
  {
    "code": 0,
    "msg": "OK",
    "data": {}
  }
```

### 10 账号注销
##### 简要描述

- 账号注销接口

##### 请求URL
- ` /user/applyLogout `

##### 请求方式
- POST application/json

##### 参数

| 参数名         | 必选                 | 类型     | 说明                                                      |
|:------------|:-------------------|:-------|:--------------------------------------------------------|
| uid         | 是                  | int    | 用户UID                                                   |
| account     | 是                  | string | 当前账号                                                    |
| token       | 是                  | string | 登录token，需要在有效期内                                         |
| third_info  | 是                  | string | 不需要使用可以传"{}"，需要解绑第三方需要json以下字段:                         
| ------      | third_name         | string | 第三方名称，Google、Facebook等，多项目要保持名称一致，如餐厅传入Facebook， ARK要一样 |
| ------      | third_username     | string | 第三方的用户名                                                 |
| ------      | third_email        | string | 第三方的email地址，无值传入空字符串                                    |
| ------      | access_token       | string | 验证第三方是否有效的token字段                                       |
| ------      | authorization_code | string | 苹果AuthorizationCode                                     |
| game_id     | 是                  | int    | 游戏ID                                                    |
| platform_id | 是                  | int    | 大区ID                                                    |
| app_id      | 是                  | int    | 分配的APPID                                                |
| sign        | 是                  | string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)                  |

##### 返回示例

``` 
  {
    "code": 0,
    "msg": "OK",
    "data": {}
  }
```

##### 错误码
见 错误码及常量
<hr>


### 11 撤销账号注销
##### 简要描述

- 撤销账号注销接口

##### 请求URL
- ` /user/undoLogout `

##### 请求方式
- POST application/json

##### 参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|uid |是  |int |用户UID     |
|account |是  |string |当前账号     |
|token |是  |string |登录token，需要在有效期内     |
|game_id     |是  |int | 游戏ID    |
|platform_id     |是  |int | 大区ID    |
|app_id     |是  |int | 分配的APPID    |
|sign     |是  |string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)    |

##### 返回示例

``` 
  {
    "code": 0,
    "msg": "OK",
    "data": {}
  }
```

##### 错误码
见 错误码及常量
<hr>




### 12 用户信息
##### 接口名
- 用户信息查询

##### 请求URL

- ` /user/getUserInfo`

##### 请求方式

- POST application/json

##### 参数

| 参数名         | 必选 | 类型   | 说明                                                         |
| :------------- | :--- | :----- | ------------------------------------------------------------ |
| uid            | 是   | int    | 用户uid                                                         |
| account            | 是   | string    | 账号                                                         |
| token    | 是   | string | 登录token                                               |
| game_id        | 是   | int    | 游戏id                                                       |
| platform_id    | 是   | int    | 大区id                                                       |
| app_id         | 是   | int    | 分配的AppID                                               |
| sign           | 是   | string | 签名，md5(用&符号按顺序拼接以上所有字段，最后拼接&SecretKey)     |

##### 返回示例
```
  {
    "code": 0,
    "msg": "OK",
    "data": {
		"binds": {
			"email": "kangyun@outlook.com",
			"mobile": "",
			"thirds": []
		},
		"card_id_info": {
			"is_real_name": 0,
			"adult": 0,
			"play_time": 0
		}
    }
  }
```
##### 返回参数说明

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|binds |json   |绑定信息  |
|binds.email |string   |绑定的email  |
|binds.mobile |string   |绑定的手机号  |
|binds.thirds |array   |绑定的第三方  |
|card_id_info |json   |身份信息  |
|card_id_info.is_real_name |int   |是否实名，1 是  |
|card_id_info.adult |int   |是否成人，1 是  |
|card_id_info.play_time |int   |反沉迷，可玩时长，秒数  |


##### 错误码
见 错误码及常量
<hr>

### 错误码及常量	

|错误码| 说明                      |
|:----    |:------------------------|
|0     | 成功                      |
|1     | 失败                      |
|2     | 不是post请求                |
|101   | 获取app_id错误,可能不存在        |
|102   | 查询appid遇到错误, 可能不存在      |
|104   | app_id类型错误              |
|105   | 签名错误                    |
|106   | 不存在的game_id             |
|107   | 请求的数据解析错误               |
|108   | 请求数据校验失败                |
|109   | 请求的数据不正确                |
|110   | 发送post请求失败              |
|111   | post请求,读取数据失败           |
|112   | 登录token解析错误             |
|113   | 登录token内的uid与传入的uid不相等  |
|114   | 触发ip 访问限制规则             |
|115   | ip访问限制规则配置错误            |
|116   | 登录限制，ip已被锁定             |
|117   | 登录限制配置错误                |
|118   | 验证码发送频率限制配置错误           |
|119   | 验证码发送频率：60秒内仅允许发送一次     |
|120   | 验证码发送频率：此ip超过了限定上线锁ip   |
|121   | 注册限制：配置错误               |
|122   | 注册限制：ip已被锁定             |
|123   | 图形验证码验证错误               |
|1201  | 验证码不存在                  |
|1202  | 验证码错误                   |
|1204  | 删除验证码出错                 |
|1205  | 生成token失败               |
|1207  | 邮箱格式错误                  |
|1208  | 手机号格式错误                 |
|1209  | 用户名长度错误                 |
|1210  | 用户名格式错误                 |
|1211  | 游客或第三方ID长度错误            |
|1212  | 游客或第三方ID格式错误            |
|2308  | 第三方账号格式错误               |
|2309  | 第三方id解析失败               |
|2310  | 不支持的第三方id               |
|2311  | 第三方账号uid为空              |
|2312  | 获取hash master db 事务错误   |
|2313  | 获取账号master db事务错误       |
|2314  | 获取项目用户master db事务错误     |
|2315  | 账号hash表插入执行错误           |
|2316  | 账号表执行插入错误               |
|2317  | 项目用户执行插入错误              |
|2318  | hash库事务提交失败             |
|2319  | 邮箱或手机号注册时，验证码和密码不允许都是空  |
|2320  | 注册成功                    |
|2321  | 生成账号uid错误               |
|3324  | 账号被禁用                   |
|3326  | 登录密码错误                  |
|3327  | 账号或密码错误                 |
|3328  | 验证码和密码不能都是空值            |
|3332  | 查询项目用户表出错               |
|3333  | 账号注销中                   |
|3337  | 主账号登录异常                 |
|3338  | 插入项目用户表时失败              |
|3339  | 登录成功                    |
|3340  | 查询项目用户表Scan时出错          |
|4340  | 忘记密码更新数据失败              |
|4341  | 账号不存在                   |
|5360  | 原密码错误                   |
|5361  | 修改密码账号不存在               |
|5362  | 修改密码更新数据失败              |
|5363  | 修改密码查询项目用户表不存在          |
|5364  | 修改密码查询项目用户表错误           |
|5365  | 账号与uid不匹配               |
|6380  | 邮件模板配置查询失败              |
|6382  | 短信配置查询失败                |
|6383  | 短信配置为空                  |
|6384  | 发送验证码账号不存在              |
|6385  | 发送注册验证码账号已存在            |
|6387  | 验证码写入失败                 |
|6388  | 验证码更新失败                 |
|6389  | 未知的短信类型                 |
|7301  | 绑定账号，账号不存在              |
|7335  | 查询已绑定账号失败，找不到           |
|7336  | 查询第三方绑定信息错误             |
|7338  | 要绑定的账号已经存在              |
|7339  | 获取hash库事务错误             |
|7340  | 获取账号、项目用户事务错误           |
|7341  | hash 账号插入事务执行错误         |
|7342  | 项目表执行插入事务错误             |
|7343  | 账号表执行更新事务错误             |
|7344  | 查询uid查询错误，找不到           |
|7345  | 查询第三方已绑定信息Scan时错误       |
|7346  | 获取项目事务失败                |
|7347  | 查询账号信息错误                |
|7348  | Email已绑定                |
|7349  | 手机号不为空，已绑定              |
|8320  | 解绑，账号不存在                |
|8323  | 邮箱已经解绑                  |
|8324  | 手机号已经解绑                 |
|8332  | 被解绑账号不存在                |
|8333  | 被解绑的账号与当前账号不属于同一个账号下    |
|8334  | 被解绑账号不存在，查询主账号          |
|8335  | 查询第三方账号时错误              |
|8336  | 查询第三方账号Scan时错误          |
|8337  | 不支持解绑注册时的类型             |
|8338  | 获取hash库事务错误             |
|8339  | 获取账号用户事务错误              |
|8340  | 获取项目用户事务错误              |
|8341  | 项目用户删除事务错误              |
|8342  | hash账号更新事务错误            |
|8343  | 账号表执行更新事务错误             |
|8344  | 不支持解绑当前账号               |
|9340  | 登录校验失败                  |
|9341  | 登录Token解析失败             |
|9342  | 登录Token与uid不匹配          |
|9343  | 请求域名不是白名单内              |
|10365 | 添加注销申请失败                |
|10366 | 账号注销修改用户状态失败            |
|10367 | 添加注销申请已存在               |
|10368 | 申请注销的账号不存在              |
|10369 | 账号不存在                   |
|10370 | 账号查询错误                  |
|10371 | 账号uid不匹配                |
|11384 | 撤销账号注销修改用户状态失败          |
|11385 | 撤销账号注销记录不存在             |
|11386 | 撤销账号不存在                 |
|11387 | 撤销账号与记录中的不匹配            |
|12300 | 姓名错误                    |
|12301 | 身份证号错误                  |
|12302 | 调用身份证实名接口错误             |
|12303 | 更新错误                    |
|12304 | 根据account获取账号id不存在      |
|12305 | 更新账号信息错误                

### 第三方账号编码
|第三方|编码|
|:----    |:---  |
|Facebook  | 1001|
|Google    | 1002|
|Twitter   | 1003|
|Youtube   | 1004|
|Instagram | 1005|
|Apple     | 1006|
|Whatsapp  | 1007|
|Skype     | 1008|
|Linkedln  | 1009|
|Line      | 1010|
|VK        | 1011|
|Weixin    | 1012|
|Reddit    | 1013|
|Weibo     | 1014|
|QQ        | 1015|
|Tiktok    | 1016|






