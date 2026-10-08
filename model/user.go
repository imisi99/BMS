// Package model
package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	GormModel

	FullName      string `gorm:"column:fullname;not null"`
	UserName      string `gorm:"column:username;not null;index"`
	Email         string `gorm:"unique;not null"`
	Password      string `gorm:"not null"`
	Notifications bool
}

type GormModel struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
