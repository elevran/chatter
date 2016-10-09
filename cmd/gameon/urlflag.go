package main

import (
	"fmt"
	"net/url"
)

//-----------------------------------------------------------------------------
// urlFlag defines command line handling for URL values (e.g., GameOn! server)
type urlFlag struct {
	url *url.URL
}

// Url returns the flag value as a parsed URL
func (uf *urlFlag) Url() *url.URL {
	return uf.url
}

// Set implements the flag.Getter/flag.Var interface requirements
func (uf *urlFlag) Set(value string) error {
	if uf.url != nil {
		return fmt.Errorf("Value already set to %s", uf.url.String())
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("invalid URL: %s", err)
	}
	uf.url = parsed
	return nil
}

// String implements the flag.Getter/flag.Var interface requirements
func (uf *urlFlag) String() string {
	if uf.url == nil {
		return "<nil>"
	}
	return uf.url.String()
}

// Get implements the flag.Getter/flag.Var interface requirements
func (uf *urlFlag) Get() interface{} {
	return (*url.URL)(uf.url)
}
