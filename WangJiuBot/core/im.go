package core

import (
	"errors"
	"fmt"

	"reflect"
	"strings"
	"sync"
	"time"
)

type (
	ImageUrl  string
	ImagePath string
	Edit      int
	Replace   int
	Notify    int
	Article   []string
)

var (
	E Edit
	R Replace
	N Notify
)

type Faker struct {
	Message string
	Type    string
	UserID  string
	ChatID  int
	BaseSender
}

func (sender *Faker) GetRawMessage() interface{} {
	return sender.Message
}
func (sender *Faker) GetUsername() string {
	return ""
}

func (sender *Faker) GetReplySenderUserID() int {
	return 0
}

func (sender *Faker) GetContent() string {
	return sender.Message
}
func (sender *Faker) GetUserID() string {
	return sender.UserID
}
func (sender *Faker) IsAdmin() bool {
	return true
}

func (sender *Faker) Reply(msgs ...interface{}) (int, error) {
	rt := ""
	var n *Notify
	for _, msg := range msgs {
		switch msg.(type) {
		case []byte:
			rt = (string(msg.([]byte)))
		case string:
			rt = (msg.(string))
		case Notify:
			v := msg.(Notify)
			n = &v
		}
	}
	if rt != "" && n != nil {
		NotifyMasters(rt)
	}
	return 0, nil
}

func (sender *Faker) Copy() Sender {
	faker := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(Faker)
	return &faker
}

type BaseSender struct {
	matches [][]string //匹配规则
	goon    bool       //标记是否跳过
	child   Sender
	Content string //内容
}

func (sender *BaseSender) SetMatch(ss []string) {

	//sender.matches = append(sender.matches,ss)
	//这样只会保持一个切片规则 [["","",...]]
	sender.matches = [][]string{ss}
}
func (sender *BaseSender) SetAllMatch(ss [][]string) {
	sender.matches = ss
}

func (sender *BaseSender) SetContent(content string) {
	sender.Content = content
}

func (sender *BaseSender) GetMatch() []string {
	return sender.matches[0]
}
func (sender *BaseSender) GetAllMatch() [][]string {
	return sender.matches
}

func (sender *BaseSender) Continue() {
	//设置标志为跳过
	sender.goon = true
}

func (sender *BaseSender) IsContinue() bool {
	return sender.goon
}

func (sender *BaseSender) ClearContinue() {
	sender.goon = false
}
func (sender *BaseSender) Get(index ...int) string {
	//不传入index时, 默认获取第一条规则
	i := 0
	if len(index) != 0 {
		i = index[0]
	}
	if len(sender.matches) == 0 {
		return ""
	}
	if len(sender.matches[0]) < i+1 {
		return ""
	}
	return sender.matches[0][i]
}

func (sender *BaseSender) Delete() error {
	return nil
}

func (sender *BaseSender) Disappear(lifetime ...time.Duration) {

}

func (sender *BaseSender) Finish() {

}

func (sender *BaseSender) IsMedia() bool {
	return false
}

func (sender *BaseSender) GetRawMessage() interface{} {
	return nil
}

func (sender *BaseSender) IsReply() bool {
	return false
}

func (sender *BaseSender) GetMessageID() int {
	return 0
}

func (sender *BaseSender) GetUserID() string {
	return ""
}
func (sender *BaseSender) GetChatID() int {
	return 0
}
func (sender *BaseSender) GetImType() string {
	return ""
}

var (
	TimeOutError   = errors.New("指令超时")
	InterruptError = errors.New("被其他指令中断")
)

// 主要的作用是：防止同一个用户发送的指令太快
var waits sync.Map

// 携带对象
type Carry struct {
	Chan    chan interface{} //管道
	Pattern string           //模式
	Result  chan interface{} //结果
	Sender  Sender           //回复对象

}

type (
	forGroup string //标记是否群发
	again    string
	YesOrNo  string
	Range    []int
	Switch   []string
)

var (
	YesNo YesOrNo = "yesno"
	Yes   YesOrNo = "yes"
	No    YesOrNo = "no"
)
var (
	ForGroup forGroup
	Again    again
	GoAgain  = func(str string) again {
		return again(str)
	}
)

func (_ *BaseSender) Await(sender Sender, callback func(Sender) interface{}, params ...interface{}) interface{} {
	c := &Carry{}
	timeout := 24 * 60 * 60 * time.Second
	var handleErr func(err error)
	var fg *forGroup

	//选择要用什么模式回复
	for _, param := range params {
		switch param.(type) {
		case string:
			c.Pattern = param.(string)
		case time.Duration:
			du := param.(time.Duration)
			if du != 0 {
				timeout = du
			}
		case func() string:
			callback = param.(func(Sender) interface{})
		case func(error):
			handleErr = param.(func(error))
		case forGroup:
			a := param.(forGroup)
			fg = &a
		}

	}
	//匹配所有内容
	if c.Pattern == "" {
		c.Pattern = `[\s\S]*`
	}
	c.Chan = make(chan interface{}, 1)
	c.Result = make(chan interface{}, 1)
	key := fmt.Sprintf("u=%v&c=%v&i=%v", sender.GetUserID(), sender.GetChatID(), sender.GetImType())
	if fg != nil {
		key += fmt.Sprintf("&t=%v&f=true", time.Now().Unix())
	}
	//说明这个用户已经发送过指令,但是还没执行结束,所以不执行当前这条命令
	//这里会把c存入map
	if oc, ok := waits.LoadOrStore(key, c); ok {
		oc.(*Carry).Chan <- InterruptError
	}
	defer func() {
		waits.Delete(key)
	}()

	for {
		select {
		case result := <-c.Chan:
			switch result.(type) {
			case Sender:
				s := result.(Sender)
				if callback == nil {
					return s.GetContent()
				}
				result := callback(s)
				if v, ok := result.(again); ok {
					if v == "" {
						c.Result <- nil
					} else {
						c.Result <- string(v)
					}
				} else if _, ok := result.(YesOrNo); ok {
					if strings.Contains(strings.ToLower(s.GetContent()), "y") {
						return Yes
					}

					if strings.Contains(strings.ToLower(s.GetContent()), "n") {
						return No
					}
					c.Result <- "Y or n ?"
				} else if vv, ok := result.(Switch); ok {
					ct := s.GetContent()
					for _, v := range vv {
						if ct == v {
							return v
						}
					}
					c.Result <- fmt.Sprintf("请从%s中选择一个。", strings.Join(vv, "、"))
				} else if vv, ok := result.(Range); ok {
					ct := s.GetContent()
					n := Int(ct)
					if fmt.Sprint(n) == ct {
						if (n >= vv[0]) && (n <= vv[1]) {

							return n
						}
					}
					c.Result <- fmt.Sprintf("请从%d~%d中选择一个整数。", vv[0], vv[1])
				} else {
					c.Result <- result
					return nil
				}
			case error:
				if handleErr != nil {
					handleErr(result.(error))
				}
				c.Result <- nil
				return nil
			}
		case <-time.After(timeout):
			if handleErr != nil {
				handleErr(TimeOutError)
			}
			c.Result <- nil
			return nil

		}
	}

}
