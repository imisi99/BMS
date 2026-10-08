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
