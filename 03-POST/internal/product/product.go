package product

import "encoding/json"

type Product struct {
	ID          int
	Name        string
	Quantity    int
	CodeValue   string
	IsPublished bool
	Expiration  string
	Price       float64
}

func CreateProductsListByBytes(data []byte) ([]Product, error) {
	newProductsList := []Product{}

	err := json.Unmarshal(data, &newProductsList)
	if err != nil {
		return nil, err
	}

	return newProductsList, nil
}
