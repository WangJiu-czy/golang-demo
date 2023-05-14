package wxmp

import (
	"WangJiuBot/core"
	"fmt"
	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"io"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	wxmp      = core.NewBucket("wxmp")
	materials = core.NewBucket("wxmpMaterial")
)

func init() {
	file_dir := "logs/wxmp/"
	os.MkdirAll(file_dir, os.ModePerm)
	logs.Debug("加载微信公众号服务")
	core.Server.Any("/wx/", func(c *gin.Context) {
		wc := wechat.NewWechat()
		memory := cache.NewMemory()
		cfg := &offConfig.Config{
			AppID:          wxmp.Get("app_id"),
			AppSecret:      wxmp.Get("app_servet"),
			Token:          wxmp.Get("token"),
			EncodingAESKey: wxmp.Get("encoding_aes_key"),
			Cache:          memory,
		}

		officialAccount := wc.GetOfficialAccount(cfg)
		server := officialAccount.GetServer(c.Request, c.Writer)
		server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
			if msg.Event == "subcribe" {
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(wxmp.Get("subscribe_reply", "感谢关注！"))}
			}
			sender := &Sender{}
			sender.Message = msg.Content
			sender.Wait = make(chan []interface{}, 1)
			sender.uid = fmt.Sprintf(string(msg.FromUserName))
			core.Senders <- sender
			//sender.Wait Reply()将内容写入到Responses中,然后调用Finish（）写入Wait管道中
			end := <-sender.Wait //管道没有数据会阻塞，等待sender处理后调用finsh(),将消息写入sender.Wait的管道中
			ss := []string{}
			url := ""
			if len(end) == 0 {
				ss = append(ss, wxmp.Get("default_reply", "无法回复该消息"))
			}
			for _, item := range end {
				switch item.(type) {
				case error:
					ss = append(ss, item.(error).Error())
				case string:
					ss = append(ss, item.(string))
				case []byte:
					ss = append(ss, string(item.([]byte)))
				case core.ImageUrl:
					url = string(item.(core.ImageUrl))
				}

			}

			mediaID := ""
			if url != "" && len(ss) == 0 {
				filename := file_dir + fmt.Sprint(time.Now().UnixNano()) + ".jpg"
				err := func() error {
					f, err := os.Create(filename)
					if err != nil {
						return err
					}
					rsp, err := httplib.Get(url).Response()
					_, err = io.Copy(f, rsp.Body)

					if err != nil {
						f.Close()
						return err
					}
					f.Close()
					logs.Info("filename=====>" + filename)
					m := officialAccount.GetMaterial()
					mediaID, _, err = m.AddMaterial(material.MediaTypeImage, filename)
					logs.Debug("mediaID===>" + mediaID)
					if err != nil {
						logs.Error("AddMaterial管理素材出错")
						logs.Error(err.Error())
						return err
					}
					logs.Debug("mediaID=======>" + mediaID)
					materials.Set(mediaID, filename)
					return nil
				}()
				if err != nil {

					ss = append(ss, err.Error())
					goto TEXT
				}
				return &message.Reply{MsgType: message.MsgTypeImage, MsgData: message.NewImage(mediaID)}
			}
		TEXT:
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(strings.Join(ss, "\n\n"))}
		})
		err := server.Serve()
		if err != nil {
			return
		}
		server.Send()
	})
}

type Sender struct {
	Message   string
	Responses []interface{}
	Wait      chan []interface{}
	uid       string
	core.BaseSender
}

func (sender *Sender) GetUsername() string {
	return ""
}
func (sender *Sender) GetImType() string {
	return "wxmp"
}
func (sender *Sender) GetReplySenderUserID() int {
	return 0
}
func (sender *Sender) GetUserID() string {
	return sender.uid
}
func (sender *Sender) GetContent() string {
	if sender.Content != "" {
		return sender.Content
	}

	return sender.Message
}
func (sender *Sender) GetRawMessage() interface{} {
	return sender.Message
}

func (sender *Sender) IsAdmin() bool {
	return strings.Contains(wxmp.Get("masters"), fmt.Sprint(sender.uid))
}

func (sender *Sender) Reply(msgs ...interface{}) (int, error) {
	sender.Responses = append(sender.Responses, msgs...)
	return 0, nil
}

func (sender *Sender) Finish() {
	if sender.Responses == nil {
		sender.Responses = []interface{}{}
	}
	sender.Wait <- sender.Responses
}
func (sender *Sender) Copy() core.Sender {
	wxmp := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(Sender)
	return &wxmp
}
