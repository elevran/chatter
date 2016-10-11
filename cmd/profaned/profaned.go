package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	h, err := newHandler()

	if err != nil {
		fmt.Println("failed to create handler", err.Error())
		os.Exit(1)
	}

	http.HandleFunc("/", h.handle)
	http.ListenAndServe(":80", nil)
}
