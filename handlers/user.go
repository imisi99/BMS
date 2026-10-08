// Package handlers
package handlers

import (
	"net/http"

	"BMS/core"
	"BMS/model"
	"BMS/schemas"

	"github.com/gin-gonic/gin"
)

func (s *Service) AddUser(ctx *gin.Context) {
	var payload schemas.AddUser

	if err := ctx.BindJSON(payload); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Failed to parse the payload."})
		return
	}

	hashedPass, err := HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to hash password."})
		return
	}

	user := model.User{
		FullName: payload.FullName,
		UserName: payload.UserName,
		Email:    payload.Email,
		Password: hashedPass,
	}

	if err := s.DB.AddUser(&user); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "User created successfully."})
}

func (s *Service) VerifyUser(ctx *gin.Context) {
	var payload schemas.VerifyUser

	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Failed to parse the payload."})
		return
	}

	var token string
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
