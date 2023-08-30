package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Notification struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Date      string `json:"date"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

var notifications = make([]Notification, 0)

func sendPaymentNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newNotification Notification
		err := json.NewDecoder(r.Body).Decode(&newNotification)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newNotification.Id = len(notifications) + 1
		newNotification.CreatedAt = time.Now().Format("2006-01-02T15:04:05")
		newNotification.UpdatedAt = time.Now().Format("2006-01-02T15:04:05")
		notifications = append(notifications, newNotification)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, notifications)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
func getNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("content-type", "application/josn")
		err := json.NewEncoder(w).Encode(notifications)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/notification", sendPaymentNotification)
	http.HandleFunc("/getNotifications", getNotifications)
	log.Println("Listening from notification controller...!!!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
