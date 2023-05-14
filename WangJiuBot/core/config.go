package core

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// 读取配置中的规则
type Yaml struct {
	Replies []Reply
}

var (
	ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	Config      = Yaml{}
)

func ReadYaml(ruleDir string, conf interface{}, _ string) {
	ph := ruleDir + "config.yaml"
	logger.Info(ph)
	if _, err := os.Stat(ruleDir); err != nil {
		os.MkdirAll(ruleDir, os.ModePerm)
	}
	f, err := os.OpenFile(ph, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(f)
	if len(bytes) == 0 {
		return
	}
	f.Close()
	content, err := ioutil.ReadFile(ph)
	if err != nil {
		logger.Warn(fmt.Sprintf("解析配置文件%s出错: %v \n", ph, err))
		return
	}
	if yaml.Unmarshal(content, conf) != nil {
		logger.Warn(fmt.Sprintf("解析配置文件%s出错: %v \n", ph, err))
		return
	}
	logger.Debug(fmt.Sprintf("解析配置文件%s\n", ph))

}
