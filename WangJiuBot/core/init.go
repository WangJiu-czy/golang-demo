package core

import (
	"bufio"
	"github.com/beego/beego/v2/core/logs"
	"os"
	"regexp"
	"time"
)

// 延迟执行的时间
var Duration time.Duration

func init() {
	logger.Info("初始化")
	killp()
	_, err := os.Stat("/etc/WangJiu/")
	if err != nil {
		os.MkdirAll("/etc/WangJiu/", os.ModePerm)
	}
	for _, arg := range os.Args {
		if arg == "-d" {
			initStore()
			Daemon()
		}
	}
	initStore()
	ReadYaml(ExecPath+"/conf/", &Config, "")
	InitReplies()
	initToHandleMessage()
	logs.Debug("已加载的规则")
	for _, function := range functions {
		logs.Debug(function.Rules)
	}
	file, err := os.Open("/etc/WangJiu/sets.conf")
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if v := regexp.MustCompile(`^\s*set\s+(\S+)\s+(\S+)\s+(\S+.*)`).FindStringSubmatch(line); len(v) > 0 {
				b := Bucket(v[1])
				if b.Get(v[2]) != v[3] {
					b.Set(v[2], v[3])
				}
			}
		}
		file.Close()
	}
	initSys()
	Duration = time.Duration(WangJiu.GetInt("duration", 5)) * time.Second
	WangJiu.Set("started_at", time.Now().Format("2006-01-02 15:04:05"))
	api_key := WangJiu.Get("api_key")
	if api_key == "" {
		api_key := time.Now().UnixNano()
		WangJiu.Set("api_key", api_key)
	}
	if WangJiu.Get("uuid") == "" {
		WangJiu.Set("uuid", GetUUID())
	}
}
