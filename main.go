package main

import (
	"strconv"

	"github.com/sajalmia381/store-api/src/api"
	"github.com/sajalmia381/store-api/src/config"
	"github.com/sajalmia381/store-api/src/dependency"
	"github.com/sajalmia381/store-api/src/v1/db"
	"github.com/sajalmia381/store-api/src/v1/dtos"
)

func main() {
	server := config.New()
	db.GetDmManager()
	go intSuperAdmin()
	go initDefaultUser()

	api.Routes(server)

	server.Logger.Fatal(server.Start(":" + config.ServerPort))
}

func initDefaultUser() {
	userService := dependency.GetUserService()
	user, err := userService.FindByEmail(config.DefaultUserEmail)
	if err == nil {
		config.DefaultUserId = &user.ID
		return
	}
	num, err := strconv.Atoi(config.SuperAdminNumber)
	if err != nil {
		num = 1234567891
	}
	var menNum uint = uint(num)
	payload := dtos.UserRegisterDTO{
		Name:     config.DefaultUserName,
		Email:    config.DefaultUserEmail,
		Password: config.DefaultUserPassword,
		Number:   &menNum,
	}
	user, _ = userService.Store(payload)
	config.DefaultUserId = &user.ID
}

func intSuperAdmin() {
	if config.SuperAdminEmail != "" {
		userService := dependency.GetUserService()
		_, err := userService.FindByEmail(config.SuperAdminEmail)
		if err == nil {
			return
		}
		num, err := strconv.Atoi(config.SuperAdminNumber)
		if err != nil {
			num = 1234567830
		}
		var menNum uint = uint(num)
		payload := dtos.UserRegisterDTO{
			Name:     config.SuperAdminName,
			Email:    config.SuperAdminEmail,
			Password: config.SuperAdminPassword,
			Number:   &menNum,
		}
		userService.StoreSuperAdmin(payload)
	}
}
