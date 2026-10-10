package handlers

import (
	"net/http"

	"BMS/core"
	"BMS/model"
	"BMS/schemas"

	"github.com/gin-gonic/gin"
)

func (s *Service) CreateStop(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	admin := ctx.GetBool("admin")
	if err := s.VerifyAdminAuth(uid, purpose, admin); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var payload schemas.CreateStop
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Failed to parse payload"})
		return
	}

	stop := model.Stop{
		Name:      payload.Name,
		Longitude: payload.Longitude,
		Latitude:  payload.Latitude,
	}

	if err := s.DB.AddStop(&stop); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"msg": "Stop Created successfully."})
}

func (s *Service) CreateRoute(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	admin := ctx.GetBool("admin")
	if err := s.VerifyAdminAuth(uid, purpose, admin); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var payload schemas.CreateRoute
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Failed to parse payload"})
		return
	}

	if len(payload.Stops) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "A route must contain at least one stop."})
		return
	}

	route := model.Route{
		Name: payload.Name,
	}

	if err := s.DB.AddRoute(&route, payload.Stops); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"msg": "Route Created successfully."})
}
