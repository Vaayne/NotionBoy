package notion

import "context"

type Notion interface {
	ParseContent()
	CreateNewRecord()
	UpdateDatabaseProperties(ctx context.Context, cfg *NotionConfig) (string, error)
}
