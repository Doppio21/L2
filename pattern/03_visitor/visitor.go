package main

import "fmt"

// Visitor pattern — это шаблон поведенческого проектирования,
// который позволяет добавлять поведение в структуру без фактического изменения структуры.

// Перенос операций в классы посетителей выгоден в следующих случаях:
// требуется множество несвязанных операций над структурой объекта,
// классы, составляющие структуру объекта, известны и не должны изменяться,
// необходимо часто добавлять новые операции,
// алгоритм включает в себя несколько классов структуры объекта, но желательно управлять им в одном единственном месте,
// Алгоритм должен работать с несколькими независимыми иерархиями классов.
// Недостаток этого шаблона, однако, заключается в том, что он затрудняет расширение иерархии классов, 
// так как новые классы обычно требуют добавления нового метода к каждому посетителю. 

// Допустим, есть 3 объекта:
// Server
// Client
// DB
// Каждая из приведенных выше структур реализует общую форму интерфейса и необходимо добавить
// новое поведение SendData()

type Service interface {
	getData() string
	accept(visitor)
}

type Server struct {
}

func (s *Server) getData() string {
	return "data"
}

func (s *Server) accept(v visitor) {
	v.visitForServer(s)
}

type Client struct {
}

func (c *Client) getData() string {
	return "data"
}

func (c *Client) accept(v visitor) {
	v.visitForClient(c)
}

type DB struct {
}

func (d *DB) getData() string {
	return "data"
}

func (d *DB) accept(v visitor) {
	v.visitForDB(d)
}

type visitor interface {
	visitForServer(*Server)
	visitForClient(*Client)
	visitForDB(*DB)
}

type SendData struct {
	Data string
}

func (s *SendData) visitForServer(srv *Server) {
	fmt.Println("server: send data")
}

func (s *SendData) visitForClient(c *Client) {
	fmt.Println("client: send data")
}

func (s *SendData) visitForDB(d *DB) {
	fmt.Println("db: send data")
}


func main(){
	server:= Server{}
	client:= Client{}
	db:= DB{}
	
	sendData:=SendData{}
	server.accept(&sendData)
	client.accept(&sendData)
	db.accept(&sendData)
}