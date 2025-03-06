package main

import "fmt"

func main() {
	firstDigit := 5
	secondDigit := 10
	sum := SumTwoElem(firstDigit, secondDigit)
	fmt.Print("Сумма: ", sum)
}

func SumTwoElem(a, b int) int {
	sum := a + b
	var x float32 = 1
	sumA := a + b + int(x)
	fmt.Println(sumA)
	return sum
}

var sum = SumTwoElem(5, 5)
