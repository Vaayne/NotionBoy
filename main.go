package main

import (
	"context"
	"fmt"
	"notionboy/config"
	"notionboy/notion"
)

func main() {
	config.LoadConfig(config.GetConfig())

	content := notion.Content{
		Text: "#notion #bot This is #second #note created by bot",
	}

	res := notion.CreateNewRecord(context.Background(), config.GetConfig().DatabaseID, content)
	fmt.Println(res)
}
