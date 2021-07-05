package main

import (
	"notionboy/config"
	"notionboy/wxgzh"
)

func main() {
	config.LoadConfig(config.GetConfig())
	wxgzh.Run()
}
