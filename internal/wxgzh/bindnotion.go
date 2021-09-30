package wxgzh

import (
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/utils"
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
Token: secret_xxx,DatabaseID: xxxx
`
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(text)}
}

func parseBindNotionConfig(text string) (string, string) {
	log.Warn("----- parseBindNotionConfig ------")
	r := regexp.MustCompile(`Token:(?P<Token>.*)\W.*DatabaseID:(?P<DatabaseID>.*)`)
	res := r.FindStringSubmatch(text)
	if len(res) < 3 {
		log.Warn("----- parseBindNotionConfig ------")
		return "", ""
	}
	log.Info("Parse TokenL ", res)
	return strings.TrimSpace(res[1]), strings.TrimSpace(res[2])
}

func unBindingNotion(c *gin.Context, msg *message.MixMessage) *message.Reply {
	db.DeleteWxAccount(msg.GetOpenID())
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("成功解除 Notion 绑定！")}
}
