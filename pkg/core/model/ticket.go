package model

import "time"

type TicketStatus string

const (
	TicketStatusError   = "ERROR"
	TicketStatusToTreat = "TO_TREAT"
	TicketStatusTreated = "TREATED"
)

type Ticket struct {
	ID           uint         `db:"id"`
	Body         []byte       `db:"body"`
	Status       TicketStatus `db:"status"`
	ErrorDetails string       `db:"error_details"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
}

func NewTicket(body []byte) *Ticket {
	return &Ticket{
		Body:   body,
		Status: TicketStatusToTreat,
	}
}
