package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"goweb/06-DotEnv/internal/domain"
	"os"
)

type Store interface {
	FindAll() (products []domain.Product, err error)
	FindByID(id int) (product domain.Product, err error)
	Insert(product domain.Product) (domain.Product, error)
	Update(product domain.Product) (domain.Product, error)
	Delete(id int) error
	save(products []domain.Product) (err error)
}

type FileManager struct {
	FilePath string
	lastId   int
}

func NewFileManager(filePath string) *FileManager {
	return &FileManager{
		FilePath: filePath,
		lastId:   500,
	}
}

func (fm *FileManager) FindAll() (products []domain.Product, err error) {
	data, err := os.ReadFile(fm.FilePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (fm *FileManager) FindByID(id int) (product domain.Product, err error) {
	products, err := fm.FindAll()
	if err != nil {
		return domain.Product{}, err
	}

	for _, p := range products {
		if p.ID == id {
			return p, nil
		}
	}

	return domain.Product{}, errors.New("product not found")
}

func (fm *FileManager) Insert(product domain.Product) (newProduct domain.Product, err error) {
	products, err := fm.FindAll()
	if err != nil {
		return domain.Product{}, err
	}

	fm.lastId++
	product.ID = fm.lastId
	fmt.Println(product)
	products = append(products, product)

	if err = fm.save(products); err != nil {
		return domain.Product{}, err
	}

	newProduct = product
	return newProduct, nil
}

func (fm *FileManager) Update(product domain.Product) (updatedProduct domain.Product, err error) {
	products, err := fm.FindAll()
	if err != nil {
		return domain.Product{}, err
	}

	for i, p := range products {
		if p.ID == product.ID {
			products[i] = product
			if err = fm.save(products); err != nil {
				return domain.Product{}, err
			}
			updatedProduct = product
			return updatedProduct, nil
		}
	}

	return domain.Product{}, errors.New("product not found")
}

func (fm *FileManager) Delete(id int) (err error) {
	products, err := fm.FindAll()
	if err != nil {
		return err
	}

	index := -1
	for i, p := range products {
		if p.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("product not found")
	}

	products = append(products[:index], products[index+1:]...)
	if err = fm.save(products); err != nil {
		return err
	}

	return nil
}

func (fm *FileManager) save(products []domain.Product) (err error) {
	jsonData, err := json.MarshalIndent(products, "", "\t")
	if err != nil {
		return err
	}

	if err = os.WriteFile("../../products.json", jsonData, os.ModePerm); err != nil {
		return err
	}

	return nil
}
