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
	// First, find the task to get its airdrop_id and name
	var task model.AirdropTask
	if err := r.DB.First(&task, id).Error; err != nil {
		return err
	}

	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Find all AccountAirdrop IDs for this airdrop
		var aaIDs []uint
		tx.Model(&model.AccountAirdrop{}).Where("airdrop_id = ?", task.AirdropID).Pluck("id", &aaIDs)

		// Delete all daily Tasks that match (same name + same account-airdrops)
		if len(aaIDs) > 0 {
			tx.Where("account_airdrop_id IN ? AND name = ?", aaIDs, task.Name).Delete(&model.Task{})
		}

		// Delete the template task itself
		return tx.Delete(&model.AirdropTask{}, id).Error
	})
}
