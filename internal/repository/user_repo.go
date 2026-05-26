package repository

import (
	"github.com/bsyrlhabibi/airdrop/internal/model"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, id).Error
	return &user, err
}
