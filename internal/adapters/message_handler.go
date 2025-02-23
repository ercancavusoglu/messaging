package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ercancavusoglu/messaging/internal/adapters/scheduler"
	"github.com/ercancavusoglu/messaging/internal/ports"
)

type MessageHandler struct {
	messageService ports.MessageService
	scheduler      *scheduler.SchedulerService
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewMessageHandler(messageService ports.MessageService, scheduler *scheduler.SchedulerService) *MessageHandler {
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
	if err := h.scheduler.Start(h.ctx); err != nil {
		fmt.Printf("[StartScheduler] Error: %v\n", err)
		h.jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

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
	messages, err := h.messageService.GetSendedMessages()
	if err != nil {
		h.jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	h.jsonResponse(w, http.StatusOK, messages)
}

func (h *MessageHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Printf("[Handler] Error encoding response: %v\n", err)
	}
}
