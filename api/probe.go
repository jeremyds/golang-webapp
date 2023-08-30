package api

import "github.com/gin-gonic/gin"

func ReadinessProbe(ctx *gin.Context) {
	ctx.Status(200)
}
