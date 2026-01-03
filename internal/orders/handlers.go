package orders

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/matthosch/go_ecommerce_api/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

// PlaceOrder handles the HTTP request to place a new order.
func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var tempOrder createOrderParams
	if err := json.Read(r, &tempOrder); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		log.Println(err)

		switch err {
		case ErrProductNoStock:
			http.Error(w, err.Error(), http.StatusConflict)
			return
		case ErrProductNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.Write(w, http.StatusCreated, createdOrder)
}

// GetOrderDetails handles the HTTP request to get order details by ID.
func (h *handler) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}
	orderDetails, err := h.service.GetOrderDetails(r.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, orderDetails)
}
