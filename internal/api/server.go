package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/asafdavid23/endoflife-datastore/internal/models"
)

func SendNotifications(products []models.Product) {
	// Notifier microservice URL
	notifierURL := "http://endoflife-notifier:8080/api/notify"

	payload, err := json.Marshal(products)
	if err != nil {
		log.Printf("Error marshaling products: %v", err)
		return
	}

	resp, err := http.Post(notifierURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error sending notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Notifier microservice responded with status: %v", resp.Status)
	} else {
		log.Println("Notifications sent successfully.")
	}
}

func TestNotifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var products []models.Product
	if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Notify the user about the products

}
