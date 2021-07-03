package notion

import (
	"context"
	"fmt"
	"os"

	"github.com/kjk/notion"
)

type Notion interface {
	ParseContent()
	CreateNewRecord()
}

var globalClient = notion.NewClient(os.Getenv("BEARER_TOKEN"), nil)

// func GetNotionClient(secret string) *notion.Client {

// 	if globalClient == nil {
// 		client := notion.NewClient(secret, nil)
// 		globalClient = client
// 	}
// 	return globalClient
// }

type Content struct {
	Tags []string
	Text string
}

func CreateNewRecord(ctx context.Context, databaseID string, content Content) {

	var multiSelect []notion.SelectOptions

	for _, tag := range content.Tags {
		selectOption := notion.SelectOptions{
			Name: tag,
		}
		multiSelect = append(multiSelect, selectOption)
	}

	params := notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   databaseID,
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Text": notion.DatabasePageProperty{
				Type: "rich_text",
				RichText: []notion.RichText{
					{
						Type: "text",
						// PlainText: content.Text,
						Text: &notion.Text{
							Content: content.Text,
						},
					},
				},
			},
			"Tags": notion.DatabasePageProperty{
				Type:        "multi_select",
				MultiSelect: multiSelect,
			},
		},
	}

	page, err := globalClient.CreatePage(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println(page.ID)
}
