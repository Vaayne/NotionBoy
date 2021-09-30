package wxgzh

import (
	"fmt"
	"notionboy/internal/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
)

func Run() {
	r := gin.Default()

	wc := wechat.NewWechat()
	account := NewOfficialAccount(wc)
	r.Any("/", account.Serve)
	//获取ak
	r.GET("/api/v1/oa/basic/get_access_token", account.GetAccessToken)
	//获取微信callback IP
	r.GET("/api/v1/oa/basic/get_callback_ip", account.GetCallbackIP)
	//获取微信API接口 IP
	r.GET("/api/v1/oa/basic/get_api_domain_ip", account.GetAPIDomainIP)
	//清理接口调用次数
	r.GET("/api/v1/oa/basic/clear_quota", account.ClearQuota)
	svc := config.GetConfig().Service
	r.Run(fmt.Sprintf("%s:%s", svc.Host, svc.Port))
}
