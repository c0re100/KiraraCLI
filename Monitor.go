package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

func echo(w http.ResponseWriter, r *http.Request) {
    c, _ := upgrader.Upgrade(w, r, nil)
    defer c.Close()
    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            break
        }
        fmt.Println(string(message))
    }
}

func main() {
    fmt.Println("Kirara SSR(5*) Probability Monitor")
	http.HandleFunc("/", echo)
    log.Fatal(http.ListenAndServe(":12345", nil))
}
