package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/charismen/home-api/internal/service"
)

// ItemHandler handles HTTP requests for items
type ItemHandler struct {
	service *service.ItemService
}

// NewItemHandler creates a new item handler
func NewItemHandler(service *service.ItemService) *ItemHandler {
	return &ItemHandler{
		service: service,
	}
}

// SyncItems handles the /sync endpoint
func (h *ItemHandler) SyncItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.service.SyncItems(r.Context())
	if err != nil {
		http.Error(w, "Failed to sync items: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Items synced successfully",
	})
}

// GetItems handles the /items endpoint
func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items, err := h.service.GetAllItems()
	if err != nil {
		http.Error(w, "Failed to get items: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}