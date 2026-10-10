package model

type Stop struct {
	GormModel
	Name      string
	Longitude float64
	Latitude  float64
}

type Route struct {
	GormModel
	Name  string
	Stops []RouteStop
}

type RouteStop struct {
	RouteID  string
	Position int
	StopID   string
	Stop     Stop
}
