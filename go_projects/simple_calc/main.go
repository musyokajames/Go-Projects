package main

import "fmt"

func calculate() {
	var num1, num2 int
	var operator string
	fmt.Print("Enter your first number:")
	fmt.Scan(&num1)
	fmt.Print("Enter operator:")
	fmt.Scan(&operator)
	fmt.Print("Enter your second number:")
	fmt.Scan(&num2)

	var sum, difference, product, quotient int
	switch operator {
	case "+":
		sum = num1 + num2
		fmt.Println("Sum:", sum)

	case "-":
		difference = num1 - num2
		fmt.Println("Difference:", difference)

	case "*":
		product = num1 * num2
		fmt.Println("Product:", product)

	case "/":
		if num2 != 0 {
			quotient = num1 / num2
			fmt.Println("Quotient:", quotient)
		} else {
			fmt.Println("Division by zero cannot be perfomed!")
		}

	default:
		fmt.Println("Invalid operator! Please use one of this: +, -, /, *")
	}

}

func main() {

	calculate()

}
