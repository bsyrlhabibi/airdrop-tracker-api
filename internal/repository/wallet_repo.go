package repository

import (
	"github.com/bsyrlhabibi/airdrop/internal/model"
	"gorm.io/gorm"
)

type WalletRepo struct {
	DB *gorm.DB
}

func NewWalletRepo(db *gorm.DB) *WalletRepo {
	return &WalletRepo{DB: db}
}

func (r *WalletRepo) Create(wallet *model.Wallet) error {
	return r.DB.Create(wallet).Error
}

func (r *WalletRepo) FindByUser(userID uint) ([]model.Wallet, error) {
	var wallets []model.Wallet
	err := r.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&wallets).Error
	return wallets, err
}

func (r *WalletRepo) FindByID(id, userID uint) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.DB.Where("id = ? AND user_id = ?", id, userID).First(&wallet).Error
	return &wallet, err
}

func (r *WalletRepo) Delete(id, userID uint) error {
	return r.DB.Where("user_id = ?", userID).Delete(&model.Wallet{}, id).Error
}
