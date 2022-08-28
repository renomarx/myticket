package ports

import "github.com/renomarx/myticket/pkg/core/model"

// TicketRepository tickets database
type TicketRepository interface {
	// CreateTicket create ticket in database
	// Return an error only if the ticket couldn't be saved
	CreateTicket(ticket *model.Ticket) error
	// UpdateTicket ticket in database
	UpdateTicket(ticket *model.Ticket) error
}

// ProductRepository products database
type ProductRepository interface {
	// SaveProduct upsert products
	SaveProducts(products []model.Product) error
}
