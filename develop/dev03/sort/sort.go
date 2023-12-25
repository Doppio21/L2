package sort

import (
	"flag"
	"os"
	"sort"
	"strconv"
	"strings"
)

type keys struct {
	n, r, u bool
	k       int
}

func ParserKeys() (*keys, error) {
	k := &keys{}
	flag.BoolVar(&k.n, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&k.r, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&k.u, "u", false, "не выводить повторяющие строки")
	flag.IntVar(&k.k, "k", 1, "указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умлчанию разделитель - пробел)")
	flag.Parse()

	k.k--
	return k, nil
}

type Sort struct {
	Data [][]string
	keys *keys
}

func NewSort() (*Sort, error) {
	k, err := ParserKeys()
	if err != nil {
		return nil, err
	}
	data := make([][]string, 0)

	file, err := os.ReadFile("input.txt")
	if err != nil {
		return nil, err
	}

	ss := strings.Split(string(file), "\n")
	for _, s := range ss {
		data = append(data, strings.Split(s, " "))
	}

	return &Sort{
		Data: data,
		keys: k,
	}, nil
}

func (s *Sort) Run() {
	switch {
	case s.keys.n:
		if s.keys.u {
			s.Data = unique(s.Data)
		}
		s.Data = valueSort(s.Data, s.keys)
	default:
		if s.keys.u {
			s.Data = unique(s.Data)
		}
		s.Data = startSort(s.Data, s.keys)
	}
}

func startSort(ss [][]string, k *keys) [][]string {
	sort.Slice(ss, func(i, j int) bool {
		if k.r {
			return ss[i][k.k] > ss[j][k.k]
		}
		return ss[i][k.k] < ss[j][k.k]
	})
	return ss
}

func unique(ss [][]string) [][]string {
	m := make(map[string]struct{})
	ret := make([][]string, 0)

	for _, s := range ss {
		str := strings.Builder{}
		for _, v := range s {
			str.WriteString(v)
		}
		if _, ok := m[str.String()]; !ok {
			m[str.String()] = struct{}{}
			ret = append(ret, s)
		}
	}
	return ret
}

func valueSort(ss [][]string, k *keys) [][]string {
	sort.Slice(ss, func(i, j int) bool {
		ok := true
		v1, err := strconv.Atoi(ss[i][k.k])
		if err != nil {
			ok = false
		}
		v2, err := strconv.Atoi(ss[j][k.k])
		if err != nil {
			ok = false
		}
		if k.r {
			if ok {
				return v1 > v2
			} else {
				return ss[i][k.k] > ss[j][k.k]
			}
		} else {
			if ok {
				return v1 < v2
			} else {
				return ss[i][k.k] < ss[j][k.k]
			}
		}
	})
	return ss
}
