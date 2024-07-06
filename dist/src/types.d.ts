export declare class uint8 {
    value: number;
    constructor(value: number);
    get(): number;
    set(value: number): void;
    valid(value: number): boolean;
    toHex(): string;
    toBit(): string;
    toStr(): string;
}
export declare class uint16 {
    value: number;
    constructor(value: number);
    get(): number;
    set(value: number): void;
    valid(value: number): boolean;
    toHex(): string;
    toBit(): string;
    toStr(): string;
}
export type MemoryWrite = {
    name: string;
    address: number;
    value: number[];
};
export type Instruction = {
    Mnemonic: string;
    Bytes: number;
    Cycles: number[];
    Operands: Operand[];
    Immediate: boolean;
    Flags: Flags;
};
export type Operand = {
    Name: string;
    Bytes: number;
    Immediate: boolean;
    Increment: boolean;
    Decrement: boolean;
};
export type Flags = {
    Z: string;
    N: string;
    H: string;
    C: string;
};
export type CpuState = {
    PC: uint16;
    SP: uint16;
    A: uint8;
    F: uint8;
    Z: boolean;
    N: boolean;
    H: boolean;
    C: boolean;
    BC: uint16;
    DE: uint16;
    HL: uint16;
    PREFIXED: boolean;
    IR: uint8;
    OPERAND_VALUE: uint16;
    IE: uint8;
    IME: boolean;
    HALTED: boolean;
};
export type GameboyState = {
    PREV_CPU_STATE: CpuState;
    CURR_CPU_STATE: CpuState;
    INSTR: Instruction;
    MEMORY_WRITES: MemoryWrite[];
};
export declare enum MessageType {
    ConnectionMessageType = 0,// initial client http connection even before upgrading to websocket
    CommandMessageType = 10,// sends a command to the server (message.data: 'step', 'reset')
    InitialMemoryMapsMessageType = 50,// notifies the client of the initial memory maps (message.data: MemoryWrite)
    GameboyStateMessageType = 70,// notifies the client of the current gameboy state (message.data: GameboyState)
    ErrorMessageType = 90
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
export type ServerMessageTypes = ConnectionMessage | CommandMessage | InitialMemoryMapsMessage | GameboyStateMessage | ErrorMessageType;
