package controller

import (
	"golang-chap46/config"
	"golang-chap46/database"
	"golang-chap46/helper"
	"golang-chap46/models"
	"golang-chap46/service"
	"golang-chap46/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	service service.Service
	log     *zap.Logger
	cacher  database.Cacher
	config  config.Configuration
}

func NewUserController(service service.Service, log *zap.Logger, cacher database.Cacher, config config.Configuration) *UserController {
	return &UserController{service, log, cacher, config}
}

func (c *UserController) Register(ctx *gin.Context) {
	req := models.User{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.log.Error("invalid request", zap.Error(err))
		helper.ResponseError(ctx, err.Error(), "invalid request", http.StatusBadRequest)
		return
	}

	if err := c.service.User.Register(req); err != nil {
		c.log.Error("failed to register user", zap.Error(err))
		helper.ResponseError(ctx, err.Error(), "failed to register user", http.StatusInternalServerError)
		return
	}

	c.log.Info("user registered successfully", zap.String("email", req.Email))
	helper.ResponseOK(ctx, nil, "user registered successfully", http.StatusCreated)
}

func (c *UserController) Login(ctx *gin.Context) {
	req := models.LoginRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.log.Error("invalid request", zap.Error(err))
		helper.ResponseError(ctx, err.Error(), "invalid request", http.StatusBadRequest)
		return
	}

	user, err := c.service.User.Login(req.Email, req.Password)
	if err != nil {
		c.log.Error("failed to login user", zap.Error(err))
		helper.ResponseError(ctx, err.Error(), "failed to login user", http.StatusUnauthorized)
		return
	}

	expTime := 12 * time.Hour
	token, err := utils.GenerateToken(user.ID, c.config.JwtSecret, expTime)
	if err != nil {
		c.log.Error("failed to generate token", zap.Error(err))
		helper.ResponseError(ctx, err.Error(), "failed to generate token", http.StatusInternalServerError)
		return
	}

	key := strconv.Itoa(user.ID)
	c.cacher.SaveToken(key, token)
	response := models.LoginResponse{
		Token: token,
	}
	c.log.Info("user logged in successfully", zap.String("email", req.Email))
	helper.ResponseOK(ctx, response, "user logged in successfully", http.StatusOK)
}
