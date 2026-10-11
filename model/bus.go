package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bus struct {
	GormModel

	Capacity int
}

func (b *Bus) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.NewString()
	}
	return err
}
