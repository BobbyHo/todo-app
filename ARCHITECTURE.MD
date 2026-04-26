# Architecture

This todo app is a small client/server system built around a gRPC API. The server owns the todo business operations, Redis stores users and todo records, and Redis pub/sub broadcasts todo activity to clients that subscribe through a streaming gRPC call.

## Current Architecture

```text
todoclient CLI
    |
    | gRPC requests
    v
todoserver
    |
    | store/read users and todo items
    v
Redis data store

todoserver
    |
    | publish todo activity
    v
Redis pub/sub
    |
    | subscribed by server-side ListenTodos stream
    v
todoclient listentodo
```

## Main Components

### gRPC API

The API contract lives in `api/v1/proto/todo.proto`. It defines the todo service, request and response messages, todo state enum, and the server-streaming `ListenTodos` RPC.

The generated Go files are:

- `api/v1/proto/todo.pb.go` for protobuf message types.
- `api/v1/proto/todo_grpc.pb.go` for gRPC client and server bindings.

The gRPC service exposes operations to add/delete users, add/update/delete/list todo items, list all todos for a user, and listen for todo activity.

### CLI Client

The CLI client lives under `cmd/client`. It uses Cobra commands to turn command-line input into gRPC calls.

Examples of supported commands:

- `adduser`
- `deleteuser`
- `addtodo`
- `updatetodo`
- `deletetodo`
- `listtodo`
- `listalltodos`
- `listentodo`

The client is intentionally thin. It handles flags, request construction, JSON-style output, and connection setup, while the server owns validation and persistence.

### gRPC Server

The server lives under `cmd/server`.

- `svc.go` handles configuration, Redis connection setup, service wiring, and gRPC server startup/shutdown.
- `svcHandler.go` implements the RPC methods generated from the proto contract.
- `main.go` handles process startup, signal handling, and service lifecycle.

The server is the central coordination point. It authenticates whether a user exists before todo operations, maps protobuf messages to internal model structs, writes data through the store interface, and publishes activity messages through the message queue interface.

### Redis Store

Redis is the persistence layer for users and todo records. The implementation lives under `internal/db/redis`.

At a high level:

- Each user has a collection of todo items.
- A user ID identifies the user record.
- A task ID identifies a todo item inside a user's collection.

The rest of the app talks to Redis through interfaces in `internal/models`, which keeps storage concerns mostly isolated from the server handlers.

### Redis Pub/Sub

Redis is also used as the message broker for todo activity. The implementation lives under `internal/msgqueue/redis`.

When todo operations happen, the server publishes a short activity message to the todo command channel. Clients using `listentodo` call the `ListenTodos` streaming RPC; the server subscribes to Redis pub/sub and forwards received messages over the gRPC stream.

## Request Flow

For a write operation such as `addtodo`:

1. The CLI parses command flags.
2. The CLI sends a gRPC `AddTodo` request to the server.
3. The server checks that the user exists.
4. The server publishes a todo activity message to Redis pub/sub.
5. The server writes the todo item to Redis.
6. The server returns a gRPC response to the CLI.

For `listentodo`:

1. The CLI opens a server-streaming gRPC request.
2. The server subscribes to the Redis todo command channel.
3. Redis pub/sub delivers todo activity messages.
4. The server forwards each message to the CLI over the gRPC stream.

## Known Tradeoffs

Redis is currently both the database and message broker. That keeps the app simple, but it also couples persistence and event delivery to one dependency.

The server uses package-level service state, which is acceptable for a small sample app but makes testing and dependency injection harder as the app grows.

The Redis-backed tests behave like integration tests because they require a live Redis instance at `127.0.0.1:6379`.

The pub/sub path is best-effort event delivery. If a listener is offline, it will miss messages published while it is disconnected.
