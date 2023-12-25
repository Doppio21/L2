package main

import (
	"L2/develop/dev05/grep"
	"fmt"
)

func main() {
	g, err := grep.NewG()
	if err != nil {
		fmt.Println(err)
		return
	}

	g.Run()
}
