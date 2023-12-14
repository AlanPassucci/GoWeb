package ticket

import (
	"fmt"

	"goweb/Desafio/internal/domain"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	GetAll(ctx *gin.Context) ([]domain.Ticket, error)
	GetAllByDestination(ctx *gin.Context, destination string) ([]domain.Ticket, error)
}

type repository struct {
	db []domain.Ticket
}

func NewRepository(db []domain.Ticket) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx *gin.Context) ([]domain.Ticket, error) {
	if len(r.db) == 0 {
		return []domain.Ticket{}, fmt.Errorf("empty list of tickets")
	}

	return r.db, nil
}

func (r *repository) GetAllByDestination(ctx *gin.Context, destination string) ([]domain.Ticket, error) {
	var ticketsDest []domain.Ticket

	if len(r.db) == 0 {
		return []domain.Ticket{}, fmt.Errorf("empty list of tickets")
	}

	for _, t := range r.db {
		if t.Country == destination {
			ticketsDest = append(ticketsDest, t)
		}
	}

	return ticketsDest, nil
}
