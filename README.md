## System Information

The project is tested based on the following:

- Go version: 1.14.4
- Ubuntu: 18.04
- Docker: 19.03.1

## How to build?

The following steps assume that Go environment (golang 1.14.1+) has been setup in the system. 

a) check out the code into the GOPATH (i.e. /go/src)

b) cd to the root directory of the project

c) run "make" 
- this will compile both the server and client binaries
- the target exec files are stored in dist/usr/bin
```bash
make
```


## how to run the todoserver?
### Using Docker Container

The following steps assume that "make" has completed successfully. 

a) From the project's root diectory, run the following command:

```bash
make docker
```

- you can replace "APP_VERSION" with a different verion 

b) Once step a) is done, a docker image is built. You could verify it by running the "docker images" command.
```bash
docker images -a | grep todo
todoapp                                    1.0                             c566404ac937        26 minutes ago      825MB       21.5MB
```

c) cd deploy and run docker-compose
```bash
cd deploy
docker-compose up
```
Note: The default listening port for the todoserver is "12345" and it is defined in the todoserverconf.json file. The docker-compose command brings up a redis container and the todoserver container. If you changed the APP_VERSIon in the Makefile, please also update the app image version in docker-compose.yml.

## how to run todoclient?
The todoclient binary is located in the dist/usr/bin directory. You could run the binary directly from there or copy it to another location. 

The todoclient CLI options are:
```bash
./todoclient --help
This Todo application allows a user to register or delete a user, add/update/delete/list Todo tasks

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
By default, the client tries to connect to "127.0.0.1:12345". If the server is running in another VM or machine, you could specify the server IP addres by using the -s flag. For example:
```bash
./todoclient adduser -u "test-user-A" -s "10.0.0.47:12345"
```

Also, -u (--username) is required for all the commands.

### adduser
The adduser command must be called first before using other todo commands. An example of adduser command is the following:
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
The listentodo command will block and wait for the server to send todo commands (entered by another todoclient) to it. To exit, press Ctrl+C.

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
### app/v1/proto
gRPC interface defintions (todo.proto) and the auto-generated client and server stubs (todo.pb.go)

### cmd
The cmd directory contains the server (cmd/server) and client (cmd/client) main execution codes.

#### client/cmd
This directory contains CLI implementations for each commands supported by the the todoclient. The root.go file defines flags and variables that are required by each child commands (i.e. addtodo, updatetodo and etc). todo-client.go defines the gRPC client stub that interacts with the gRPC server.

#### server
The file svc.go defines functions that handles server start up and stop. svcHandler.go defines the gRPC server functions that receives commands from the todoclient. It also interacts with the database and message queues for storing todo tasks and publishing todo events.

### internal
#### models
The models package contains interface definitions for core_services (i.e. the Todo Service), database access methods (dbstore) and message publishing (msgpublish.go)

#### db
At the moment, redis is the database for storing the TODO records. todoRecord.go implements the interface methods defined in the models package. 

##### data model
Each user (TodoUser)contains a collection of todo tasks (TodoData). The key to access the record is the userId and inside a user record, taskId is the key to retrieve individual todo task.

#### msgqueue
#### msgredis
Redis is also used as a message broker for pubsub. The msgredis implements methods defined in the models package.

### config
sample configuration file for the todoserver

### deploy
Dockerfile and Docker-Compose files


