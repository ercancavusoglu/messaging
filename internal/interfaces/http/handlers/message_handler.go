package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ercancavusoglu/messaging/internal/domain/message"
	"github.com/ercancavusoglu/messaging/internal/domain/scheduler"
)

type MessageHandler struct {
	messageService *message.Service
	scheduler      *scheduler.Scheduler
}

func NewMessageHandler(messageService *message.Service, scheduler *scheduler.Scheduler) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		scheduler:      scheduler,
	}
}

func (h *MessageHandler) StartScheduler(w http.ResponseWriter, r *http.Request) {
	if h.scheduler.IsRunning() {
		h.jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Scheduler is already running",
		})
		return
	}

	ctx := r.Context()
	go func() {
		if err := h.scheduler.Start(ctx); err != nil {
			// Log error
		}
	}()

	h.jsonResponse(w, http.StatusOK, map[string]string{
		"message": "Scheduler started successfully",
	})
}

func (h *MessageHandler) StopScheduler(w http.ResponseWriter, r *http.Request) {
	if !h.scheduler.IsRunning() {
		h.jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Scheduler is not running",
		})
		return
	}

	h.scheduler.Stop()
	h.jsonResponse(w, http.StatusOK, map[string]string{
		"message": "Scheduler stopped successfully",
	})
}

func (h *MessageHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.messageService.List()
	if err != nil {
		h.jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch messages",
		})
		return
	}

	h.jsonResponse(w, http.StatusOK, messages)
}

func (h *MessageHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}
