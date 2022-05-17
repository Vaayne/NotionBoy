package main

import (
	"notionboy/internal/app"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db"
)

func main() {
	config.LoadConfig(config.GetConfig())
	db.InitDB()
	app.Run()
}
