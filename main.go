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
	_, err := userService.FindByEmail("anonymous@gmail.com")
	if err == nil {
		return
	}
	var num uint = 1234567891
	payload := dtos.UserRegisterDTO{
		Name:     "Anonymous User",
		Email:    "anonymous@gmail.com",
		Password: "simple_password",
		Number:   &num,
	}
	userService.Store(payload)
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
			num = 1234567891
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
