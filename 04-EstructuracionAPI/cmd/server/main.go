package main

import (
	"goweb/04-EstructuracionAPI/cmd/server/handler"
	"goweb/04-EstructuracionAPI/internal/domain"
	"goweb/04-EstructuracionAPI/internal/product"
	"goweb/04-EstructuracionAPI/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	var products []domain.Product
	pkg.FillDB[[]domain.Product]("../../products.json", &products)

	repository := product.NewRepositorySlice(products, len(products))
	service := product.NewDefaultService(repository)
	handler := handler.NewProductHandler(service)

	server := gin.Default()

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Pong!")
	})

	productsGroup := server.Group("/products")
	productsGroup.GET("/", handler.GetAllProducts())
	productsGroup.GET("/:id", handler.GetProductByID())
	productsGroup.GET("/search", handler.GetAllProductsAbovePrice())
	productsGroup.POST("/", handler.CreateProduct())

	server.Run()
}
