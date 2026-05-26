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
	// Delete tasks for all account-airdrops in this account
	r.DB.Where("account_airdrop_id IN (SELECT id FROM account_airdrops WHERE account_id = ?)", id).Delete(&model.Task{})

	// Delete account-airdrops
	r.DB.Where("account_id = ?", id).Delete(&model.AccountAirdrop{})

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

// CloneAccount duplicates an account with its account-airdrops and tasks (not wallets).
func (r *AccountRepo) CloneAccount(sourceID, userID uint, name, color string) (*model.Account, error) {
	// Find source account
	var source model.Account
	if err := r.DB.Where("id = ? AND user_id = ?", sourceID, userID).First(&source).Error; err != nil {
		return nil, err
	}

	// Create new account
	newAccount := &model.Account{
		UserID: userID,
		Name:   name,
		Color:  color,
		Notes:  "Cloned from: " + source.Name,
	}
	if err := r.DB.Create(newAccount).Error; err != nil {
		return nil, err
	}

	// Clone account-airdrops with their tasks
	var sourceAAs []model.AccountAirdrop
	r.DB.Where("account_id = ?", sourceID).Preload("Tasks").Find(&sourceAAs)

	for _, sa := range sourceAAs {
		newAA := model.AccountAirdrop{
			AccountID: newAccount.ID,
			AirdropID: sa.AirdropID,
			Status:    sa.Status,
			Notes:     sa.Notes,
		}
		if err := r.DB.Create(&newAA).Error; err != nil {
			continue
		}

		// Clone tasks
		for _, t := range sa.Tasks {
			newTask := model.Task{
				AccountAirdropID: newAA.ID,
				Description:      t.Description,
				Frequency:        t.Frequency,
				IsCompleted:      false,
				CompletedAt:      nil,
				GasSpent:         0,
				TxHash:           "",
			}
			r.DB.Create(&newTask)
		}
	}

	// Reload with preloads
	r.DB.Preload("Wallets").
		Preload("AccountAirdrops").
		Preload("AccountAirdrops.Airdrop").
		Preload("AccountAirdrops.Tasks").
		First(newAccount, newAccount.ID)

	return newAccount, nil
}
