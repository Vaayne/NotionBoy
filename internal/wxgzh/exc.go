package wxgzh

import (
	"fmt"
	"net/http"
	"notionboy/internal/pkg/config"
	notion "notionboy/internal/pkg/notion"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

//OfficialAccount 公众号操作样例
type OfficialAccount struct {
	wc              *wechat.Wechat
	officialAccount *officialaccount.OfficialAccount
}

//OfficialAccount new
func NewOfficialAccount(wc *wechat.Wechat) *OfficialAccount {
	//init config
	wechatConfig := config.GetConfig().Wechat
	offCfg := &offConfig.Config{
		AppID:          wechatConfig.AppID,
		AppSecret:      wechatConfig.AppSecret,
		Token:          wechatConfig.Token,
		EncodingAESKey: wechatConfig.EncodingAESKey,
	}
	log.Debug("AppID: ", offCfg.AppID)
	wc.SetCache(cache.NewMemory())
	officialAccount := wc.GetOfficialAccount(offCfg)
	return &OfficialAccount{
		wc:              wc,
		officialAccount: officialAccount,
	}
}

func transformToNotionContent(msg *message.MixMessage) *notion.Content {
	content := notion.Content{
		Text: msg.Content,
	}
	return &content
}

//Serve 处理消息
func (ex *OfficialAccount) Serve(c *gin.Context) {
	// 传入request和responseWriter
	server := ex.officialAccount.GetServer(c.Request, c.Writer)
	server.SkipValidate(true)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		return messageHandler(c, msg)
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		log.Errorf("Serve Error, err=%v", err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		log.Errorf("Send Error, err=%v", err)
		return
	}
}

//GetAccessToken 获取ak
func (ex *OfficialAccount) GetAccessToken(c *gin.Context) {
	ak, err := ex.officialAccount.GetAccessToken()
	if err != nil {
		log.Errorf("get ak error, err=%v", err)
		RenderError(c, err)
		return
	}
	RenderSuccess(c, ak)
}

//GetCallbackIP 获取微信callback IP地址
func (ex *OfficialAccount) GetCallbackIP(c *gin.Context) {
	ipList, err := ex.officialAccount.GetBasic().GetCallbackIP()
	if err != nil {
		log.Errorf("GetCallbackIP error, err=%v", err)
		RenderError(c, err)
		return
	}
	RenderSuccess(c, ipList)
}

//GetAPIDomainIP 获取微信callback IP地址
func (ex *OfficialAccount) GetAPIDomainIP(c *gin.Context) {
	ipList, err := ex.officialAccount.GetBasic().GetAPIDomainIP()
	if err != nil {
		log.Errorf("GetAPIDomainIP error, err=%v", err)
		RenderError(c, err)
		return
	}
	RenderSuccess(c, ipList)
}

//GetAPIDomainIP  清理接口调用次数
func (ex *OfficialAccount) ClearQuota(c *gin.Context) {
	err := ex.officialAccount.GetBasic().ClearQuota()
	if err != nil {
		log.Errorf("ClearQuota error, err=%v", err)
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "success")
}

//RenderError render error
func RenderError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
}

//RenderSuccess render success
func RenderSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
