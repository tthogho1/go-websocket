package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

type Client struct {
	Name string
	Conn *websocket.Conn
}

type Hub struct {
	Clients map[string]*Client
}

var hub = Hub{
	Clients: make(map[string]*Client),
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/2", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index2.html")
	})

	http.Handle("/ws", websocket.Handler(msgHandler))

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func handleErr(err error) {
	switch err.(type) {
	case *os.PathError:
		log.Println("File Path Error:", err)
	default:
		log.Println("Unknown Error:", err)
	}
}

type Message struct {
	Type_   string `json:"type"`
	From    string `json:"name"`
	To      string `json:"to"`
	Message string `json:"message"`
}

func msgHandler(ws *websocket.Conn) {
	defer ws.Close()

	query := ws.Request().URL.Query()
	if query != nil {
		log.Print("connect from " + query.Get("name"))
	}

	defer func() {
		deleteName := query.Get("name")
		delete(hub.Clients, deleteName)
		log.Println("length of hub.Clients: ", len(hub.Clients))
		log.Println("WebSocket connection closed", deleteName)
	}()

	name := query.Get("name")
	client := &Client{
		Name: name,
		Conn: ws,
	}
	hub.Clients[name] = client

	decoder := json.NewDecoder(ws)
	for {
		//msg := ""
		//err := websocket.Message.Receive(ws, &msg)
		var msg map[string]string
		err := decoder.Decode(&msg)

		if err != nil {
			handleErr(err)
			log.Println("break from loop:", err)
			break
		}

		message := Message{
			Type_:   msg["type"],
			From:    msg["from"],
			To:      msg["to"],
			Message: msg["message"],
		}

		jsonBytes, err := json.Marshal(message)
		if err != nil {
			fmt.Println("encode failed:", err)
			return
		}

		to := msg["to"]
		for _, c := range hub.Clients {
			if c.Name == to {
				log.Printf("send message from %q to %q ", msg["from"], msg["to"])
				err = websocket.Message.Send(c.Conn, string(jsonBytes))
				if err != nil {
					handleErr(err)
				}
				break
			}
		}
	}
}
