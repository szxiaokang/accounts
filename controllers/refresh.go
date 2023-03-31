/**
 * @project Accounts
 * @filename refresh.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/3/20 11:02
 * @version 1.0
 * @description
 * 相关定时刷新
 */

package controllers

import (
	"accounts/base"
	"accounts/models"
	"net/http"
	"time"
)

// TimingRefresh 定时刷新的信息
func TimingRefresh() {
	//定时刷新接口key
	go refreshAppKey()

	//定时读取game_config表数据,写入内存中
	go refreshGameConfig()

	//定时处理账号注销和账号恢复
	go refreshUserDelete()

	//定时刷新节假日信息
	go refreshHoliday()
}

// 定时刷新接口key
func refreshAppKey() {
	//先初始化一次接口key
	base.RefreshApiAppKey(true)

	//定时刷新接口key
	for range time.Tick(time.Second * time.Duration(base.GConf.RefreshTime.AppKeyRefreshTime)) {
		base.RefreshApiAppKey(false)
	}
}

// 定时刷新game_config
func refreshGameConfig() {
	//先初始化一次
	models.RefreshGameConfig()

	//定时刷新
	for range time.Tick(time.Second * time.Duration(base.GConf.RefreshTime.GameConfigRefreshTime)) {
		models.RefreshGameConfig()
	}
}

// 定时处理
func refreshUserDelete() {
	//定时处理用户注销请求,遍历game_config中的数据
	for range time.Tick(time.Second * time.Duration(base.GConf.RefreshTime.UserDeleteRefreshTime)) {
		models.RefreshUserDelete()
	}
}

func refreshHoliday() {
	//先初始化一次
	models.RefreshHoliday()

	//定时刷新
	for range time.Tick(time.Second * time.Duration(base.GConf.RefreshTime.HolidayRefreshTime)) {
		models.RefreshHoliday()
	}
}

// 心跳监测，正常返回200
func Heartbeat(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(200)
	return
}
