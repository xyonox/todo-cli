package main

import (
	"bufio"
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	list := map[int]string{}

	list[0] = "Elisabeth oder so"

	for {
		fmt.Print("Befehl eingeben (list, add, done, quit): ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				_, err := fmt.Fprintln(os.Stderr, "Eingabe konnte nicht gelesen werden:", err)
				if err != nil {
					fmt.Println("Eingabe konnte nicht gelesen werden.")
					continue
				}
			}
			continue
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			fmt.Println("Kein Befehl eingegeben.")
			continue
		}

		switch input {
		case "list":
			listTodos(list)
		case "add":
			fmt.Print("Neue Todo eingeben: 	")
			if !scanner.Scan() {
				_, err := fmt.Fprintln(os.Stderr, "Eingabe konnte nicht gelesen werden:", err)
				if err != nil {
					fmt.Println("Eingabe konnte nicht gelesen werden.")
					continue
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
				_, err := fmt.Fprintln(os.Stderr, "Eingabe konnte nicht gelesen werden:", err)
				if err != nil {
					fmt.Println("Eingabe konnte nicht gelesen werden.")
					continue
				}
			}

			secondInput := strings.TrimSpace(scanner.Text())

			index, err := strconv.Atoi(secondInput)
			if err != nil {
				fmt.Println("Could not convert to int.")
				continue
			}

			if index < 1 || index > len(list) {
				fmt.Println("Index out of range.")
				continue
			}

			index--

			doneTodo(&list, index)
		case "quit":
			return
		}

		fmt.Println("cmd:", input)
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
