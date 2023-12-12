/*
Crear una ruta /ping que debe respondernos con un string que contenga pong con el status 200 OK.
Crear una ruta /products que nos devuelva la lista de todos los productos en la slice.
Crear una ruta /products/:id que nos devuelva un producto por su id.
Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo precio sean mayor a un valor priceGt.
*/

package main

import (
	"goweb/03-POST/internal/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong!")
}

func main() {
	router := gin.Default()
	hd := handler.NewDefaultProducts()
	hd.InitialiceProductsList()

	router.GET("/ping", Ping)

	products := router.Group("/products")
	products.GET("/", hd.GetAllProducts())
	products.GET("/:id", hd.GetProductByID())
	products.GET("/search", hd.GetProductsAbovePrice())
	products.POST("/", hd.CreateProduct())

	router.Run()
}
