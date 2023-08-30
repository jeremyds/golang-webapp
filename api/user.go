package api

import (
	"test_jump/database"

	"github.com/gin-gonic/gin"
)

type ListUsersResponse struct {
	ID        uint    `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balance   float64 `json:"balance"`
}

func ListUsers(ctx *gin.Context) {
	var data []ListUsersResponse
	database.DB.Model(&database.User{}).Find(&data)

	ctx.JSON(200, data)
}
