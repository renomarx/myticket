package service

import (
	"github.com/renomarx/myticket/pkg/core/model"
	"github.com/renomarx/myticket/pkg/core/ports"
	"github.com/sirupsen/logrus"
)

// FallbackQueueSize maximum queue size for tickets to treat again
// Should be big enough to store enough tickets during a database cut, but small enough to avoid consuming too much RAM
const FallbackQueueSize = 10000

// TicketFallback service handling tickets that coudln't be saved in database
//
// For this implementation, I choosed an in-memory queue with a go channel
// It would be better to use an external queue like rabbitmq or redis streams,
// and to read from it in the worker part of this service: that would assure us
// to keep the tickets even in case of a crash of this service,
// but it's a little too sophisticated for this exercise
type TicketFallback struct {
	queue         chan *model.Ticket
	TicketHandler ports.TicketHandler
}

// NewTicketFallback TicketFallback constructor
func NewTicketFallback(ticketHandler ports.TicketHandler) *TicketFallback {
	return &TicketFallback{
		queue:         make(chan *model.Ticket, FallbackQueueSize),
		TicketHandler: ticketHandler,
	}
}

func (f *TicketFallback) Push(ticket *model.Ticket) {
	f.queue <- ticket
}

func (f *TicketFallback) Run() {
	logrus.Println("Tickets fallback running")
	for ticket := range f.queue {
		f.TicketHandler.Handle(ticket)
		// TODO: increment metric of tickets queue size to be able to monitor it
	}
}
