package main

import (
	"fmt"
	"goweb/06-DotEnv/cmd/server/handler"
	"goweb/06-DotEnv/internal/product"
	"goweb/06-DotEnv/pkg/store"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err.Error())
	}
	fileManager := store.NewFileManager("../../products.json")
	repository := product.NewRepositorySlice(fileManager)
	service := product.NewDefaultService(repository)
	handler := handler.NewProductHandler(service)

	server := gin.Default()

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Pong!")
	})

	productsGroup := server.Group("/products")
	{
		productsGroup.GET("/", handler.GetAllProducts())
		productsGroup.GET("/:id", handler.GetProductByID())
		productsGroup.GET("/search", handler.GetAllProductsAbovePrice())
		productsGroup.POST("/", handler.CreateProduct())
		productsGroup.PUT("/:id", handler.UpdateProduct())
		productsGroup.PATCH("/:id", handler.PartialUpdateProduct())
		productsGroup.DELETE("/:id", handler.DeleteProduct())
	}

	server.Run()
}
