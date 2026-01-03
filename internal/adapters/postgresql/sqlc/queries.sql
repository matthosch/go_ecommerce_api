-- name: ListProducts :many
SELECT
  * 
FROM 
  products;

-- name: FindProductByID :one
SELECT
  *
FROM
  products
WHERE
  id = $1;

-- name: CreateOrder :one
INSERT INTO orders (
  customer_id
) VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
  order_id,
  product_id,
  quantity,
  price_cents
) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateProductQuantity :one
UPDATE products
SET quantity = $2
WHERE id = $1
RETURNING *;

-- name: GetOrderByID :one
SELECT
  id,
  customer_id,
  created_at
FROM
  orders
WHERE
  id = $1;

-- name: GetProductsByOrderID :many
SELECT
  p.id,
  p.name,
  oi.quantity,
  oi.price_cents
FROM
  order_items oi
JOIN
  products p ON oi.product_id = p.id
WHERE
  oi.order_id = $1;