package wxgzh

import (
	"fmt"
	"notionboy/internal/pkg/db"
	notion "notionboy/internal/pkg/notion"
	"notionboy/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

func messageHandler(c *gin.Context, msg *message.MixMessage) *message.Reply {

	if msg.MsgType == message.MsgType(message.EventSubscribe) {
		return bindNotion(c, msg)
	}

	userID := msg.GetOpenID()
	content := transformToNotionContent(msg)
	memCache := utils.GetCache()
	userCache := memCache.Get(userID)
	log.Infof("UserID: %s, content: %s, msgType: %s, userCache: %s", userID, content, msg.MsgType, userCache)

	if msg.Content == "绑定" {
		return bindNotion(c, msg)
	} else if msg.Content == "解绑" {
		return unBindingNotion(c, msg)
	}

	if memCache.Get(userID) != nil {
		token, databaseID := parseBindNotionConfig(content.Text)
		log.Infof("Token: %s,\tDatabaseID: %s", token, databaseID)
		if token == "" || databaseID == "" {
			text := `
错误的 Token 和 DatabaseID，请按如下格式回复：
Token: secret_xxx,DatabaseID: xxxx
`
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(text)}
		} else {
			flag, err := notion.BindNotion(c, token, databaseID)
			if flag {
				log.Debug("Token is valid, saving account.")
				db.SaveAccount(&db.Account{
					NtDatabaseID: databaseID,
					NtToken:      token,
					WxUserID:     userID,
				})
				memCache.Delete(userID)
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("恭喜 🎉 成功绑定 Notion！")}
			} else {
				msg := fmt.Sprintf("绑定 Notion 失败, 请检查后重新绑定！ 失败原因: %v", err)
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(msg)}
			}
		}
	}

	// 获取用户信息
	accountInfo := db.QueryAccountByWxUser(msg.GetOpenID())
	if accountInfo.ID == 0 {
		return bindNotion(c, msg)
	}

	res, _ := notion.CreateNewRecord(c, &notion.NotionConfig{BearerToken: accountInfo.NtToken, DatabaseID: accountInfo.NtDatabaseID}, content)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
}
