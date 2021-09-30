package main

import (
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db"
	"notionboy/internal/wxgzh"
)

func main() {
	config.LoadConfig(config.GetConfig())
	db.InitDB()
	wxgzh.Run()
}
