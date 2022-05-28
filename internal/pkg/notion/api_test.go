package notion

import (
	"context"
	"os"
	"testing"

	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func shouldRunTest() bool {
	if os.Getenv("BEARER_TOKEN") == "" || os.Getenv("DATABASE_ID") == "" {
		return false
	}
	return true
}

func TestUpdateDatabase(t *testing.T) {
	if !shouldRunTest() {
		logrus.Info("Skip test: TestUpdateDatabase")
		return
	}

	ctx := context.Background()
	notionConfig := &NotionConfig{
		DatabaseID:  os.Getenv("DATABASE_ID"),
		BearerToken: os.Getenv("BEARER_TOKEN"),
	}

	respMsg, err := updateDatabase(ctx, notionConfig, &notionapi.DatabaseUpdateRequest{})
	assert.Nil(t, err, respMsg)
}

func TestCreateNewRecord(t *testing.T) {
	if !shouldRunTest() {
		logrus.Info("Skip test: TestCreateNewRecord")
		return
	}

	ctx := context.Background()
	notionConfig := &NotionConfig{
		DatabaseID:  os.Getenv("DATABASE_ID"),
		BearerToken: os.Getenv("BEARER_TOKEN"),
	}
	content := &Content{
		Tags: []string{"test"},
		Text: "This is test",
	}

	respMsg, err := CreateNewRecord(ctx, notionConfig, content)
	assert.Nil(t, err, respMsg)
}
