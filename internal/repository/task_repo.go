package repository

import (
	"github.com/bsyrlhabibi/airdrop/internal/model"
	"gorm.io/gorm"
)

type TaskRepo struct {
	DB *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{DB: db}
}

func (r *TaskRepo) Create(task *model.Task) error {
	return r.DB.Create(task).Error
}

func (r *TaskRepo) FindByAccountAirdrop(accountAirdropID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := r.DB.Preload("Category").Where("account_airdrop_id = ?", accountAirdropID).Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepo) FindByID(id uint) (*model.Task, error) {
	var task model.Task
	err := r.DB.Preload("Category").First(&task, id).Error
	return &task, err
}

func (r *TaskRepo) Update(task *model.Task) error {
	return r.DB.Save(task).Error
}

func (r *TaskRepo) Delete(id uint) error {
	return r.DB.Delete(&model.Task{}, id).Error
}

// FindTodayByUser finds all tasks for a user through account-airdrops and accounts.
func (r *TaskRepo) FindTodayByUser(userID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := r.DB.
		Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
		Joins("JOIN accounts ON accounts.id = account_airdrops.account_id").
		Where("accounts.user_id = ?", userID).
		Preload("Category").
		Preload("AccountAirdrop").
		Preload("AccountAirdrop.Airdrop").
		Order("tasks.created_at DESC").
		Find(&tasks).Error
	return tasks, err
}
