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
		s.setState(Sleep)
	case Sleep:
		fmt.Println("State: Sleep")
		s.setState(Off)
	case Off:
		fmt.Println("State: Off")
		s.setState(Working)
	}
}



func main() {
	systemState := System{}
	systemState.setState(SystemState(0))
	for i := 0; i < 6; i++ {
		systemState.ChangeState()
	}
}
