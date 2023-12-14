package pkg

import (
	"encoding/json"
	"goweb/05-PUT.PATCH.DELETE/internal/domain"
	"os"
	"time"
)

func FillDB[T any](path string, slice *T) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal(data, slice); err != nil {
		panic(err.Error())
	}
}

func CheckCodeValueAlreadyExist(products []domain.Product, code string) bool {
	for _, p := range products {
		if p.CodeValue == code {
			return true
		}
	}
	return false
}

func IsValidDate(dateString string) bool {
	layout := "02/01/2006"
	location, _ := time.LoadLocation("UTC")

	if _, err := time.ParseInLocation(layout, dateString, location); err != nil {
		return false
	}

	return true
}
