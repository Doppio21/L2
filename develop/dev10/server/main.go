package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Config struct {
	Address string
}

type Server struct {
	cfg Config

	mu    sync.Mutex
	conns map[int]net.Conn
}

func New(cfg Config) *Server {
	return &Server{
		cfg:   cfg,
		conns: make(map[int]net.Conn),
	}
}

func (s *Server) Run(ctx context.Context) error {
	listen, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()
	defer s.closeAllConns()
	defer listen.Close()

	connID := 0
	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		connID++
		s.addConn(connID, conn)

		wg.Add(1)
		go func(conn net.Conn, connID int) {
			defer wg.Done()
			defer s.delConn(connID)
			if err := s.handleConn(conn); err != nil {
				fmt.Println(err)
			}

		}(conn, connID)
	}

}

func (s *Server) addConn(id int, conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conns[id] = conn
}

func (s *Server) delConn(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, ok := s.conns[id]
	if !ok {
		return
	}

	conn.Close()
	delete(s.conns, id)
}

func (s *Server) closeAllConns() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, conn := range s.conns {
		conn.Close()
		delete(s.conns, id)
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		} else if errors.Is(err, io.EOF) {
			return nil
		}

		buff = buff[:n]
		fmt.Println(string(buff))

		_, err = conn.Write(buff)
		if err != nil {
			return err
		}
	}
}

func main() {
	var (
		address = "localhost:12312"
	)

	s := New(Config{Address: address})
	if err := s.Run(context.Background()); err != nil {
		return
	}
}
