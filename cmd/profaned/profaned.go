package main

import (
	"fmt"
	"net/http"
	"os"
)

type ProfanityFinder interface {
	Find(*http.Request) (bool, error)
}

func main() {
	finder, err := newProfanityFinder()

	if err != nil {
		fmt.Println("failed to create handler", err.Error())
		os.Exit(1)
	}

	h := &handler{
		profanity: finder,
	}

	http.HandleFunc("/", h.handle)
	http.ListenAndServe(":80", nil)
}

type handler struct {
	profanity ProfanityFinder
}

func (h *handler) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil { // no data -> no profanity found
		w.WriteHeader(http.StatusOK)
		return
	}

	found, err := h.profanity.Find(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if found {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("profanity detected"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
