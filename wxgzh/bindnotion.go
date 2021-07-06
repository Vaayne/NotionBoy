package wxgzh

import (
	"notionboy/config"
	"notionboy/db"
	"notionboy/notion"
	"notionboy/utils"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

func bindNotion(c *gin.Context, msg *message.MixMessage) *message.Reply {
	log.Warn("----- bindNotion ------")
	memCache := utils.GetCache()
	memCache.Set(msg.GetOpenID(), []string{msg.Content}, 60*time.Second)
	text := `
欢迎使用 Notion Boy，您还未绑定。
如需绑定，请在 1 分钟之内，按照下面的格式回复 Notion 的 Token 和 DatabaseID 来绑定

获取 Token 和 DatabaseID 的相关方法，请参考官方文档 https://developers.notion.com/docs/getting-started
Token 是以 "secret_" 开头的字符串，
DatabaseID 则是分享的 Page 链接的后半部分

--- 下面是具体的格式 ---
Token: secret_xxx
DatabaseID: xxxx
`
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(text)}
}

func parseBindNotionConfig(text string) (string, string) {
	log.Warn("----- parseBindNotionConfig ------")
	r := regexp.MustCompile(`Token: (?P<Token>.*) .*DatabaseID: (?P<DatabaseID>.*)`)
	res := r.FindStringSubmatch(text)
	log.Info("Parse TokenL ", res)
	return res[1], res[2]
}

func checkNotionBinding(c *gin.Context, token, databaseID string) bool {
	content := notion.Content{Text: "#NotionBoy 欢迎🎉使用 Notion Boy!"}
	res := notion.CreateNewRecord(c, config.Notion{BearerToken: token, DatabaseID: databaseID}, content)
	return strings.Contains(res, "创建 Note 成功")
}

func unBindingNotion(c *gin.Context, msg *message.MixMessage) *message.Reply {
	db.DeleteWxAccount(msg.GetOpenID())
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("成功解除 Notion 绑定！")}
}
