package controller

import (
	"golang-chap46/config"
	"golang-chap46/database"
	"golang-chap46/service"

	"go.uber.org/zap"
)

type Controller struct {
	User    UserController
	Product ProductController
}

func NewController(service service.Service, log *zap.Logger, cacher database.Cacher, config config.Configuration) *Controller {
	return &Controller{
		User:    *NewUserController(service, log, cacher, config),
		Product: *NewProductController(service, log),
	}

}
