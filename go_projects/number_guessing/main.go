package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func number_guesser() {
	for {
		luckyNum := rand.Intn(50)
		num_of_guesses := 0

		for {
			var input int
			fmt.Println("Enter a random number between 1 and 50:")
			fmt.Scan(&input)
			num_of_guesses++

			if input < 1 || input > 50 {
				fmt.Println("This number is out of range choose another one in the range")
				num_of_guesses--
				continue

			} else if input == luckyNum {
				fmt.Println("BINGO, you got it!")
				fmt.Println("Number of guesses:", num_of_guesses)
				break

			} else if input < luckyNum {
				fmt.Println("Oops your guess is LOW")
			} else if input > luckyNum {
				fmt.Println("Oh no your guess is HIGH")
			}
		}
		var choice string
		fmt.Println("Would you like to continue? Y/N:")
		fmt.Scan(&choice)
		if strings.ToLower(choice) != "y" {
			return
		}

	}
}
func main() {
	number_guesser()
}
