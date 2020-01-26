package users

import (
	"github.com/coby9241/frontend-service/internal/models/users"
	"github.com/jinzhu/gorm"
)

// Repository is
type Repository interface {
	GetUserByUID(uid string) (*users.User, error)
}

// UserRepositoryImpl is
type UserRepositoryImpl struct {
	DB *gorm.DB
}

// NewUserRepositoryImpl is
func NewUserRepositoryImpl(storage *gorm.DB) Repository {
	return &UserRepositoryImpl{
		DB: storage,
	}
}

// GetUserByUID is
func (r *UserRepositoryImpl) GetUserByUID(uid string) (*users.User, error) {
	i := users.User{}
	err := r.DB.Where("UID = ?", uid).First(&i).Error
	return &i, err
}
