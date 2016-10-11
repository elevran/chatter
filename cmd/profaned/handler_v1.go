// +build !v2 !v3

package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

type handler struct {
	re *regexp.Regexp
}

func newHandler() (*handler, error) {
	var err error
	h := &handler{}

	h.re, err = regexp.Compile("(boogers|snot|poop|shucks|argh)") // our list of nasty words
	return h, err
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if h.re.Match(body) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("profanity detected"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
