package repository

import (
	"github.com/bsyrlhabibi/airdrop/internal/model"
	"gorm.io/gorm"
)

type AccountRepo struct {
	DB *gorm.DB
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{DB: db}
}

func (r *AccountRepo) Create(account *model.Account) error {
	return r.DB.Create(account).Error
}

func (r *AccountRepo) FindByUser(userID uint) ([]model.Account, error) {
	var accounts []model.Account
	err := r.DB.Where("user_id = ?", userID).
		Preload("Wallets").
		Preload("Airdrops").
		Order("created_at ASC").
		Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepo) FindByID(id, userID uint) (*model.Account, error) {
	var account model.Account
	err := r.DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("Wallets").
		Preload("Airdrops").
		First(&account).Error
	return &account, err
}

func (r *AccountRepo) Update(account *model.Account) error {
	return r.DB.Save(account).Error
}

func (r *AccountRepo) Delete(id, userID uint) error {
	// Check if account has airdrops or wallets
	var airdropCount int64
	r.DB.Model(&model.Airdrop{}).Where("account_id = ?", id).Count(&airdropCount)

	var walletCount int64
	r.DB.Model(&model.Wallet{}).Where("account_id = ?", id).Count(&walletCount)

	if airdropCount > 0 || walletCount > 0 {
		return gorm.ErrRecordNotFound // Will use custom error in handler
	}

	return r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Account{}).Error
}

func (r *AccountRepo) DeleteCascade(id, userID uint) error {
	// Delete tasks for all airdrops in this account
	r.DB.Where("airdrop_id IN (SELECT id FROM airdrops WHERE account_id = ?", id).Delete(&model.Task{})

	// Delete airdrops
	r.DB.Where("account_id = ?", id).Delete(&model.Airdrop{})

	// Delete wallets
	r.DB.Where("account_id = ?", id).Delete(&model.Wallet{})

	// Delete account
	return r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Account{}).Error
}

func (r *AccountRepo) FindOrCreateDefault(userID uint) (*model.Account, error) {
	var account model.Account
	err := r.DB.Where("user_id = ? AND name = ?", userID, "Default").First(&account).Error
	if err == gorm.ErrRecordNotFound {
		account = model.Account{
			UserID: userID,
			Name:   "Default",
			Color:  "#3B82F6",
			Notes:  "Auto-created default account",
		}
		err = r.DB.Create(&account).Error
	}
	return &account, err
}
