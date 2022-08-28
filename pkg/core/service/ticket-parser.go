package service

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/renomarx/myticket/pkg/core/model"
	"github.com/sirupsen/logrus"
)

type TicketParser struct {
}

func NewTicketParser() *TicketParser {
	return &TicketParser{}
}

func (parser *TicketParser) Parse(ticket *model.Ticket) (*model.Order, error) {
	order := model.Order{}

	// Headers parsing
	scanner := bufio.NewScanner(bytes.NewReader(ticket.Body))
	seek := 0
	for scanner.Scan() {
		bytes := scanner.Bytes()
		seek += len(bytes) + 1
		line := string(bytes)
		if line == "" {
			// Line break, we stop custom parsing to use next csv parser
			break
		}
		parser.parseHeader(line, &order)
	}

	// CSV parsing
	csvData := ticket.Body[seek:]
	reader := csv.NewReader(bytes.NewBuffer(csvData))
	_, err := reader.Read() // skip first line
	if err != nil {
		if err != io.EOF {
			return &order, err
		}
	}
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			logrus.Error(err)
			continue
		}
		parser.parseCSVLine(record, &order)
	}
	return &order, nil
}

func (parser *TicketParser) parseHeader(line string, order *model.Order) {
	header := strings.Split(line, ":")
	if len(header) != 2 {
		// Header format not recognized
		return
	}
	key := strings.ToLower(strings.TrimSpace(header[0]))
	value := strings.TrimSpace(header[1])
	switch key {
	case "order":
		order.ID = value
	case "vat":
		fvalue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logrus.Error(err)
			return
		}
		order.VAT = fvalue
	case "total":
		fvalue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logrus.Error(err)
			return
		}
		order.Total = fvalue
	}
}

func (parser *TicketParser) parseCSVLine(record []string, order *model.Order) error {
	if len(record) < 2 {
		return fmt.Errorf("Record length < 2, impossible to parse %v", record)
	}
	product := model.Product{
		ID:        strings.TrimSpace(record[1]),
		Name:      strings.TrimSpace(record[0]),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if len(record) >= 3 {
		price, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
		if err != nil {
			logrus.Error(err)
		}
		product.Price = price
	}
	order.Products = append(order.Products, product)
	return nil
}
