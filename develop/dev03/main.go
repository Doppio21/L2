package main

import (
	"L2/develop/dev03/sort"
	"fmt"
)

func main() {
	sort, err := sort.NewSort()
	if err != nil {
		fmt.Println(err)
	}
	sort.Run()
	for _, s := range sort.Data {
		fmt.Println(s)
	}
}
