package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

type Config struct {
	Host    string
	Port    string
	Timeout time.Duration
}

type Telnet struct {
	cfg Config

	conn net.Conn
}

func New(cfg Config) *Telnet {
	return &Telnet{
		cfg: cfg,
	}
}

func (t *Telnet) Dial() error {
	d := net.Dialer{Timeout: t.cfg.Timeout}
	conn, err := d.Dial("tcp", net.JoinHostPort(t.cfg.Host, t.cfg.Port))
	if err != nil {
		return err
	}

	t.conn = conn
	return nil
}

func (t *Telnet) SendAndRecv(buff []byte) ([]byte, error) {
	left := len(buff)
	for left != 0 {
		n, err := t.conn.Write(buff)
		if err != nil {
			return nil, err
		}

		left -= n
		buff = buff[n:]
	}
	b := make([]byte, 1024)
	n, err := t.conn.Read(b)
	if err != nil {
		return nil, err
	}
	b = b[:n]

	return b, nil
}

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "таймаут на подключение к серверу")
	flag.Parse()

	host := os.Args[len(os.Args)-2]
	port := os.Args[len(os.Args)-1]

	t := New(Config{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	})

	if err := t.Dial(); err != nil {
		fmt.Println(err)
		return
	}
	defer t.conn.Close()

	for {
		msg := ""
		_, err := fmt.Scan(&msg)
		if err != nil {
			break
		}

		if msg == "quit" {
			break
		}

		resp, err := t.SendAndRecv([]byte(msg))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(resp))
	}
}
