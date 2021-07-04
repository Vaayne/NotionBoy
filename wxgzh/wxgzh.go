package wxgzh

import (
	"fmt"
	"net/http"

	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	// "github.com/silenceper/wechat/v2/officialaccount/message"
)

func serveWechat(rw http.ResponseWriter, req *http.Request) {
	wc := wechat.NewWechat()
	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:     "xxx",
		AppSecret: "xxx",
		Token:     "xxx",
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)
	//设置接收消息的处理方法
	// server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
	// 	//回复消息：演示回复用户发送的消息
	// 	text := message.NewText(msg.Content)
	// 	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	// })

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}

func main() {
	http.HandleFunc("/", serveWechat)
	fmt.Println("wechat server listener at", ":8001")
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
