package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type TodoStatus struct {
	key         string
	translation string
}

var (
	TodoStatusDone    = TodoStatus{key: "TODO-STATUS-DONE", translation: "Done"}
	TodoStatusWorking = TodoStatus{key: "TODO-STATUS-WORKING", translation: "Working on"}
	TodoStatusNone    = TodoStatus{key: "TODO-STATUS-NONE", translation: "Not started"}
)

const (
	CommandList = "list, add, working, done, remove, quit"
)

// listTodos prints the todos in the map
func listTodos(todos map[int]string) {
	fmt.Println("TODOs:")
	fmt.Println("----------------------------------------\n")
	category := map[string]string{}
	for key, value := range todos {
		splitted := strings.SplitN(value, "!DATA-INFORMATION!", 2)

		if len(splitted) != 2 {
			category["unknown"] = fmt.Sprintf("%v[%v] %v\n", category["unknown"], key+1, value)
			continue
		}

		status := splitted[0]
		text := splitted[1]

		switch status {
		case TodoStatusNone.key, TodoStatusWorking.key, TodoStatusDone.key:
			category[status] = fmt.Sprintf("%v[%v] %v\n", category[status], key+1, text)
		default:
			category["unknown"] = fmt.Sprintf("%v[%v] %v\n", category["unknown"], key+1, value)
		}
	}

	printed := false

	if category[TodoStatusNone.key] != "" {
		fmt.Printf("%v:\n%v\n", TodoStatusNone.translation, category[TodoStatusNone.key])
		printed = true
	}
	if category[TodoStatusWorking.key] != "" {
		fmt.Printf("%v:\n%v\n", TodoStatusWorking.translation, category[TodoStatusWorking.key])
		printed = true
	}
	if category[TodoStatusDone.key] != "" {
		fmt.Printf("%v:\n%v\n", TodoStatusDone.translation, category[TodoStatusDone.key])
		printed = true
	}
	if !printed {
		fmt.Println("Keine Todos vorhanden.")
	}
	fmt.Println("----------------------------------------")
	if category["unkown"] != "" {
		fmt.Printf("Nicht zu einer Kategorien rückzuführen: \n%v\n", category["unkown"])
	}
}

// removeTodoFromList removes a todo item from the list at the specified index and reorders the remaining items.
func removeTodoFromList(list *map[int]string, index int) {
	newList := map[int]string{}
	i := 0
	for key, value := range *list {
		if key == index {
			continue
		} else {
			newList[i] = value
			i++
		}
	}
	*list = newList
}

func doneTodo(list *map[int]string, index int) {

	splitted := strings.Split((*list)[index], "!DATA-INFORMATION!")

	var todo string

	switch len(splitted) {
	case 1:
		todo = splitted[0]
	case 2:
		todo = splitted[1]
	}

	(*list)[index] = TodoStatusDone.key + "!DATA-INFORMATION!" + todo
}

func workingTodo(list *map[int]string, index int) {

	splitted := strings.Split((*list)[index], "!DATA-INFORMATION!")

	var todo string

	switch len(splitted) {
	case 1:
		todo = splitted[0]
	case 2:
		todo = splitted[1]
	}

	(*list)[index] = TodoStatusWorking.key + "!DATA-INFORMATION!" + todo
}

func addTodo(list *map[int]string, todo string) {
	(*list)[len(*list)] = TodoStatusNone.key + "!DATA-INFORMATION!" + todo
}

