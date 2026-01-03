package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repo "github.com/matthosch/go_ecommerce_api/internal/adapters/postgresql/sqlc"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock  = errors.New("product has insufficient stock")
	ErrOrderNotFound   = errors.New("order not found")
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
	GetOrderDetails(ctx context.Context, orderID int64) (orderDetails, error)
}

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{repo: repo, db: db}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)

	// create an order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}
	// look for the product if exists
	for _, item := range tempOrder.Items {
		product, err := qtx.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrProductNoStock
		}

		// create order items
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, err
		}
		// deduct the product stock quantity
		_, err = qtx.UpdateProductQuantity(ctx, repo.UpdateProductQuantityParams{
			ID:       item.ProductID,
			Quantity: product.Quantity - item.Quantity,
		})
		if err != nil {
			return repo.Order{}, err
		}
	}

	tx.Commit(ctx)
	return order, nil
}

// GET /orders/{id} to retrieve order details
func (s *svc) GetOrderDetails(ctx context.Context, orderID int64) (orderDetails, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return orderDetails{}, ErrOrderNotFound
		}
		return orderDetails{}, err
	}

	products, err := s.repo.GetProductsByOrderID(ctx, orderID)
	if err != nil {
		return orderDetails{}, err
	}

	return orderDetails{
		OrderId:    order.ID,
		CustomerID: order.CustomerID,
		CreatedAt:  order.CreatedAt,
		Products:   products,
	}, nil
}
