// +build !v2 !v3

package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

type regexSearch struct {
	re *regexp.Regexp
}

func newProfanityFinder() (*regexSearch, error) {
	var err error
	s := &regexSearch{}

	s.re, err = regexp.Compile("(boogers|snot|poop|shucks|argh)") // our list of nasty words
	return s, err
}

func (s *regexSearch) Find(r *http.Request) (bool, error) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return true, err // be pessimistic on errors
	}

	return s.re.Match(body), nil
}
