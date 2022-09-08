package main

import (
	"net/http"

	"github.com/aidenappl/celcat-gcal/routers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/getCalendar", routers.HandleGetCalendar).Methods(http.MethodGet)

	http.ListenAndServe(":8000", r)
}
