package core

import "time"

type AppStore interface {
	Set(key, val interface{})

	Get(key interface{}, defaultVal ...string) string
	GetInt(key interface{}, defaultVal ...int) int
	GetBool(key interface{}, defaultVal ...bool) bool
	Foreach(f func(k, v []byte) error)
	Create(i interface{}) error
	//TODO:还有
}
type Sender interface {
	GetUserID() string
	GetChatID() int
	GetImType() string
	GetMessageID() int
	GetUsername() string
	IsReply() bool
	GetReplySenderUserID() int
	GetRawMessage() interface{}
	SetMatch([]string)
	SetAllMatch([][]string)
	GetMatch() []string
	GetAllMatch() [][]string
	Get(...int) string
	GetContent() string
	SetContent(string)
	IsAdmin() bool
	IsMedia() bool
	Reply(...interface{}) (int, error)
	Delete() error
	Disappear(lifetime ...time.Duration) //消息处理的等待时间
	Finish()
	Continue()
	IsContinue() bool
	ClearContinue()
	Await(Sender, func(Sender) interface{}, ...interface{}) interface{}
	Copy() Sender
}
