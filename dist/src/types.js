"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MessageType = exports.uint16 = exports.uint8 = void 0;
const padStart = (str, length, pad) => {
    const len = str.length;
    if (len >= length)
        return str;
    else
        return pad.repeat(length - len) + str;
};
class uint8 {
    constructor(value) {
        if (this.valid(value))
            this.value = value;
        else
            throw new Error(`Invalid uint8 value: ${value}`);
    }
    get() {
        return this.value;
    }
    set(value) {
        if (this.valid(value))
            this.value = value;
        else
            throw new Error(`Invalid uint8 value: ${value}`);
    }
    valid(value) {
        return this.value >= 0 && this.value < Math.pow(2, 8);
    }
    toHex() {
        return padStart(this.value.toString(16).toUpperCase(), 2, "0");
    }
    toBit() {
        return padStart(this.value.toString(2).toUpperCase(), 8, "0");
    }
    toStr() {
        return this.value.toString();
    }
}
exports.uint8 = uint8;
class uint16 {
    constructor(value) {
        if (this.valid(value))
            this.value = value;
        else
            throw new Error(`Invalid uint16 value: ${value}`);
    }
    get() {
        return this.value;
    }
    set(value) {
        if (this.valid(value))
            this.value = value;
        else
            throw new Error(`Invalid uint16 value: ${value}`);
    }
    valid(value) {
        return this.value >= 0 && this.value < Math.pow(2, 16);
    }
    toHex() {
        return padStart(this.value.toString(16).toUpperCase(), 4, "0");
    }
    toBit() {
        return padStart(this.value.toString(2).toUpperCase(), 16, "0");
    }
    toStr() {
        return this.value.toString();
    }
}
exports.uint16 = uint16;
var MessageType;
(function (MessageType) {
    // client -> server
    MessageType[MessageType["ConnectionMessageType"] = 0] = "ConnectionMessageType";
    MessageType[MessageType["CommandMessageType"] = 10] = "CommandMessageType";
    // server -> client
    MessageType[MessageType["InitialMemoryMapsMessageType"] = 50] = "InitialMemoryMapsMessageType";
    MessageType[MessageType["GameboyStateMessageType"] = 70] = "GameboyStateMessageType";
    MessageType[MessageType["ErrorMessageType"] = 90] = "ErrorMessageType";
})(MessageType || (exports.MessageType = MessageType = {}));
