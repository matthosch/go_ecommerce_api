package products

import (
	"context"

	repo "github.com/matthosch/go_ecommerce_api/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductByID(ctx context.Context, id int64) (repo.Product, error)
	CreateProduct(ctx context.Context, tempProduct Product) (repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) FindProductByID(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductByID(ctx, id)
}

func (s *svc) CreateProduct(ctx context.Context, tempProduct Product) (repo.Product, error) {
	return s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:         tempProduct.Name,
		PriceInCents: tempProduct.PriceInCents,
		Quantity:     tempProduct.Quantity,
	})
}
