package main

import (
	"github.com/iris-contrib/middleware/jwt"
)

var j = jwt.New(
	jwt.Config{
		// 从请求头的Authorization字段中提取，这个是默认值
		Extractor: jwt.FromAuthHeader,
		// 设置一个函数返回秘钥，关键在于return []byte("这里设置秘钥")
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("czy"), nil
		},
		// 设置一个加密方法
		SigningMethod: jwt.SigningMethodHS256,
	},
)
