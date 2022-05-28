package db

import (
	"fmt"
	"notionboy/internal/pkg/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	UserTypeWechat   = "wechat"
	UserTypeTelegram = "telegram"
)

type Account struct {
	gorm.Model
	// UserType wechat or telegram or ...
	UserType string `gorm:"type:varchar(20);not null"`
	UserID   string `gorm:"unique,index,column:user_id"`

	// Notion database id
	DatabaseID string `gorm:"uniqueIndex,not null,column:database_id"`
	// Notion access token
	AccessToken string `gorm:"not null,column:access_token"`
	// IsLatestSchema is latest schema
	IsLatestSchema bool `gorm:"column:is_latest_schema;default:false"`
}

type NotionOauthInfo struct {
	gorm.Model
	AccessToken   string `gorm:"not null,column:access_token"`
	WorkspaceID   string `gorm:"not null,column:workspace_id"`
	WorkspaceName string `gorm:"column:workspace_name"`
	WorkspaceIcon string `gorm:"column:workspace_icon"`
	BotID         string `gorm:"not null,column:bot_id"`
	UserID        string `gorm:"unique,not null,column:user_id"`
	UserName      string `gorm:"column:user_name"`
	UserEmail     string `gorm:"column:user_email"`
	UserInfo      string `gorm:"column:user_info"`
}

var conn *gorm.DB

func GetDBConn() *gorm.DB {
	if conn == nil {
		database := config.GetConfig().Databases
		if database.MySQL.Host != "" {
			m := database.MySQL
			fmt.Printf("Connect Mysql %s:%d", m.Host, m.Port)
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
				m.User, m.Pass, m.Host, m.Port, m.Database)
			c, err := gorm.Open(mysql.New(mysql.Config{
				DSN:                       dsn,   // data source name
				DefaultStringSize:         256,   // default size for string fields
				DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
				DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
				DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
				SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
			}), &gorm.Config{})
			if err != nil {
				panic("failed to connect database")
			}
			conn = c
		} else {
			c, err := gorm.Open(sqlite.Open(config.GetConfig().Databases.Sqlite.File), &gorm.Config{})
			if err != nil {
				panic("failed to connect database")
			}
			conn = c
		}
	}
	return conn
}

func InitDB() {
	db := GetDBConn()
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&NotionOauthInfo{})
}

func QueryAccountByWxUser(wxUserID string) *Account {
	var account Account

	db := GetDBConn()
	db.Where(&Account{UserID: wxUserID, UserType: UserTypeWechat}).First(&account)
	return &account
}

func SaveAccount(account *Account) {
	db := GetDBConn()
	db.Create(account)
}

func DeleteWxAccount(wxUserID string) {
	db := GetDBConn()
	db.Where(&Account{UserID: wxUserID, UserType: UserTypeWechat}).Delete(&Account{})
}

func SaveNotionOauth(info *NotionOauthInfo) {
	db := GetDBConn()
	db.Create(info)
}

func UpdateIsLatestSchema(databaseID string, isLatest bool) {
	db := GetDBConn()
	db.Model(&Account{}).Where("database_id = ?", databaseID).Update("is_latest_schema", isLatest)
	logrus.Debugf("UpdateIsLatestSchema %s %t", databaseID, isLatest)
}
