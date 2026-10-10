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
	FetchRoute(rid string, route *model.Route) error
	SaveRoute(uid, rid string) error
	RemoveRoute(uid, rid string) error
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
	route := model.User{GormModel: model.GormModel{ID: rid}}
	if err := db.DB.Model(&user).Association("Routes").Append(&route); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return &CustomError{Code: http.StatusNotFound, Message: "User or route not found."}
		}
		log.Println("[GormDB] Failed to save route -> ", err.Error())
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
