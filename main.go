package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func msgHandler(ws *websocket.Conn) {
	defer ws.Close()

	query := ws.Request().URL.Query()
	if query != nil {
		log.Print("connect from " + query.Get("name"))
	}

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
			log.Fatalln(err)
		}

		to := msg["to"]
		for _, c := range hub.Clients {
			if c.Name == to {
				err = websocket.Message.Send(c.Conn, fmt.Sprintf(`%q から %q というメッセージを受け取りました。`, msg["from"], msg["message"]))
				if err != nil {
					log.Fatalln(err)
				}

				break
			}
		}
	}

}
