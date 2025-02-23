package http

import (
	"net/http"

	"github.com/ercancavusoglu/messaging/internal/interfaces/http/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(messageHandler *handlers.MessageHandler) http.Handler {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/messages", messageHandler.GetMessages).Methods("GET")
	api.HandleFunc("/scheduler/start", messageHandler.StartScheduler).Methods("GET")
	api.HandleFunc("/scheduler/stop", messageHandler.StopScheduler).Methods("GET")

	return router
}
