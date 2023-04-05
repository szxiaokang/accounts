/**
 * @project accounts
 * @filename init.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/21 15:27
 * @version 1.0
 * @description
 * 系统启动初始化信息
 */

package base

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// InitAppService 系统启动入口
func InitAppService() {

	//初始化配置
	initConfig()

	//初始化日志配置
	initLogs()

	//初始化Redis
	initRedis()

	//初始化mysql, 有错误则停止
	initMysql()

	//其它初始化
	initOther()

	//随机字符串初始化
	rand.Seed(time.Now().UnixNano())
}

// 初始化参数配置
func initConfig() {
	confFile := flag.String("conf", "../config-file-example.toml", "the config file path")
	host := flag.String("host", "", "set host")
	id := flag.Int64("id", 8720, "set id")
	buildSql := flag.String("buildDdl", "", "build table sql, value mainUser or gameUser")
	flag.Parse()
	ConfFile = *confFile
	GConf.Server.Host = *host
	GConf.Server.ID = *id
	if *buildSql != "" {
		buildDdl(*buildSql)
	}

	if ConfFile == "" {
		log.Fatal().Msg("the config file empty, please check")
	}

	err := UnmarshalToml(ConfFile, GConf)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to load conf file")
	}
}

// 初始化日志配置
func initLogs() {
	//创建数据上报日志目录
	if GConf.Server.DataLogsPath[len(GConf.Server.DataLogsPath)-1:] != "/" {
		GConf.Server.DataLogsPath = GConf.Server.DataLogsPath + "/"
	}
	_, pathErr := os.Stat(GConf.Server.DataLogsPath)
	if pathErr != nil {
		makeErr := os.MkdirAll(GConf.Server.DataLogsPath, 0755)
		if makeErr != nil {
			log.Fatal().Err(makeErr).Msg("data log path make error")
		}
	}
	//记录上报数据的日志
	LoginDataLog = zerolog.New(NewFileWriter(GConf.Server.DataLogsPath, "login", true))
	RegisterDataLog = zerolog.New(NewFileWriter(GConf.Server.DataLogsPath, "register", true))

	//初始化普通日志
	if GConf.Server.LogRoot[len(GConf.Server.LogRoot)-1:] != "/" {
		GConf.Server.LogRoot = GConf.Server.LogRoot + "/"
	}
	_, pathErr = os.Stat(GConf.Server.LogRoot)
	if pathErr != nil {
		makeErr := os.MkdirAll(GConf.Server.LogRoot, 0755)
		if makeErr != nil {
			log.Fatal().Err(makeErr).Msg("server log path make error")
		}
	}

	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z07:00"
	zerolog.SetGlobalLevel(zerolog.Level(GConf.Server.LogLevel))
	accountLog := NewFileWriter(GConf.Server.LogRoot, "account", false)
	log.Logger = zerolog.New(accountLog).With().Timestamp().Logger()

	//输出到控制台+日志文件的日志
	multi := zerolog.MultiLevelWriter(accountLog, os.Stdout)
	MultipleLog = zerolog.New(multi).With().Timestamp().Logger()

	MultipleLog.Info().Msg("Start account service...")
	MultipleLog.Info().Msgf("Account host (%s)", GConf.Server.Host)
}

// 初始化redis
func initRedis() {
	MultipleLog.Info().Msg("init redis...")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         GConf.RedisConfig.Addr,
		Password:     GConf.RedisConfig.Pass,
		DB:           GConf.RedisConfig.Db,
		PoolSize:     GConf.RedisConfig.PoolSize,
		MinIdleConns: GConf.RedisConfig.MinIdle,
		MaxConnAge:   time.Duration(GConf.RedisConfig.MaxLifetime) * time.Second,
	})

	// 需要使用context库
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping().Result()
	if err != nil {
		MultipleLog.Fatal().Msgf("init redis failure: %s", err.Error())
	}

	//初始化账号uid自增值
	_, err = RedisClient.Get(RedisUidAutoIncrementKey).Int64()
	if err != nil {
		err = RedisClient.Set(RedisUidAutoIncrementKey, GConf.Base.AutoIncrementUid, 0).Err()
		if err != nil {
			MultipleLog.Fatal().Msgf("init redis set AutoIncrementUid error: %s", err.Error())
		}
	}

	MultipleLog.Info().Msg("redis connect success")
}

// 初始化数据库
func initMysql() {
	//如果EnabledGameList配置为空，则停止
	if len(GConf.Base.EnabledGameList) == 0 {
		MultipleLog.Info().Msg("EnabledGameList empty, stop")
		os.Exit(1)
	}

	//base 库
	AccountBaseDb = connMysql(GConf.MysqlAccountBase)
	MultipleLog.Info().Msg("AccountBaseDb connect success")

	//主账号master库
	for id, mysqlConf := range GConf.MysqlAccountMasterList {
		AccountMasterDbMap.Store(id, connMysql(mysqlConf))
		MultipleLog.Info().Msgf("Account master db %s connect success", id)
	}

	//主账号slave库
	for id, mysqlConf := range GConf.MysqlAccountSlaveList {
		AccountSlaveDbMap.Store(id, connMysql(mysqlConf))
		MultipleLog.Info().Msgf("Account slave db %s connect success", id)
	}

	//项目用户master库
	for id, mysqlConf := range GConf.MysqlGameUserMasterList {
		info := strings.Split(id, "_")
		gameId, _ := strconv.Atoi(info[2])
		exists := false
		for _, v := range GConf.Base.EnabledGameList {
			if v == gameId {
				exists = true
				break
			}
		}
		if exists == false {
			MultipleLog.Info().Msgf("game user db master init, EnabledGameList not exists %d, please check config", gameId)
			os.Exit(1)
		}
		GameUserMasterDbMap.Store(id, connMysql(mysqlConf))
		MultipleLog.Info().Msgf("game user master db %s connect success", id)
	}

	//项目用户slave库
	for id, mysqlConf := range GConf.MysqlGameUserSlaveList {
		info := strings.Split(id, "_") // id: master|slave_db_gid_pid_N
		gameId, _ := strconv.Atoi(info[2])
		exists := false
		for _, v := range GConf.Base.EnabledGameList {
			if v == gameId {
				exists = true
				break
			}
		}
		if exists == false {
			MultipleLog.Info().Msgf("game user db slave init, EnabledGameList not exists %d, please check config", gameId)
			os.Exit(1)
		}
		GameUserSlaveDbMap.Store(id, connMysql(mysqlConf))
		MultipleLog.Info().Msgf("game user slave db %s connect success", id)
	}
}

