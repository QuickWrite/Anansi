package main

import (
	"log"
	"net/http"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := r.PathValue("data")

	w.Write([]byte(data))
	log.Printf("Received the data: \"%s\"", data)
}

func main() {
	http.HandleFunc("/{data}", viewHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
