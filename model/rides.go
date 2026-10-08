package model

type Rides struct {
	GormModel

	BusID     string
	RouteID   string
	longitude string
	latitude  string
	Capacity  int
}
