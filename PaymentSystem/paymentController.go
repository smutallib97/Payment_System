package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PaymentModel struct {
	Id        int     `json:"id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

var payments = make([]PaymentModel, 0)

func createPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newPayment PaymentModel
		err := json.NewDecoder(r.Body).Decode(&newPayment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newPayment.Id = len(payments) + 1
		newPayment.CreatedAt = time.Now().Format("2006-01-02T15:04:05")
		newPayment.UpdatedAt = time.Now().Format("2006-01-02T15:04:05")
		payments = append(payments, newPayment)
		fmt.Fprintln(w, newPayment)
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(payments)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func processPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		var updatedPayment PaymentModel
		err := json.NewDecoder(r.Body).Decode(&updatedPayment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i, payment := range payments {
			if payment.Id == updatedPayment.Id {
				updatedPayment.CreatedAt = payments[i].CreatedAt
				updatedPayment.UpdatedAt = time.Now().Format("2006-01-02T15:04:05")
				payments[i] = updatedPayment
				fmt.Fprintln(w, updatedPayment)
				w.WriteHeader(http.StatusOK)
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
func getPaymentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		paymentIDStr := r.URL.Query().Get("id")
		paymentID, err := strconv.Atoi(paymentIDStr)
		if err != nil {
			http.Error(w, "Invalid payment ID", http.StatusNotFound)
			return
		}
		var foundPayment *PaymentModel
		for _, payment := range payments {
			if payment.Id == paymentID {
				foundPayment = &payment
				break
			}
		}
		if foundPayment != nil {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Payment Status : %s", foundPayment.Status)
		} else {
			http.Error(w, "Payment not found", http.StatusNotFound)
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}
func getByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		paymentIDStr := r.URL.Query().Get("id")
		paymentID, err := strconv.Atoi(paymentIDStr)
		if err != nil {
			http.Error(w, "Invalid payment Id", http.StatusNotFound)
			return
		}
		var foundPayment *PaymentModel
		for _, payment := range payments {
			if payment.Id == paymentID {
				foundPayment = &payment
				w.Header().Set("Content-Type", "application/json")
				err := json.NewEncoder(w).Encode(payment)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				break
			}
		}
		if foundPayment != nil {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Payment not found", http.StatusNotFound)
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}
func main() {
	http.HandleFunc("/createPayment", createPayment)
	http.HandleFunc("/index", index)
	http.HandleFunc("/getById", getByID)
	http.HandleFunc("/processPayment", processPayment)
	http.HandleFunc("/getPaymentStatus", getPaymentStatus)
	log.Println("Listening...!!!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
