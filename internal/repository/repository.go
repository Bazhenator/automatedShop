package repository

import (
	"automatedShop/internal/dataprovider"
	db "automatedShop/internal/repository/psql"
)

type Repository struct {
	AuthRepo IAuthRepository
	ShopRepo IShopRepository
}

func NewRepository(provider *dataprovider.Provider) *Repository {
	return &Repository{
		AuthRepo: db.NewAuthProvider(provider),
		ShopRepo: db.NewShopProvider(provider),
	}
}
