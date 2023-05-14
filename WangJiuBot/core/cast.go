package core

import "time"
import "github.com/beego/beego/v2/adapter/httplib"

var channels = []string{}

func init() {
	go func() {
		time.Sleep(time.Second * 20)
		for {
			for _, channel := range channels {
				msg, _ := httplib.Get(channel).String()
				if msg != "" && WangJiu.Get(channel) != msg {
					NotifyMasters(msg)
					WangJiu.Set(channel, msg)
				}
			}
			time.Sleep(time.Minute)
		}

	}()
}
