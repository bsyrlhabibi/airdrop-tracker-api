package repository

import (
	"github.com/bsyrlhabibi/airdrop/internal/model"
	"gorm.io/gorm"
)

type AccountAirdropRepo struct {
	DB *gorm.DB
}

func NewAccountAirdropRepo(db *gorm.DB) *AccountAirdropRepo {
	return &AccountAirdropRepo{DB: db}
}

func (r *AccountAirdropRepo) Create(aa *model.AccountAirdrop) error {
	return r.DB.Create(aa).Error
}

func (r *AccountAirdropRepo) FindByID(id uint) (*model.AccountAirdrop, error) {
	var aa model.AccountAirdrop
	err := r.DB.Preload("Airdrop").Preload("Tasks").First(&aa, id).Error
	return &aa, err
}

func (r *AccountAirdropRepo) FindByAccount(accountID uint) ([]model.AccountAirdrop, error) {
	var aas []model.AccountAirdrop
	err := r.DB.Where("account_id = ?", accountID).
		Preload("Airdrop").
		Preload("Tasks").
		Order("created_at DESC").
		Find(&aas).Error
	return aas, err
}

func (r *AccountAirdropRepo) FindByAccountAndAirdrop(accountID, airdropID uint) (*model.AccountAirdrop, error) {
	var aa model.AccountAirdrop
	err := r.DB.Where("account_id = ? AND airdrop_id = ?", accountID, airdropID).
		First(&aa).Error
	return &aa, err
}

func (r *AccountAirdropRepo) Update(aa *model.AccountAirdrop) error {
	return r.DB.Save(aa).Error
}

func (r *AccountAirdropRepo) Delete(id uint) error {
	// Delete tasks first
	r.DB.Where("account_airdrop_id = ?", id).Delete(&model.Task{})
	return r.DB.Delete(&model.AccountAirdrop{}, id).Error
}

func (r *AccountAirdropRepo) GetTasks(accountAirdropID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := r.DB.Where("account_airdrop_id = ?", accountAirdropID).
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

func (r *AccountAirdropRepo) AddTask(accountAirdropID uint, task *model.Task) error {
	task.AccountAirdropID = accountAirdropID
	return r.DB.Create(task).Error
}

// AssignAirdrop creates an AccountAirdrop link and optionally creates tasks from a template.
func (r *AccountAirdropRepo) AssignAirdrop(accountID, airdropID uint, status, notes string, tasks []model.Task) (*model.AccountAirdrop, error) {
	// Check if already assigned
	var existing model.AccountAirdrop
	err := r.DB.Where("account_id = ? AND airdrop_id = ?", accountID, airdropID).First(&existing).Error
	if err == nil {
		return &existing, nil // Already assigned
	}

	aa := &model.AccountAirdrop{
		AccountID: accountID,
		AirdropID: airdropID,
		Status:    status,
		Notes:     notes,
	}
	if aa.Status == "" {
		aa.Status = "active"
	}

	if err := r.DB.Create(aa).Error; err != nil {
		return nil, err
	}

	// Create tasks if provided
	for _, t := range tasks {
		t.AccountAirdropID = aa.ID
		r.DB.Create(&t)
	}

	// Reload with preloads
	r.DB.Preload("Airdrop").Preload("Tasks").First(aa, aa.ID)

	return aa, nil
}

// RemoveAirdrop removes the account-airdrop link and its tasks.
func (r *AccountAirdropRepo) RemoveAirdrop(accountID, airdropID uint) error {
	var aa model.AccountAirdrop
	if err := r.DB.Where("account_id = ? AND airdrop_id = ?", accountID, airdropID).First(&aa).Error; err != nil {
		return err
	}
	return r.DB.Delete(aa).Error
}
