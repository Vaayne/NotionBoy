package notion

import (
	"context"

	"fmt"
	"regexp"
	"strings"

	notionapi "github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

type Notion interface {
	ParseContent()
	CreateNewRecord()
}

type NotionConfig struct {
	DatabaseID  string `json:"database_id"`
	BearerToken string `json:"bearer_token"`
}

func GetNotionClient(token string) *notionapi.Client {
	return notionapi.NewClient(notionapi.Token(token), func(c *notionapi.Client) {})
}

type Content struct {
	Tags []string `json:"tags"`
	Text string   `json:"text"`
}

func (c *Content) parseTags() {
	regexp, _ := regexp.Compile(`#.*? `)
	match := regexp.FindAllString(c.Text, -1)
	if len(match) > 0 {
		tags := make([]string, 0)
		for _, m := range match {
			tag := strings.Trim(m, "# ")
			tags = append(tags, tag)
		}
		c.Tags = tags
	}
}

func CreateNewRecord(ctx context.Context, notionConfig *NotionConfig, content *Content) (string, error) {

	content.parseTags()

	var multiSelect []notionapi.Option

	for _, tag := range content.Tags {
		selectOption := notionapi.Option{
			Name: tag,
		}
		multiSelect = append(multiSelect, selectOption)
	}

	databasePageProperties := notionapi.Properties{
		"Text": notionapi.RichTextProperty{
			Type: "rich_text",
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: notionapi.Text{
						Content: content.Text,
					},
				},
			},
		},
		// todo 后续支持 Name 字段
		// "Name": notionapi.TitleProperty{
		// 	Type: "title",
		// 	Title: []notionapi.RichText{
		// 		{
		// 			Type: "text",
		// 			Text: notionapi.Text{
		// 				Content: content.Text,
		// 			},
		// 		},
		// 	},
		// },
	}

	if multiSelect != nil {
		databasePageProperties["Tags"] = notionapi.MultiSelectProperty{
			Type:        "multi_select",
			MultiSelect: multiSelect,
		}
	}
	pageCreateRequest := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(notionConfig.DatabaseID),
		},
		Properties: databasePageProperties,
	}
	client := notionapi.NewClient(notionapi.Token(notionConfig.BearerToken), func(c *notionapi.Client) {})
	page, err := client.Page.Create(ctx, pageCreateRequest)
	var msg string
	if err != nil {
		msg = fmt.Sprintf("创建 Note 失败，失败原因, %v", err)
		logrus.Error(msg)
	} else {
		pageID := strings.Replace(page.ID.String(), "-", "", -1)
		msg = fmt.Sprintf("创建 Note 成功，如需编辑更多，请前往 https://www.notion.so/%s", pageID)
		logrus.Info(msg)
	}
	return msg, err
}

func UpdateDatabase(ctx context.Context, notionConfig *NotionConfig) (string, error) {
	databaseUpdateRequest := &notionapi.DatabaseUpdateRequest{
		Properties: notionapi.PropertyConfigs{
			"Tags": notionapi.MultiSelectPropertyConfig{
				Type: notionapi.PropertyConfigTypeMultiSelect,
				MultiSelect: notionapi.Select{
					Options: []notionapi.Option{},
				},
			},
			"Text": notionapi.RichTextPropertyConfig{
				Type: notionapi.PropertyConfigTypeRichText,
			},
			// todo 后续支持 Name 字段
			// "Name": notionapi.TitlePropertyConfig{
			// 	Type: notionapi.PropertyConfigTypeTitle,
			// },
		},
	}

	client := notionapi.NewClient(notionapi.Token(notionConfig.BearerToken), func(c *notionapi.Client) {})
	database, err := client.Database.Update(ctx, notionapi.DatabaseID(notionConfig.DatabaseID), databaseUpdateRequest)
	var msg string
	if err != nil {
		msg = fmt.Sprintf("Update Database 失败，失败原因, %v", err)
		logrus.Error(msg)
	} else {
		msg = fmt.Sprintf("成功更新 Database: %s", database.ID.String())
	}
	return msg, err
}

func BindNotion(ctx context.Context, token string) (string, error) {
	// 获取用户绑定的 Database ID，如果有多个，只取找到的第一个
	databaseID, err := getDatabaseID(ctx, token)
	if err != nil {
		return "", err
	}

	// 第一次绑定的时候自动建立 Text 和 Tags，确保绑定成功
	cfg := &NotionConfig{BearerToken: token, DatabaseID: databaseID}
	msg, err := UpdateDatabase(ctx, cfg)
	logrus.Infof("Update database: %s", msg)
	if err != nil {
		return "", err
	}

	content := &Content{Text: "#NotionBoy 欢迎🎉使用 Notion Boy!"}
	msg, err = CreateNewRecord(ctx, cfg, content)
	logrus.Infof("CreateNewRecord: %s", msg)
	if err != nil {
		return "", err
	}
	return databaseID, nil
}

func getDatabaseID(ctx context.Context, token string) (string, error) {
	logrus.Debug("Token is: ", token)
	cli := GetNotionClient(token)
	searchFilter := make(map[string]string)
	searchFilter["property"] = "object"
	searchFilter["value"] = "database"
	searchReq := notionapi.SearchRequest{
		PageSize: 1,
		Filter: map[string]string{
			"property": "object",
			"value":    "database",
		},
	}
	res, err := cli.Search.Do(ctx, &searchReq)
	if err != nil {
		return "", err
	}
	databases := res.Results
	if len(databases) == 0 {
		return "", fmt.Errorf("至少需要绑定一个 Database")
	}
	database := databases[0].(*notionapi.Database)
	logrus.Debugf("Find Database: %#v", database)
	databaseId := database.ID.String()
	return databaseId, nil
}
