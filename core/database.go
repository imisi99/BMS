package core

import (
	"errors"
	"log"
	"net/http"

	"BMS/model"

	"gorm.io/gorm"
)

type DB interface {
	AddUser(user *model.User) error
	VerifyUser(user *model.User, username string) error
	FetchUser(uid string, user *model.User) error
	AddStop(stop *model.Stop) error
	AddRoute(route *model.Route, stops []string) error
	GetRoutes(uid string, user *model.User) error
	GetStops(uid string, user *model.User) error
	FetchRoute(rid string, route *model.Route) error
	SaveRoute(uid, rid string) error
	SaveStop(uid, sid string) error
	RemoveRoute(uid, rid string) error
	AddRide(ride *model.Ride) error
	UpdateRide(rid string, cap int, long, lat float64) error
	GetRide(rid string, ride *model.Ride) error
	FindRideByRoute(rid string, rides *[]model.Ride) error
	FindRideByStops(start, stop string, rides *[]model.Ride) error
}

type GormDB struct {
	DB *gorm.DB
}

func (db *GormDB) AddUser(user *model.User) error {
	if err := db.DB.Create(user).Error; err != nil {
		log.Println("[GormDB] Failed to create user -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "Failed to create new user."}
	}
	return nil
}

func (db *GormDB) VerifyUser(user *model.User, username string) error {
	if err := db.DB.Where("username = ? OR email = ?", username).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &CustomError{Code: http.StatusNotFound, Message: "Invalid username or password"}
		}
		log.Println("[GormDB] Failed to verify user -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "server error."}
	}
	return nil
}

func (db *GormDB) FetchUser(uid string, user *model.User) error {
	if err := db.DB.Where("id = ?", uid).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &CustomError{Code: http.StatusNotFound, Message: "User not found."}
		}
		log.Println("[GormDB] Failed to fetch user -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "server error."}
	}
	return nil
}

func (db *GormDB) AddStop(stop *model.Stop) error {
	if err := db.DB.Create(stop).Error; err != nil {
		log.Println("[GormDB] Failed to create stop -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "Failed to create new stop"}
	}
	return nil
}

func (db *GormDB) AddRoute(route *model.Route, stops []string) error {
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Stops").Create(route).Error; err != nil {
			return err
		}

		links := make([]model.RouteStop, len(stops))
		for i, sid := range stops {
			links[i] = model.RouteStop{
				RouteID:  route.ID,
				Position: i + 1,
				StopID:   sid,
			}
		}

		if err := tx.Omit("Stop").Create(&links).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return &CustomError{Code: http.StatusBadRequest, Message: "One or more stops do not exist"}
		}
		log.Println("[GormDB] Failed to create route -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "Failed to create new route."}

	}
	return nil
}

func (db *GormDB) GetRoutes(uid string, user *model.User) error {
	if err := db.DB.Preload("Routes.Stops.Stop").Where("id = ?", uid).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &CustomError{Code: http.StatusNotFound, Message: "User not found."}
		}
		log.Println("[GormDB] Failed to fetch routes -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "server error."}
	}

	return nil
}

func (db *GormDB) GetStops(uid string, user *model.User) error {
	if err := db.DB.Preload("Stop").Where("id = ?", uid).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &CustomError{Code: http.StatusNotFound, Message: "User not found."}
		}
		log.Println("[GormDB] Failed to fetch stops -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "server error."}
	}
	return nil
}

func (db *GormDB) FetchRoute(rid string, route *model.Route) error {
	if err := db.DB.Where("id = ?", rid).First(route).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &CustomError{Code: http.StatusNotFound, Message: "Route not found."}
		}
		log.Println("[GormDB] Failed to fetch route -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "server error."}
	}
	return nil
}

func (db *GormDB) SaveRoute(uid, rid string) error {
	user := model.User{GormModel: model.GormModel{ID: uid}}
	route := model.Route{GormModel: model.GormModel{ID: rid}}
	if err := db.DB.Model(&user).Association("Routes").Append(&route); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return &CustomError{Code: http.StatusNotFound, Message: "User or route not found."}
		}
		log.Println("[GormDB] Failed to save route -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "server error."}
	}

	return nil
}

func (db *GormDB) SaveStop(uid, sid string) error {
	user := model.User{GormModel: model.GormModel{ID: uid}}
	stop := model.Stop{GormModel: model.GormModel{ID: sid}}
	if err := db.DB.Model(&user).Association("Stops").Append(&stop); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return &CustomError{Code: http.StatusNotFound, Message: "User or stop not found."}
		}
		log.Println("[GormDB] Failed to save stop -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "server error."}

	}
	return nil
}

func (db *GormDB) RemoveRoute(uid, rid string) error {
	user := model.User{GormModel: model.GormModel{ID: uid}}
	route := model.User{GormModel: model.GormModel{ID: rid}}
	if err := db.DB.Model(&user).Association("Routes").Delete(&route); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return &CustomError{Code: http.StatusNotFound, Message: "User or route not found."}
		}
		log.Println("[GormDB] Failed to save route -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "server error."}
	}
	return nil
}

func (db *GormDB) AddRide(ride *model.Ride) error {
	if err := db.DB.Create(ride).Error; err != nil {
		log.Println("[GormDB] Failed to create stop -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "Failed to create new ride"}
	}
	return nil
}

func (db *GormDB) UpdateRide(rid string, cap int, long, lat float64) error {
	if res := db.DB.Model(&model.Ride{}).Where("id = ?", rid).Updates(map[string]any{
		"latitude": lat, "longitude": long, "cap": cap,
	}); res.RowsAffected == 0 || res.Error != nil {
		if res.RowsAffected == 0 {
			return &CustomError{Code: http.StatusNotFound, Message: "Ride not found."}
		}
		log.Println("[GormDB] Failed to update ride -> ", res.Error.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "server error."}
	}
	return nil
}

func (db *GormDB) GetRide(rid string, ride *model.Ride) error {
	if err := db.DB.Preload("Bus").Where("id = ?", rid).First(ride).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &CustomError{Code: http.StatusNotFound, Message: "Ride not found."}
		}
		log.Println("[GormDB] Failed to get ride -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "server error."}
	}
	return nil
}

func (db *GormDB) FindRideByRoute(rid string, rides *[]model.Ride) error {
	if err := db.DB.Preload("Bus").Where("routeid = ?", rid).Find(rides).Error; err != nil {
		log.Println("[GormDB] Failed to find rides -> ", err.Error())
		return &CustomError{Code: http.StatusInternalServerError, Message: "Failed to find a ride"}
	}
	return nil
}

func (db *GormDB) FindRideByStops(start, stop string, rides *[]model.Ride) error {
	return nil
}
