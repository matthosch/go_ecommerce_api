package orders

import (
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/matthosch/go_ecommerce_api/internal/adapters/postgresql/sqlc"
)

type orderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParams struct {
	CustomerID int64       `json:"customerId"`
	Items      []orderItem `json:"items"`
}

type orderDetails struct {
	OrderId    int64                          `json:"orderId"`
	CustomerID int64                          `json:"customerId"`
	CreatedAt  pgtype.Timestamptz             `json:"createdAt"`
	Products   []repo.GetProductsByOrderIDRow `json:"products"`
}
