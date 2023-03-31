## 账号系统 v1.0

**一 项目说明**

1. 该系统使用场景主要为游戏公司账号体系，旨在让创业的游戏公司减少账号、支付等基础开发
2. 功能特点
   - 分库分表存储账号信息，容量可大大提高
   - 可指定分库分表数量
   - 支持读写分离
   - 支持苹果删账号政策
   - 服务器登录校验
   - 支持多种注册方式（email, 手机号, 用户名, 第三方, 游客）
   - 支持email、手机号+验证码或密码方式登录、注册
   - 日志结构化，方便二次处理、查询
   - 支持注册、登录透传数据，方便数据收集
   - 支持测试白名单
   - 限制访问策略，有效保证系统安全
   - 错误码定义清晰，能够快速定位问题
3. 系统接口:
    - /user/login         登录
    - /user/register      注册、登录，适用于游客、第三方，减少请求
    - /user/forgetPassword 忘记密码
    - /user/changePassword 修改密码
    - /user/sendSmsCode   发送验证码
    - /user/bindAccount   绑定账号
    - /user/unBindAccount 解绑
    - /user/loginAuth     服务器登录校验
    - /user/applyLogout  账号注销申请
    - /user/undoLogout    撤销账号注销
    - /user/whiteList     白名单校验
    - /captcha/image     图形验证码展示，根据返回的captcha_id+.png，
    - /captcha/verify    图形验证码验证
4. 接口文档见 document.md

**二 环境要求**

```
    go 1.19+
    mysql 5.7+
    redis 6.0+
```

**三 源文件及目录说明**

        ├── base                 # 基础目录, 包括协议定义, 常量, 常用函数, 错误处理等
        │   ├── constants.go     # 错误常量定义
        │   ├── captcha.go       # 图形验证码
        │   ├── defs.go          # 全局变量、常量、结构体定义       
        │   ├── error.go         # 错误处理
        │   ├── init.go          # 启动初始化  
        │   ├── mail.go          # 邮件发送
        │   ├── middleware.go    # http服务中间件
        │   ├── sms.go           # 短信发送
        │   ├── utils.go         # 常用基础函数
        │   ├── logs.go          # 日志
        │   └── limit.go         # 限制访问方法
        ├── controllers          # 控制器目录
        │   ├── users.go         # 账号控制器实体
        │   ├── captcha.go       # 图形验证码
        │   ├── refresh.go       # 配置刷新
        │   └── users_test.go    # 账号控制器单元测试
        ├── models               # 数据库操作model
        │   ├── users_model.go   # 账号操作model
        │   └── refresh_model.go # 定时刷新操作model
        ├── routers              # 路由
        │   └── routers.go       # 登录路由
        ├── go.mod              
        ├── go.sum
        ├── main.go
        ├── config-file-example.toml    #配置文件示例
        ├── gameUserDeleteTpl.sql       #项目用户申请删除SQL模板文件
        ├── gameUserTpl.sql             #项目用户SQL模板文件
        ├── mainUserTpl.sql             #主账号SQL模板文件
        └── README.md

