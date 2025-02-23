package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ercancavusoglu/messaging/internal/domain/message"
	"github.com/ercancavusoglu/messaging/internal/domain/scheduler"
)

type MessageHandler struct {
	messageService *message.Service
	scheduler      *scheduler.Scheduler
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewMessageHandler(messageService *message.Service, scheduler *scheduler.Scheduler) *MessageHandler {
	ctx, cancel := context.WithCancel(context.Background())
	return &MessageHandler{
		messageService: messageService,
		scheduler:      scheduler,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (h *MessageHandler) StartScheduler(w http.ResponseWriter, r *http.Request) {
	if h.scheduler.IsRunning() {
		h.jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Scheduler is already running",
		})
		return
	}

	fmt.Println("[StartScheduler] Scheduler is not running")
	go func() {
		if err := h.scheduler.Start(h.ctx); err != nil {
			fmt.Printf("[StartScheduler] Error: %v\n", err)
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
	h.cancel()
	ctx, cancel := context.WithCancel(context.Background())
	h.ctx = ctx
	h.cancel = cancel

	h.jsonResponse(w, http.StatusOK, map[string]string{
		"message": "Scheduler stopped successfully",
	})
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.messageService.List(100)
	if err != nil {
		h.jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	h.jsonResponse(w, http.StatusOK, messages)
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		To      string `json:"to"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg, err := h.messageService.Create(req.To, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func (h *MessageHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Printf("[Handler] Error encoding response: %v\n", err)
	}
}
