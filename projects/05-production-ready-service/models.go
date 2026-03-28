package main

import "time"

type OrderStatus string

var (
	StatusPending    OrderStatus = "pending"
	StatusProcessing OrderStatus = "processing"
	StatusCompleted  OrderStatus = "completed"
	StatusFailed     OrderStatus = "failed"
)

type Order struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	Status     OrderStatus `json:"status"`
	Total      float64     `json:"total"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Error      string      `json:"error,omitempty"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CreateOrderRequest struct {
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
}

type InventoryResponse struct {
	ProductID string `json:"product_id"`
	Available bool   `json:"available"`
	Stock     int    `json:"stock"`
}
