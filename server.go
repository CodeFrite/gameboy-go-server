package main

/* server.go is the API server that serves the gameboy emulator to the client
 * - it is a websocket server that listens for incoming connections on port :8080
 * - it serves the endpoint /gameboy that the client connects to
 * - upon connection, a new instance of the gameboy emulator hosted on github.com/codefrite/gameboy-go/gameboy is created and mapped to the connection
 * - any following message from a client is then processed on the corresponding gameboy emulator, assuring that each client has its own instance of the emulator

 * Each connection should have its own instance of the gameboy as well as the cpu and ppu state channels. This is naturally achieved by the fact
 * that http.HandleFunc("/gameboy", handleWSMessage) will create a new goroutine for each connection.
 * TODO: Check if the map of emulators is still necessary, since each connection will run in its own goroutine

 */

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/codefrite/gameboy-go/gameboy"
	"github.com/gorilla/websocket"
)

/* package level variables declaration */
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

/* closes the connection and deletes it from the map */
func cleanup(conn *websocket.Conn) {
	fmt.Printf("Cleaning up conn %v\n", emulators)
	conn.Close()
	mutex.Lock()
	delete(emulators, conn) // safe because no-op if conn is not in the map
	mutex.Unlock()
}

/* sends a message to the client
 * conn is the websocket connection to send the message to
 * message is the message to be marshalled and sent through the connection (must be a text message)
 */
func sendMessage(conn *websocket.Conn, message *Message) {
	payload, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Err: Couldn't marshal message:", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, payload)
}

/* handler for the /gameboy route
 * - upgrades the HTTP connection to a websocket connection
 * - creates a new instance of the gameboy emulator and maps it to the connection
 * - sends the initial state of the CPU to the client
 * - then enters in a loop that listens for incoming messages from the client
 * - processes the messages and sends the updated state of the CPU to the client

 * Depending on the message content received from the client, the server will:
 * - step the CPU 1 cycle and send the updated state to the client
 * - run the CPU and send the updated state to the client (TODO: must decide how to handle the different intermediate states of the CPU and also process the users joypad inputs)
 */

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
	// create cpu & ppu state channels
	cpuStateChannel := make(chan *gameboy.CpuState)
	ppuStateChannel := make(chan *gameboy.PpuState)
	apuStateChannel := make(chan *gameboy.ApuState)
	memoryStateChannel := make(chan *[]gameboy.MemoryWrite)
	joypadStateChannel := make(chan *gameboy.JoypadState)

	db := gameboy.NewDebugger(cpuStateChannel, ppuStateChannel, apuStateChannel, memoryStateChannel, joypadStateChannel)
	db.LoadRom("tetris.gb")
	// return the initial state of the CPU
	emulators[conn] = db
	mutex.Unlock()

	// log new connection
	fmt.Println("New connection")

	// send the initial memory maps (name, address, data dump) to the client
	attachedMemories := db.GetAttachedMemories()
	sendMessage(conn, InitialMemoryMapsMessage(&attachedMemories))

	// send the initial state of the CPU to the client
	go func() {
		db.Step()
		cpuInitialState := <-cpuStateChannel
		sendMessage(conn, CPUStateMessage(cpuInitialState))
		ppuInitialState := <-ppuStateChannel
		sendMessage(conn, PPUStateMessage(ppuInitialState))
		apuInitialState := <-apuStateChannel
		sendMessage(conn, APUStateMessage(apuInitialState))
		memoryInitialState := <-memoryStateChannel
		sendMessage(conn, MemoryStateMessage(memoryInitialState))
	}()

	fmt.Println("Initial state sent ... listening to incoming messages")

	// now that the connection has been established, we can start listening for messages
	for {
		// read the incoming message
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Err: Couldn't read incoming message", err)
			break
		}
		fmt.Println("received msg", p)

		// unmarshall the message
		var msg Message = Message{}
		json.Unmarshal(p, &msg)
		fmt.Println("")
		fmt.Println("> received msg", p)
		fmt.Println("--------------")
		fmt.Println("- msg.Type", msg.Type)
		fmt.Println("- msg.Data", msg.Data)

		// Process the message
		switch msg.Type {
		// gameboy related messages
		case StepMessageType:
			fmt.Println("+ step request")
			emulators[conn].Step()
			cpuState := <-cpuStateChannel
			sendMessage(conn, CPUStateMessage(cpuState))
			ppuState := <-ppuStateChannel
			sendMessage(conn, PPUStateMessage(ppuState))
			apuState := <-apuStateChannel
			sendMessage(conn, APUStateMessage(apuState))
			memoryState := <-memoryStateChannel
			sendMessage(conn, MemoryStateMessage(memoryState))
		case RunMessageType:
			fmt.Println("+ received run request")
			go func() {
				emulators[conn].Run()
				for {
					select {
					case cpuState := <-cpuStateChannel:
						sendMessage(conn, CPUStateMessage(cpuState))
					case ppuState := <-ppuStateChannel:
						sendMessage(conn, PPUStateMessage(ppuState))
					case apuState := <-apuStateChannel:
						sendMessage(conn, APUStateMessage(apuState))
					case memoryState := <-memoryStateChannel:
						sendMessage(conn, MemoryStateMessage(memoryState))
					}
				}
			}()
		case JoypadStateType:
			fmt.Println("+ received joypad request")
			var joypadState gameboy.JoypadState
			json.Unmarshal(p, &joypadState)
			// forward the joypad state to the emulator
			joypadStateChannel <- &joypadState
		// debugger related messages
		case AddBreakpointType:
			var message AddBreakpointMessage
			json.Unmarshal(p, &message)
			fmt.Printf("+ received add breakpoint request: 0x%04X\n", message.Data)
			db.AddBreakPoint(uint16(message.Data))
		case RemoveBreakpointType:
			var message RemoveBreakpointMessage
			json.Unmarshal(p, &message)
			fmt.Printf("+ received delete breakpoint request: 0x%04X\n", message.Data)
			db.RemoveBreakPoint(uint16(message.Data))
		default:
			fmt.Println("Err: Unknown message type", string(p))
			return
		}
	}
}

/* registers the handler for the different routes
 * '/gameboy' is the route that the client will connect to from gameboy-go.codefrite.dev
 */
func route() {
	http.HandleFunc("/gameboy", handleWSMessage)
}

/* setup the API routes, serves and listen on port :8080 */
func main() {
	route()
	fmt.Println("Server running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
