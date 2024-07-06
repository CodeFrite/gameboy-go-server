package main

import "github.com/codefrite/gameboy-go/gameboy"

type MessageType uint

const (
	// client -> server
	ConnectionMessageType MessageType = 0  // initial client http connection even before upgrading to websocket
	CommandMessageType    MessageType = 10 // sends a command to the server (message.data: 'step', 'reset')
	// server -> client
	ErrorMessageType        MessageType = 60
	GameboyStateMessageType MessageType = 70 // notifies the client of the current gameboy state (message.data: GameboyState)
	CPUStateMessageType     MessageType = 71 // notifies the client of the current CPU state (message.data: CpuState)
	MemoryStateMessageType  MessageType = 72 // notifies the client of a memory write (message.data: MemoryWrite) // TODO! should be an array []MemoryWrite that retrieves only changes and not whole memory state
)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"` // TODO! should be a map[string]interface{} instead of
}

func newGameboyStateMessage(data *gameboy.GameboyState) *Message {
	message := &Message{}
	message.Type = GameboyStateMessageType
	message.Data = data
	return message
}
