package handler

import (
	"fmt"
	"goweb/04-EstructuracionAPI/internal/domain"
	"goweb/04-EstructuracionAPI/internal/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	s product.Service
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

type BodyRequestCreateProduct struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func NewProductHandler(s product.Service) *ProductHandler {
	return &ProductHandler{
		s: s,
	}
}

func ConvertProductSliceToProductJSONSlice(products []domain.Product) (productsJSON []ProductJSON) {
	for _, p := range products {
		productJSON := ConvertProductToProductJSON(p)
		productsJSON = append(productsJSON, productJSON)
	}
	return productsJSON
}

func ConvertProductToProductJSON(product domain.Product) (productJSON ProductJSON) {
	return ProductJSON{
		ID:          product.ID,
		Name:        product.Name,
		Quantity:    product.Quantity,
		CodeValue:   product.CodeValue,
		IsPublished: product.IsPublished,
		Expiration:  product.Expiration,
		Price:       product.Price,
	}
}

func (h *ProductHandler) GetAllProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := h.s.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
		productsJSON := ConvertProductSliceToProductJSONSlice(products)
		ctx.JSON(http.StatusOK, productsJSON)
	}
}

func (h *ProductHandler) GetProductByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid product identifier",
			})
			return
		}

		productFound, err := h.s.GetByID(productID)
		if err != nil {
			switch err {
			case product.ErrNotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			default:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "an unexpected error occurred",
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, ConvertProductToProductJSON(productFound))
	}
}

func (h *ProductHandler) GetAllProductsAbovePrice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productPriceGt, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid or non-existent query parameter",
			})
			return
		}

		productsWithPriceGt, err := h.s.GetAllAbovePrice(productPriceGt)
		if err != nil {
			switch err {
			case product.ErrNotGt:
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "an unexpected error occurred",
				})
				return
			}
		}

		productsJSON := ConvertProductSliceToProductJSONSlice(productsWithPriceGt)
		ctx.JSON(http.StatusOK, productsJSON)
	}
}

func (h *ProductHandler) CreateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body BodyRequestCreateProduct
		var m map[string]any
		if err := ctx.ShouldBindJSON(&m); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		name, ok := m["name"].(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "name is a obligatory attribute"})
			return
		}
		qtyFloat, ok := m["quantity"].(float64)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "quantity is a obligatory attribute"})
			return
		}
		qty, err := strconv.Atoi(fmt.Sprintf("%.0f", qtyFloat))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, product.ErrInvalidQuantity)
			return
		}
		code, ok := m["code_value"].(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "code_value is a obligatory attribute"})
			return
		}
		exp, ok := m["expiration"].(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "expiration is a obligatory attribute"})
			return
		}
		price, ok := m["price"].(float64)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "price is a obligatory attribute"})
			return
		}

		body = BodyRequestCreateProduct{
			Name:       name,
			Quantity:   qty,
			CodeValue:  code,
			Expiration: exp,
			Price:      price,
		}
		isPublished, ok := m["is_published"].(bool)
		if ok {
			body.IsPublished = isPublished
		}

		p := domain.Product{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		newProduct, err := h.s.Create(p)
		if err != nil {
			switch {
			case err == product.ErrInvalidName || err == product.ErrInvalidCodeValue || err == product.ErrInvalidQuantity || err == product.ErrInvalidExpiration || err == product.ErrInvalidPrice:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "an unexpected error occurred",
				})
				return
			}
		}

		newProductJSON := ConvertProductToProductJSON(newProduct)
		ctx.JSON(http.StatusOK, newProductJSON)
	}
}
