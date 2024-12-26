package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	var ip [4]byte
	ip[0], ip[1], ip[2], ip[3] = 127, 0, 0, 1
	fmt.Println(formatIP(ip))

	a := 2
	b := 20
	fmt.Println(listEven(a, b))

	input := "Hello, world!"
	result := countChars(input)
	for char, count := range result {
		fmt.Printf("'%c' встречается %d раз\n", char, count)
	}

	p1 := Point{X: 1.0, Y: 4.0}
	p2 := Point{X: 4.0, Y: 6.0}

	segment := Segment{Start: p1, End: p2}
	segmentLength := segment.Length()
	fmt.Println("Длина отрезка: ", segmentLength)

	t := Triangle{A: Point{X: 0, Y: 0}, B: Point{X: 3, Y: 0}, C: Point{X: 0, Y: 4}}
	tArea := t.Area()
	fmt.Println("Площадь треугольника: ", tArea)

	c := Circle{Center: Point{X: 0, Y: 0}, Radius: Segment{Start: Point{X: 0, Y: 0}, End: Point{X: 0, Y: 4}}}
	cArea := c.Area()
	fmt.Println("Площадь круга: ", cArea)

	printArea(t)
	printArea(c)

	numbers := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	fmt.Println("Исходный срез:", numbers)
	squareFunc := square
	resultFunc := Map(numbers, squareFunc)
	fmt.Println("Срез после применения функции:", resultFunc)
}

func formatIP(ip [4]byte) string {
	parts := make([]string, len(ip))
	for i := 0; i < len(ip); i++ {
		parts[i] = strconv.Itoa(int(ip[i]))
	}
	result := strings.Join(parts, ".")
	return result
}

func listEven(a, b int) ([]int, error) {
	var even []int
	if a <= b {
		for i := a; i <= b; i++ {
			if i%2 == 0 {
				even = append(even, i)
			}
		}
		return even, nil
	}
	return even, fmt.Errorf("левая граница больше правой")
}

func countChars(s string) map[rune]int {
	counts := make(map[rune]int)
	s = strings.ToLower(s)
	for _, char := range s {
		counts[char]++
	}
	return counts
}

type Point struct {
	X float64
	Y float64
}

type Segment struct {
	Start Point
	End   Point
}

func (s Segment) Length() float64 {
	dx := s.End.X - s.Start.X
	dy := s.End.Y - s.Start.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type Triangle struct {
	A Point
	B Point
	C Point
}

type Circle struct {
	Center Point
	Radius Segment
}

func (t Triangle) Area() float64 {
	aLength := Segment{Start: t.A, End: t.B}.Length()
	bLength := Segment{Start: t.B, End: t.C}.Length()
	cLength := Segment{Start: t.C, End: t.A}.Length()
	p := (aLength + bLength + cLength) / 2
	return math.Sqrt(p * (p - aLength) * (p - bLength) * (p - cLength))
}

func (c Circle) Area() float64 {
	radiusLength := c.Radius.Length()
	return math.Pi * radiusLength * radiusLength
}

type Shape interface {
	Area() float64
	Perimitr() float64
}

func printArea(s Shape) {
	result := s.Area()
	fmt.Printf("Площадь фигуры: %.2f\n", result)
}

func Map(slice []float64, f func(float64) float64) []float64 {
	newSlice := make([]float64, len(slice))
	copy(newSlice, slice)
	for i, v := range newSlice {
		newSlice[i] = f(v)
	}
	return newSlice
}

func square(x float64) float64 {
	return x * x
}
