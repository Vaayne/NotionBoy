package wxgzh

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/r2"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	oactx "github.com/silenceper/wechat/v2/officialaccount/context"
	"github.com/sirupsen/logrus"
)

const GET_MEDIA_URL = "https://api.weixin.qq.com/cgi-bin/media/get"

var (
	httpClient *resty.Client
	r2Client   r2.R2
	reg        *regexp.Regexp
)

func init() {
	httpClient = resty.New()
	httpClient.SetTimeout(60 * time.Second)
	r2Client = r2.NewR2Client(config.GetConfig().R2Config.Token, config.GetConfig().R2Config.Url)
	logrus.Debugf("init r2Client: %v", r2Client)
	reg, _ = regexp.Compile(`filename="(.*)"`)
}

//Media 临时素材管理
type Media struct {
	*oactx.Context
}

func NewMedia(ctx *oactx.Context) *Media {
	return &Media{ctx}
}

type MediaResp struct {
	Errcode  int    `json:"errcode"`
	Errmsg   string `json:"errmsg"`
	VideoURL string `json:"video_url"`
}

type GetMediaResp struct {
	MediaResp
	Data        []byte
	ContentType string
	FileName    string
	R2URL       string
}

func (m *Media) getMedia(ctx context.Context, mediaId, userID string) (*GetMediaResp, error) {
	getMediaResp, err := m.downloadMediaFromWx(ctx, mediaId)
	if err != nil {
		return nil, err
	}
	err = m.uploadMediaToR2(ctx, getMediaResp, userID)
	if err != nil {
		return nil, err
	}
	return getMediaResp, nil
}

func (m *Media) downloadMediaFromWx(ctx context.Context, mediaId string) (*GetMediaResp, error) {
	getMediaResp := &GetMediaResp{}

	accessToken, err := m.GetAccessToken()
	if err != nil {
		getMediaResp.Errmsg = fmt.Sprintf("get access token error: %v", err)
		return getMediaResp, err
	}

	uri := fmt.Sprintf("%s?access_token=%s&media_id=%s", GET_MEDIA_URL, accessToken, mediaId)
	resp, err := httpClient.R().
		SetContext(ctx).
		Get(uri)
	if err != nil {
		getMediaResp.Errmsg = fmt.Sprintf("get media from wx error: %v", err)
		return getMediaResp, err
	}

	respContentType := resp.Header().Get("Content-Type")
	if strings.Contains(respContentType, "application/json") {
		var res MediaResp
		err = json.Unmarshal(resp.Body(), &res)
		if err != nil {
			getMediaResp.Errmsg = fmt.Sprintf("get media Unmarshal resp error: %v", err)
			return getMediaResp, err
		}
		getMediaResp.Errmsg = fmt.Sprintf("get media from wx error: errcode=%v , errmsg=%v", res.Errcode, res.Errmsg)
		return getMediaResp, err
	}

	filename := reg.FindStringSubmatch(resp.Header().Get("Content-disposition"))[1]
	getMediaResp.ContentType = respContentType
	getMediaResp.FileName = filename
	getMediaResp.Data = resp.Body()
	return getMediaResp, nil
}

func (m *Media) uploadMediaToR2(ctx context.Context, r *GetMediaResp, userID string) error {
	objectName := fmt.Sprintf("%s-%s", userID, r.FileName)
	log.Printf("objectName: %s\nfilename: %s", objectName, r.FileName)
	url, err := r2Client.Upload(ctx, objectName, r.ContentType, r.Data)
	if err != nil {
		return fmt.Errorf("upload media to r2 error: %v", err)
	}
	r.R2URL = url
	return nil
}
