package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/elevran/chatter/pkg/gameon"
)

type Client struct {
	httpClient *http.Client
	serverURL  string
}

func newClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (cl *Client) doHello(hello *gameon.Hello) (*gameon.MessageCollection, error) {
	return cl.doRequest("/hello", hello.UserID, hello)
}

func (cl *Client) doGoodbye(goodbye *gameon.Goodbye) (*gameon.MessageCollection, error) {
	return cl.doRequest("/goodbye", goodbye.UserID, goodbye)
}

func (cl *Client) doCommand(command *gameon.Command) (*gameon.MessageCollection, error) {
	return cl.doRequest("/room", command.UserID, command)
}

func (cl *Client) doRequest(path string, userID string, body interface{}) (*gameon.MessageCollection, error) {
	url := cl.serverURL + path

	reqBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reqBuf := bytes.NewBuffer(reqBytes)

	req, err := http.NewRequest("POST", url, reqBuf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(gameon.UserIDHeader, userID)

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var msgs gameon.MessageCollection
	err = json.Unmarshal(respBytes, &msgs)
	if err != nil {
		return nil, err
	}

	return &msgs, nil
}
