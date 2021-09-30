package notion

import (
	"context"
	"fmt"
	"notionboy/config"
	"regexp"
	"strings"

	notionapi "github.com/kjk/notion"
	"github.com/sirupsen/logrus"
)

type Notion interface {
	ParseContent()
	CreateNewRecord()
}

func GetNotionClient(token string) *notionapi.Client {
	return notionapi.NewClient(token, nil)
}

type Content struct {
	Tags []string
	Text string
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

func CreateNewRecord(ctx context.Context, notionConfig config.Notion, content Content) string {

	content.parseTags()

	var multiSelect []notionapi.SelectOptions

	for _, tag := range content.Tags {
		selectOption := notionapi.SelectOptions{
			Name: tag,
		}
		multiSelect = append(multiSelect, selectOption)
	}

	databasePageProperties := notionapi.DatabasePageProperties{
		"Text": notionapi.DatabasePageProperty{
			Type: "rich_text",
			RichText: []notionapi.RichText{
				{
					Type: "text",
					// PlainText: content.Text,
					Text: &notionapi.Text{
						Content: content.Text,
					},
				},
			},
		},
	}

	if multiSelect != nil {
		databasePageProperties["Tags"] = notionapi.DatabasePageProperty{
			Type:        "multi_select",
			MultiSelect: multiSelect,
		}
	}
	params := notionapi.CreatePageParams{
		ParentType:             notionapi.ParentTypeDatabase,
		ParentID:               notionConfig.DatabaseID,
		DatabasePageProperties: &databasePageProperties,
	}
	client := notionapi.NewClient(notionConfig.BearerToken, nil)
	page, err := client.CreatePage(ctx, params)
	var msg string
	if err != nil {
		msg = fmt.Sprintf("创建 Note 失败，失败原因, %v", err)
		logrus.Error(msg)
	} else {
		pageID := strings.Replace(page.ID, "-", "", -1)
		msg = fmt.Sprintf("创建 Note 成功，如需编辑更多，请前往 https://www.notion.so/%s", pageID)
		logrus.Info(msg)
	}
	return msg
}
