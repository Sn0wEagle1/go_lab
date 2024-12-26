package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")

	var name string
	var greeting string
	name = "Student"
	greeting = hello(name)
	fmt.Println(greeting)

	var a, b int64
	var err error
	a = 2
	b = 20
	fmt.Println("Все четные элементы: ")
	err = printEven(a, b)
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}

	var num1, num2 float64 = 10, 0
	var operator string = "/"
	var num3, num4 float64 = 24, 2
	var operator1 string = "*"
	result, err := apply(num1, num2, operator)
	result1, err1 := apply(num3, num4, operator1)
	fmt.Println("Результат: ", result)
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}
	fmt.Println("Результат: ", result1)
	if err1 != nil {
		fmt.Println("Ошибка: ", err1)
	}
}

func hello(name string) string {
	return "Hello, " + name + "!"
}

func printEven(a, b int64) error {
	if a <= b {
		for i := a; i <= b; i++ {
			if i%2 == 0 {
				fmt.Println(i)
			}
		}
		return nil
	}
	return fmt.Errorf("левая граница больше правой")
}

func apply(num1, num2 float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return (num1 + num2), nil
	case "-":
		return (num1 - num2), nil
	case "*":
		return (num1 * num2), nil
	case "/":
		if num2 != 0 {
			return (num1 / num2), nil
		}
		return 0, fmt.Errorf("деление на 0 невозможно")
	}
	return 0, fmt.Errorf("неверный оператор")
}
