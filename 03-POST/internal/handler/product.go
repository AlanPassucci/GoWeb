package handler

import (
	"fmt"
	"goweb/03-POST/internal/filemanager"
	"goweb/03-POST/internal/product"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DefaultProducts struct {
	db     []product.Product
	lastID int
}

type BodyRequestCreate struct {
	Name        string
	Quantity    int
	CodeValue   string
	IsPublished bool
	Expiration  string
	Price       float64
}

type ProductJSON struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func NewDefaultProducts() *DefaultProducts {
	return &DefaultProducts{
		db:     make([]product.Product, 0),
		lastID: 500,
	}
}

func (d *DefaultProducts) InitialiceProductsList() {
	fileManager := filemanager.NewFileManager("./products.json")

	err := fileManager.ReadFile()
	if err != nil {
		panic(err)
	}

	d.db, err = product.CreateProductsListByBytes(fileManager.Data)
	if err != nil {
		panic(err)
	}
}

func ConvertProductToProductJSON(p product.Product) ProductJSON {
	return ProductJSON{
		ID:          p.ID,
		Name:        p.Name,
		Quantity:    p.Quantity,
		CodeValue:   p.CodeValue,
		IsPublished: p.IsPublished,
		Expiration:  p.Expiration,
		Price:       p.Price,
	}
}

func (d *DefaultProducts) CheckCodeValueAlreadyExist(code string) bool {
	for _, p := range d.db {
		if p.CodeValue == code {
			return true
		}
	}
	return false
}

func isValidDate(dateString string) bool {
	parsedTime, err := time.Parse("02/01/2006", dateString)
	if err != nil {
		return false
	}

	return parsedTime.Year() >= 2023 &&
		parsedTime.Month() >= 1 && parsedTime.Month() <= 12 &&
		parsedTime.Day() >= 1 && parsedTime.Day() <= daysInMonth(parsedTime.Year(), parsedTime.Month())
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func (d *DefaultProducts) GetAllProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productListJSON := []ProductJSON{}
		for _, p := range d.db {
			pJSON := ConvertProductToProductJSON(p)
			productListJSON = append(productListJSON, pJSON)
		}
		ctx.JSON(http.StatusOK, productListJSON)
	}
}

func (d *DefaultProducts) GetProductByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		for _, p := range d.db {
			if p.ID == productId {
				productJSON := ProductJSON{
					ID:          p.ID,
					Name:        p.Name,
					Quantity:    p.Quantity,
					CodeValue:   p.CodeValue,
					IsPublished: p.IsPublished,
					Expiration:  p.Expiration,
					Price:       p.Price,
				}
				ctx.JSON(http.StatusOK, productJSON)
				return
			}
		}

		ctx.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
	}
}

func (d *DefaultProducts) GetProductsAbovePrice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
		for _, product := range d.db {
			if product.Price > productPriceGt {
				productsWithPriceGt = append(productsWithPriceGt, product)
			}
		}

		if len(productsWithPriceGt) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("There are not any products with price greater than %.2f", productPriceGt)})
			return
		}

		productListJSON := []ProductJSON{}
		for _, p := range productsWithPriceGt {
			pJSON := ConvertProductToProductJSON(p)
			productListJSON = append(productListJSON, pJSON)
		}
		ctx.JSON(http.StatusOK, productListJSON)
	}
}

func (d *DefaultProducts) CreateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body BodyRequestCreate
		var m map[string]any
		if err := ctx.ShouldBindJSON(&m); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
			return
		}

		name, ok := m["name"].(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "name is a obligatory attribute"})
			return
		}
		qtyFloat, ok := m["quantity"].(float64)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "quantity is a obligatory attribute"})
			return
		}
		qty, err := strconv.Atoi(fmt.Sprintf("%.0f", qtyFloat))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if qty < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "quantity must be greater than 0"})
			return
		}
		code, ok := m["code_value"].(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "code_value is a obligatory attribute"})
			return
		}
		if d.CheckCodeValueAlreadyExist(code) {
			ctx.JSON(http.StatusConflict, gin.H{"message": "code already exist"})
			return
		}
		exp, ok := m["expiration"].(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "expiration is a obligatory attribute"})
			return
		}
		if !isValidDate(exp) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid expiration date"})
			return
		}
		price, ok := m["price"].(float64)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "price is a obligatory attribute"})
			return
		}
		if price < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "price must be greater than 0"})
			return
		}

		body = BodyRequestCreate{
			Name:        name,
			Quantity:    qty,
			CodeValue:   code,
			IsPublished: false,
			Expiration:  exp,
			Price:       price,
		}

		p := product.Product{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}
		d.lastID++
		p.ID = d.lastID
		d.db = append(d.db, p)

		ctx.JSON(http.StatusCreated, map[string]any{
			"message": "product created",
			"data": ProductJSON{
				ID:          p.ID,
				Name:        p.Name,
				Quantity:    p.Quantity,
				CodeValue:   p.CodeValue,
				IsPublished: p.IsPublished,
				Expiration:  p.Expiration,
				Price:       p.Price,
			},
		})
	}
}
