package cmd

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients   = make(map[*websocket.Conn]string) 
	broadcast = make(chan string)
	mutex     = &sync.Mutex{}
)


func StartServer(port string) {
	InitializeDatabase()

	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server started on ws://localhost%s/ws\n", serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	var authData struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}

	if err := conn.ReadJSON(&authData); err != nil {
		log.Println("Failed to read authentication data:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Authentication failed"))
		return
	}

	if err := AuthenticateUser(authData.Username, authData.Token); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid token"))
		return
	}

	mutex.Lock()
	clients[conn] = authData.Username
	mutex.Unlock()

	log.Printf("%s connected", authData.Username)
	conn.WriteMessage(websocket.TextMessage, []byte("Welcome, "+authData.Username))

	// Send last 10 messages as chat history
	history := GetLastMessages(10)
	for _, msg := range history {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			log.Printf("%s disconnected", authData.Username)
			break
		}
		fullMsg := fmt.Sprintf("%s: %s", authData.Username, string(msg))
		SaveMessage(authData.Username, string(msg))
		broadcast <- fullMsg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			client.WriteMessage(websocket.TextMessage, []byte(msg))
		}
		mutex.Unlock()
	}
}
