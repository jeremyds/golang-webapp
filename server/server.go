package server

import (
	"test_jump/api"
	"test_jump/config"

	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var Ws *gin.Engine = nil

func InitServer() {
	gin.SetMode(config.Config.WsConfig.Mode)

	r := gin.Default()

	r.GET("/ready", api.ReadinessProbe)

	r.GET("/users", api.ListUsers)
	r.POST("/invoice", api.AddInvoice)
	r.POST("/transaction", api.CreateTransaction)

	Ws = r
}

func RunServer() {
	addr := fmt.Sprintf("%s:%d", config.Config.WsConfig.Host, config.Config.WsConfig.Port)
	if err := Ws.Run(addr); err != nil {
		fmt.Printf("Failed to run the server: %v", err)
		os.Exit(1)
	}
}
