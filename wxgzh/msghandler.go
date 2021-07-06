package wxgzh

import (
	"notionboy/config"
	"notionboy/db"
	"notionboy/notion"
	"notionboy/utils"

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
Token: secret_xxx
DatabaseID: xxxx
`
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(text)}
		} else {
			if checkNotionBinding(c, token, databaseID) {
				log.Debug("Token is valid, saving account.")
				db.SaveAccount(&db.Account{
					NtDatabaseID: databaseID,
					NtToken:      token,
					WxUserID:     userID,
				})
				memCache.Delete(userID)
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("恭喜 🎉 成功绑定 Notion！")}
			} else {
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("绑定 Notion 失败，无效的 Token 或 DatabaseID， 请重新绑定！")}
			}
		}
	}

	// 获取用户信息
	accountInfo := db.QueryAccountByWxUser(msg.GetOpenID())
	if accountInfo.ID == 0 {
		return bindNotion(c, msg)
	}

	res := notion.CreateNewRecord(c, config.Notion{BearerToken: accountInfo.NtToken, DatabaseID: accountInfo.NtDatabaseID}, *content)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
}
