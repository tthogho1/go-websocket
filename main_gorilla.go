package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
)

/*type Client struct {
	Name string
	Conn *websocket.Conn
}

type Hub struct {
	Clients map[string]*Client
}*/

func main_bk() {

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	/*	hub := Hub{
		Clients: make(map[string]*Client),
	}*/

	// return html page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// start websocket server
	r.HandleFunc("/ws", websocketHandler)
	http.ListenAndServe(":8080", nil)
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			return
		}

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			return
		}
	}
}

func wsHandler(hub *Hub) http.HandlerFunc {
	// output log
	log.Println("new connection")
	return func(w http.ResponseWriter, r *http.Request) {
		/*
			//upGrader := websocket.Upgrader{}
			//ws, err := upGrader.Upgrade(w, r, nil)
			if err != nil {
				http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
				return
			}
			defer ws.Close()

			name := r.URL.Query().Get("name")
			if name == "" {
				http.Error(w, "Please provide a name parameter", http.StatusBadRequest)
				return
			}

			client := &Client{
				Name: name,
				Conn: ws,
			}
			hub.Clients[name] = client
			defer delete(hub.Clients, name)

			var saveError error = nil
			for {
				messageType, message, err := ws.ReadMessage()
				if err != nil {
					saveError = err
					break
				}

				switch messageType {
				case websocket.TextMessage:
					if len(message) > 0 && message[0] == '@' {
						to := message[1:]
						if client, ok := hub.Clients[string(to)]; ok {
							err = client.Conn.WriteMessage(messageType, message)
							if err != nil {
								saveError = err
								break
							}
						}
					}
				}

				if saveError != nil {
					break
				}
			}
			// output err to log
			log.Println("read message error:", saveError.Error())
		*/
	}
}
