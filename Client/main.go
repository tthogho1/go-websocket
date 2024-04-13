package main

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

// WebSocket Client Sample
func main() {
	log.SetFlags(log.Lmicroseconds)

	// WebSocket Dial
	ws, dialErr := websocket.Dial("ws://localhost:3000/ws?name=client2", "", "http://localhost:3000/")
	if dialErr != nil {
		log.Fatal(dialErr)
	}
	defer ws.Close()

	// Send Message
	sendRestMsg(ws, "Hello")

	// Receive Message Ligic
	var recvMsg string
	for {
		recvErr := websocket.Message.Receive(ws, &recvMsg)
		if recvErr != nil {
			log.Fatal(recvErr)
			break
		}
		log.Println("Receive : " + recvMsg + ", from Server")
	}
}

// Send Message Ligic
type Message struct {
	Type_   string `json:"type"`
	From    string `json:"name"`
	To      string `json:"to"`
	Message string `json:"message"`
}

func sendRestMsg(ws *websocket.Conn, msg string) {

	message := Message{
		Type_:   "SDP",
		From:    "client2",
		To:      "client1",
		Message: msg,
	}

	jsonBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("encode failed:", err)
		return
	}

	sendErr := websocket.Message.Send(ws, string(jsonBytes))
	if sendErr != nil {
		log.Fatal(sendErr)
	}

	log.Println("Send : " + msg + ", to Server")
}
