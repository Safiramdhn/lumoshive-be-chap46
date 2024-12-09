package infra

import (
	"golang-chap46/config"
	"golang-chap46/controller"
	"golang-chap46/database"
	"golang-chap46/helper"
	"golang-chap46/middleware"
	"golang-chap46/repository"
	"golang-chap46/service"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Cfg        config.Configuration
	DB         *gorm.DB
	Ctl        controller.Controller
	Log        *zap.Logger
	Cacher     database.Cacher
	Middleware *middleware.AuthMiddleware
}

func NewServiceContext() (*ServiceContext, error) {

	handlerError := func(err error) (*ServiceContext, error) {
		return nil, err
	}

	// instance config
	config, err := config.ReadConfig()
	if err != nil {
		handlerError(err)
	}

	// instance looger
	log, err := helper.InitZapLogger()
	if err != nil {
		handlerError(err)
	}

	// instance database
	db, err := database.InitDB(config)
	if err != nil {
		handlerError(err)
	}

	rdb := database.NewCacher(config, 60*60)

	// middleware := middleware.NewMiddleware(log, rdb)

	// instance repository
	repository := repository.NewRepository(db, log)

	// instance service
	service := service.NewService(*repository)

	// instance controller
	Ctl := controller.NewController(*service, log, rdb, config)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(log, rdb, config.JwtSecret)

	// Return service context
	return &ServiceContext{
		Cfg:        config,
		DB:         db,
		Ctl:        *Ctl,
		Log:        log,
		Cacher:     rdb,
		Middleware: authMiddleware,
	}, nil
}
