package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Create a command-line to-do list where you can add, remove, and list tasks.
func addTasks(toDoList *[]string) {
	scanner := bufio.NewScanner(os.Stdin)

	for {

		var task string
		fmt.Print("Enter Task to do:")
		scanner.Scan()
		task = scanner.Text()
		*toDoList = append(*toDoList, task)

		var choice string
		fmt.Println("Add more tasks? Y/N:")
		fmt.Scan(&choice)
		if strings.ToLower(choice) != "y" {
			fmt.Println("Exiting...")
			break
		}

	}
	fmt.Println("Your ToDo List:", *toDoList)
}

func removeTasks(toDoList *[]string) {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		if len(*toDoList) == 0 {
			fmt.Println("No items to remove!")
			break
		}
		fmt.Println("Your ToDo List:", *toDoList)
		fmt.Println("Enter task number to remove (starting from 1):")
		scanner.Scan()
		var userInput string
		userInput = scanner.Text()

		input, err := strconv.Atoi(userInput)
		if err != nil || input < 1 || input > len(*toDoList) {
			fmt.Println("Invalid input, please enter a valid task number:")
			continue
		}

		*toDoList = append((*toDoList)[:input-1], (*toDoList)[input:]...)
		fmt.Println("Updated toDo List:", *toDoList)

		var choice string
		fmt.Println("Remove more tasks? Y/N:")
		fmt.Scan(&choice)
		if strings.ToLower(choice) != "y" {
			fmt.Println("Exiting...")
			break
		}

	}
	fmt.Println("New ToDo List:", *toDoList)

}

func viewTasks(toDoList *[]string) {
	if len(*toDoList) == 0 {
		fmt.Println("No tasks in the to-do list")
		return
	}
	for key, task := range *toDoList {
		//fmt.Println(key, task)
		fmt.Printf("%d. %s\n", key+1, task)
	}

}

func main() {
	toDoList := []string{}

	for {
		fmt.Println("--------------")
		fmt.Println("1.ADD TASKS")
		fmt.Println("2. DELETE TASKS")
		fmt.Println("3. VIEW TASK")
		fmt.Println("4. EXIT")
		fmt.Print("What action do you want to perform?:")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:

			addTasks(&toDoList)

		case 2:

			removeTasks(&toDoList)

		case 3:
			fmt.Println("Viewing Tasks")
			viewTasks(&toDoList)

		case 4:
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid Choice!")
		}
	}

}
