//
package main

import (
	"fmt"
	"net/http"
)

func main() {
	room := newChatRoom()

	http.HandleFunc("/hello", room.hello)
	http.HandleFunc("/goodbye", room.goodbye)
	http.HandleFunc("/room", room.message)
	err := http.ListenAndServe(":80", nil)

	if err != nil {
		fmt.Println(err)
	}
}
