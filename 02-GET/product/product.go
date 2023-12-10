package product

import "encoding/json"

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func DecerealizeJSON(data []byte) ([]Product, error) {
	newProductsList := []Product{}

	err := json.Unmarshal(data, &newProductsList)
	if err != nil {
		return nil, err
	}

	return newProductsList, nil
}
