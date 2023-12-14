package handler

import (
	"goweb/Desafio/internal/ticket"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	sv ticket.Service
}

type TicketJSON struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Country string  `json:"country"`
	Time    string  `json:"time"`
	Price   float64 `json:"price"`
}

func NewHandler(s ticket.Service) *Handler {
	return &Handler{
		sv: s,
	}
}

func (h *Handler) GetTicketsByDestination() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		destination := ctx.Param("dest")

		tickets, err := h.sv.GetTicketsByDestination(ctx, destination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "an unexpected error occurred",
			})
			return
		}

		var ticketsJSON []TicketJSON
		for _, t := range tickets {
			ticketJSON := TicketJSON{
				Id:      t.Id,
				Name:    t.Name,
				Email:   t.Email,
				Country: t.Country,
				Time:    t.Time,
				Price:   t.Price,
			}
			ticketsJSON = append(ticketsJSON, ticketJSON)
		}
		ctx.JSON(http.StatusOK, ticketsJSON)
	}
}

func (h *Handler) GetPercentageByDestination() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		destination := ctx.Param("dest")

		percentage, err := h.sv.GetPercentageByDestination(ctx, destination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "an unexpected error occurred",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"percentage": percentage,
		})
	}
}
