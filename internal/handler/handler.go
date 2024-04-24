package handler

import (
	"html/template"
	"log"
	"net/http"
	//
	"github.com/mishaRomanov/wb-l0/internal/storage/cache"
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
	log.Println("New get by id request")
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
	tmpl, err := template.ParseFiles("/Users/misha/coding/GoProjects/wb-l0/internal/handler/template.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(w, "orders", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
