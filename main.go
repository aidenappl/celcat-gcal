package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aidenappl/celcat-gcal/routers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// AWS Healthcheck Handler
	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET", "POST", "OPTIONS")

	// Get Calendar Method
	r.HandleFunc("/getCalendar", routers.HandleGetCalendar).Methods(http.MethodGet)

	// Host & Serve the API

	// Launch API Listener
	fmt.Printf("✅ APLB Celcat Gcal Service running on port %s\n", "8000")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Origin", "Authorization", "Accept", "X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
