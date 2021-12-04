package db

import (
	"fmt"
	"notionboy/internal/pkg/config"

	"gorm.io/driver/mysql"
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
