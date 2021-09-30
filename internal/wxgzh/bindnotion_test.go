package wxgzh

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBindNotionConfig(t *testing.T) {

	type testData struct {
		Text       string
		Token      string
		DatabaseID string
	}

	testDatas := []testData{
		{
			Text:       `Token: secret_token\nDatabaseID: DatabaseID`,
			Token:      `secret_token`,
			DatabaseID: `DatabaseID`,
		},
		{
			Text:       `Token: secret_token \nDatabaseID: DatabaseID`,
			Token:      `secret_token`,
			DatabaseID: `DatabaseID`,
		},
		{
			Text:       `Token: secret_token DatabaseID: DatabaseID`,
			Token:      `secret_token`,
			DatabaseID: `DatabaseID`,
		},
		{
			Text:       `Token: secret_token\tDatabaseID: DatabaseID`,
			Token:      `secret_token`,
			DatabaseID: `DatabaseID`,
		},
		{
			Text:       `Token: secret_token     DatabaseID: DatabaseID`,
			Token:      `secret_token`,
			DatabaseID: `DatabaseID`,
		},
	}

	for _, data := range testDatas {
		fmt.Println(data.Text)
		token, databaseID := parseBindNotionConfig(data.Text)
		assert.Equal(t, data.Token, token, "token is invalid")
		assert.Equal(t, data.DatabaseID, databaseID, "databaseID is invalid")
	}
}
