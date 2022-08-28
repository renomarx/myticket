package main

import (
	"github.com/renomarx/myticket/pkg/controller"
	"github.com/renomarx/myticket/pkg/core/service"
	"github.com/renomarx/myticket/pkg/repository"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Println("APP STARTING")

	db := repository.NewPostgresDB()

	// Injecting nil as TicketFallback because we'll inject it later (circular reference)
	ticketHandler := service.NewRawTicketHandler(db, db, nil)

	ticketFallback := service.NewTicketFallback(ticketHandler)

	// Injectiong TicketFallback
	ticketHandler.TicketFallback = ticketFallback

	api := controller.NewRestAPI(ticketHandler)

	// We could launch multiple workers, but one should be enough
	go ticketFallback.Run()

	api.Serve()

}
