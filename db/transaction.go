package db

import (
	"context"
	"gorm.io/gorm"
)

func WithTransaction(fn func(tx *gorm.DB) error) error {
	context.Background()

	tx := database.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
