package main

import (
	"context"
	"notionboy/notion"
	"os"
)

func main() {

	content := notion.Content{
		Tags: []string{"notion", "bot"},
		Text: "#notion #bot This is second note created by bot",
	}

	databaseID := os.Getenv("DATABASE_ID")
	notion.CreateNewRecord(context.Background(), databaseID, content)
}
