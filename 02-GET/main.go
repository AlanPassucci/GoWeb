/*
Crear una ruta /ping que debe respondernos con un string que contenga pong con el status 200 OK.
Crear una ruta /products que nos devuelva la lista de todos los productos en la slice.
Crear una ruta /products/:id que nos devuelva un producto por su id.
Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo precio sean mayor a un valor priceGt.
*/

package main

import (
	"fmt"
	"goweb/02-GET/filemanager"
	"goweb/02-GET/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var productsList []product.Product

func InitialiceProductsList() {
	fileManager := filemanager.NewFileManager("./products.json")

	data, err := fileManager.ReadFile()
	if err != nil {
		panic(err)
	}

	productsList, err = product.DecerealizeJSON(data)
	if err != nil {
		panic(err)
	}
}

func Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong!")
}

func GetAllProducts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, productsList)
}

func GetProductByID(ctx *gin.Context) {
	productId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	for _, product := range productsList {
		if product.ID == productId {
			ctx.JSON(http.StatusOK, product)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

func GetProductsAbovePrice(ctx *gin.Context) {
	productPriceGtAsString := ctx.Query("priceGt")
	if productPriceGtAsString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid or non-existent query parameter"})
		return
	}

	productPriceGt, err := strconv.ParseFloat(productPriceGtAsString, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	productsWithPriceGt := []product.Product{}
	for _, product := range productsList {
		if product.Price > productPriceGt {
			productsWithPriceGt = append(productsWithPriceGt, product)
		}
	}

	if len(productsWithPriceGt) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("There are not any products with price greater than %.2f", productPriceGt)})
		return
	}

	ctx.JSON(http.StatusOK, productsWithPriceGt)
}

func main() {
	router := gin.Default()

	InitialiceProductsList()

	router.GET("/ping", Ping)

	products := router.Group("/products")
	products.GET("/", GetAllProducts)
	products.GET("/:id", GetProductByID)
	products.GET("/search", GetProductsAbovePrice)

	router.Run()
}
