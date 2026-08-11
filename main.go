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

const (
	TODO_STATUS_DONE    = "TODO-STATUS-DONE"
	TODO_STATUS_WORKUNG = "TODO-STATUS-WORKING"
	TODO_STATUS_NONE    = "TODO-STATUS-NONE"
	COMMAND_LIST        = "list, add, working, done, remove, quit"
)

// listTodos prints the todos in the map
func listTodos(todos map[int]string) {
	fmt.Println("TODOs:")
	for key, value := range todos {
		fmt.Printf("[%v] %v \n", key+1, value)
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

	(*list)[index] = TODO_STATUS_DONE + "!DATA-INFORMATION!" + todo
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

	(*list)[index] = TODO_STATUS_WORKUNG + "!DATA-INFORMATION!" + todo
}

func getTodoStatus(index int) string {
	return ""
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
		list[0] = "Einkaufen"
		list[1] = "Auto Mieten"
		list[2] = "Lego Set zu ende bauen"
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
		fmt.Printf("Befehl eingeben (%v): ", COMMAND_LIST)
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
			fmt.Println(secondInput)
			list[len(list)] = TODO_STATUS_NONE + "!DATA-INFORMATION!" + secondInput
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

			fmt.Printf("Bearbeitung eingeben (%v): ", list[index])
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
			list[index] = secondInput

		case "quit", "exit":
			err := save(list)
			if err != nil {
				return
			}

			// End the loop
			return
		// When the user types an unknown command, print him the commands
		default:
			fmt.Printf("Unbekannter Befehl. Bekannte Befehle: %v\n", COMMAND_LIST)
		}
	}
}
