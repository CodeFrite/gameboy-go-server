package main

import "github.com/codefrite/gameboy-go/gameboy"

type MessageType uint

const (

	/* client -> server */

	// init
	ConnectionMessageType MessageType = 0 // initial client http connection even before upgrading to websocket

	// execution
	StepMessageType    MessageType = 10 // step the gameboy
	RunMessageType     MessageType = 11 // run the gameboy
	ResetMessageType   MessageType = 12 // restart the gameboy with the current rom
	LoadRomMessageType MessageType = 13 // load a given rom on the gameboy
	Unused_0           MessageType = 14 // unused (to leave some space in case i need to change to add an event related to the current series)
	SaveStateType      MessageType = 15 // save the current state of the gameboy
	LoadStateType      MessageType = 16 // load the saved state to the gameboy
	JoypadStateType    MessageType = 17 // joypad keypresses (message.data: JoypadState)
	Unused_2           MessageType = 18 // unused
	Unused_3           MessageType = 19 // unused

	// breakpoints
	AddBreakpointType         MessageType = 20 // add a breakpoint at a certain address on the gameboy
	RemoveBreakpointType      MessageType = 21 // remove a breakpoint at a certain address on the gameboy
	DisableBreakPointType     MessageType = 22 // disable a breakpoint at a certain address
	EnableBreakPointType      MessageType = 23 // enable a breakpoint at a certain address
	DisableAllBreakPointsType MessageType = 24 // disable all breakpoint (but does not delete them)
	EnableAllBreakPointsType  MessageType = 25 // enable all breakpoint

	/* server -> client */
	InitialMemoryMapsMessageType MessageType = 50 // notifies the client of the initial memory maps (message.data: MemoryWrite)
	CPUStateMessageType          MessageType = 70 // notifies the client of the current CPU state (message.data: CpuState)
	PPUStateMessageType          MessageType = 71 // notifies the client of the current PPU state (message.data: PpuState)
	APUStateMessageType          MessageType = 72 // notifies the client of the current APU state (message.data: ApuState)
	MemoryStateMessageType       MessageType = 73 // notifies the client of a memory write (message.data: MemoryWrite) // TODO! should be an array []MemoryWrite that retrieves only changes and not whole memory state
	ErrorMessageType             MessageType = 90 // notifies the client of an error (message.data: string)
)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"` // TODO! should be a map[string]interface{} instead of
}

type AddBreakpointMessage struct {
	Type MessageType `json:"type"`
	Data uint        `json:"data"`
}

type RemoveBreakpointMessage struct {
	Type MessageType `json:"type"`
	Data uint        `json:"data"`
}

func InitialMemoryMapsMessage(data *[]gameboy.MemoryWrite) *Message {
	message := &Message{}
	message.Type = InitialMemoryMapsMessageType
	message.Data = data
	return message
}

func CPUStateMessage(data *gameboy.CpuState) *Message {
	message := &Message{}
	message.Type = CPUStateMessageType
	message.Data = data
	return message
}

func PPUStateMessage(data *gameboy.PpuState) *Message {
	message := &Message{}
	message.Type = PPUStateMessageType
	message.Data = data
	return message
}

func APUStateMessage(data *gameboy.ApuState) *Message {
	message := &Message{}
	message.Type = APUStateMessageType
	message.Data = data
	return message
}

func MemoryStateMessage(data *[]gameboy.MemoryWrite) *Message {
	message := &Message{}
	message.Type = MemoryStateMessageType
	message.Data = data
	return message
}

func JoypadStateMessage(data *gameboy.JoypadState) *Message {
	message := &Message{}
	message.Type = JoypadStateType
	message.Data = data
	return message
}
