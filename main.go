package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	receiptStore = make(map[string]Receipt)
	mutex        sync.Mutex
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ReceiptResponse struct {
	ID string `json:"id"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/receipts/process", receiptHandler)
	r.HandleFunc("/receipts/{id}/points", pointsHandler)

	port := ":8080"
	fmt.Println("Server running on", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func calculatePoints(receipt Receipt) (int, error) {
	points := 0

	reg := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(reg.FindAllString(receipt.Retailer, -1))

	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil {
		if total == math.Floor(total) {
			points += 50
		} else if math.Mod(total, 0.25) == 0 {
			points += 25
		}
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
			if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}

	if purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		if purchaseDate.Day()%2 != 0 {
			points += 6
		}
	}

	if purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		hour := purchaseTime.Hour()
		if hour == 14 || (hour == 15 && purchaseTime.Minute() == 0) {
			points += 10
		}
	}
	return points, nil
}

func receiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt

	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	mutex.Lock()
	receiptStore[id] = receipt
	mutex.Unlock()

	response := ReceiptResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func pointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	mutex.Lock()
	receipt, ok := receiptStore[id]
	mutex.Unlock()

	if !ok {
		http.Error(w, "receipt not found", http.StatusNotFound)
		return
	}

	points, err := calculatePoints(receipt)
	if err != nil {
		http.Error(w, "unable to calculate points", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
