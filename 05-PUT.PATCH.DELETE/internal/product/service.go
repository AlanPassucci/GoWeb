package product

import (
	"goweb/05-PUT.PATCH.DELETE/internal/domain"
	"goweb/05-PUT.PATCH.DELETE/pkg"
)

type Service interface {
	GetAll() (products []domain.Product, err error)
	GetByID(id int) (product domain.Product, err error)
	GetAllAbovePrice(priceGt float64) (product []domain.Product, err error)
	Create(product domain.Product) (domain.Product, error)
	Update(product domain.Product, patch bool) (domain.Product, error)
	Delete(id int) error
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

func (s *DefaultService) Update(product domain.Product, patch bool) (domain.Product, error) {
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
	if !patch {
		if product.CodeValue == "" || pkg.CheckCodeValueAlreadyExist(products, product.CodeValue) {
			return domain.Product{}, ErrInvalidCodeValue
		}
	}
	if product.Expiration == "" || !pkg.IsValidDate(product.Expiration) {
		return domain.Product{}, ErrInvalidExpiration
	}
	if product.Price < 1 {
		return domain.Product{}, ErrInvalidPrice
	}

	updatedProduct, err := s.rp.Update(product)
	if err != nil {
		return domain.Product{}, err
	}
	return updatedProduct, nil
}

func (s *DefaultService) Delete(id int) error {
	if err := s.rp.Delete(id); err != nil {
		return err
	}
	return nil
}
