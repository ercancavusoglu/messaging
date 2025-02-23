package http

import (
    "net/http"
    
    "github.com/gorilla/mux"
    "github.com/ercancavusoglu/messaging/internal/interfaces/http/handlers"
)

type Router struct {
    *mux.Router
    messageHandler *handlers.MessageHandler
}

func NewRouter(messageHandler *handlers.MessageHandler) *Router {
    r := mux.NewRouter()
    router := &Router{
        Router:         r,
        messageHandler: messageHandler,
    }
    
    router.setupRoutes()
    return router
}

func (r *Router) setupRoutes() {
    api := r.PathPrefix("/api/v1").Subrouter()
    
    api.HandleFunc("/scheduler/start", r.messageHandler.StartScheduler).Methods(http.MethodPost)
    api.HandleFunc("/scheduler/stop", r.messageHandler.StopScheduler).Methods(http.MethodPost)
    api.HandleFunc("/messages", r.messageHandler.ListMessages).Methods(http.MethodGet)
}

