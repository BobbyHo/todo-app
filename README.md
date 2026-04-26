## System Information

The project is tested based on the following:

- Go version: 1.26.1
- Ubuntu: 18.04
- Docker: 19.03.1

## How to build

The following steps assume that a Go environment (Go 1.26.1+) has been set up on the system.

a) Check out the code into the GOPATH (i.e. /go/src).

b) Change to the root directory of the project.

c) Run `make`.
- This compiles both the server and client binaries.
- The binaries are stored in `dist/usr/bin`.
```bash
make
```


## How to run the todoserver
### Using Docker

The following steps assume that `make` has completed successfully.

a) From the project's root directory, run the following command:

```bash
make docker
```

- You can replace `APP_VERSION` with a different version.

b) Once step a) is done, a Docker image is built. You can verify it by running the `docker images` command.
```bash
docker images -a | grep todo
todoapp                                    1.0                             c566404ac937        26 minutes ago      825MB       21.5MB
```

c) Change to the `deploy` directory and run Docker Compose.
```bash
cd deploy
docker-compose up
```
Note: The default listening port for the todoserver is `12345`, which is defined in `todoserverconf.json`. The Docker Compose command starts both the Redis container and the todoserver container. If you change `APP_VERSION` in the Makefile, also update the app image version in `docker-compose.yml`.

## How to run todoclient
The todoclient binary is located in the `dist/usr/bin` directory. You can run the binary directly from there or copy it to another location.

The todoclient CLI options are:
```bash
./todoclient --help
This Todo application allows a user to register or delete a user and add, update, delete, or list Todo tasks.

Usage:
  todoclient [flags]
  todoclient [command]

Available Commands:
  addtodo      add a todo task
  adduser      add a user to Todo App
  deletetodo   delete a todo task
  deleteuser   delete a user
  help         Help about any command
  listalltodos list all todos
  listentodo   listen to todo commands
  listtodo     list a todo task
  updatetodo   update an existing todo task

Flags:
      --config string        config file (default is todoclient.yaml)
  -d, --description string   Task Description (default "-1-1")
  -e, --duedate string       Task Due Date (default "-1-1")
  -h, --help                 help for todoclient
  -p, --progress int32       Task Progress 0: Todo; 1: In Progress; 2: Done (default -1)
  -s, --serverip string      Server Address (default "127.0.0.1:12345")
  -t, --task string          Task Title
  -u, --username string      User Name

Use "todoclient [command] --help" for more information about a command.
```
By default, the client tries to connect to `127.0.0.1:12345`. If the server is running in another VM or on another machine, specify the server address by using the `-s` flag. For example:
```bash
./todoclient adduser -u "test-user-A" -s "10.0.0.47:12345"
```

The `-u` (`--username`) flag is required for all commands.

### adduser
The `adduser` command must be called before using other todo commands. For example:
```bash
./todoclient adduser -u "test-user-A"
Added user successfully: 
{
    "userId": "test-user-A"
}
```

### deleteuser
```bash
./todoclient deleteuser -u "test-user-B"
Deleted user successfully:
{
    "userId": "test-user-B"
}
```

### addtodo
```bash
./todoclient addtodo -u "test-user-A" -t "test-user-A-task-1"
Added todo successfully:
{
    "userId": "test-user-A",
    "taskId": "test-user-A-task-1",
    "progress": 0
}
```

### updatetodo
```bash
./todoclient updatetodo -u "test-user-A" -t "test-user-A-task-1" -d "test-user-A-example" -e "2020-06-08" -p 1
Updated todo successfully: 
{
    "userId": "test-user-A",
    "taskId": "test-user-A-task-1",
    "description": "test-user-A-example",
    "duedate": "2020-06-08",
    "progress": 1
}
```

### listtodo
```bash
./todoclient listtodo -u "test-user-A" -t "test-user-A-task-1"
{
    "userId": "test-user-A",
    "taskId": "test-user-A-task-1",
    "description": "test-user-A-example",
    "duedate": "2020-06-08",
    "progress": 1
}
```

### listalltodos
```bash
./todoclient listalltodos -u "test-user-A" 
{
    "userId": "test-user-A",
    "taskId": "test-user-A-task-1",
    "description": "test-user-A-example",
    "duedate": "2020-06-08",
    "progress": 1
}
```

### listentodo
The `listentodo` command blocks and waits for the server to send todo commands entered by another todoclient. To exit, press Ctrl+C.

```bash
./todoclient listentodo -u "test-user-1"
Received Listen Reply: action:"test-user-A list all todos"
Received Listen Reply: action:"test-user-A update todo task: test-user-A-task-1"
```

### deletetodo
```bash
./todoclient deletetodo -u "test-user-A" -t "test-user-A-task-1" 
2020/06/08 12:09:48 Setting state to UNDEFIND
Deleted todo successfully:
{
    "userId": "test-user-A",
    "taskId": "test-user-A-task-1"
}
```

## Project Structure
### api/v1/proto
gRPC interface definitions (`todo.proto`) and the generated protobuf and gRPC files (`todo.pb.go` and `todo_grpc.pb.go`).

### cmd
The `cmd` directory contains the server (`cmd/server`) and client (`cmd/client`) entry points.

#### client/cmd
This directory contains CLI implementations for each command supported by todoclient. The `root.go` file defines flags and variables shared by child commands such as `addtodo` and `updatetodo`. The `todo-client.go` file defines the gRPC client code that interacts with the gRPC server.

#### server
The `svc.go` file defines functions that handle server startup and shutdown. The `svcHandler.go` file defines the gRPC server functions that receive commands from todoclient. It also interacts with the database and message queue for storing todo tasks and publishing todo events.

### internal
#### models
The `models` package contains interface definitions for core services, database access (`dbstore`), and message publishing (`msgpublish.go`).

#### db
Redis is the database for storing TODO records. The `todoRecord.go` file implements the interface methods defined in the `models` package.

##### data model
Each user (`TodoUser`) contains a collection of todo tasks (`TodoData`). The `userId` identifies the user record, and the `taskId` identifies an individual todo task inside that user record.

#### msgqueue
#### msgredis
Redis is also used as a message broker for pub/sub. The `msgredis` package implements methods defined in the `models` package.

### config
Sample configuration file for the todoserver.

### deploy
Dockerfile and Docker Compose files.
