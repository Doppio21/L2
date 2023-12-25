package main

import "fmt"

// Command pattern используется для инкапсуляции запроса в виде объекта 
// и позволяет параметризовать клиентов с помощью различных запросов, 
// ставить запросы в очередь или журнал, а также поддерживать невыполнимые операции. 


// интерфейс, определяющий метод выполнения.
type Command interface {
	Execute() string
}

// структура, реализующая интерфейс и содержащая объект-приемник.
type SimpleCommand struct {
	receiver Receiver
}

func (sc *SimpleCommand) Execute() string {
	return sc.receiver.Action()
}

// структура, которая инициирует выполнение команды.
type Invoker struct {
	command Command
}

func (i *Invoker) ExecuteCommand() string {
	return i.command.Execute()
}

// структура, которая получает запрос и выполняет действие.
type Receiver struct{}

func (r *Receiver) Action() string {
	return "Action Performed"
}

func main() {
	receiver := &Receiver{}
	concreteCommand := &SimpleCommand{receiver: *receiver}
	invoker := &Invoker{command: concreteCommand}
	result := invoker.ExecuteCommand()
	fmt.Println(result)
}
