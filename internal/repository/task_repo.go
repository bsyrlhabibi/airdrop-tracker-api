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
// For recurring tasks (daily/weekly/monthly), it auto-creates new entries
// for the requested date if they don't exist yet. This ensures each day
// has independent task entries with their own status.
func (r *TaskRepo) FindTodayByAccount(accountID uint, date string) ([]model.Task, error) {
	// First: auto-expand recurring tasks for this date
	r.expandRecurringTasks(accountID, date)

	// Then: return all tasks for this date
	var tasks []model.Task
	err := r.DB.
		Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
		Where("account_airdrops.account_id = ?", accountID).
		Where("tasks.date = ?", date).
		Preload("Category").
		Preload("AccountAirdrop").
		Preload("AccountAirdrop.Airdrop").
		Order("tasks.status ASC, tasks.created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

// expandRecurringTasks creates new task entries for recurring tasks
// that don't have an entry for the given date yet.
func (r *TaskRepo) expandRecurringTasks(accountID uint, date string) {
	// Parse the target date
	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return
	}

	// Find all account-airdrops for this account
	var accountAirdrops []model.AccountAirdrop
	r.DB.Where("account_id = ?", accountID).Find(&accountAirdrops)

	for _, aa := range accountAirdrops {
		// Find recurring tasks (daily/weekly/monthly) for this account-airdrop
		// that are "seed" tasks (original templates)
		var recurringTasks []model.Task
		r.DB.Where("account_airdrop_id = ? AND frequency != 'once'", aa.ID).Find(&recurringTasks)

		for _, rt := range recurringTasks {
			// Check if a task already exists for this date with this name+account_airdrop
			var existing int64
			r.DB.Model(&model.Task{}).
				Where("account_airdrop_id = ? AND name = ? AND date = ?", aa.ID, rt.Name, date).
				Count(&existing)

			if existing > 0 {
				continue // Already has entry for this date
			}

			// Check if this recurring task should appear on the target date
			if !shouldRecurOnDate(rt, targetDate) {
				continue
			}

			// Create new task entry for this date with "pending" status
			newTask := model.Task{
				AccountAirdropID: aa.ID,
				CategoryID:       rt.CategoryID,
				Name:             rt.Name,
				Status:           "pending",
				Frequency:        rt.Frequency,
				Date:             &targetDate,
			}
			r.DB.Create(&newTask)
		}
	}
}

// shouldRecurOnDate checks if a recurring task should appear on the target date.
func shouldRecurOnDate(task model.Task, target time.Time) bool {
	if task.Date == nil {
		return true // No start date, always show
	}

	start := *task.Date

	// Normalize to date only
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	targetDate := time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, time.UTC)

	// Target must be >= start date
	if targetDate.Before(startDate) {
		return false
	}

	switch task.Frequency {
	case "daily":
		return true
	case "weekly":
		// Same day of week
		return targetDate.Weekday() == startDate.Weekday()
	case "monthly":
		// Same day of month
		return targetDate.Day() == startDate.Day()
	default:
		return false
	}
}
