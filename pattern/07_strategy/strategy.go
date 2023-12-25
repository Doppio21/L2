package main

import "fmt"

// Strategy pattern это шаблон проектирования поведенческого программного обеспечения, 
// который позволяет выбрать алгоритм во время выполнения. 
// Вместо того, чтобы реализовывать один алгоритм напрямую, 
// код получает инструкции во время выполнения о том, какой 
// из семейств алгоритмов следует использовать. 

type Operator interface {
	Apply(int, int) int
}

type Operation struct {
	Operator Operator
}

func (o *Operation) Operate(a, b int) int {
	return o.Operator.Apply(a, b)
}

type Addition struct{}

func (Addition) Apply(a, b int) int {
	return a + b
}

type Subtraction struct{}

func (Subtraction) Apply(a, b int) int {
	return a - b
}

func main() {
	add := Operation{Addition{}}
	fmt.Println(add.Operator.Apply(3, 2))

	sub := Operation{Subtraction{}}
	fmt.Println(sub.Operator.Apply(3, 2))
}
