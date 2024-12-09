package controller

import (
	"golang-chap46/helper"
	"golang-chap46/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductController struct {
	service service.Service
	log     *zap.Logger
}

func NewProductController(service service.Service, log *zap.Logger) *ProductController {
	return &ProductController{service, log}
}

func (c *ProductController) GetAllProductsController(ctx *gin.Context) {
	products, err := c.service.Product.GetAllProducts()
	if err != nil {
		c.log.Error("Error fetching products", zap.Error(err))
		helper.ResponseError(ctx, err.Error(), "Error fetching products", http.StatusInternalServerError)
		return
	}

	c.log.Info("Products successfully retrieved", zap.Any("products", products))
	helper.ResponseOK(ctx, products, "Products successfully retrieved", http.StatusOK)
}
