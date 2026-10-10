package schemas

type CreateRoute struct {
	Name  string   `json:"name" binding:"required"`
	Stops []string `json:"stops" binding:"required"`
}

type CreateStop struct {
	Name      string  `json:"name" binding:"required"`
	Longitude float64 `json:"long" binding:"required"`
	Latitude  float64 `json:"lat" binding:"required"`
}

type UpdateRoute struct {
	RouteID string `json:"routeid" binding:"required"`
}
