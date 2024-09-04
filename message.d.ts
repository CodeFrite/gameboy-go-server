declare enum MessageType {
  // client -> server
  ConnectionMessageType = 0, // initial client http connection even before upgrading to websocket
  CommandMessageType = 10, // sends a command to the server (message.data: 'step', 'reset')
  // server -> client
  ErrorMessageType = 60,
  GameboyStateMessageType = 70, // notifies the client of the current gameboy state (message.data: GameboyState)
  CPUStateMessageType = 71, // notifies the client of the current CPU state (message.data: CpuState)
  MemoryStateMessageType = 72, // notifies the client of a memory write (message.data: MemoryWrite) // TODO: should be an array []MemoryWrite that retrieves only changes and not whole memory state
}

interface Message {
  type: MessageType;
  data: any; // TODO: Consider using a more specific type or a union type for better type safety.
}
