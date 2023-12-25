package grep

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type G struct {
	file     string
	pattern  string
	count    bool
	exclude  bool
	caseSens bool
	fixed    bool
	numbers  bool
	after    int
	before   int
	context  int
}

func NewG() (*G, error) {
	g := &G{}
	flag.BoolVar(&g.count, "c", false, "count (количество строк)")
	flag.BoolVar(&g.caseSens, "i", false, "ignore-case (игнорировать регистр)")
	flag.BoolVar(&g.exclude, "v", false, "invert (вместо совпадения, исключать)")
	flag.BoolVar(&g.fixed, "F", false, "fixed, точное совпадение со строкой, не паттерн")
	flag.BoolVar(&g.numbers, "n", false, "line num, напечатать номер строки")
	flag.IntVar(&g.after, "A", 0, "after печатать +N строк после совпадения")
	flag.IntVar(&g.before, "B", 0, "before печатать +N строк до совпадения")
	flag.IntVar(&g.context, "C", 0, "context (A+B) печатать ±N строк вокруг совпадения")
	flag.CommandLine.Parse(os.Args[2 : len(os.Args)-1])

	g.file = os.Args[len(os.Args)-1]
	g.pattern = os.Args[1]

	return g, nil
}

func (g *G) resolveMatcher() func(string) bool {
	if g.caseSens {
		g.pattern = strings.ToLower(g.pattern)
	}

	matcher := regexp.MustCompile(g.pattern).MatchString
	if g.fixed {
		matcher = func(s string) bool {
			return strings.Contains(s, g.pattern)
		}
	}

	if g.caseSens {
		old := matcher
		matcher = func(s string) bool {
			return old(strings.ToLower(s))
		}
	}

	if g.exclude {
		return func(s string) bool {
			return !matcher(s)
		}
	}

	return matcher
}

func (g *G) print(before []string, s string, idx int) {
	if g.count {
		return
	}

	for i, b := range before {
		if g.numbers {
			fmt.Printf("%d: ", idx-len(before)+i)
		}
		fmt.Println(b)
	}

	if g.numbers {
		fmt.Printf("%d: ", idx)
	}
	fmt.Println(s)
}

func (g *G) Run() error {
	f, err := os.Open(g.file)
	if err != nil {
		return err
	}
	defer f.Close()

	var (
		num           int
		matches       int
		before        = max(g.before, g.context)
		after         = max(g.after, g.context)
		leftAfter     = 0
		beforeStrings = make([]string, 0, before)
		m             = g.resolveMatcher()
	)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		num++

		str := sc.Text()
		if m(str) {
			matches++
			g.print(beforeStrings, str, num)
			beforeStrings = beforeStrings[:0]
			leftAfter = after
			continue
		} else if leftAfter != 0 {
			g.print(nil, str, num)
			leftAfter--
		}

		if before != 0 {
			if len(beforeStrings) == before {
				beforeStrings = beforeStrings[1:]
			}

			beforeStrings = append(beforeStrings, str)
		}
	}

	if g.count {
		fmt.Println(matches)
	}

	return nil
}
