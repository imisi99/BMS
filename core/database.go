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
	FetchRoute(rid string, route *model.Route) error
	SaveRoute(uid, rid string) error
	RemoveRoute(uid, rid string) error
}

type GormDB struct {
	DB *gorm.DB
}

func (db *GormDB) AddUser(user *model.User) error {
	if err := db.DB.Create(user).Error; err != nil {
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
