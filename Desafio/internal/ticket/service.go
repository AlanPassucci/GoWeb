package ticket

import (
	"errors"
	"goweb/Desafio/internal/domain"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetTicketsByDestination(ctx *gin.Context, destination string) ([]domain.Ticket, error)
	GetPercentageByDestination(ctx *gin.Context, destination string) (float64, error)
}

type service struct {
	rp Repository
}

func NewService(r Repository) *service {
	return &service{
		rp: r,
	}
}

func (s *service) GetTicketsByDestination(ctx *gin.Context, destination string) ([]domain.Ticket, error) {
	if destination == "" {
		return nil, errors.New("invalid destination")
	}

	ticketsDest, err := s.rp.GetAllByDestination(ctx, destination)
	if err != nil {
		return nil, err
	}

	return ticketsDest, nil
}

func (s *service) GetPercentageByDestination(ctx *gin.Context, destination string) (float64, error) {
	if destination == "" {
		return 0.0, errors.New("invalid destination")
	}

	tickets, err := s.rp.GetAll(ctx)
	if err != nil {
		return 0.0, err
	}

	ticketsDest, err := s.rp.GetAllByDestination(ctx, destination)
	if err != nil {
		return 0.0, err
	}

	percentage := float64(len(ticketsDest)) * 100 / float64(len(tickets))
	return float64(percentage), nil
}
