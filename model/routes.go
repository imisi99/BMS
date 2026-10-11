package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stop struct {
	GormModel
	Name      string  `gorm:"not null"`
	Longitude float64 `gorm:"not null"`
	Latitude  float64 `gorm:"not null"`
}

type Route struct {
	GormModel
	Name string `gorm:"not null"`

	// Relatons
	Stops []*RouteStop `gorm:"many2many:route_stops"`
}

type RouteStop struct {
	RouteID  string `gorm:"primaryKey"`
	Position int    `gorm:"primaryKey"`
	StopID   string `gorm:"column:stopid; not null"`
	Stop     Stop
}

func (s *Stop) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return err
}

func (r *Route) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return err
}
