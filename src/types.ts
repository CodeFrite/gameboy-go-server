const padStart = (str: string, length: number, pad: string) => {
  const len = str.length;
  if (len >= length) return str;
  else return pad.repeat(length - len) + str;
};

export class uint8 {
  value: number;
  constructor(value: number) {
    if (this.valid(value)) this.value = value;
    else throw new Error(`Invalid uint8 value: ${value}`);
  }
  get() {
    return this.value;
  }
  set(value: number) {
    if (this.valid(value)) this.value = value;
    else throw new Error(`Invalid uint8 value: ${value}`);
  }
  valid(value: number) {
    return this.value >= 0 && this.value < 2 ** 8;
  }
  toHex(): string {
    return padStart(this.value.toString(16).toUpperCase(), 2, "0");
  }

  toBit(): string {
    return padStart(this.value.toString(2).toUpperCase(), 8, "0");
  }
  toStr(): string {
    return this.value.toString();
  }
}

export class uint16 {
  value: number;
  constructor(value: number) {
    if (this.valid(value)) this.value = value;
    else throw new Error(`Invalid uint16 value: ${value}`);
  }
  get() {
    return this.value;
  }
  set(value: number) {
    if (this.valid(value)) this.value = value;
    else throw new Error(`Invalid uint16 value: ${value}`);
  }
  valid(value: number) {
    return this.value >= 0 && this.value < 2 ** 16;
  }
  toHex(): string {
    return padStart(this.value.toString(16).toUpperCase(), 4, "0");
  }

  toBit(): string {
    return padStart(this.value.toString(2).toUpperCase(), 16, "0");
  }

  toStr(): string {
    return this.value.toString();
  }
}

export type MemoryWrite = {
  name: string;
  address: number;
  data: number[];
};

export type Instruction = {
  Mnemonic: string; // instruction mnemonic
  Bytes: number; // number of bytes the instruction takes
  Cycles: number[]; // number of cycles the instruction takes to execute. The first element is the number of cycles the instruction takes when the condition is met, the second element is the number of cycles the instruction takes when the condition is not met (see RETZ for example)
  Operands: Operand[]; // instruction operands used as function arguments
  Immediate: boolean; // is the operand an immediate value or should it be fetched from memory
  Flags: Flags; // cpu flags affected by the instruction
};

export type Operand = {
  Name: string; // operand name: register, n8/n16 (immediate unsigned value), e8 (immediate signed value), a8/a16 (memory location)
  Bytes: number; // number of bytes the operand takes (optional)
  Immediate: boolean; // is the operand an immediate value or should it be fetched from memory
  Increment: boolean; // should the program counter be incremented after fetching the operand
  Decrement: boolean; // should the program counter be decreased after fetching the operand
};

export type Flags = {
  Z: string; // Zero flag: set if the result is zero (all bits are 0)
  N: string; // Subtract flag: set if the instruction is a subtraction
  H: string; // Half carry flag: set if there was a carry from bit 3 (result is 0x0F)
  C: string; // Carry flag: set if there was a carry from bit 7 (result is 0xFF)
};

export type CpuState = {
  // Special registers
  PC: uint16; // Program Counter
  SP: uint16; // Stack Pointer
  A: uint8; // Accumulator
  F: uint8; // Flags: Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4)
  Z: boolean; // Zero flag
  N: boolean; // Subtraction flag
  H: boolean; // Half Carry flag
  C: boolean; // Carry flag

  // 16-bits general purpose registers
  BC: uint16;
  DE: uint16;
  HL: uint16;

  // Instruction
  PREFIXED: boolean; // Is the current instruction prefixed with 0xCB
  IR: uint8; // Instruction Register
  OPERAND_VALUE: uint16; // Current operand fetched from memory (this register doesn't physically exist in the CPU)

  // Interrupts
  IE: uint8; // Interrupt Enable
  IME: boolean; // interrupt master enable
  HALTED: boolean; // is the CPU halted
};

export type GameboyState = {
  PREV_CPU_STATE: CpuState;
  CURR_CPU_STATE: CpuState;
  INSTR: Instruction;
  MEMORY_WRITES: MemoryWrite[];
};

export enum MessageType {
  // client -> server
  ConnectionMessageType = 0, // initial client http connection even before upgrading to websocket
  CommandMessageType = 10, // sends a command to the server (message.data: 'step', 'reset')
  // server -> client
  InitialMemoryMapsMessageType = 50, // notifies the client of the initial memory maps (message.data: MemoryWrite)
  GameboyStateMessageType = 70, // notifies the client of the current gameboy state (message.data: GameboyState)
  ErrorMessageType = 90, // notifies the client of an error (message.data: string)
}

export type Message = {
  type: MessageType;
  data: unknown;
};

export type ConnectionMessage = Message & {
  type: MessageType.ConnectionMessageType;
  data: undefined;
};

export type CommandMessage = Message & {
  type: MessageType.CommandMessageType;
  data: "step" | "run";
};

export type InitialMemoryMapsMessage = Message & {
  type: MessageType.InitialMemoryMapsMessageType;
  data: MemoryWrite[];
};

export type GameboyStateMessage = Message & {
  type: MessageType.GameboyStateMessageType;
  data: GameboyState;
};

export type ErrorMessageType = Message & {
  type: MessageType.ErrorMessageType;
  data: string;
};

export type ServerMessageTypes =
  | ConnectionMessage
  | CommandMessage
  | InitialMemoryMapsMessage
  | GameboyStateMessage
  | ErrorMessageType;
