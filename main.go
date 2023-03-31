/**
 * @project Accounts
 * @filename main.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/15 10:02
 * @version 1.0
 * @description
 *
 */

package main

import (
	"accounts/base"
	"accounts/controllers"
	"accounts/routers"
)

func main() {
	//系统初始化,启动
	base.InitAppService()

	//定时刷新, 将数据加载到内存
	controllers.TimingRefresh()

	//加载路由以及启动服务
	routers.InitRouterService()
}
