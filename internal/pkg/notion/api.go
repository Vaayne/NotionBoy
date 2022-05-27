package notion

import (
	"context"
	"time"

	"fmt"
	"regexp"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

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
		"Name": notionapi.TitleProperty{
			Type: "title",
			Title: []notionapi.RichText{
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

func CreateNewMediaRecord(ctx context.Context, notionConfig *NotionConfig, mediaURL, mediaType string) (string, error) {
	databasePageProperties := notionapi.Properties{
		"Name": notionapi.TitleProperty{
			Type: "title",
			Title: []notionapi.RichText{
				{
					Type: "text",
					Text: notionapi.Text{
						Content: "Media " + time.Now().Local().Format(time.RFC3339),
					},
				},
			},
		},
	}
	var mediaBlock notionapi.Block
	if strings.HasPrefix(mediaType, "image") {
		mediaBlock = notionapi.ImageBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeImage,
			},
			Image: notionapi.Image{
				Type: "external",
				External: &notionapi.FileObject{
					URL: mediaURL,
				},
			},
		}
	} else if strings.HasPrefix(mediaType, "video") {
		mediaBlock = notionapi.VideoBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeVideo,
			},
			Video: notionapi.Video{
				Type: "external",
				External: &notionapi.FileObject{
					URL: mediaURL,
				},
			},
		}
	} else {
		mediaBlock = notionapi.FileBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeFile,
			},
			File: notionapi.BlockFile{
				Type: "external",
				External: &notionapi.FileObject{
					URL: mediaURL,
				},
			},
		}
	}

	pageCreateRequest := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(notionConfig.DatabaseID),
		},
		Properties: databasePageProperties,
		Children: []notionapi.Block{
			mediaBlock,
			notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							Type: "text",
							Text: notionapi.Text{
								Content: mediaURL,
								Link: &notionapi.Link{
									Url: mediaURL,
								},
							},
						},
					},
				},
			},
		},
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