// 其它初始化
func initOther() {
	if GConf.RequestLimitRule.Enabled {
		//将白名单list写入到map, 方便比较
		GConf.RequestLimitRule.WhiteListMap = make(map[string]int)
		for _, v := range GConf.RequestLimitRule.WhiteList {
			GConf.RequestLimitRule.WhiteListMap[v] = 1
		}

		//将服务器登录校验白名单写入map
		GConf.RequestLimitRule.LoginAuthWhiteListMap = make(map[string]int)
		for _, v := range GConf.RequestLimitRule.LoginAuthWhiteList {
			GConf.RequestLimitRule.LoginAuthWhiteListMap[v] = 1
		}
	}
}

// 生成主账号、hash、项目用户 库、表,
// 依据 mainUserTpl.sql生成主账号、hash库表，命令： go run main.go --buildDdl mainUser
// 依据 gameUserTpl.sql、gameUserDeleteTpl.sql 生成项目用户库表，命令： go run main.go --buildDdl gameUser-16-1 #16-1：16 代表项目，1 代表大区
func buildDdl(table string) {
	fmt.Printf("buildDdl params: %s\n", table)
	if table != "mainUser" {
		buildGameUser(table)
	}
	content, _ := os.ReadFile("./sql/main_user_tpl.sql")

	dbNumber := MainAccountDbNumber
	tbNumber := MainAccountTableNumber
	dbNameFormat := "account_info_%d"
	tableDdl := string(content)

	for dbId := 1; dbId <= dbNumber; dbId++ {
		dbName := fmt.Sprintf(dbNameFormat, dbId)
		ddl := fmt.Sprintf("create database %s default character set utf8mb4 collate utf8mb4_0900_ai_ci;\nuse %s;\n", dbName, dbName)
		dbFile := dbName + ".sql"
		mainUserDdl := ""
		for tableId := 1; tableId <= tbNumber; tableId++ {
			mainUserDdl += fmt.Sprintf(tableDdl, tableId, tableId)
		}
		err := os.WriteFile(dbFile, []byte(ddl+mainUserDdl), os.ModePerm)
		if err != nil {
			fmt.Printf("write %s error: %s", dbName, err.Error())
			os.Exit(-1)
		}
	}

	fmt.Printf("\nbuild main user ddl success :)\n")
	os.Exit(1)
}

// 生成项目用户表、申请删除表
func buildGameUser(param string) {
	params := strings.Split(param, "-")
	gameId := params[1]
	platformId := params[2]

	content, err := os.ReadFile("./sql/game_user_tpl.sql")
	if err != nil {
		fmt.Printf("read gameUserTpl.sql error: %s", err.Error())
		os.Exit(-1)
	}
	dbNumber := GameUserDbNumber
	tbNumber := GameUserTableNumber
	dbNameFormat := fmt.Sprintf("game_user_%s_%s", gameId, platformId) + "_%d" //16 代表项目编号，1 代表大区编号
	//fmt.Printf("dbNameFormat value: " + dbNameFormat)
	//return
	tableDdl := string(content)

	//申请删除表&删除记录表
	deleteTableContent, _ := os.ReadFile("./sql/game_user_delete_tpl.sql")
	deleteDdl := string(deleteTableContent)

	for dbId := 1; dbId <= dbNumber; dbId++ {
		dbName := fmt.Sprintf(dbNameFormat, dbId)
		ddl := fmt.Sprintf("create database %s default character set utf8mb4 collate utf8mb4_0900_ai_ci;\nuse %s;\n", dbName, dbName)
		dbFile := dbName + ".sql"
		userDdl := ""
		for tableId := 1; tableId <= tbNumber; tableId++ {
			userDdl += fmt.Sprintf(tableDdl, tableId)
		}
		err = os.WriteFile(dbFile, []byte(ddl+userDdl), os.ModePerm)
		if err != nil {
			fmt.Printf("write %s error: %s", dbName, err.Error())
			os.Exit(-1)
		}

		file, _ := os.OpenFile(dbFile, os.O_APPEND, 0777)

		//注销申请与删除记录表
		ddl = ""
		for tableId := 1; tableId <= GameUserDeleteTableNumber; tableId++ {
			ddl += fmt.Sprintf(deleteDdl, tableId)
		}
		file.WriteString(ddl)
		file.Close()
	}
	fmt.Printf("build game user ddl success\n")
	os.Exit(1)
}
