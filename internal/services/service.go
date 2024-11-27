package services

import (
	"automatedShop/internal/repository"
	authService "automatedShop/internal/services/auth"
	shopService "automatedShop/internal/services/shop"
)

type Service struct {
	AuthService IAuthService
	ShopService IShopService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		AuthService: authService.NewAuthService(repos.AuthRepo),
		ShopService: shopService.NewShopService(repos.ShopRepo),
	}
}
