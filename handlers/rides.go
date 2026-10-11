package handlers

import (
	"net/http"

	"BMS/core"
	"BMS/model"
	"BMS/schemas"

	"github.com/gin-gonic/gin"
)

// TODO:
// A Tracking Endpoint
// The View Ride Also Fetches the Route

func (s *Service) CreateRide(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	admin := ctx.GetBool("admin")
	if err := s.VerifyAdminAuth(uid, purpose, admin); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var payload schemas.CreateRide
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Failed to parse payload."})
		return
	}

	ride := model.Ride{
		RouteID: payload.RouteID,
		BusID:   payload.BusID,
	}

	if err := s.DB.AddRide(&ride); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"msg": "Ride created successfully.", "id": ride.ID})
}

func (s *Service) UpdateRide(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	admin := ctx.GetBool("admin")
	if err := s.VerifyAdminAuth(uid, purpose, admin); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var payload schemas.UpdateRide
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Failed to parse payload."})
		return
	}

	if err := s.DB.UpdateRide(payload.RideID, payload.Capacity, payload.Longitude, payload.Latitude); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"msg": "Ride updated successfully."})
}

func (s Service) GetRide(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	if err := s.VerifyAuth(uid, purpose); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var payload schemas.GetRide
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Failed to parse payload."})
		return
	}

	var ride model.Ride
	if err := s.DB.GetRide(payload.RideID, &ride); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	response := schemas.FindRidesRespsonse{
		RouteID:   ride.RouteID,
		Available: (ride.Bus.Capacity - ride.Capacity),
	}

	ctx.JSON(http.StatusOK, gin.H{"ride": response})
}

func (s *Service) FindRideByRoute(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	if err := s.VerifyAuth(uid, purpose); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var payload schemas.FindRideByRoute
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Failed to parse payload."})
		return
	}

	var rides []model.Ride
	if err := s.DB.FindRideByRoute(payload.RouteID, &rides); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var response []schemas.FindRidesRespsonse
	for _, ride := range rides {
		response = append(response, schemas.FindRidesRespsonse{
			RouteID:   ride.ID,
			Available: (ride.Bus.Capacity - ride.Capacity),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"rides": response})
}

func (s *Service) FindRideByStops(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	purpose := ctx.GetString("purpose")
	if err := s.VerifyAuth(uid, purpose); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var payload schemas.FindRideByStops
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Failed to parse payload."})
		return
	}

	var rides []model.Ride
	if err := s.DB.FindRideByStops(payload.Start, payload.Stop, &rides); err != nil {
		cm := err.(*core.CustomError)
		ctx.JSON(cm.Code, gin.H{"msg": cm.Message})
		return
	}

	var response []schemas.FindRidesRespsonse
	for _, ride := range rides {
		response = append(response, schemas.FindRidesRespsonse{
			RouteID:   ride.ID,
			Available: (ride.Bus.Capacity - ride.Capacity),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"rides": response})
}
