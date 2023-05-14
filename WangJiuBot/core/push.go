package core

import (
	"fmt"
	"strings"
)

//消息推送模块

// 存储tg,qq,wx,等通知功能的函数实现
var Pushs = map[string]func(interface{}, string, interface{}){}

// 针对用户的推送组
var GroupPushs = map[string]func(interface{}, interface{}, string){}

type Chat struct {
	Class  string
	ID     int
	UserID int
}

func (ct *Chat) Push(content interface{}) {
	switch content.(type) {
	case string:
		if push, ok := GroupPushs[ct.Class]; ok {
			push(ct.ID, ct.UserID, content.(string))
		}
	case error:
		if push, ok := GroupPushs[ct.Class]; ok {
			push(ct.ID, ct.UserID, content.(error).Error())
		}
	}
}

// 通知管理员
func NotifyMasters(content string) {
	go func() {
		content = strings.Trim(content, " ")
		if WangJiu.GetBool("ignore_notify", false) == true {
			return
		}
		if token := WangJiu.Get("Bark"); token != "" {
			GetPush(fmt.Sprintf("https://api.day.app/%s/%s/%s", token, "WangJiu-Bot", content))

		}
		for _, class := range []string{"tg", "qq", "wx"} {
			//从通知类型的桶中, 获取推知用户的token
			notify := Bucket(class).Get("notifiers")
			if notify == "" {
				notify = Bucket(class).Get("masters")

			}
			for _, v := range strings.Split(notify, "&") {
				if push, ok := Pushs[class]; ok {
					push(v, content, nil)
				}
			}

		}

	}()

}
