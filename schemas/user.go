// Package schemas
package schemas

type AddUser struct {
	FullName string `json:"fullname" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyUser struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ViewSavedRoutes struct {
	RouteID string
	Name    string
	Start   string
	Stop    string
}

type ViewSavedStops struct {
	StopID string
	Name   string
}
