package product

import (
	"goweb/04-EstructuracionAPI/internal/domain"
	"goweb/04-EstructuracionAPI/pkg"
)

type Service interface {
	GetAll() (products []domain.Product, err error)
	GetByID(id int) (product domain.Product, err error)
	GetAllAbovePrice(priceGt float64) (product []domain.Product, err error)
	Create(product domain.Product) (domain.Product, error)
}

type DefaultService struct {
	rp Repository
}

func NewDefaultService(rp Repository) *DefaultService {
	return &DefaultService{
		rp: rp,
	}
}

func (s *DefaultService) GetAll() (products []domain.Product, err error) {
	products, err = s.rp.FindAll()
	if err != nil {
		return nil, err
	}
	return
}

func (s *DefaultService) GetByID(id int) (product domain.Product, err error) {
	product, err = s.rp.FindByID(id)
	if err != nil {
		return
	}

	return product, nil
}

func (s *DefaultService) GetAllAbovePrice(priceGt float64) ([]domain.Product, error) {
	productsWithPriceGt := []domain.Product{}
	products, err := s.rp.FindAll()
	if err != nil {
		return nil, err
	}

	for _, p := range products {
		if p.Price > priceGt {
			productsWithPriceGt = append(productsWithPriceGt, p)
		}
	}

	if len(productsWithPriceGt) == 0 {
		return nil, ErrNotGt
	}

	return productsWithPriceGt, nil
}

func (s *DefaultService) Create(product domain.Product) (domain.Product, error) {
	if product.Name == "" {
		return domain.Product{}, ErrInvalidName
	}
	if product.Quantity < 1 {
		return domain.Product{}, ErrInvalidQuantity
	}
	products, err := s.GetAll()
	if err != nil {
		return domain.Product{}, err
	}
	if product.CodeValue == "" || pkg.CheckCodeValueAlreadyExist(products, product.CodeValue) {
		return domain.Product{}, ErrInvalidCodeValue
	}
	if product.Expiration == "" || !pkg.IsValidDate(product.Expiration) {
		return domain.Product{}, ErrInvalidExpiration
	}
	if product.Price < 1 {
		return domain.Product{}, ErrInvalidPrice
	}

	insertedProduct, err := s.rp.Insert(product)
	if err != nil {
		return domain.Product{}, err
	}
	return insertedProduct, nil
}
