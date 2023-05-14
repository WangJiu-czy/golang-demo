package core

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"strings"
)

var (
	BeforeStop = []func(){}
	pidf       = "/var/run/WangJiu.pid"
)

func Daemon() {
	for _, bs := range BeforeStop {
		bs()
	}
	args := os.Args[1:]
	execArgs := make([]string, 0)
	for i := 0; i < len(args); i++ {
		if strings.Index(args[i], "-d") == 0 {
			continue
		}
		execArgs = append(execArgs, args[i])

	}
	proc := exec.Command(os.Args[0], execArgs...)
	err := proc.Start()
	if err != nil {
		panic(err)
	}
	logger.Info(WangJiu.Get("name", "枉久") + "以静默形式运行")
	//将运行起来的程序写入pid文件, 然后退出当前程序
	err = os.WriteFile(pidf, []byte(fmt.Sprintf("%d", proc.Process.Pid)), 0o644)
	if err != nil {
		logger.Warn(err.Error())
	}
	os.Exit(0)

}

func GitPull(filename string) (bool, error) {
	//TODO:GitPull
	return false, nil
}

func CompileCode() error {
	//TODO:CompileCode
	return nil
}
func Download() error {
	//TODO:Download
	return nil
}
func killp() {

	pids, err := ppid()
	if err == nil {
		if len(pids) == 0 {
			return
		} else {
			exec.Command("sh", "-c", "kill -9 "+strings.Join(pids, " ")).Output()
		}
	}

}

// 找到不是当前程序的pid,方便后面kill 之前运行的程序
func ppid() ([]string, error) {
	pid := fmt.Sprint(os.Getpid())
	pids := []string{}
	rtn, err := exec.Command("sh", "-c", "pidof "+pname).Output()
	if err != nil {
		return pids, err
	}

	re := regexp.MustCompile(`[\d]+`)
	for _, v := range re.FindAll(rtn, -1) {
		if string(v) != pid {
			pids = append(pids, string(v))
		}
	}
	return pids, nil
}

func initSys() {
	AddCommand("", []Function{
		{
			Rules: []string{"raw ^name$"},
			Handle: func(s Sender) interface{} {
				s.Disappear()
				return name()
			},
		},
		{
			Rules: []string{"raw ^升级$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				//TODO:升级指令待实现
				return nil
			},
		},
		{
			Rules: []string{"raw ^编译$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				s.Reply("正在编译程序...", E)
				if err := CompileCode(); err != nil {
					return err.Error() + "编译个🐔8"

				}
				s.Reply("编译程序完毕。", E)
				return nil
			},
		},
		{
			Rules: []string{"raw ^重启$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				s.Disappear()
				WangJiu.Set("rebootInfo", fmt.Sprintf("%v %v %v", s.GetImType(), s.GetChatID(), s.GetUserID()))
				s.Reply("即将重启！", E)
				Daemon()
				return nil
			},
		},
		{
			Rules: []string{"raw ^命令$"},
			Handle: func(s Sender) interface{} {
				s.Disappear()
				ss := []string{}
				ruless := [][]string{}
				for _, f := range functions {
					if len(f.Rules) > 0 {
						rules := []string{}
						for i := range f.Rules {
							rules = append(rules, f.Rules[i])
						}
						ruless = append(ruless, rules)
					}
				}
				for j := range ruless {
					for i := range ruless[j] {
						ruless[j][i] = strings.Trim(ruless[j][i], "^$")
						ruless[j][i] = strings.Replace(ruless[j][i], `(\S+)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `(\S*)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `(.+)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `(.*)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `\s+`, " ", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `\s*`, " ", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `.+`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `.*`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `[(]`, "(", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `[)]`, ")", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `([\s\S]+)`, "?", -1)
					}
					ss = append(ss, strings.Join(ruless[j], "\n"))
				}

				return strings.Join(ss, "\n")
			},
		},
	})
}
