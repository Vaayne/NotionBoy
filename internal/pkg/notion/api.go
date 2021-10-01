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
	return notionapi.NewClient(notionapi.Token(token), nil)
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
	}

	if multiSelect != nil {
		databasePageProperties["Tags"] = notionapi.MultiSelectProperty{
			Type:        "multi_select",
			MultiSelect: multiSelect,
		}
	}
	pageCreateRequest := &notionapi.PageCreateRequest{
		Parent: notionapi.PageCreateDatabaseParent{
			DatabaseID: notionConfig.DatabaseID,
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

func BindNotion(ctx context.Context, token, databaseID string) (bool, error) {
	// 第一次绑定的时候自动建立 Text 和 Tags，确保绑定成功
	cfg := &NotionConfig{BearerToken: token, DatabaseID: databaseID}
	_, err := UpdateDatabase(ctx, cfg)
	if err != nil {
		return false, err
	}

	content := &Content{Text: "#NotionBoy 欢迎🎉使用 Notion Boy!"}
	_, err = CreateNewRecord(ctx, cfg, content)
	if err != nil {
		return false, err
	}
	return true, nil
}
