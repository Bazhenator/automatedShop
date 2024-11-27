package main

import (
	"automatedShop/configs"
	"automatedShop/internal/app/shop"
	"fmt"
	"log"
)

const (
	configFile string = "./configs/config.yaml"
)

func main() {
	conf, err := configs.ReadConfigFromYAML[configs.ShopConfig](configFile)
	if err != nil {
		panic(fmt.Errorf("read of config from '%s' failed: %w", configFile, err))
	}

	err = configs.ValidateConfig(conf)
	if err != nil {
		panic(fmt.Errorf("'%s' parsing failed: %w", configFile, err))
	}

	err = shop.ProcessApp(conf)
	if err != nil {
		log.Fatal(err)
	}
}
