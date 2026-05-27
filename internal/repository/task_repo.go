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

// FindTodayByAccount finds all tasks for an account on a specific date.
// Includes: tasks with matching date, daily/weekly/monthly tasks (recurring), and tasks with no date.
func (r *TaskRepo) FindTodayByAccount(accountID uint, date string) ([]model.Task, error) {
	var tasks []model.Task
	err := r.DB.
		Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
		Where("account_airdrops.account_id = ?", accountID).
		Where("tasks.date = ? OR tasks.frequency != 'once'", date).
		Preload("Category").
		Preload("AccountAirdrop").
		Preload("AccountAirdrop.Airdrop").
		Order("tasks.status ASC, tasks.created_at DESC").
		Find(&tasks).Error
	return tasks, err
}
