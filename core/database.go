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
