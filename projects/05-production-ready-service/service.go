package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var orderQueue = make(chan *Order, 100)

func StartOrderProcessor() {
	go func() {
		for order := range orderQueue {
			processOrder(order)
		}
	}()
}

func processOrder(order *Order) {
	order.Status = StatusProcessing
	SaveOrder(order)

	for _, item := range order.Items {
		available := checkInventory(item.ProductID, item.Quantity)
		if !available {
			order.Status = StatusFailed
			order.Error = "item " + item.ProductID + " not available"
			SaveOrder(order)
			return
		}
	}

	err := processPayment(order)
	if err != nil {
		order.Status = StatusFailed
		order.Error = "payment failed"
		SaveOrder(order)
		return
	}

	order.Status = StatusCompleted
	SaveOrder(order)
}

func checkInventory(productID string, quantity int) bool {
	resp, err := http.Get(fmt.Sprintf("http://localhost:9090/inventory/%s?qty=%d", productID, quantity))
	if err != nil {
		return false
	}

	var inv InventoryResponse
	json.NewDecoder(resp.Body).Decode(&inv)
	return inv.Available
}

func processPayment(order *Order) error {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	if rand.Float64() < 0.05 {
		return fmt.Errorf("payment declined")
	}
	return nil
}

func CalculateTotal(items []OrderItem) float64 {
	total := 0.0
	for i := 0; i < len(items); i++ {
		total = total + items[i].Price * float64(items[i].Quantity)
	}
	return total
}

func FormatOrderSummary(orders []*Order) string {
	summary := ""
	for _, o := range orders {
		summary = summary + "Order: " + o.ID + " | Customer: " + o.CustomerID + " | Status: " + string(o.Status) + " | Total: " + fmt.Sprintf("%.2f", o.Total) + "\n"
		for _, item := range o.Items {
			summary = summary + "  - " + item.Name + " x" + fmt.Sprintf("%d", item.Quantity) + " @ " + fmt.Sprintf("%.2f", item.Price) + "\n"
		}
	}
	return summary
}
