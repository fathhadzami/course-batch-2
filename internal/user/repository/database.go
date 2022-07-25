package repository

import (
	"context"
	"course/internal/domain"

	"gorm.io/gorm"
)

type DatabaseRepo struct {
	db *gorm.DB
}

func NewDatabaseRepo(db *gorm.DB) *DatabaseRepo {
	return &DatabaseRepo{db}
}

func (dr DatabaseRepo) IsUserExist(ctx context.Context, userID int) bool {
	var user domain.User
	err := dr.db.WithContext(ctx).Where("id = ?", userID).Take(&user).Error
	return err != nil && user.ID > 0
}

func (dr DatabaseRepo) GetUserByID(ctx context.Context, userID int) bool {
	return true
}
