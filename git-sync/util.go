package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

//TODO:准备

type Config struct {
	//将在提交title前面添加前缀, 默认:"backup: %date%"
	CommitFormat string `json:"commit_format"`

	//指定日期格式，默认为："2006-01-02 15:04:05"
	CommitDate string `json:"commit_date"`
	//列出提交主体中受提交影响的文件名，默认值：true
	AddAffectedFiles bool `json:"add_affected_files"`

	//备份之间的时间间隔（以秒为单位），默认值：300
	BackupInterval int `json:"backup_interval"`

	//提交命令，默认："git commit -m"
	CommitCommand string `json:"commit_command"`

	//启用调试模式（详细日志记录、额外信息等），默认值：false
	DebugMode bool `json:"debug"`

	//启动时启用从远程拉取，默认值：true
	PullOnStart bool `json:"pull_on_start"`
}

/*
读取配置文件
Loads and parses config from:
 - On Unix systems, $XDG_CONFIG_HOME or $HOME/.config
 - On Windows, it returns %AppData%
配置文件位置取决于 os.UserConfigDir()

 if config is not found the fallback config is:

		Config{
	     AutoCommitPrefix:      "backup: ",
	     BackupInterval:        300,
	     CommitCommand:         "git commit -m",
	     AddAffectedFiles:      true,
	     CommitTitleDateFormat: "2023-01-01 15:15:15",
		 DebugMode:             false,
         PullOnStart:           true,
		}
*/
func getConfig() Config {
	fallbackConf := Config{
		CommitFormat:     "backup:%date%",
		CommitDate:       "2006-01-02 15:04:05",
		AddAffectedFiles: true,
		BackupInterval:   300,
		CommitCommand:    "git commit -m",
		DebugMode:        false,
		PullOnStart:      true,
	}

	confDir, _ := os.UserConfigDir()
	confFile := path.Join(confDir, "gitsync.json")
	confContent, err := os.ReadFile(confFile)
	if err != nil {
		log.Println("[WARNING] gitsync config not found: ", err)
		log.Println("using fallback config...")
		return fallbackConf
	}
	resConfig := Config{}
	err = json.Unmarshal(confContent, &resConfig)
	if err != nil {
		log.Println("[WARNING] couldn't parse config", err)
		log.Println("using fallback config")
		return fallbackConf
	}
	return resConfig
}

//检查路径中的 git 可执行文件
func CheckForGit(conf Config) bool {
	DebugLog(conf, "checking for git executable in path ...")
	//LookPath 找到这个命令的位置,返回绝对路径, 而不是 ./go 或者 .\go.exe
	_, err := exec.LookPath("git")
	return err == nil
}

func DebugLog(conf Config, msg string) {
	if conf.DebugMode {
		log.Println("[DEBUG] ", msg)
	}

}

func runCmd(cmd []string) (val string, err error) {
	//cmd[0] 命令 cmd[1:]...都是参数
	command := exec.Command(cmd[0], cmd[1:]...)
	//CombinedOutput 运行命令并返回其组合的标准输出和标准错误。
	out, err := command.CombinedOutput()
	if err != nil {
		return "", err
	}
	//TrimSpace 返回字符串 s 的一部分，删除了所有前导和尾随空格
	return strings.TrimSpace(string(out)), nil
}
