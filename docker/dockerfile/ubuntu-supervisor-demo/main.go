package main

import (
	"fmt"
	"log"
	"net/http"
)

func _URI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.URL.String())
}

func main() {
	http.HandleFunc("/", _URI)
	if err := http.ListenAndServe(":7000", nil); err != nil {
		log.Fatal(err)
	}
}
