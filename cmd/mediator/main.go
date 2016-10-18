package main

import (
	"fmt"
	"net/http"
)

func main() {
	s := newServer()

	http.HandleFunc("/", s.handleHTTP)
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println(err)
	}
}
