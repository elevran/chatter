package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
)

// The subcommand interface defines methods to allow the main function to work with sub-command generically.
type subcommand interface {
	Keyword() string                               // sub-command keyword
	Usage() string                                 // usage text
	Parse([]string) error                          // parse command line arguments
	Process(Doer, Signer, io.Writer) (bool, error) // do work
	// Uses the given Doer, Signer and io.Writer (for tracing). if Doer or Signer are nil, will internally
	// create a default implementation.
	// TODO: Returns?
}

//-----------------------------------------------------------------------------
// The actual sub-commands defined by our gameon cli
//
// All sub-commands follow exactly the same implementation scheme so only
// the first command is documented.
//-----------------------------------------------------------------------------

//-- delete a room by its identifier
type deleteRoom struct {
	server url.URL // GameOn! server
	roomID string  // identifier of room to delete
}

// create a new delete sub-command instance, along with its flags set
func deleteRoomSubcommand(server url.URL, _ string) *deleteRoom {
	return &deleteRoom{
		server: server,
	}
}

// implement GoStringer, mostly for debugging/tracing in verbose mode
func (del *deleteRoom) GoString() string {
	return fmt.Sprintf("%s - roomid:%s server:%s", del.Keyword(), del.roomID, del.server.String())
}

// sub-command's keyword
func (del *deleteRoom) Keyword() string {
	return "delete"
}

// usage string
func (del *deleteRoom) Usage() string {
	return fmt.Sprintf("%s %s", del.Keyword(), "<room-identifier>")
}

// parse and validate sub-command and its flags
func (del *deleteRoom) Parse(argv []string) error {
	if len(argv) != 1 {
		return fmt.Errorf("%s", del.Usage())
	}

	del.roomID = argv[0]
	return nil
}

// Attempts to delete the room denoted by roomID from a Game On! server at server.
// The delete room functionality uses the given Doer. If Doer is nil, it'll create a default implementation.
// Flow is reported via the tracer, which can be set by the caller.
func (del *deleteRoom) Process(client Doer, auth Signer, trace io.Writer) (bool, error) {
	fmt.Fprintln(trace, location(), del, "Doer", client, "Signer", auth)

	// validate requirements
	if del.server.String() == "" || del.roomID == "" {
		return false, fmt.Errorf("%#v", del)
	}

	httpClient, authenticator := client, auth
	if httpClient == nil { // using default as a fallback, caller explicitly said it doesn't care
		httpClient = newDefaultDoer(del.server)
	}
	if authenticator == nil { // caller doesn't care, use minimal to get by...
		authenticator = nullSigner{}
	}

	u := del.server
	u.Path = "/map/v1/sites/" + del.roomID

	fmt.Fprintln(trace, "DELETE", u.String())
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		fmt.Fprintln(trace, location(), "NewRequest.Error", err.Error())
		return false, err
	}

	_ = authenticator.Sign(req) // TODO: handle error
	resp, err := httpClient.Do(req)
	defer resp.Body.Close() // TODO: even if err != nil?
	if err != nil {
		fmt.Fprintln(trace, location(), "DELETE.Error", err.Error())
		return false, err
	}

	fmt.Fprintln(trace, "Status", resp.StatusCode)
	switch resp.StatusCode {
	// XXX handle cases
	}

	/*defer r.Body.Close()
	b, e := ioutil.ReadAll(r.Body)
	if e == nil {
		body = string(b)
	}*/
	/* io.Copy(ioutil.Discard, resp.Body) or copy to trace?
	resp.Body.Close() */

	/*body, err := extractBody(resp)
	if err != nil {
		checkpoint(locus, fmt.Sprintf("Body.Error err=%s", err.Error()))
		return
	}
	checkpoint(locus, fmt.Sprintf("Status=%s", resp.Status))

	switch resp.StatusCode {
	case http.StatusNoContent:
		checkpoint(locus, "Deleted")
		return
	case http.StatusOK, http.StatusForbidden, http.StatusNotFound:
		checkpoint(locus, "Sigh. There is no use trying any more.")
		printResponseBody(locus, resp, body)
		stopTrying = true
		return
	default:
		err = RegError{fmt.Sprintf("Unhandled Status: %s", resp.Status)}
		checkpoint(locus, fmt.Sprintf("Unhandled Status=%s", resp.Status))
		printResponseBody(locus, resp, body)
		return
	}*/

	return true, nil // TODO figure out return values
}

//-- list registered rooms by owner and (optionally) name
type listRooms struct {
	server   url.URL // GameOn! server
	gameonID string  // GameOn! id (owner field in queries, e.g., email)
	roomName string  // room name, if multiple rooms have been registered
}

func listRoomsSubcommand(server url.URL, goid string) *listRooms {
	return &listRooms{
		server:   server,
		gameonID: goid,
	}
}

func (ls *listRooms) GoString() string {
	return fmt.Sprintf("%s - owner:%s name:%s server:%s", ls.Keyword(), ls.gameonID, ls.roomName,
		ls.server.String())
}

func (ls *listRooms) Keyword() string {
	return "list"
}

func (ls *listRooms) Usage() string {
	return fmt.Sprintf("%s %s", ls.Keyword(), "[room-name-to-match]")
}

func (ls *listRooms) Parse(argv []string) error {
	argc := len(argv)
	if argc > 2 {
		return fmt.Errorf("%s", ls.Usage())
	} else if argc == 1 {
		ls.roomName = argv[0]
	}
	return nil
}

func (ls *listRooms) Process(client Doer, _ Signer, trace io.Writer) (bool, error) {
	return true, nil
}

//-- register a new room
type registerRoom struct {
	server      url.URL  // GameOn! server
	gameonID    string   // GameOn! id
	roomName    string   // room name
	callbackURL *url.URL // room callback URL (must be accessible from GameOn! server)
}

func registerRoomSubcommand(server url.URL, goid string) *registerRoom {
	return &registerRoom{
		server:   server,
		gameonID: goid,
	}
}

func (reg *registerRoom) GoString() string {
	return fmt.Sprintf("%s - name:%s cb:%s server:%s", reg.Keyword(), reg.roomName, reg.callbackURL.String(),
		reg.server.String())
}

func (reg *registerRoom) Keyword() string {
	return "register"
}

func (reg *registerRoom) Usage() string {
	return fmt.Sprintf("%s %s %s", reg.Keyword(), "<room-name>", "<room-callback-URL>")
}

func (reg *registerRoom) Parse(argv []string) error {
	if len(argv) != 2 {
		return fmt.Errorf("%s", reg.Usage())
	}

	var err error
	reg.callbackURL, err = url.Parse(argv[1])

	if err != nil {
		return err
	}

	reg.roomName = argv[0]
	return nil
}

func (reg *registerRoom) Process(client Doer, auth Signer, trace io.Writer) (bool, error) {
	return true, nil
}

// current program location (file, line and function name)
func location() string {
	pc, file, line, _ := runtime.Caller(1)
	funcname := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s:%d %s", file, line, funcname)
}
