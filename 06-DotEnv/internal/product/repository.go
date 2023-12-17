package product

import (
	"goweb/06-DotEnv/internal/domain"
	"goweb/06-DotEnv/pkg/store"
)

type Repository interface {
	FindAll() (products []domain.Product, err error)
	FindByID(id int) (product domain.Product, err error)
	Insert(product domain.Product) (domain.Product, error)
	Update(product domain.Product) (domain.Product, error)
	Delete(id int) error
}

type RepositorySlice struct {
	fm store.Store
}

func NewRepositorySlice(fm store.Store) *RepositorySlice {
	return &RepositorySlice{
		fm: fm,
	}
}

func (r *RepositorySlice) FindAll() (products []domain.Product, err error) {
	products, err = r.fm.FindAll()
	if err != nil {
		return nil, err
	}

	return
}

func (r *RepositorySlice) FindByID(id int) (product domain.Product, err error) {
	product, err = r.fm.FindByID(id)
	if err != nil {
		return domain.Product{}, ErrNotFound
	}
	return
}

func (r *RepositorySlice) Insert(product domain.Product) (domain.Product, error) {
	newProduct, err := r.fm.Insert(product)
	if err != nil {
		return domain.Product{}, err
	}
	return newProduct, nil
}

func (r *RepositorySlice) Update(product domain.Product) (domain.Product, error) {
	updatedProduct, err := r.fm.Update(product)
	if err != nil {
		return domain.Product{}, ErrNotFound
	}
	return updatedProduct, nil
}

func (r *RepositorySlice) Delete(id int) error {
	if err := r.fm.Delete(id); err != nil {
		return ErrNotFound
	}
	return nil
}
