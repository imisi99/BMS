// Package handlers
package handlers

import (
	"log"
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

	hashedPass, err := s.HashPassword(payload.Password)
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

	token, err := s.GenerateJWT(user.ID, "login", JWTExpiry)
	if err != nil {
		log.Println("Failed to generate token -> ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create access token."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "User created successfully.", "token": token})
}

func (s *Service) VerifyUser(ctx *gin.Context) {
	var payload schemas.VerifyUser

	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Failed to parse the payload."})
		return
	}

	var token string

	token, err := s.Authorization(payload)
	if err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Service) AddRoute(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	if err := s.VerifyAuth(uid, purpose); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var payload schemas.UpdateRoute

	if err := s.DB.SaveRoute(uid, payload.RouteID); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Route saved successfully."})
}

func (s *Service) RemoveRoute(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	if err := s.VerifyAuth(uid, purpose); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var payload schemas.UpdateRoute

	if err := s.DB.RemoveRoute(uid, payload.RouteID); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Route removed successfully."})
}
