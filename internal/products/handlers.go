package products

import (
	"fmt"
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

// ListProducts handles the HTTP request to list all products.
func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

// FindProductByID handles the HTTP request to find a product by its ID.
func (h *handler) FindProductByID(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from URL
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.service.FindProductByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, product)
}

// CreateProduct handles the HTTP request to create a new product.
func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var tempProduct Product
	if err := json.Read(r, &tempProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate payload
	if err := validateProductInput(tempProduct); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.service.CreateProduct(r.Context(), tempProduct)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, product)
}

func validateProductInput(p Product) error {
	if p.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if p.PriceInCents <= 0 {
		return fmt.Errorf("product price must be greater than zero")
	}
	if p.Quantity < 0 {
		return fmt.Errorf("product quantity cannot be negative")
	}
	return nil
}
