package service

import (
	"time"

	"github.com/renomarx/myticket/pkg/core/model"
	"github.com/renomarx/myticket/pkg/core/ports"
)

type RawTicketHandler struct {
	TicketRepo     ports.TicketRepository
	TicketFallback ports.TicketFallback
	TicketParser   ports.TicketParser
	ProductRepo    ports.ProductRepository
}

func NewRawTicketHandler(
	ticketRepo ports.TicketRepository,
	productRepo ports.ProductRepository,
	ticketFallback ports.TicketFallback,
) *RawTicketHandler {
	return &RawTicketHandler{
		TicketRepo:     ticketRepo,
		TicketParser:   NewTicketParser(),
		ProductRepo:    productRepo,
		TicketFallback: ticketFallback,
	}
}

func (handler *RawTicketHandler) Handle(ticket *model.Ticket) {
	// We save the ticket in database before even treating it, to avoid a maximum losses
	ticket.CreatedAt = time.Now()
	ticket.UpdatedAt = time.Now()
	err := handler.TicketRepo.CreateTicket(ticket)
	if err != nil {
		// In case of saving error at this point, we push to a service that will retry as much as it can to treat the ticket
		handler.TicketFallback.Push(ticket)
		return
	}
	// Then we parse it in order to populate our products' database
	order, err := handler.TicketParser.Parse(ticket)
	if err != nil {
		// Update ticket to status error, for it to be eventually treated manually later
		ticket.Status = model.TicketStatusError
		ticket.ErrorDetails = err.Error()
		handler.TicketRepo.UpdateTicket(ticket)
	}
	// We could save the whole order, for now we only consider products
	handler.ProductRepo.SaveProducts(order.Products)
}
