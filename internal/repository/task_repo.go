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

// FindTodayByAccount finds tasks for an account on a specific date.
// It auto-creates account task entries for airdrop tasks whose
// start_date → end_date range covers the requested date.
func (r *TaskRepo) FindTodayByAccount(accountID uint, date string) ([]model.Task, error) {
	// Parse target date
	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	// Auto-expand: ensure account tasks exist for today
	r.expandTasksFromAirdrop(accountID, date, targetDate)

	// Return all tasks for this date (use DATE() for SQLite compatibility)
	var tasks []model.Task
	err = r.DB.
		Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
		Where("account_airdrops.account_id = ?", accountID).
		Where("DATE(tasks.date) = ?", date).
		Preload("Category").
		Preload("AccountAirdrop").
		Preload("AccountAirdrop.Airdrop").
		Order("tasks.status ASC, tasks.created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

// expandTasksFromAirdrop ensures account task entries exist for each
// airdrop task whose start_date→end_date range covers the target date.
func (r *TaskRepo) expandTasksFromAirdrop(accountID uint, date string, targetDate time.Time) {
	target := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)

	// Find all account-airdrops for this account
	var accountAirdrops []model.AccountAirdrop
	r.DB.Where("account_id = ?", accountID).Find(&accountAirdrops)

	for _, aa := range accountAirdrops {
		// Get airdrop tasks (global template) for this airdrop
		var airdropTasks []model.AirdropTask
		r.DB.Where("airdrop_id = ?", aa.AirdropID).Find(&airdropTasks)

		for _, at := range airdropTasks {
			// Check if target date falls within start_date → end_date range
			if at.StartDate == nil {
				continue // No start date, skip
			}

			start := time.Date(at.StartDate.Year(), at.StartDate.Month(), at.StartDate.Day(), 0, 0, 0, 0, time.UTC)

			// Target must be >= start_date
			if target.Before(start) {
				continue
			}

			// If end_date is set, target must be <= end_date
			if at.EndDate != nil {
				end := time.Date(at.EndDate.Year(), at.EndDate.Month(), at.EndDate.Day(), 0, 0, 0, 0, time.UTC)
				if target.After(end) {
					continue
				}
			}

			// Check if account task already exists for this date + airdrop task name
			var existing int64
			r.DB.Model(&model.Task{}).
				Where("account_airdrop_id = ? AND name = ? AND DATE(date) = ?", aa.ID, at.Name, date).
				Count(&existing)

			if existing > 0 {
				continue // Already exists
			}

			// Create new account task entry for this date
			newTask := model.Task{
				AccountAirdropID: aa.ID,
				CategoryID:       at.CategoryID,
				Name:             at.Name,
				Status:           "pending",
				Frequency:        "daily",
				Date:             &target,
			}
			r.DB.Create(&newTask)
		}
	}
}
