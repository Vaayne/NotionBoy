package wxgzh

import (
	"notionboy/internal/pkg/db"
	notion "notionboy/internal/pkg/notion"
	"notionboy/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

func (ex *OfficialAccount) messageHandler(c *gin.Context, msg *message.MixMessage) *message.Reply {

	if msg.MsgType == message.MsgType(message.EventSubscribe) {
		return bindNotion(c, msg)
	}

	if msg.MsgType == message.MsgType(message.EventUnsubscribe) {
		return unBindingNotion(c, msg)
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

	// 获取用户信息
	accountInfo := db.QueryAccountByWxUser(msg.GetOpenID())
	if accountInfo.ID == 0 {
		return bindNotion(c, msg)
	}
	notionConfig := &notion.NotionConfig{BearerToken: accountInfo.AccessToken, DatabaseID: accountInfo.DatabaseID}

	// 如果不是最新的 Scheam，更新 Schema
	if !accountInfo.IsLatestSchema {
		notion.UpdateDatabaseProperties(c, notionConfig)
		db.UpdateIsLatestSchema(accountInfo.DatabaseID, true)
	}

	switch msg.MsgType {
	case message.MsgTypeText:
		// 保存文本信息到 Notion
		res, _ := notion.CreateNewRecord(c, notionConfig, content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
	case message.MsgTypeImage, message.MsgTypeVideo, message.MsgTypeVoice:
		// 保存媒体信息到 Notion
		media := NewMedia(ex.officialAccount.GetContext())
		getMediaResp, err := media.getMedia(c, msg.MediaID, accountInfo.DatabaseID)
		if err != nil {
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(err.Error())}
		}
		res, _ := notion.CreateNewMediaRecord(c, notionConfig, getMediaResp.R2URL, getMediaResp.ContentType)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
	}
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("Unsupport Message!")}
}
