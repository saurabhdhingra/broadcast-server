package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)


func ConnectClient(serverAddr, username, token string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", serverAddr), nil)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer conn.Close()


	authData := map[string]string{
		"username": username,
		"token":    token,
	}
	authBytes, _ := json.Marshal(authData)
	conn.WriteMessage(websocket.TextMessage, authBytes)


	_, authResponse, err := conn.ReadMessage()
	if err != nil || string(authResponse) == "Invalid token" {
		log.Println("Authentication failed")
		return
	}

	fmt.Println(string(authResponse)) 

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Disconnected from server")
				os.Exit(0)
			}
			fmt.Println("\n" + string(msg))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := conn.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}
}
