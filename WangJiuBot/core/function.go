package core

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	cron "github.com/robfig/cron/v3"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var c *cron.Cron

func init() {
	c = cron.New()
	c.Start()
}

type Function struct {
	Rules   []string
	FindAll bool
	Admin   bool
	Handle  func(s Sender) interface{}
	Cron    string
}

var (
	functions = []Function{}
	name      = func() string {
		return WangJiu.Get("name", "枉久")
	}
	//进程名
	pname = regexp.MustCompile(`([^/\s]+)$`).FindStringSubmatch(os.Args[0])[1]
	//会有协程来消费这个回复对象通道
	Senders chan Sender
)

// 回复监听...执行handler函数
// 监听管道,负责执行关键字匹配的逻辑
func initToHandleMessage() {

	Senders = make(chan Sender)
	go func() {
		for {
			go handleMessage(<-Senders)
		}
	}()
}

// 这个是读取js脚本中的注释规则,提取后解析成能用的rule
func AddCommand(prefix string, cmds []Function) {
	for j := range cmds {
		for i := range cmds[j].Rules {
			if strings.Contains(cmds[j].Rules[i], "raw") {
				cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "raw", "", -1)
				continue
			}
			cmds[j].Rules[i] = strings.ReplaceAll(cmds[j].Rules[i], `\r\a\w`, "raw")
			if strings.Contains(cmds[j].Rules[i], "$") {
				continue
			}
			if prefix != "" {
				cmds[j].Rules[i] = prefix + `\s+` + cmds[j].Rules[i]
			}
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "(", `[(]`, -1)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], ")", `[)]`, -1)
			cmds[j].Rules[i] = regexp.MustCompile(`\?$`).ReplaceAllString(cmds[j].Rules[i], `([\s\S]+)`)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], " ", `\s+`, -1)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "?", `(\S+)`, -1)
			cmds[j].Rules[i] = "^" + cmds[j].Rules[i] + "$"
		}
		functions = append(functions, cmds[j])
		if cmds[j].Cron != "" {
			cmd := cmds[j]
			if _, err := c.AddFunc(cmds[j].Cron, func() {
				cmd.Handle(&Faker{})
			}); err != nil {
				logs.Warn("任务%v添加失败%v", cmds[j].Rules[0], err)
			}
		}
	}
}

// 逻辑处理
func handleMessage(sender Sender) {
	content := TrimHiddenCharacter(sender.GetContent())

	defer sender.Finish()
	defer func() {
		logs.Info("%v ==> %v", sender.GetContent(), "finished")
	}()

	u, g, i := fmt.Sprintf(sender.GetUserID()), fmt.Sprint(sender.GetChatID()), fmt.Sprint(sender.GetImType())
	con := true
	mtd := false
	//遍历sync.map
	//主要的作用是：等待那些需要回复是或者否的指令
	//这个主要是执行程序系统相关的命令
	//第一条handle中包含Await（） 系统指令过来后,会在waits添加key,   然后在这里Range中判断
	waits.Range(func(k, v interface{}) bool {
		c := v.(*Carry)
		vs, _ := url.ParseQuery(k.(string))
		userID := vs.Get("u")
		chatID := vs.Get("c")
		imType := vs.Get("i")
		fGroup := vs.Get("f")
		if imType != i || chatID != g || (userID != u && fGroup == "") {
			//跳过
			return true
		}
		if m := regexp.MustCompile(c.Pattern).FindString(content); m != "" {
			mtd = true
			//会在这里阻塞到sender被消费
			c.Chan <- sender
			//等待到Result管道有数据
			sender.Reply(<-c.Result)

			//如果不继续
			if !sender.IsContinue() {
				con = false
				return false
			}

			content = TrimHiddenCharacter(sender.GetContent())

		}
		return true

	})

	//不继续...结束执行
	if mtd && !con {
		return
	}
	reply.Foreach(func(k []byte, v []byte) error {
		if string(v) == "" {
			return nil
		}
		//解析出一个正则表达式对象
		reg, err := regexp.Compile(string(k))
		if err == nil && reg.FindString(content) != "" {
			//接受到的消息匹配上正则，则把正则对应的回复内容返回
			sender.Reply(string(v))
		}
		return nil

	})

	for _, function := range functions {
		for _, rule := range function.Rules {
			var matched bool

			if function.FindAll {
				if res := regexp.MustCompile(rule).FindAllStringSubmatch(content, -1); len(res) > 0 {
					tmp := [][]string{}
					for i := range res {
						tmp = append(tmp, res[i][1:])
					}
					//添加所有匹配的关键字
					sender.SetAllMatch(tmp)
					matched = true
				}

			} else {
				if res := regexp.MustCompile(rule).FindStringSubmatch(content); len(res) > 0 {
					sender.SetMatch(res[1:])
					matched = true
				}
			}

			if matched {
				log.Println(fmt.Sprintf("%v===========>%v", content, rule))
				if function.Admin && !sender.IsAdmin() {
					sender.Delete()
					sender.Disappear()
					sender.Reply("再捣乱我就报警啦～")
					return
				}
				rt := function.Handle(sender)
				if rt != nil {
					sender.Reply(rt)
				}
				if sender.IsContinue() {
					sender.ClearContinue()
					content = TrimHiddenCharacter(sender.GetContent())
					//如果继续就跳到第一层循环,继续匹配下一个handler函数
					goto goon
				}
				return
			}

		}
	goon:
	}

	//屏蔽关键字
	recall := WangJiu.Get("recall")
	if recall != "" {
		recalled := false
		for _, v := range strings.Split(recall, "&") {
			reg, err := regexp.Compile(v)
			if err == nil && reg.FindString(content) != "" {
				if !sender.IsAdmin() {
					sender.Delete()
					sender.Reply("枉久清除了不好的消息~", time.Duration(time.Second))
					recalled = true
					break
				}
			}
		}
		if recalled {
			return
		}
	}
}

// 通过key, 获取cookies中key对应的值
func FetchCookieValue(ps ...string) string {
	var key, cookies string
	if len(ps) == 2 {
		if len(ps[0]) > len(ps[1]) {
			key, cookies = ps[1], ps[0]
		} else {
			key, cookies = ps[0], ps[1]
		}
	}
	match := regexp.MustCompile(key + `=([^;]*);{0,1}`).FindStringSubmatch(cookies)
	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}
