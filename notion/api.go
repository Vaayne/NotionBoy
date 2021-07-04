package notion

import (
	"context"
	"fmt"
	"notionboy/config"
	"regexp"
	"strings"

	notionapi "github.com/kjk/notion"
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

func CreateNewRecord(ctx context.Context, databaseID string, content Content) string {

	content.parseTags()

	var multiSelect []notionapi.SelectOptions

	for _, tag := range content.Tags {
		selectOption := notionapi.SelectOptions{
			ID:   "",
			Name: tag,
		}
		multiSelect = append(multiSelect, selectOption)
	}

	params := notionapi.CreatePageParams{
		ParentType: notionapi.ParentTypeDatabase,
		ParentID:   databaseID,
		DatabasePageProperties: &notionapi.DatabasePageProperties{
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
			"Tags": notionapi.DatabasePageProperty{
				Type:        "multi_select",
				MultiSelect: multiSelect,
			},
		},
	}
	client := GetNotionClient(config.GetConfig().BearerToken)
	page, err := client.CreatePage(ctx, params)
	if err != nil {
		return fmt.Sprintf("创建 Note 失败，失败原因, %s", err)
	}
	pageID := strings.Replace(page.ID, "-", "", -1)
	return fmt.Sprintf("成功创建 Note，如需编辑更多，请前往 https://www.notion.so/%s to edit.", pageID)
}
