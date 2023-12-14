package ticket

import (
	"testing"

	"goweb/Desafio/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var tickets = []domain.Ticket{
	{
		Id:      "1",
		Name:    "Tait Mc Caughan",
		Email:   "tmc0@scribd.com",
		Country: "Finland",
		Time:    "17:11",
		Price:   785.00,
	},
	{
		Id:      "2",
		Name:    "Padget McKee",
		Email:   "pmckee1@hexun.com",
		Country: "China",
		Time:    "20:19",
		Price:   537.00,
	},
	{
		Id:      "3",
		Name:    "Yalonda Jermyn",
		Email:   "yjermyn2@omniture.com",
		Country: "China",
		Time:    "18:11",
		Price:   579.00,
	},
}

var ticketsByDestination = []domain.Ticket{
	{
		Id:      "2",
		Name:    "Padget McKee",
		Email:   "pmckee1@hexun.com",
		Country: "China",
		Time:    "20:19",
		Price:   537.00,
	},
	{
		Id:      "3",
		Name:    "Yalonda Jermyn",
		Email:   "yjermyn2@omniture.com",
		Country: "China",
		Time:    "18:11",
		Price:   579.00,
	},
}

type stubRepo struct {
	db *DbMock
}

type DbMock struct {
	db  []domain.Ticket
	spy bool
	err error
}

func NewRepositoryTest(dbM *DbMock) Repository {
	return &stubRepo{dbM}
}

func (r *stubRepo) GetAll(ctx *gin.Context) ([]domain.Ticket, error) {
	r.db.spy = true
	if r.db.err != nil {
		return []domain.Ticket{}, r.db.err
	}
	return tickets, nil
}

func (r *stubRepo) GetAllByDestination(ctx *gin.Context, destination string) ([]domain.Ticket, error) {

	var tkts []domain.Ticket

	r.db.spy = true
	if r.db.err != nil {
		return []domain.Ticket{}, r.db.err
	}

	for _, t := range r.db.db {
		if t.Country == destination {
			tkts = append(tkts, t)
		}
	}

	return tkts, nil
}

func TestGetTicketsByDestination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)

	dbMock := &DbMock{
		db:  tickets,
		spy: false,
		err: nil,
	}
	repo := NewRepositoryTest(dbMock)
	service := NewService(repo)

	tkts, err := service.GetTicketsByDestination(ctx, "China")

	assert.Nil(t, err)
	assert.True(t, dbMock.spy)
	assert.Equal(t, len(ticketsByDestination), tkts)
}

func TestGetGetPercentageByDestination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)

	dbMock := &DbMock{
		db:  tickets,
		spy: false,
		err: nil,
	}
	repo := NewRepositoryTest(dbMock)
	service := NewService(repo)

	percentage, err := service.GetPercentageByDestination(ctx, "China")

	assert.Nil(t, err)
	assert.NotNil(t, percentage)
	assert.True(t, dbMock.spy)
}
