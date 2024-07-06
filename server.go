package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/codefrite/gameboy-go/gameboy"
	"github.com/gorilla/websocket"
)

var (
	emulators = make(map[*websocket.Conn]*gameboy.Debugger)
	mutex     = &sync.Mutex{}
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024 * 8,
		WriteBufferSize: 1024 * 8,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Close the connection and remove it from the map
func cleanup(conn *websocket.Conn) {
	fmt.Printf("Cleaning up conn %v\n", emulators)
	conn.Close()
	mutex.Lock()
	delete(emulators, conn) // safe because no-op if conn is not in the map
	mutex.Unlock()
}

func sendMessage(conn *websocket.Conn, message *Message) {
	payload, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Err: Couldn't marshal message:", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, payload)
}

func handleWSMessage(w http.ResponseWriter, r *http.Request) {
	// upgrade the HTTP connection to a websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Err: Couldn't upgrade HTTP connection to ws:", err)
		return
	}
	// defer cleanup of the connection on error or exit when this function returns
	defer cleanup(conn)

	// since maps are not concurrent-safe, we need use a mutex to protect the map
	mutex.Lock()
	db := gameboy.NewDebugger()
	db.Init("tetris.gb")
	// return the initial state of the CPU
	emulators[conn] = db
	mutex.Unlock()

	// log new connection
	fmt.Println("New connection")

	// send the initial memory maps (name, address, data dump) to the client
	attachedMemories := db.GetAttachedMemories()
	sendMessage(conn, InitialMemoryMapsMessage(&attachedMemories))

	// send the initial state of the CPU to the client
	gameboyInitialState := db.Step()
	sendMessage(conn, GameboyStateMessage(gameboyInitialState))

	// now that the connection has been established, we can start listening for messages
	for {
		// read the incoming message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Err: Couldn't read incoming message", err)
			break
		}

		// Accept only text messages
		if messageType != websocket.TextMessage {
			fmt.Println("Err: Received non-text message")
			break
		}

		// Process the message
		switch string(p) {
		case "step":
			newState := emulators[conn].Step()
			sendMessage(conn, GameboyStateMessage(newState))
		case "run":
		default:
			fmt.Println("Err: Unknown message type", string(p))
			return
		}
	}

}

func route() {
	http.HandleFunc("/gameboy", handleWSMessage)
}

func main() {
	route()
	fmt.Println("Server running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
