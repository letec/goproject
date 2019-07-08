package main

import (
	"goproject/src/controller"
	"goproject/src/model"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// HallConfig 房间配置
var HallConfig map[string]string

func init() {
	HallConfig = map[string]string{"tableNumbers": "35"}
}

// Message Define our message object
type Message struct {
	Oid     string `json:"Oid"`
	Code    string `json:"Code"`
	Info    string `json:"Info"`
	Message string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 初始化房间数据
func initHallInfo() {
	tableNumbers, _ := strconv.Atoi(HallConfig["tableNumbers"])
	conn := model.Mongo.Copy()
	defer conn.Close()
	conn.DB("online").C("HALL_ChineseChess_Tables").DropCollection()
	conn.DB("online").C("HALL_ChineseChess_Users").DropCollection()
	for i := 1; i <= tableNumbers; i++ {
		index := strconv.Itoa(i)
		t := controller.Table{
			ID:            index,
			PlayerAStatus: "0",
			PlayerBStatus: "0",
			PlayerCStatus: "0",
			Status:        "0",
		}
		conn.DB("online").C("HALL_ChineseChess_Tables").Insert(t)
	}
}

func main() {
	// connect to redis server
	model.MongoDBConnection()

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
			ws.Close()
			delete(clients, ws)
		}
		ws.SetCloseHandler(func(code int, text string) error {
			ws.Close()
			delete(clients, ws)
			return nil
		})

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
	}
}
