package repo

import (
	"github.com/nczz/phishguard/internal/model"
	"gorm.io/gorm"
)

type UserRepo struct{ DB *gorm.DB }

func (r *UserRepo) Create(u *model.User) error {
	return r.DB.Create(u).Error
}

func (r *UserRepo) FindByID(id int64) (*model.User, error) {
	var u model.User
	err := r.DB.First(&u, id).Error
	return &u, err
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.DB.Where("email = ?", email).First(&u).Error
	return &u, err
}

func (r *UserRepo) FindAllByTenant(tenantID int64) ([]model.User, error) {
	var users []model.User
	err := r.DB.Where("tenant_id = ?", tenantID).Find(&users).Error
	return users, err
}

func (r *UserRepo) Update(u *model.User) error {
	return r.DB.Save(u).Error
}

func (r *UserRepo) Delete(id int64) error {
	return r.DB.Delete(&model.User{}, id).Error
}
