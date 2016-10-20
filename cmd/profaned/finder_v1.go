// +build !v2 !v3

package main

import (
	"io"
	"io/ioutil"
	"regexp"
)

type regexProfanityFinder struct {
	re *regexp.Regexp
}

func newProfanityFinder() (*regexProfanityFinder, error) {
	var err error
	s := &regexProfanityFinder{}

	s.re, err = regexp.Compile("(boogers|snot|poop|shucks|argh)") // our list of nasty words
	return s, err
}

func (s *regexProfanityFinder) Find(input io.Reader) (bool, error) {
	content, err := ioutil.ReadAll(input)

	if err != nil {
		return true, err // be pessimistic on errors
	}

	return s.re.Match(content), nil
}
