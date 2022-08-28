package ports

import (
	"github.com/renomarx/myticket/pkg/core/model"
)

// TicketHandler main ticket handler
type TicketHandler interface {
	// Handle treat a ticket
	// must not loose the ticket
	Handle(ticket *model.Ticket)
}

// TicketParser ticket parser
type TicketParser interface {
	// Parse a ticket to a model.Order
	Parse(ticket *model.Ticket) (*model.Order, error)
}

// TicketFallback service handling tickets that coudln't be saved in database
type TicketFallback interface {
	// Push a ticket that needs to be treated again
	Push(ticket *model.Ticket)
}
