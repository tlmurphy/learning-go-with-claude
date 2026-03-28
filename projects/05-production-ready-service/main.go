package main

import (
	"fmt"
	"log"
	"net/http"
)

var version = "0.1.0"

func main() {
	fmt.Printf("Order Service v%s starting...\n", version)

	StartOrderProcessor()

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleListOrders(w, r)
		case "POST":
			handleCreateOrder(w, r)
		default:
			http.Error(w, "method not allowed", 405)
		}
	})

	http.HandleFunc("/orders/get", handleGetOrder)
	http.HandleFunc("/orders/delete", handleDeleteOrder)
	http.HandleFunc("/orders/summary", handleOrderSummary)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
