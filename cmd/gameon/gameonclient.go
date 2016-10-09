package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Doer provides an interface to execute HTTP requests using an "injected" client.
// The use of Doer is to allow customization of the (http/https) client used (e.g., for testing).
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Signer provides an interface for signing request properties to validate their authenticity to the server
type Signer interface {
	Sign(*http.Request) error
}

//-----------------------------------------------------------------------------
// the default Doer is an http.Client customized based on the target protocol
func newDefaultDoer(target url.URL) Doer {
	if target.String() != "" && target.Scheme == "https" {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // TODO: GameOn! server certificate is self-signed?
				},
			},
		}
	}
	return &http.Client{}
}

//-----------------------------------------------------------------------------
type nullSigner struct{}

func (ns nullSigner) Sign(_ *http.Request) error {
	return nil
}

//-----------------------------------------------------------------------------
type gameonSigner struct {
	id     string // GameOn! identifier
	secret string // GameOn! shared secret
}

func newSigner(id, secret string) *gameonSigner {
	return &gameonSigner{
		id:     id,
		secret: secret,
	}
}

// TODO: ignoring errors is not good for you...
func (gs gameonSigner) Sign(req *http.Request) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}

	body, _ := copyBodyOf(req)
	hash := sha256.New()
	_, _ = hash.Write([]byte(body))
	bodyHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	ts := time.Now().UTC().Format(time.RFC3339Nano) // TODO: should we account for time shifts?

	tokens := make([]string, 0)
	tokens = append(tokens, []string{gs.id, ts}...)
	if len(body) > 0 {
		tokens = append(tokens, bodyHash)
	}

	msgauth := hmac.New(sha256.New, []byte(gs.secret))
	_, _ = msgauth.Write([]byte(strings.Join(tokens, "")))
	signature := base64.StdEncoding.EncodeToString(msgauth.Sum(nil))

	// Set the required headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json,text/plain")
	req.Header.Set("gameon-id", gs.id)
	req.Header.Set("gameon-date", ts)
	req.Header.Set("gameon-sig-body", bodyHash)
	req.Header.Set("gameon-signature", signature)
	return nil
}

func copyBodyOf(req *http.Request) (string, error) {
	if req.Body == nil {
		return "", nil
	}

	// read the content and restore the io.ReadCloser
	body, _ := ioutil.ReadAll(req.Body) // TODO: is it safe to ignore the error?
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return string(body), nil
}
