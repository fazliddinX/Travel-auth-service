package handler

import (
	"Auth-service/genproto/auth_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) Register(c *gin.Context) {
	user := auth_service.RegisterUser{}
	err := c.ShouldBindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	h.User.Create(&user)

}