**四 配置说明**
```toml
[Base]
    Env="dev"
    LoginTokenExpires = 7200 #秒
    RefreshTokenExpires = 31536000 #秒
    ServerKey = "R]3(Fo#N}WVksCb@kang*Kang&230219" #token加密key, 与token服务要一致
    CodeExpires = 10 #10分钟, 验证码有效时间
    CodeCheckMax = 10 #30秒内验证码最大校验次数
    CallRestfulTimeout = 60 #seconds，过期时间
    AppIdConfPath = "/data/service/accounts/appid" #AppId 目录
    EnabledGameList = [16,18] #可以使用此系统的游戏id， 项目编号最大6位，大区编号最大3位位
    AutoIncrementUid = 100000 #主账号自增起始uid, 项目起始uid=game_id(最大6位) + platform_id(最大3位) + 0000000000 加上 主账号uid， 如 100100000100000，每个项目用户支持无限制增长，设计的在99亿内很好区分

[RefreshTime]
    AppKeyRefreshTime = 7200 #单位秒，app key的定时刷新时间
    UserStatusRefreshTime = 60 #单位秒，用户状态的定时刷新时间
    HolidayRefreshTime = 86400 #单位秒， 节假日刷新时间
[HttpTimeout]
    ReadTimeout = 300 #http Server ReadTimeout
    WriteTimeout = 300 #http Server WriteTimeout
    IdleTimeout = 300 #http Server IdleTimeout
[MysqlTimeout]
    MysqlTimeout = "30s" #Timeout for establishing connections
    MysqlReadTimeout = "30s" #I/O read timeout
    MysqlWriteTimeout = "30s" #I/O write timeout

#请求限制规则
[RequestLimitRule]
    Enabled = true #是否启用访问限制，true代表启用, false关闭
    WhiteList = ["/user/loginAuth", "/user/heartbeat", "/captcha/image", "/captcha/verify"] #不限制访问的白名单
    Ip = [3, 100, 1800] #公用, 3秒内，一个ip最多允许请求100次，超过则锁定1800秒，需要输入验证码正确才能继续，此处如果使用了Nginx, 最好用Nginx那一环来做
    Login = [3600, 10, 86400] #登录，一个账号，1小时(3600秒)内，登录失败10次，锁定此ip 24小时(86400秒)不能登录
    VerifyCode = [60, 1, 10] #验证码发送频率，一个账号，60秒内仅允许发送1次。 一个ip, 60秒内只能发送10次，超过锁此ip 60秒,需要输入验证码正确才能继续
    Register = [10, 10, 86400] #一个ip 10秒内注册成功10个后, 锁定此ip 24小时，再注册时需要输入验证码正确才能继续
    LoginAuthWhiteList = [] #服务器登录校验白名单，ip或域名，空值 代表不限制
[Server]
    Name = "accounts"                                    #服务名称
    LogRoot = "/data/logs/service/account/"              #日志地址
    DataLogsPath = "/data/logs/service/account_data/"    #数据日志地址, 数据部门需要的数据, 不做任何处理, 客户端上报什么就写入什么
    Host = ":8080"

#短信业务
[AliSmsConfig]
    RegionId = "cn-hangzhou"
    AccessId = "xx"
    SecretKey = "xx"

#邮件
[MailConfig]
    Hostname = "smtp.qq.com"
    Port = 465
    Username = "xx@qq.com"
    Password = "xx"
    Charset = "utf-8"

#账号基本功能库
[MysqlAccountBase]
    Host        = "127.0.0.1"
    User        = "root"
    Port        = 3306
    Pass        = "123456"
    MaxConn     = 200
    MaxIdle     = 1000
    MaxLifetime = 30
    DbName      = "account_base_info"

#Redis
[RedisConfig]
    Addr = "127.0.0.1:6379"
    Pass = "123456"
    Db = 8
    PoolSize     = 200
    MinIdle     = 10
    MaxLifetime = 30

#账号MasterDB
[MysqlAccountMasterList]
    [MysqlAccountMasterList.account_master_db_1]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "account_info_1"
    [MysqlAccountMasterList.account_master_db_2]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "account_info_2"
    [MysqlAccountMasterList.account_master_db_3]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "account_info_3"

#账号SlaveDB
[MysqlAccountSlaveList]
    [MysqlAccountSlaveList.account_slave_db_1]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "account_info_1"
    [MysqlAccountSlaveList.account_slave_db_2]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "account_info_2"
    [MysqlAccountSlaveList.account_slave_db_3]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "account_info_3"

#项目用户Slave DB
[MysqlGameUserMasterList]
    # 16 项目，1区
    [MysqlGameUserMasterList.master_db_16_1_1] #16_1_1, 16代表项目编号， 1代表大区， 1 代表第一个数据库
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_16_1_1"
    [MysqlGameUserMasterList.master_db_16_1_2]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_16_1_2"
    # 16项目，2区
    [MysqlGameUserMasterList.master_db_16_2_1]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_16_2_1"
    [MysqlGameUserMasterList.master_db_16_2_2]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_16_2_2"
    # 18项目， 1区
    [MysqlGameUserMasterList.master_db_18_1_1]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_18_1_1"
    [MysqlGameUserMasterList.master_db_18_1_2]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 10000
        MaxLifetime = 30
        DbName      = "game_user_18_1_2"
# 从库
[MysqlGameUserSlaveList]
    # 16项目，1区
    [MysqlGameUserSlaveList.slave_db_16_1_1]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_16_1_1"
    [MysqlGameUserSlaveList.slave_db_16_1_2]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_16_1_2"

    # 16项目，2区
    [MysqlGameUserSlaveList.slave_db_16_2_1]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_16_2_1"
    [MysqlGameUserSlaveList.slave_db_16_2_2]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_16_2_2"

    # 18项目，1区
    [MysqlGameUserSlaveList.slave_db_18_1_1]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_18_1_1"
    [MysqlGameUserSlaveList.slave_db_18_1_2]
        Host        = "127.0.0.1"
        User        = "root"
        Port        = 3306
        Pass        = "123456"
        MaxConn     = 200
        MaxIdle     = 1000
        MaxLifetime = 30
        DbName      = "game_user_18_1_2"

```


**注意事项**

​	以下配置项，注意登录、支付、token刷新的配置要相同

```
ServerKey          		#生成token的key
LoginTokenExpires  		#登录token过期时间，秒数
RefreshTokenExpires	    #刷新token过期时间，秒数
```

**五 编译&启动**

    1 编译
    	go build
    2 启动
    	./accounts

**六 数据库生成及配置**
   
   1 生成数据库、表sql文件，并执行生成后的sql

      生成主账号、hash库表，命令： go run main.go --buildDdl mainUser
      生成项目用户库表，命令： go run main.go --buildDdl gameUser-16-1 #16-1：16 代表项目，1 代表大区

   2 设置配置文件中的 Mysql、Redis连接信息，
     

**七 数据库说明及定义**

1.account_base_db 基础库，存储短信、邮件模板、节假日等

主账号库，默认有3个，分表为 account_info_1、2、3， 每个库内默认有100个账号表和100个hash表.

