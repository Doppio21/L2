package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

const (
	_ = iota
	ErrorFprintln
	ErrorQuery
	ErrorPrintln
)

func main() {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, err.Error())
		if err != nil {
			os.Exit(ErrorFprintln)
		}
		os.Exit(ErrorQuery)
	}
	_, err = fmt.Println(time)
	if err != nil {
		os.Exit(ErrorPrintln)
	}
}
