# Gameboy Server

This is a go server exposing the gameboy emulator as a websocket

context: rest-api
technology: websockets
endpoints: GET /gameboy

messages:

- user -> server:
  - GET /gameboy: instantiates a new gameboy instance for the new user connection
- server -> user:

# Endpoints

## GET /gameboy

Only one endpoint is available for this server. When hit by a user, it will instantiate a new gameboy and map it to the user's connection. It will then listen to incoming messages and react to it accordingly either by sending interrupts to the gameboy, actions to the emulator or by sending messages to the user like gameboy state notifications.

## Messages

### step (user->server)

The step message is sent by the user to the server to request the gameboy to execute a single instruction. The server will then send an interrupt to the gameboy to execute the next instruction.
