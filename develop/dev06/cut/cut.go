package cut

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Cut struct {
	f int
	d string
	s bool
}

func NewCut() *Cut {
	cut := &Cut{}
	flag.IntVar(&cut.f, "f", 0, "fields - выбрать поля (колонки)")
	flag.StringVar(&cut.d, "d", "", "delimiter - использовать другой разделитель")
	flag.BoolVar(&cut.s, "s", false, "separated - только строки с разделителем")
	flag.Parse()

	return cut
}

func (c *Cut) Run() error {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		s := sc.Text()
		p := strings.Split(s, c.d)
		if len(p) == 1 && c.s || len(p) < c.f {
			continue
		}

		fmt.Println(p[c.f-1])
	}

	return os.Stdout.Sync()
}
