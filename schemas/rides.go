package schemas

type CreateRide struct {
	BusID   string `json:"busid" binding:"required"`
	RouteID string `json:"routeid" binding:"required"`
	RiderID string `json:"riderid" binding:"required"`
}

type UpdateRide struct {
	RideID    string  `json:"rideid" binding:"required"`
	Capacity  int     `json:"capacity" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
}

type GetRide struct {
	RideID string `json:"rideid" binding:"required"`
}

type FindRideByRoute struct {
	RouteID string `json:"routeid" binding:"required"`
}

type FindRideByStops struct {
	Start string `json:"start" binding:"required"`
	Stop  string `json:"stop" binding:"required"`
}

type FindRidesRespsonse struct {
	RouteID   string
	Available int
}
