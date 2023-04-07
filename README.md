## 账号系统 v1.0

**一 项目说明**

1. 该系统使用场景主要为游戏公司账号体系，旨在让创业的游戏公司减少账号、支付等基础开发
2. 功能特点
   - 分库分表存储账号信息，容量可大大提高
   - 可指定分库分表数量
   - 支持读写分离
   - 支持苹果删账号政策
   - 支持多种注册方式（email, 手机号, 用户名, 第三方, 游客）
   - 支持email、手机号+验证码或密码方式登录、注册
   - 日志结构化，方便二次处理、查询
   - 支持注册、登录透传数据，方便数据收集
   - 支持测试白名单
   - 限制访问策略，有效保证系统安全
   - 错误码定义清晰，能够快速定位问题
   - 实名认证、防沉迷等
3. 系统接口
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
4. 接口文档见 Document.md

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
        ├── appid              # appid 样例目录
        ├── sql              # 基础数据库sql及账号模板文件
        │   ├── account_base_info.sql       # 基础库sql,包含创建数据库DDL
        │   ├── game_user_delete_tpl.sql       # 项目用户删除申请表
        │   ├── game_user_tpl.sql       # 项目用户表
        │   └── main_user_tpl.go       # 主账号及hash表
        ├── go.mod              
        ├── go.sum
        ├── main.go
        ├── config-file-example.toml    #配置文件示例
        ├── Document.sql             #接口文档
        └── README.md

**四 配置说明**

见 config-file-example.toml，内有注释说明


**注意事项**

​	以下配置项，注意登录、token刷新的配置要相同

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

1.account_base_info 基础库，存储短信、邮件模板、节假日等；

表:
- white_user_list  白名单表
- mail_tpl 邮件模板表
- sms_tpl 短信模板表
- user_delete_config 账号注销配置表
- holiday 节假日表，用于防沉迷日期判断

2.account_info_1、2、3 主账号及hash库（默认3个库），每个库内默认有100个账号表和100 hash表；

3.game_user_1、2（默认2个），项目用户表（默认2个），每个库内默认有100个用户表及5个删除申请表；

详情见文件，内有注释说明

**八 接口文档**

见 Document.md