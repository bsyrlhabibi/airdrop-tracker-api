package repository

import (
	"github.com/bsyrlhabibi/airdrop/internal/model"
	"gorm.io/gorm"
)

type AirdropTaskRepo struct {
	DB *gorm.DB
}

func NewAirdropTaskRepo(db *gorm.DB) *AirdropTaskRepo {
	return &AirdropTaskRepo{DB: db}
}

func (r *AirdropTaskRepo) FindByAirdropID(airdropID uint) ([]model.AirdropTask, error) {
	var tasks []model.AirdropTask
	err := r.DB.Preload("Category").Where("airdrop_id = ?", airdropID).Order("sort_order ASC, created_at ASC").Find(&tasks).Error
	return tasks, err
}

func (r *AirdropTaskRepo) FindByID(id uint) (*model.AirdropTask, error) {
	var task model.AirdropTask
	err := r.DB.Preload("Category").First(&task, id).Error
	return &task, err
}

func (r *AirdropTaskRepo) Create(task *model.AirdropTask) error {
	return r.DB.Create(task).Error
}

func (r *AirdropTaskRepo) Update(task *model.AirdropTask) error {
	return r.DB.Save(task).Error
}

func (r *AirdropTaskRepo) Delete(id uint) error {
	return r.DB.Delete(&model.AirdropTask{}, id).Error
}
