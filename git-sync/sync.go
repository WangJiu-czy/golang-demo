package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//使用 `git status -s` 获取更改的文件，使用给定状态代码的正确动词为文件添加前缀，修剪文件名
func gitAffectedFiles(conf Config) []string {
	DebugLog(conf, "parsing files and their state")
	out, _ := runCmd([]string{"git", "status", "-s"})
	r := strings.Split(out, "\n")
	res := make([]string, 0)
	for _, file := range r {
		file = strings.TrimSpace(file)
		if len(file) == 0 {

			break
		}
		change := ""
		switch file[0] {
		case 'M':
			//修改
			change = "modified"
		case 'A':
			change = "added"
		case 'D':
			change = "deleted"
		case 'R':
			change = "renamed"
		case 'C':
			//复制
			change = "copied"
		case 'U':
			change = "update but unmerged"
		case '?':
			continue

		}
		//获取文件名字,拼接具体操作
		if strings.Contains(file, "\"") {
			//Unquote 将 s 解释为单引号、双引号或反引号的 Go 字符串文字，返回 s 引号的字符串值
			v, err := strconv.Unquote(strings.TrimSpace(file[1:]))
			if err != nil {
				log.Fatalln("[ERROR] couldn't parse encoded characters: ", err)
			}
			file = " " + v
		}
		res = append(res, strings.TrimSpace(file[1:])+"("+change+")")
	}
	var c rune
	if len(res) > 1 {
		c = 's'

	}
	DebugLog(conf, fmt.Sprintf("parsed '%d' changed file%c...", len(res), c))
	return res
}
func GitPull(conf Config) {
	_, err := runCmd([]string{"git", "pull"})
	if err != nil {
		log.Println("[WARNING] pulling changes from remote failed: ", err)
	}
	DebugLog(conf, "pull changes from remote")
}
func GitRepoHashChanges(conf Config) bool {
	DebugLog(conf, "checking if repo has changes...")
	out, err := runCmd([]string{"git", "status", "-s"})
	return err == nil && len(out) != 0
}

//检查当前目录是否是 git 存储库
func CheckIfGitRepo(config Config) bool {
	DebugLog(config, "checking if current directory is a git repository...")
	_, err := runCmd([]string{"git", "status", "-s"})
	return err == nil
}

//将所有更改添加到暂存区域
func GitAdd(config Config) {
	DebugLog(config, "adding all changes to the staged area...")
	_, err := runCmd([]string{"git", "add", "-A"})
	if err != nil {
		log.Println("[WARNING] adding to staging area failed: ", err)

	}
	DebugLog(config, "added changes")

}

//推送提交到远程仓库
func GitPush(config Config) {
	DebugLog(config, "pushing commit to remote...")
	_, err := runCmd([]string{"git", "push"})
	if err != nil {
		log.Println("[WARNING] push to remote failed :", err)
		return
	}
	log.Println("pushed commmits to remote...")
}

//提交
func GitCommit(config Config) {
	commitContent := generateCommitContent(config)
	log.Println("new commit:", strconv.Quote(strings.Join(commitContent, " ")))
	_, err := runCmd(commitContent)
	if err != nil {
		log.Println("[WARNING] commiting failed: ", err)
		return
	}
	DebugLog(config, "made commit")
}

//根据用户在Config中的配置生成commit内容
func generateCommitContent(config Config) []string {
	DebugLog(config, "generating commit content...")
	commitTime := time.Now().Format(config.CommitDate)
	commitContent := strings.ReplaceAll(config.CommitFormat, "%date%", commitTime)
	commit := make([]string, 0)
	if config.AddAffectedFiles {
		affectedFiles := gitAffectedFiles(config)
		commitContent += "\n\n" + "Affected files:\n" + strings.Join(affectedFiles, "\n")
		commit = append(commit, strings.Split(config.CommitCommand, " ")...)
	}
	commit = append(commit, commitContent)
	return commit
}
