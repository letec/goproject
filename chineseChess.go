package main

import (
	"goproject/src/controller"
	"goproject/src/model"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message Define our message object
type Message struct {
	Oid     string `json:"Oid"`
	Code    string `json:"Code"`
	Info    string `json:"Info"`
	Message string `json:"message"`
}

// 初始化房间数据
func initHallInfo() {
	info := controller.HallConfig["chineseChess"]
	tableNumbers, _ := strconv.Atoi(info["tableNumbers"])
	for i := 1; i <= tableNumbers; i++ {
		index := strconv.Itoa(i)
		model.RedisSetHash("HALL_chineseHall_"+index, "PlayerA", "")
		model.RedisSetHash("HALL_chineseHall_"+index, "PlayerB", "")
		model.RedisSetHash("HALL_chineseHall_"+index, "PlayerAStatus", "")
		model.RedisSetHash("HALL_chineseHall_"+index, "PlayerBStatus", "")
		model.RedisSetHash("HALL_chineseHall_"+index, "Status", "0")
		model.RedisSetHash("HALL_chineseHall_"+index, "Step", "")
		model.RedisSetHash("HALL_chineseHall_"+index, "ChessPanel", "")
	}
}

func main() {
	// connect to redis server
	model.RedisConnect()

	// Hall data install
	initHallInfo()

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8082")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		switch msg.Code {
		case "chat":
			if msg.Info == "C" {

			} else if msg.Info == "P" {

			} else {

			}
			break
		}
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
