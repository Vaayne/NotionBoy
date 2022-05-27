package main

import (
	"notionboy/internal/app"
	"notionboy/internal/pkg/db"
)

func main() {
	db.InitDB()
	app.Run()
}
