package service

import (
	"test_jump/config"
	"test_jump/database"
	"test_jump/server"
)

func InitService() {
	config.InitConfig()
	database.InitDb()
	server.InitServer()
}

func RunService() {
	server.RunServer()
}
