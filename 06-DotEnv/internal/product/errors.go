package product

import "errors"

var (
	ErrInvalidName       = errors.New("invalid product name")
	ErrInvalidQuantity   = errors.New("invalid product quantity")
	ErrInvalidCodeValue  = errors.New("invalid or already in use product code value")
	ErrInvalidExpiration = errors.New("invalid product expiration")
	ErrInvalidPrice      = errors.New("invalid product price")
	ErrNotFound          = errors.New("product not found")
	ErrNotGt             = errors.New("there are not any products with price greater than")
)
