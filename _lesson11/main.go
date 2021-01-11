package main

import (
	"fmt"
	"gopl.io/ch12/params"
	"log"
	"net/http"
)

func search(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Label      []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	data.MaxResults = 10
	if err := params.Unpack(r, &data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Search: %+v\n", data)
}

func main() {
	http.HandleFunc("/search", search)
	log.Fatal(http.ListenAndServe(":12345", nil))
}
