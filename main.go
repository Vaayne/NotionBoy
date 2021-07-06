package main

import (
	"notionboy/config"
	"notionboy/db"
	"notionboy/wxgzh"
)

func main() {
	config.LoadConfig(config.GetConfig())
	db.InitDB()
	wxgzh.Run()
}
