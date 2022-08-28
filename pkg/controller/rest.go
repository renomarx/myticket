package controller

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/renomarx/myticket/pkg/core/model"
	"github.com/renomarx/myticket/pkg/core/ports"

	"github.com/julienschmidt/httprouter"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.AddHook(filename.NewHook())
}

// RestAPI Main service handling http requests
type RestAPI struct {
	MetricsController *metricsController
	ticketHandler     ports.TicketHandler
}

// RestAPIError error struct to send a json result of an error
type RestAPIError struct {
	Error string `json:"error"`
}

// NewRestAPI RestAPI constructor with dependencies injected
func NewRestAPI(ticketHandler ports.TicketHandler) *RestAPI {
	return &RestAPI{
		MetricsController: NewMetricsController(),
		ticketHandler:     ticketHandler,
	}
}

// Serve Http listen to REST_PORT
func (api *RestAPI) Serve() {
	router := httprouter.New()
	api.Route(router)
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		logrus.Fatalf("No HTTP_PORT env variable found")
	}

	logrus.Infof("API listening on %s", port)
	logrus.Fatal(http.ListenAndServe(port, router))
}

// Route configure http router
func (api *RestAPI) Route(r *httprouter.Router) {
	// Handlers
	r.GET("/ping", api.Ping)

	// Because I like to have some metrics on every service, and I like prometheus
	r.Handler("GET", "/metrics", api.MetricsController.HTTPHandler())

	logrus.Infof("Serving metrics on /metrics")
	r.NotFound = http.HandlerFunc(api.NotFound)

	r.POST("/ticket", api.TicketWebhook)
	logrus.Infof("Serving webhook /ticket")
}

// Ping handle /ping http requests (for health checks)
// Always useful for monitoring tools or orchestrators like K8s or Nomad
func (api *RestAPI) Ping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pong"))
}

// NotFound handle http routes not found by router - only to trace bad 404 http calls
func (api *RestAPI) NotFound(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("HTTP Not found: %s", r.URL.Path)
	metrics := api.MetricsController.GetMetrics()
	metrics.RouterHTTPNotFound.Inc()
	w.WriteHeader(http.StatusNotFound)
}

// TicketWebHook handle tickets
func (api *RestAPI) TicketWebhook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		// No data posted, no need to continue
		w.WriteHeader(http.StatusOK)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		// If we cannot read request body, we cannot do anything more
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ticket := model.NewTicket(body)
	api.ticketHandler.Handle(ticket)
	// TODO: also send ticket to an external queue like redis streams, rabbitMQ or Kafka, for external processes
	w.WriteHeader(http.StatusOK)
}
