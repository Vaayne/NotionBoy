package db

import (
	"notionboy/internal/pkg/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	WxUserID     string `gorm:"unique,index,column:wx_user_id"`
	TgUserID     string `gorm:"unique,index,column:tg_user_id"`
	NtDatabaseID string `gorm:"uniqueIndex,not null,column:nt_database_id"`
	NtToken      string `gorm:"not null,column:nt_token"`
}

var conn *gorm.DB

func GetDBConn() *gorm.DB {
	if conn == nil {
		c, err := gorm.Open(sqlite.Open(config.GetConfig().Databases.Sqlite.File), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		conn = c
	}
	return conn
}

func InitDB() {
	db := GetDBConn()
	db.AutoMigrate(&Account{})
}

func QueryAccountByWxUser(wxUserID string) *Account {
	var account Account

	db := GetDBConn()
	db.Where(&Account{WxUserID: wxUserID}).First(&account)
	return &account
}

func SaveAccount(account *Account) {
	db := GetDBConn()
	db.Create(account)
}

func DeleteWxAccount(wxUserID string) {
	db := GetDBConn()
	db.Where(&Account{WxUserID: wxUserID}).Delete(&Account{})
}
