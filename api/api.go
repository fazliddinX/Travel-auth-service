package api

import (
	"Auth-service/api/handler"
	"Auth-service/api/middleware"
	"github.com/gin-gonic/gin"
)

func Router(handler *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
	router.POST("/auth/refresh", middleware.MiddleWareRefresh(), handler.NewAccessToken)

	r := router.Group("/api")
	r.Use(middleware.MiddleWareAcces())

	r.GET("/user/profile", handler.GetById)
	r.PUT("/user/profile", handler.Update)
	r.GET("/user/", handler.Get)
	r.DELETE("/user/:id", handler.Delete)
	r.POST("/auth/reset-password", handler.PasswordRecovery) //mashi mazgi
	r.POST("/auth/login", handler.Login)
	r.GET("/user/:id/activity", handler.Activity)
	r.POST("/user/:id/follow", handler.Follow)
	r.PUT("user/:id/unfollow", handler.Unfollow)
	r.GET("user/:id/followers", handler.GetFollowers)

	return router
}
