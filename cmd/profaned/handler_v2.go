// +build v2

package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"
)

type mapFind struct {
	profanities map[string]struct{}
}

func newProfanityFinder() (*mapFind, error) {
	mf := &mapFind{
		profanities: make(map[string]struct{}),
	}
	err := mf.load("./list_of_dirty_naughty_obscene_and_otherwise_bad_words.txt")
	return mf, nil
}

func (mf *mapFind) Find(r *http.Request) (bool, error) {
	scanner := bufio.NewScanner(r.Body)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		_, found := mf.profanities[word]
		if found {
			return true, nil
		}
	}

	return false, scanner.Err()
}

func (mf *mapFind) load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		mf.profanities[scanner.Text()] = struct{}{}
	}

	return scanner.Err()
}

func variations(word string) []string {
	// add all possible variations to the map, including common change techniques (e.g., leet speak, see
	// http://www.robertecker.com/hp/research/leet-converter.php)
	return []string{strings.ToLower(word)} // demo - only add the word itself with no variations
}
