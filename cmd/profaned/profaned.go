package main

import (
	"io"
	"net/http"
)

type ProfanityFinder interface {
	Find(io.Reader) (bool, error)
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

	found, err := h.profanity.Find(r.Body)

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
