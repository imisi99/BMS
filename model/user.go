// Package model
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	GormModel

	FullName      string `gorm:"column:fullname;not null"`
	UserName      string `gorm:"column:username;not null;index"`
	Email         string `gorm:"unique;not null"`
	Password      string `gorm:"not null"`
	Notifications bool

	// Relatons
	Routes []*Route `gorm:"many2many:user_routes"`
	Stops  []*Stop  `gorm:"many2many:user_stops"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}

	return err
}

type GormModel struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
