package repository

import (
	"time"

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

func (r *TaskRepo) FindByAirdrop(airdropID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := r.DB.Where("airdrop_id = ?", airdropID).Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepo) FindByID(id uint) (*model.Task, error) {
	var task model.Task
	err := r.DB.First(&task, id).Error
	return &task, err
}

func (r *TaskRepo) Complete(id uint) error {
	now := time.Now()
	return r.DB.Model(&model.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_completed": true,
		"completed_at": now,
	}).Error
}

func (r *TaskRepo) Reset(id uint) error {
	return r.DB.Model(&model.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_completed": false,
		"completed_at": nil,
	}).Error
}

func (r *TaskRepo) Delete(id uint) error {
	return r.DB.Delete(&model.Task{}, id).Error
}

func (r *TaskRepo) FindTodayByUser(userID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := r.DB.Joins("JOIN airdrops ON airdrops.id = tasks.airdrop_id").
		Where("airdrops.user_id = ? AND tasks.is_completed = ?", userID, false).
		Preload("Airdrop").
		Order("tasks.created_at DESC").
		Find(&tasks).Error
	return tasks, err
}
