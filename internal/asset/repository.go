package asset

import "secureasset-manager/internal/database"

type AssetRepository interface {
	Create(asset *Asset) error
	FindAll() ([]Asset, error)
	FindByID(id string) (*Asset, error)
	Update(asset *Asset) error
	Delete(asset *Asset) error
}

type Repository struct{}

func NewRepository() AssetRepository {
	return &Repository{}
}

func (r *Repository) Create(asset *Asset) error {
	return database.DB.Create(asset).Error
}

func (r *Repository) FindAll() ([]Asset, error) {
	var assets []Asset

	err := database.DB.Find(&assets).Error

	return assets, err
}

func (r *Repository) FindByID(id string) (*Asset, error) {
	var asset Asset

	err := database.DB.First(&asset, id).Error

	return &asset, err
}

func (r *Repository) Update(asset *Asset) error {
	return database.DB.Save(asset).Error
}

func (r *Repository) Delete(asset *Asset) error {
	return database.DB.Delete(asset).Error
}
