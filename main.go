package main

import (
	"groupie-tracker/handler"
	"log"
	"net/http"
)

func main() {
	port := ":8080"
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/about/", handler.About)
	mux.HandleFunc("/dates/", handler.Dates)
	mux.HandleFunc("/locations/", handler.Locations)
	mux.HandleFunc("/location/", handler.Location)
	mux.HandleFunc("/date/", handler.Date)
	mux.HandleFunc("/relation/", handler.Relation)

	fs := http.FileServer(http.Dir("./statics/"))
	mux.Handle("/statics/", http.StripPrefix("/statics/", fs))

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Printf("Listening on http://localhost%v", port)
	server.ListenAndServe()
}
