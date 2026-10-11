package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ride struct {
	GormModel

	BusID     string `gorm:"column:busid;not null"`
	RouteID   string `gorm:"column:routeid;not null"`
	Longitude float64
	Latitude  float64
	Capacity  int

	Bus   *Bus   `gorm:"foreignKey:BusID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Route *Route `gorm:"foreignKey:RouteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (r *Ride) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return err
}
