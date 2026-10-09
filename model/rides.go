package model

type Rides struct {
	GormModel

	BusID     string
	RouteID   string
	Longitude string
	Latitude  string
	Capacity  int
}
