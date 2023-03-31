/**
 * @project Accounts
 * @filename error.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/3/19 16:06
 * @version 1.0
 * @description
 * 自定义错误
 */

package base

type MyError struct {
	Code int
	Msg  string
	Log  string
	Data interface{}
}

func (e *MyError) ErrorCode() int {
	return e.Code
}

func (e *MyError) Error() string {
	if e.Msg == "" {
		return ErrorMsg[e.Code]
	}
	return e.Msg
}
