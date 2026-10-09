package schemas

type CreateRide struct {
	BusID   string `json:"busid" binding:"required"`
	RouteID string `json:"routeid" binding:"required"`
	RiderID string `json:"riderid" binding:"required"`
}
