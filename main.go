package main

import (
	"fmt"
	"go-trading/bitflyer"
	"go-trading/config"
	"go-trading/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	fmt.Println(apiClient.GetBalance())
}
