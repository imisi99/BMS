package schemas

type UpdateRoute struct {
	RouteID string `json:"routeid" binding:"required"`
}
