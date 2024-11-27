package shop

import (
	"automatedShop/configs"
	"automatedShop/internal/dataprovider"
	"automatedShop/internal/graphics"
	"automatedShop/internal/repository"
	"automatedShop/internal/services"
	"fmt"
)

func ProcessApp(config *configs.ShopConfig) error {
	provider, err := dataprovider.NewPsqlProvider(config.DbConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize db with error: %w", err)
	}

	r := repository.NewRepository(provider)
	s := services.NewService(r)
	g := graphics.NewAppManager(s)

	g.Run()
	return nil
}