hash表存储的是账号对应的uid，如 a123456@gmail.com 对应的uid为 100018，算法是md5账号，然后在crc32 取模得到库、表

账号表为存储账号信息，如email、手机号、uid、密码等，根据uid取模得到库、表以及项目用户的库、表

  
基础数据库表:
- white_user_list  白名单表
- mail_tpl 邮件模板表
- sms_tpl 短信模板表
- user_delete_config 账号注销配置表
- holiday 节假日表，用于防沉迷日期判断

账号及hash表:

```sql
CREATE TABLE `account_%d`
(
    `uid`             bigint NOT NULL,
    `username`        varchar(64)           DEFAULT NULL COMMENT '用户名',
    `email`           varchar(128)          DEFAULT NULL COMMENT '用户的email',
    `guest`           varchar(128)          DEFAULT NULL COMMENT '游客',
    `third`           varchar(128)          DEFAULT NULL COMMENT '第三方账号，thirdName_ThirdId,如fb_112257954430192',
    `mobile`          varchar(32)           DEFAULT NULL COMMENT '手机号',
    `password`        char(32)     NOT NULL DEFAULT '' COMMENT '密码',
    `created_time`    int(11) NOT NULL DEFAULT '0' COMMENT '注册时间',
    `created_ip`      varchar(16)  NOT NULL DEFAULT '' COMMENT '注册时ip',
    `updated_time`    int(11) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
    `last_login_time` int(11) NOT NULL DEFAULT '0' COMMENT '最后登录时间',
    `name`            varchar(32)  NOT NULL DEFAULT '' COMMENT '实名认证',
    `card_id`         varchar(32)  NOT NULL DEFAULT '' COMMENT '身份证id',
    `type`            tinyint(1) NOT NULL DEFAULT '0' COMMENT '注册类型 1: email, 2:手机号, 3: 用户名, 4: 游客, 5: 第三方',
    `salt`            char(16)     NOT NULL DEFAULT '' COMMENT '密码盐',
    `device_type`     tinyint(1) NOT NULL DEFAULT '0' COMMENT '1 ios, 2 android, 3 pc/h5, 4 其他',
    `lang`            varchar(32)  NOT NULL DEFAULT '' COMMENT '用户语言',
    `extra`           varchar(255) NOT NULL DEFAULT '' COMMENT '其他扩展部分',
    PRIMARY KEY (`uid`) USING BTREE,
    UNIQUE KEY `email` (`email`) USING BTREE,
    UNIQUE KEY `username` (`username`) USING BTREE,
    UNIQUE KEY `guest` (`guest`) USING BTREE,
    UNIQUE KEY `third` (`third`) USING BTREE,
    UNIQUE KEY `mobile` (`mobile`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='账号表';

CREATE TABLE `account_hash_%d`
(
    `account` varchar(128) NOT NULL COMMENT 'Hash的账号',
    `uid`     bigint DEFAULT NULL COMMENT '账号id',
    PRIMARY KEY (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='账号Hash表';

```

项目用户表

```sql
CREATE TABLE `user_%d`
(
    `account`      varchar(128) NOT NULL COMMENT '账号account,确保数据唯一性',
    `uid`          bigint       NOT NULL COMMENT '项目uid',
    `main_uid`     bigint       NOT NULL COMMENT '主账号uid',
    `status`       tinyint(1) NOT NULL DEFAULT '1' COMMENT '1：正常状态, 2：被禁用，3：注销中，4：已注销',
    `type`         tinyint(4) NOT NULL DEFAULT '0' COMMENT '注册类型 1: email, 2:手机号, 3: 用户名, 4: 游客, 5: 第三方',
    `created_time` int          NOT NULL COMMENT '创建时间',
    `ext`          varchar(255) DEFAULT '' COMMENT '扩展备用信息',
    PRIMARY KEY (`account`),
    KEY            `main_uid` (`main_uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='项目用户表';
```

项目用户删除表
```sql
CREATE TABLE `user_delete_apply_%d`
(
    `uid`                 bigint       NOT NULL COMMENT '项目uid',
    `main_uid`            bigint       NOT NULL COMMENT '账号uid',
    `account`             varchar(128) NOT NULL COMMENT '账号account',
    `type`                tinyint(4) NOT NULL DEFAULT '0' COMMENT '注册类型 1: email, 2:手机号, 3: 用户名, 4: 游客, 5: 第三方',
    `status`              tinyint(1) NOT NULL DEFAULT '1' COMMENT '1：申请注销冷静期中或立即注销，2：系统处理中，3：删除成功, 4: 申请恢复，5：恢复成功',
    `apply_time`          int          NOT NULL COMMENT '申请时间',
    `execute_delete_time` int          NOT NULL COMMENT '执行注销时间',
    `ext_info`            text         NOT NULL COMMENT '其他信息，json格式，例如包含调用苹果revoke请求的refresh_token',
    `ext`                 varchar(255) DEFAULT '' COMMENT '扩展备用信息, 信息来自于项目用户表ext字段',
    PRIMARY KEY (`account`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='注销申请表';
```