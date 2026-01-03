package orders

import (
	"errors"
	"log/slog"
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
		slog.Error("invalid JSON in order request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate payload
	if err := validateOrderInput(tempOrder); err != nil {
		slog.Error("invalid order input", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		slog.Error("failed to place order", "error", err)

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
	if err != nil || id <= 0 {
		slog.Error("invalid order ID", "error", err)
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}
	orderDetails, err := h.service.GetOrderDetails(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		slog.Error("failed to get order details", "order_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, orderDetails)
}

func validateOrderInput(o createOrderParams) error {
	if o.CustomerID <= 0 {
		return errors.New("customer ID must be positive")
	}
	if len(o.Items) == 0 {
		return errors.New("at least one order item is required")
	}
	return nil
}
