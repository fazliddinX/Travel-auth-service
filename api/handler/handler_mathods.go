package handler

import (
	pb "Auth-service/genproto/auth_service"
	"Auth-service/token"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	loginReq := pb.LoginRequest{}
	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error("error in ShouldBindJSON", "error", err)
		return
	}
	res, err := h.Server.Login(context.Background(), &loginReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetById(c *gin.Context) {
	tk := c.GetHeader("Authorization")
	clams, err := token.ExtractClaimAcces(tk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Server.GetProfile(context.Background(), &pb.Id{Id: clams.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) Update(c *gin.Context) {
	profile := pb.UpdateUser{}
	err := c.ShouldBindJSON(&profile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error("error in ShouldBindJSON", "error", err)
		return
	}
	user, err := h.Server.UpdateProfile(context.Background(), &profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) Get(c *gin.Context) {
	a := c.Query("limit")
	b := c.Query("offset")
	d := c.Query("countries_visited")

	limit, err := strconv.Atoi(a)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error("error in Get", "error", err)
		return
	}
	offset, err := strconv.Atoi(b)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error("error in Get", "error", err)
		return
	}
	countriesVisited, err := strconv.Atoi(d)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error("error in Get", "error", err)
		return
	}

	userName := c.Query("user_name")
	fullName := c.Query("full_name")

	filter := pb.FilterGet{
		Limit:            int64(limit),
		Offset:           int64(offset),
		UserName:         userName,
		FullName:         fullName,
		CountriesVisited: int64(countriesVisited),
	}
	users, err := h.Server.GetUsers(context.Background(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) Delete(c *gin.Context) {
	tk := c.GetHeader("Authorization")
	clams, err := token.ExtractClaimAcces(tk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	seccess, err := h.Server.Delete(context.Background(), &pb.Id{Id: clams.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seccess)
}

func (h *Handler) PasswordRecovery(c *gin.Context) {
	email := pb.Email{}
	err := c.ShouldBindJSON(&email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error("error in ShouldBindJSON", "error", err)
		return
	}
	success, err := h.Server.PasswordRecovery(context.Background(), &email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, success)
}

func (h *Handler) NewAccessToken(c *gin.Context) {
	tk := c.GetHeader("Authorization")
	token, err := h.Server.TokenRenewal(context.Background(), &pb.RefreshToken{RefreshToken: tk})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access Token": token})
}

func (h *Handler) Activity(c *gin.Context) {
	id := c.Query("id")
	user, err := h.Server.ActivityProfile(context.Background(), &pb.Id{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) Follow(c *gin.Context) {
	tk := c.GetHeader("Authorization")
	clams, err := token.ExtractClaimAcces(tk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id := c.Query("id")
	fw, err := h.Server.Follow(context.Background(), &pb.FollowRequest{FollowerId: clams.Id, FollowingId: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fw)
}

func (h *Handler) Unfollow(c *gin.Context) {
	var id string
	err := c.ShouldBindJSON(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error("error in ShouldBindJSON", "error", err)
		return
	}
	succes, err := h.Server.Unfollow(context.Background(), &pb.FollowingId{FollowingId: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, succes)
}

func (h *Handler) GetFollowers(c *gin.Context) {
	id := c.Query("id")

	a := c.Query("limit")
	b := c.Query("offset")

	limit, err := strconv.Atoi(a)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error("error in Get", "error", err)
		return
	}
	offset, err := strconv.Atoi(b)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error("error in Get", "error", err)
		return
	}
	userName := c.Query("user_name")
	email := c.Query("email")

	filter := pb.FilterFollowers{
		Limit:    int64(limit),
		Offset:   int64(offset),
		UserName: userName,
		Email:    email,
		Id:       id,
	}
	followers, err := h.Server.GetFollowers(context.Background(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"followers": followers, "limit": int64(limit), "offset": int64(offset)})
}
