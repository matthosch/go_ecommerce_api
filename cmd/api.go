package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	repo "github.com/matthosch/go_ecommerce_api/internal/adapters/postgresql/sqlc"
	"github.com/matthosch/go_ecommerce_api/internal/orders"
	"github.com/matthosch/go_ecommerce_api/internal/products"
)

// mount sets up the HTTP routes and middleware.
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // for rate limiting and analytics/tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	// Products routes
	productService := products.NewService(repo.New(app.db))
	productsHandler := products.NewHandler(productService)
	r.Get("/products", productsHandler.ListProducts)
	r.Get("/products/{id}", productsHandler.FindProductByID)

	// Orders routes
	orderService := orders.NewService(repo.New(app.db), app.db)
	ordersHandler := orders.NewHandler(orderService)
	r.Post("/orders", ordersHandler.PlaceOrder)

	return r
}

// run starts the HTTP server.
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server starting on %s", app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string // user,pw,db name,etc
}
