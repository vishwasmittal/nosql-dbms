# NoSQL Redis type database

## Features
### Basic Functions

- **Phase 1:** Implementation of a simple inmemory database that can get commands from remote clients
- **Phase 2:** Adding a persistence layer to store data in files which can help in cache eviction when data to be stored in memory is too large

### Advanced Features
- **Phase 3:** Topic wise storing (Implemting something like collections in MongoDB)
- **Phase 4:** Introduction of Auth and permissions

## Usage

### Server
```sh
$ go run main.go
```

### Client
```sh
$ go run main.go --connect 127.0.0.1
```