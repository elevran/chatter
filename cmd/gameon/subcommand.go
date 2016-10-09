package main

import (
	"fmt"
	"net/url"
)

// The subcommand interface defines methods to allow the main function to work with sub-command generically.
type subcommand interface {
	Keyword() string           // sub-command keyword
	Usage() string             // usage text
	Parse(argv []string) error // parse command line arguments
	Process(client Doer) error // do actual work, using the given Doer (if nil, creates a default one)
}

//-----------------------------------------------------------------------------
// The actual sub-commands defined by our gameon cli
//
// All sub-commands follow exactly the same implementation scheme so only
// the first command is documented.
//-----------------------------------------------------------------------------

//-- delete a room by its identifier
type deleteRoom struct {
	roomID string // identifier of room to delete
}

// create a new delete sub-command instance, along with its flags set
func deleteRoomSubcommand() *deleteRoom {
	return &deleteRoom{}
}

// implement GoStringer, mostly for debugging/tracing in verbose mode
func (del *deleteRoom) GoString() string {
	return fmt.Sprintf("%s - roomid:%s", del.Keyword(), del.roomID)
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

// execute the delete room functionality using the given Doer. If Doer is nil, Process will create a default
// implementation. The use of Doer is to allow customization of the (http/https) client used (e.g., for testing)
func (del *deleteRoom) Process(client Doer) error {
	return nil
}

//-- list registered rooms by owner and (optionally) name
type listRooms struct {
	name  string // room owner (e.g., email)
	owner string // room name, if multiple rooms have been registered
}

func listRoomsSubcommand() *listRooms {
	return &listRooms{}
}

func (ls *listRooms) GoString() string {
	return fmt.Sprintf("%s - owner:%s name:%s", ls.Keyword(), ls.owner, ls.name)
}

func (ls *listRooms) Keyword() string {
	return "list"
}

func (ls *listRooms) Usage() string {
	return fmt.Sprintf("%s %s %s", ls.Keyword(), "<room-owner>", "[room-name-to-match]")
}

func (ls *listRooms) Parse(argv []string) error {
	argc := len(argv)
	if argc < 1 || argc > 2 {
		return fmt.Errorf("%s", ls.Usage())
	}

	ls.owner = argv[0]
	if argc == 2 {
		ls.name = argv[1]
	}
	return nil
}

func (ls *listRooms) Process(client Doer) error {
	return nil
}

//-- register a new room
type registerRoom struct {
	name     string
	callback *url.URL
}

func registerRoomSubcommand() *registerRoom {
	return &registerRoom{}
}

func (reg *registerRoom) GoString() string {
	return fmt.Sprintf("%s - name: %s cb:%s %#v", reg.Keyword(), reg.name, reg.callback)
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
	reg.callback, err = url.Parse(argv[1])

	if err != nil {
		return err
	}

	reg.name = argv[0]
	return nil
}

func (reg *registerRoom) Process(client Doer) error {
	return nil
}
