package main

import (
	"github.com/sajalmia381/store-api/src/api"
	"github.com/sajalmia381/store-api/src/config"
	"github.com/sajalmia381/store-api/src/dependency"
	"github.com/sajalmia381/store-api/src/v1/db"
	"github.com/sajalmia381/store-api/src/v1/dtos"
)

func main() {
	server := config.New()
	db.GetDmManager()

	go initDefaultUser()

	api.Routes(server)

	server.Logger.Fatal(server.Start(":" + config.ServerPort))
}

func initDefaultUser() {
	userService := dependency.GetUserService()
	num := 1816785381
	payload := dtos.UserRegisterDTO{
		Name:     "Anonymous User",
		Email:    "anonymous@gmail.com",
		Password: "simple_password",
		Number:   &num,
	}
	userService.Store(payload)
}