// save saves the list to a file
func save(list map[int]string) (err error) {
	// Marshal to encode to JSON and Unmarshal to decode from JSON
	out, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		fmt.Println("Could not marshal the todo list.")
		fmt.Printf("The Error: %v\n", err)
		return err
	}

	err = os.WriteFile("todos.json", out, 0644)
	if err != nil {
		fmt.Println("Error while writing the file.", err)
		return err
	}
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	list := map[int]string{}

	file, err := os.ReadFile("todos.json")
	if err != nil {
		// When the file does not exist
		addTodo(&list, "Einkaufen")
		addTodo(&list, "Auto Mieten")
		addTodo(&list, "Lego Set zu ende bauen")
	} else {
		err = json.Unmarshal(file, &list)
		if err != nil {
			fmt.Println("unexpect error: ", err)
			return
		}
	}

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		if err := save(list); err != nil {
			fmt.Println("Error while saving the file.", err)
		}
		os.Exit(0)
	}()

	// Main loop
	for {
		fmt.Printf("Befehl eingeben (%v): ", CommandList)
		// Trying to read input from the user
		if !scanner.Scan() {
			// If there is an error, print it to the user
			if err := scanner.Err(); err != nil {
				_, err := fmt.Fprintln(os.Stderr, "Eingabe konnte nicht gelesen werden:", err)
				if err != nil {
					fmt.Println("Eingabe konnte nicht gelesen werden.")
				}
			}
			// repeat the loop
			continue
		}
		// Trim the input (remove whitespaces)
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			fmt.Println("Kein Befehl eingegeben.")
			// repeat the loop
			continue
		}

		// check which command the user typed
		switch input {
		case "list":
			listTodos(list)
		case "add":
			fmt.Print("Neue Todo eingeben: ")
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					_, err := fmt.Fprintln(os.Stderr, "Eingabe konnte nicht gelesen werden:", err)
					if err != nil {
						fmt.Println("Eingabe konnte nicht gelesen werden.")
						// repeat the loop
						continue
					}
				}
			}
			secondInput := strings.TrimSpace(scanner.Text())
			if secondInput == "" {
				fmt.Println("Keine neue Todo eingeben")
			}
			addTodo(&list, secondInput)
		case "done":
			fmt.Print("Index eingeben: ")
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					_, err := fmt.Fprintln(os.Stderr, "Eingabe konnte nicht gelesen werden:", err)
					if err != nil {
						fmt.Println("Eingabe konnte nicht gelesen werden.")
						// repeat the loop
						continue
					}
				}
			}

			secondInput := strings.TrimSpace(scanner.Text())

			index, err := strconv.Atoi(secondInput)
			if err != nil {
				fmt.Println("Konnte den Index nicht zur Zahl nicht konvertieren.")
				// repeat the loop
				continue
			}

			if index < 1 || index > len(list) {
				fmt.Println("Index ist nicht in der Liste.")
				// repeat the loop
				continue
			}

			index--

			doneTodo(&list, index)
		case "remove":
			fmt.Print("Index eingeben: ")
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					_, err := fmt.Fprintln(os.Stderr, "Eingabe konnte nicht gelesen werden:", err)
					if err != nil {
						fmt.Println("Eingabe konnte nicht gelesen werden.")
						// repeat the loop
						continue
					}
				}
			}

			secondInput := strings.TrimSpace(scanner.Text())

			index, err := strconv.Atoi(secondInput)
			if err != nil {
				fmt.Println("Konnte den Index nicht zur Zahl nicht konvertieren.")
				// repeat the loop
				continue
			}

			if index < 1 || index > len(list) {
				fmt.Println("Index ist nicht in der Liste.")
				// repeat the loop
				continue
			}

			index--

			removeTodoFromList(&list, index)
		case "working":
			fmt.Print("Index eingeben: ")
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					_, err := fmt.Fprintln(os.Stderr, "Eingabe konnte nicht gelesen werden:", err)
					if err != nil {
						fmt.Println("Eingabe konnte nicht gelesen werden.")
						// repeat the loop
						continue
					}
				}
			}

			secondInput := strings.TrimSpace(scanner.Text())

			index, err := strconv.Atoi(secondInput)
			if err != nil {
				fmt.Println("Konnte den Index nicht zur Zahl nicht konvertieren.")
				// repeat the loop
				continue
			}

			if index < 1 || index > len(list) {
				fmt.Println("Index ist nicht in der Liste.")
				// repeat the loop
				continue
			}

			index--

			workingTodo(&list, index)

		case "edit":
			fmt.Print("Index eingeben: ")
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					_, err := fmt.Fprintln(os.Stderr, "Eingabe konnte nicht gelesen werden:", err)
					if err != nil {
						fmt.Println("Eingabe konnte nicht gelesen werden.")
						// repeat the loop
						continue
					}
				}
			}

			secondInput := strings.TrimSpace(scanner.Text())

			index, err := strconv.Atoi(secondInput)
			if err != nil {
				fmt.Println("Konnte den Index nicht zur Zahl nicht konvertieren.")
				// repeat the loop
				continue
			}

			if index < 1 || index > len(list) {
				fmt.Println("Index ist nicht in der Liste.")
				// repeat the loop
				continue
			}

			index--

			splitted := strings.SplitN(list[index], "!DATA-INFORMATION!", 2)

			fmt.Printf("Bearbeitung eingeben (%v): ", splitted[1])
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					_, err := fmt.Fprintln(os.Stderr, "Eingabe konnte nicht gelesen werden:", err)
					if err != nil {
						fmt.Println("Eingabe konnte nicht gelesen werden.")
						// repeat the loop
						continue
					}
				}
			}

			secondInput = strings.TrimSpace(scanner.Text())
			list[index] = splitted[0] + "!DATA-INFORMATION!" + secondInput

		case "quit", "exit":
			err := save(list)
			if err != nil {
				return
			}

			// End the loop
			return
		// When the user types an unknown command, print him the commands
		default:
			fmt.Printf("Unbekannter Befehl. Bekannte Befehle: %v\n", CommandList)
		}
	}
}
