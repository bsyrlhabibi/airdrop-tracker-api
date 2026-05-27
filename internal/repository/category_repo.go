package repository

import (
	"github.com/bsyrlhabibi/airdrop/internal/model"
	"gorm.io/gorm"
)

type CategoryRepo struct {
	DB *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{DB: db}
}

func (r *CategoryRepo) FindByUser(userID uint) ([]model.Category, error) {
	var categories []model.Category
	err := r.DB.Where("user_id = ?", userID).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepo) FindByID(id uint, userID uint) (*model.Category, error) {
	var cat model.Category
	err := r.DB.Where("id = ? AND user_id = ?", id, userID).First(&cat).Error
	return &cat, err
}

func (r *CategoryRepo) Create(cat *model.Category) error {
	return r.DB.Create(cat).Error
}

func (r *CategoryRepo) Update(cat *model.Category) error {
	return r.DB.Save(cat).Error
}

func (r *CategoryRepo) Delete(id uint, userID uint) error {
	return r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Category{}).Error
}
