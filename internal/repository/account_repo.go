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
		Preload("AccountAirdrops").
		Preload("AccountAirdrops.Airdrop").
		Preload("AccountAirdrops.Tasks").
		Preload("AccountAirdrops.Tasks.Category").
		Order("created_at ASC").
		Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepo) FindByID(id, userID uint) (*model.Account, error) {
	var account model.Account
	err := r.DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("Wallets").
		Preload("AccountAirdrops").
		Preload("AccountAirdrops.Airdrop").
		Preload("AccountAirdrops.Tasks").
		Preload("AccountAirdrops.Tasks.Category").
		First(&account).Error
	return &account, err
}

func (r *AccountRepo) Update(account *model.Account) error {
	return r.DB.Save(account).Error
}

func (r *AccountRepo) Delete(id, userID uint) error {
	var aaCount int64
	r.DB.Model(&model.AccountAirdrop{}).Where("account_id = ?", id).Count(&aaCount)

	var walletCount int64
	r.DB.Model(&model.Wallet{}).Where("account_id = ?", id).Count(&walletCount)

	if aaCount > 0 || walletCount > 0 {
		return gorm.ErrRecordNotFound
	}

	return r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Account{}).Error
}

func (r *AccountRepo) DeleteCascade(id, userID uint) error {
	r.DB.Where("account_airdrop_id IN (SELECT id FROM account_airdrops WHERE account_id = ?)", id).Delete(&model.Task{})
	r.DB.Where("account_id = ?", id).Delete(&model.AccountAirdrop{})
	r.DB.Where("account_id = ?", id).Delete(&model.Wallet{})
	return r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Account{}).Error
}

