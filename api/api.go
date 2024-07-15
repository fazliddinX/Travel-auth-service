package api

import (
	"Auth-service/api/handler"
	"github.com/gin-gonic/gin"
)

func Router(handler *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/api/register", handler.Register)
}
