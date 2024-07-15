package handler

import (
	pb "Auth-service/genproto/auth_service"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) Register(c *gin.Context) {
	userReq := pb.RegisterUserReq{}
	err := c.ShouldBindJSON(&userReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error("error in ShouldBindJSON", "error", err)
		return
	}
	userRes, err := h.Server.Register(context.Background(), &userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userRes)
}

func (h *Handler) Login(c *gin.Context) {

}
