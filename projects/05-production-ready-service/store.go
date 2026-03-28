package main

import (
	"fmt"
	"time"
)

// Global order storage. Accessible from anywhere.
var orders = map[string]*Order{}

func SaveOrder(order *Order) {
	order.UpdatedAt = time.Now()
	orders[order.ID] = order
}

func GetOrder(id string) (*Order, error) {
	order := orders[id]
	if order == nil {
		return nil, fmt.Errorf("not found")
	}
	return order, nil
}

func GetAllOrders() []*Order {
	result := []*Order{}
	for _, o := range orders {
		result = append(result, o)
	}
	return result
}

func GetOrdersByCustomer(customerID string) []*Order {
	result := []*Order{}
	for _, o := range orders {
		if o.CustomerID == customerID {
			result = append(result, o)
		}
	}
	return result
}
