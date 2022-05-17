package wxgzh

import (
	"fmt"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db"
	notion "notionboy/internal/pkg/notion"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

func unBindingNotion(c *gin.Context, msg *message.MixMessage) *message.Reply {
	db.DeleteWxAccount(msg.GetOpenID())
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("成功解除 Notion 绑定！")}
}

func bindNotion(c *gin.Context, msg *message.MixMessage) *message.Reply {
	log.Warn("----- bindNotion ------")
	userID := msg.GetOpenID()
	userType := db.UserTypeWechat
	stage := fmt.Sprintf("%s:%s", userType, userID)
	oauthMgr := notion.GetOauthManager()
	url := oauthMgr.OAuthURL(stage)
	log.Info("OAuthURL: ", url)
	text := config.BindNotionText
	text += url
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(text)}
}
