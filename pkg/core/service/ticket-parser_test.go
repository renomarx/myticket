package service

import (
	"testing"

	"github.com/renomarx/myticket/pkg/core/model"
	"github.com/stretchr/testify/assert"
)

func TestTicketParser(t *testing.T) {
	body := []byte(`Order: 123456
VAT: 3.10
Total: 16.90

product,product_id,price
Formule(s) midi,aZde,14.90
Café,IZ8z,2
`)
	ticket := model.Ticket{
		Body: body,
	}

	parser := NewTicketParser()
	order, err := parser.Parse(&ticket)
	if err != nil {
		t.Error(err)
	}
	// TODO: other assertions
	assert.Equal(t, 2, len(order.Products))
}
