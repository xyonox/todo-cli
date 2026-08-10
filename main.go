package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func listTodos(todos map[int]string) {
	fmt.Println("TODOs:")
	for key, value := range todos {
		fmt.Printf("[%v] %v \n", key+1, value)
	}
}

func doneTodo(list *map[int]string, index int) {
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

	// Main loop
	for {
		fmt.Print("Befehl eingeben (list, add, done, quit): ")
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
			fmt.Print("Neue Todo eingeben: 	")
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
			list[len(list)] = secondInput
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
				fmt.Println("Could not convert to int.")
				// repeat the loop
				continue
			}

			if index < 1 || index > len(list) {
				fmt.Println("Index out of range.")
				// repeat the loop
				continue
			}

			index--

			doneTodo(&list, index)
		case "quit":
			// Marshal to encode to JSON and Unmarshal to decode from JSON
			out, err := json.MarshalIndent(list, "", "  ")
			if err != nil {
				fmt.Println("Could not marshal the todo list.")
				fmt.Printf("The Error: %v\n", err)
				// repeat the loop
				continue
			}
			fmt.Println(string(out))

			err = os.WriteFile("todos.json", out, 0644)
			if err != nil {
				fmt.Println("Error while writing the file.", err)
				return
			}

			// End the loop
			return
		// When the user types an unknown command, print him the commands
		default:
			fmt.Println("Unbekannter Befehl. Bekannte Befehle: list, add, done, quit")
		}
	}
}
