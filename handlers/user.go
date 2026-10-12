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

	token, err := s.GenerateJWT(user.ID, "login", false, JWTExpiry)
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

	token, err := s.Authorization(payload, false)
	if err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Service) SaveRoute(ctx *gin.Context) {
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

func (s *Service) SaveStop(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	if err := s.VerifyAuth(uid, purpose); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var payload schemas.UpdateStop

	if err := s.DB.SaveStop(uid, payload.StopID); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Stop saved successfully."})
}

func (s *Service) GetRoutes(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	if err := s.VerifyAuth(uid, purpose); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var user model.User
	if err := s.DB.GetRoutes(uid, &user); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var routes []schemas.ViewSavedRoutes
	for _, route := range user.Routes {
		routes = append(routes, schemas.ViewSavedRoutes{
			RouteID: route.ID,
			Name:    route.Name,
			Start:   route.Stops[0].Stop.Name,
			Stop:    route.Stops[len(route.Stops)-1].Stop.Name,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"routes": routes})
}

func (s *Service) GetStops(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	if err := s.VerifyAuth(uid, purpose); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var user model.User
	if err := s.DB.GetStops(uid, &user); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var stops []schemas.ViewSavedStops
	for _, stop := range user.Stops {
		stops = append(stops, schemas.ViewSavedStops{
			StopID: stop.ID,
			Name:   stop.Name,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"stops": stops})
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
