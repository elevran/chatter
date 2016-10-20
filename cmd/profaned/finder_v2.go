// +build v2

package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type mapFinder struct {
	profanities map[string]struct{}
}

func newProfanityFinder() (*mapFinder, error) {
	mf := &mapFinder{
		profanities: make(map[string]struct{}),
	}
	err := mf.load("./list_of_dirty_naughty_obscene_and_otherwise_bad_words.txt")
	return mf, nil
}

func (mf *mapFinder) Find(input io.Reader) (bool, error) {
	scanner := bufio.NewScanner(input)
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

func (mf *mapFinder) load(path string) error {
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
