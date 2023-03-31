/**
 * @project Accounts
 * @filename middleware.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/10 10:22
 * @version 1.0
 * @description
 * 中间件
 */

package base

import "net/http"

type middlewareFunc func(http.Handler) http.Handler
type Middleware struct {
	m []middlewareFunc
}

func (a Middleware) Append(m middlewareFunc) Middleware {
	a.m = append(a.m, m)
	return a
}
func (a Middleware) Then(h http.Handler) http.Handler {

	for i := range a.m {
		h = a.m[len(a.m)-1-i](h)
	}
	return h
}
