package model

type Route struct {
	GormModel

	Locations []string
}
