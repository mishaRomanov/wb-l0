package handler

import (
	"encoding/json"
	"github.com/mishaRomanov/wb-l0/internal/storage/cache"
	"net/http"
)

// handler struct
type Handler struct {
	Cache *cache.OrdersCache
}

// func that creates a new handler
func NewHandler(inMemory *cache.OrdersCache) Handler {
	return Handler{
		Cache: inMemory,
	}
}

// Handler that returns an order with the given id
func (h *Handler) GetByID(w http.ResponseWriter, req *http.Request) {
	//Only GET method is allowed
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid method. Only GET method is supported"))
		return
	}

	//getting an id from path
	data, exists := h.Cache.Get(req.PathValue("id"))
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Order not found"))
		return
	}
	//encoding map to json
	bidata, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error while marshalling map to json"))
		return
	}
	//writing header and actual data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bidata)
}
