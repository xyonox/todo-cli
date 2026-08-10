# Todo CLI

This project is a small command-line todo application written in Go. It was
created as a learning project to understand two basic Go topics:

- working with maps
- encoding and decoding JSON files

The application stores todos locally in `todos.json` and provides a simple
interactive command loop.

## Requirements

- Go 1.25 or newer

## Running the application

From the project directory, run:

```bash
go run .
```

The program then waits for a command. The prompts and status messages in the
current implementation are in German, while this documentation is in English.

## Available commands

| Command | Description |
| --- | --- |
| `list` | Displays all current todos and their indexes. |
| `add` | Asks for a new todo and adds it to the list. |
| `done` | Asks for an index and removes the corresponding todo. |
| `quit` | Prints the current list, saves it to `todos.json`, and exits. |

Example session:

```text
list
add
Read about Go maps
done
2
quit
```

## How it works

### 1. Loading the todos

At startup, the program creates an empty map with integer keys and string
values:

```go
list := map[int]string{}
```

It then tries to read `todos.json` with `os.ReadFile`. If the file exists,
`json.Unmarshal` converts the JSON data into the Go map. If the file does not
exist, a small default list is created in memory.

The JSON file represents the map like this:

```json
{
  "0": "Buy groceries",
  "1": "Rent a car"
}
```

JSON object keys are strings, but Go converts them to integer map keys when
unmarshalling into `map[int]string`.

### 2. Reading commands

The program uses `bufio.Scanner` to read input from the terminal. Each input is
trimmed with `strings.TrimSpace` and processed by a `switch` statement.

The loop continues until the user enters `quit`. Invalid commands and input
errors are reported in the terminal.

### 3. Adding a todo

When `add` is selected, the entered text is stored in the map using the current
map length as the key:

```go
list[len(list)] = todo
```

This works for the normal add flow, but it is intentionally simple. After a
todo is removed, map indexes can be rebuilt, so this is not a permanent unique
ID system.

### 4. Completing a todo

When `done` is selected, the user enters the displayed index. The input is
converted from a string to an integer with `strconv.Atoi`.

The `doneTodo` function receives a pointer to the map. It creates a new map,
copies every todo except the selected one, and assigns the new map back to the
original variable. This also rebuilds the indexes from zero.

### 5. Saving the todos

The list is only saved when `quit` is entered. `json.MarshalIndent` converts
the Go map to readable JSON, and `os.WriteFile` writes the result to
`todos.json` with file permissions `0644`.

## Project structure

```text
.
├── main.go     # CLI logic, map operations, and JSON handling
├── todos.json  # Persisted todo data
├── go.mod      # Go module definition
└── README      # Project documentation
```

## What this project demonstrates

- declaring and modifying a `map[int]string`
- iterating over a map with `range`
- passing a map to a function
- passing a pointer to a map when replacing the map itself
- reading and writing files
- converting between Go values and JSON with `encoding/json`
- validating and converting terminal input
- handling basic errors in Go

## Current limitations

- Map iteration order is not guaranteed, so the display order may vary.
- A map is used as an indexed list, which is useful for learning but not ideal
  for a production todo application.
- The program does not support editing todos or marking them with a separate
  completed status.
