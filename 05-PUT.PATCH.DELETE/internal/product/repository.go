package product

import (
	"goweb/05-PUT.PATCH.DELETE/internal/domain"
)

type Repository interface {
	FindAll() (products []domain.Product, err error)
	FindByID(id int) (product domain.Product, err error)
	Insert(product domain.Product) (domain.Product, error)
	Update(product domain.Product) (domain.Product, error)
	Delete(id int) error
}

type RepositorySlice struct {
	db     []domain.Product
	lastId int
}

func NewRepositorySlice(db []domain.Product, lastId int) *RepositorySlice {
	return &RepositorySlice{
		db:     db,
		lastId: lastId,
	}
}

func (r *RepositorySlice) FindAll() (products []domain.Product, err error) {
	products = make([]domain.Product, len(r.db))
	copy(products, r.db)
	return
}

func (r *RepositorySlice) FindByID(id int) (product domain.Product, err error) {
	for i := range r.db {
		if r.db[i].ID == id {
			product = r.db[i]
			return
		}
	}
	return domain.Product{}, ErrNotFound
}

func (r *RepositorySlice) Insert(product domain.Product) (domain.Product, error) {
	r.lastId++
	product.ID = r.lastId
	r.db = append(r.db, product)
	return product, nil
}

func (r *RepositorySlice) Update(product domain.Product) (domain.Product, error) {
	for i := range r.db {
		if r.db[i].ID == product.ID {
			r.db[i] = product
			return product, nil
		}
	}
	return domain.Product{}, ErrNotFound
}

func (r *RepositorySlice) Delete(id int) error {
	index := -1
	for i := range r.db {
		if r.db[i].ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return ErrNotFound
	}

	r.db = append(r.db[:index], r.db[index+1:]...)

	return nil
}
