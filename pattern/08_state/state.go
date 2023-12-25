package main

import "fmt"

// State pattern — это поведенческий шаблон проектирования программного 
// обеспечения, который позволяет позволяет объекту изменять свое поведение 
// при изменении его внутреннего состояния.  Паттерн состояния можно интерпретировать как шаблон стратегии, 
// который способен переключать стратегию через вызовы методов, определенных в интерфейсе шаблона.

type SystemState int

const (
	Working SystemState = iota
	Sleep
	Off
)

type System struct {
	state SystemState
}

func (s *System) setState(newState SystemState) {
	s.state = newState
}

func (s *System) ChangeState() {
	switch s.state {
	case Working:
		fmt.Println("State: Working")
		s.setState(Working)
	case Sleep:
		fmt.Println("State: Sleep")
		s.setState(Sleep)
	case Off:
		fmt.Println("State: Off")
		s.setState(Off)
	}
}



func main() {
	systemState := System{state: Working}

	for i := 0; i < 3; i++ {
		systemState.setState(SystemState(i))
		systemState.ChangeState()
	}
}
