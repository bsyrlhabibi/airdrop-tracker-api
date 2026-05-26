package repository

import (
	"github.com/bsyrlhabibi/airdrop/internal/model"
	"gorm.io/gorm"
)

type AirdropRepo struct {
	DB *gorm.DB
}

func NewAirdropRepo(db *gorm.DB) *AirdropRepo {
	return &AirdropRepo{DB: db}
}

func (r *AirdropRepo) Create(airdrop *model.Airdrop) error {
	return r.DB.Create(airdrop).Error
}

func (r *AirdropRepo) FindByUser(userID uint) ([]model.Airdrop, error) {
	var airdrops []model.Airdrop
	err := r.DB.Where("user_id = ?", userID).Preload("Tasks").Order("created_at DESC").Find(&airdrops).Error
	return airdrops, err
}

func (r *AirdropRepo) FindByID(id, userID uint) (*model.Airdrop, error) {
	var airdrop model.Airdrop
	err := r.DB.Where("id = ? AND user_id = ?", id, userID).Preload("Tasks").First(&airdrop).Error
	return &airdrop, err
}

func (r *AirdropRepo) Update(airdrop *model.Airdrop) error {
	return r.DB.Save(airdrop).Error
}

func (r *AirdropRepo) Delete(id, userID uint) error {
	return r.DB.Where("user_id = ?", userID).Delete(&model.Airdrop{}, id).Error
}
