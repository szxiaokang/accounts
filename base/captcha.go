/**
 * @project AP2
 * @filename captcha.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/3/29 10:18
 * @version 1.0
 * @description
 * 图形验证码，基于 dchest/captcha
 */

package base

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/go-redis/redis"
	"time"
)

// Redis存储验证码
type StoreImpl struct {
	RDB        *redis.Client
	Expiration time.Duration
}

func (impl *StoreImpl) Set(id string, digits []byte) {
	impl.RDB.Set(fmt.Sprintf(CaptchaFormat, id), string(digits), impl.Expiration)
}

func (impl *StoreImpl) Get(id string, clear bool) (digits []byte) {
	bytes, _ := impl.RDB.Get(fmt.Sprintf(CaptchaFormat, id)).Bytes()
	return bytes
}

// BuildCaptchaId 生成验证码图片id
func BuildCaptchaId() string {
	//需要在New之前进行指定
	captcha.SetCustomStore(&StoreImpl{
		RDB:        RedisClient,
		Expiration: time.Second * CaptchaExpire,
	})
	return captcha.New()
}
