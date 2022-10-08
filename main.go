package main

import (
	"go-trading/config"
	"go-trading/utils"
	"log"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	log.Println("test")
}
