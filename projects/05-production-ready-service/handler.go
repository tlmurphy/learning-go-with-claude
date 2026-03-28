package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	order := &Order{
		ID:         fmt.Sprintf("ord_%d", time.Now().UnixNano()),
		CustomerID: req.CustomerID,
		Items:      req.Items,
		Status:     StatusPending,
		Total:      CalculateTotal(req.Items),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	SaveOrder(order)

	go func() {
		orderQueue <- order
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(order)
}

func handleGetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", 400)
		return
	}

	order, err := GetOrder(id)
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func handleListOrders(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")

	var result []*Order
	if customerID != "" {
		result = GetOrdersByCustomer(customerID)
	} else {
		result = GetAllOrders()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleOrderSummary(w http.ResponseWriter, r *http.Request) {
	allOrders := GetAllOrders()
	summary := FormatOrderSummary(allOrders)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(summary))
}

func handleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, err := GetOrder(id)
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}

	delete(orders, id)
	w.WriteHeader(200)
	fmt.Fprint(w, "deleted")
}
