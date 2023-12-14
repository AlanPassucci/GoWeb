package router

import (
	"goweb/Desafio/cmd/server/handler"
	"goweb/Desafio/internal/domain"
	"goweb/Desafio/internal/ticket"

	"github.com/gin-gonic/gin"
)

type Router struct {
	rt   *gin.Engine
	list []domain.Ticket
}

func NewRouter(rt *gin.Engine, list []domain.Ticket) *Router {
	return &Router{
		rt:   rt,
		list: list,
	}
}

func (r *Router) MapRoutes() {
	rp := ticket.NewRepository(r.list)
	sv := ticket.NewService(rp)
	hd := handler.NewHandler(sv)

	ticketGroup := r.rt.Group("/ticket")
	{
		ticketGroup.GET("/getByCountry/:dest", hd.GetTicketsByDestination())
		ticketGroup.GET("/getPercentage/:dest", hd.GetPercentageByDestination())
	}
}
