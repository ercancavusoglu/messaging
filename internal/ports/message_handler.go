package ports

import (
	"net/http"
)

type MessageHandler interface {
	StartScheduler(w http.ResponseWriter, r *http.Request)
	StopScheduler(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
}
