package notion

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

func bindNotion(ctx context.Context, token string) (string, error) {
	// 获取用户绑定的 Database ID，如果有多个，只取找到的第一个
	databaseID, err := getDatabaseID(ctx, token)
	if err != nil {
		return "", err
	}

	// 第一次绑定的时候自动建立 Text 和 Tags，确保绑定成功
	cfg := &NotionConfig{BearerToken: token, DatabaseID: databaseID}
	msg, err := updateDatabase(ctx, cfg)
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

func updateDatabase(ctx context.Context, notionConfig *NotionConfig) (string, error) {
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
			"Name": notionapi.TitlePropertyConfig{
				Type: notionapi.PropertyConfigTypeTitle,
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
