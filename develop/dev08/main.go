package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

type UnixProcess struct {
	pid   int
	ppid  int
	state rune
	pgrp  int
	sid   int

	binary string
}

func (p *UnixProcess) Refresh() error {
	statPath := fmt.Sprintf("/proc/%d/stat", p.pid)
	dataBytes, err := os.ReadFile(statPath)
	if err != nil {
		return err
	}

	data := string(dataBytes)
	binStart := strings.IndexRune(data, '(') + 1
	binEnd := strings.IndexRune(data[binStart:], ')')
	p.binary = data[binStart : binStart+binEnd]

	data = data[binStart+binEnd+2:]
	_, err = fmt.Sscanf(data,
		"%c %d %d %d",
		&p.state,
		&p.ppid,
		&p.pgrp,
		&p.sid)

	return err
}

func processes() ([]*UnixProcess, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	results := make([]*UnixProcess, 0, 50)
	for {
		names, err := d.Readdirnames(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, name := range names {
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}

			p, err := newUnixProcess(int(pid))
			if err != nil {
				continue
			}

			results = append(results, p)
		}
	}

	return results, nil
}

func newUnixProcess(pid int) (*UnixProcess, error) {
	p := &UnixProcess{pid: pid}
	return p, p.Refresh()
}

func cd(cur string, path string) (string, error) {
	var to string

	path = strings.TrimSpace(path)
	switch {
	case strings.HasPrefix(path, "/"):
		to = path
	case strings.HasPrefix(path, ".."),
		strings.HasPrefix(path, "./"):
		fallthrough
	default:
		to = filepath.Join(cur, path)
	}

	stat, err := os.Stat(to)
	if err != nil {
		return cur, err
	}

	if !stat.IsDir() {
		return cur, errors.New("not directory")
	}

	return to, nil
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("\033[34m", dir, "\033[32m", "$", "\033[0m ")
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		var (
			in  io.ReadWriter = os.Stdin
			out io.ReadWriter = bytes.NewBuffer(make([]byte, 1024))
		)

		str := sc.Text()
		pipelineParts := strings.Split(str, "|")
	loop:
		for idx, p := range pipelineParts {
			cmdParts := strings.Split(p, " ")

			switch cmdParts[0] {
			case "cd":
				if len(cmdParts) < 2 {
					fmt.Println("provide path")
					continue
				}

				var err error
				dir, err = cd(dir, cmdParts[1])
				if err != nil {
					fmt.Println(err)
				}
			case "pwd":
				if idx != len(pipelineParts)-1 {
					in = bytes.NewBufferString(dir)
				} else {
					fmt.Println(dir)
				}
			case "echo":
				if len(cmdParts) < 2 {
					fmt.Println("provide argument")
					continue
				}

				if idx != len(pipelineParts)-1 {
					in = bytes.NewBufferString(cmdParts[1])
				} else {
					fmt.Println(cmdParts[1])
				}
			case "kill":
				if len(cmdParts) < 2 {
					fmt.Println("provide argument")
					continue
				}

				pid, err := strconv.ParseInt(cmdParts[1], 10, 32)
				if err != nil {
					fmt.Println(err)
					continue
				}

				ps, err := os.FindProcess(int(pid))
				if err != nil {
					fmt.Println(err)
					continue
				}

				if err := ps.Kill(); err != nil {
					fmt.Println(err)
				}
			case "ps":
				pses, err := processes()
				if err != nil {
					fmt.Println(err)
					continue
				}

				buff := bytes.NewBuffer(make([]byte, 1024))
				writer := tabwriter.NewWriter(buff, 2, 4, 1, byte(' '), 0)
				fmt.Fprintln(writer, "PID\t PPID\t STATE\t PGRP\t SID\t BIN\t")
				for _, ps := range pses {
					fmt.Fprintf(writer, "%d\t %d\t %c\t %d\t %d\t %s\t\n",
						ps.pid, ps.ppid, ps.state, ps.pgrp, ps.sid, ps.binary)
				}

				writer.Flush()
				if idx != len(pipelineParts)-1 {
					in = buff
				} else {
					fmt.Println(buff.String())
				}
			case "\\quit":
				return
			default:
				if len(cmdParts) < 1 {
					fmt.Println("provide argument")
					continue
				}

				for i := 0; i < len(cmdParts); i++ {
					cmdParts[i] = strings.TrimSpace(cmdParts[i])
					if cmdParts[i] == "" {
						cmdParts = append(cmdParts[:i], cmdParts[i+1:]...)
						i--
					}
				}

				cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
				cmd.Stdin = in
				cmd.Stdout = out

				if err := cmd.Run(); err != nil {
					if cmd.ProcessState.ExitCode() != 1 {
						fmt.Println(err)
					}
					break loop
				}

				if idx == len(pipelineParts)-1 {
					data, err := io.ReadAll(out)
					if err != nil {
						fmt.Println(err)
						continue
					}

					fmt.Print(string(data))
					continue
				}

				in = out
				out = bytes.NewBuffer(make([]byte, 1024))
			}
		}

		fmt.Print("\033[34m", dir, "\033[32m", "$", "\033[0m ")
	}
}
