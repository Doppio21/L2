package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	strs := []string{"пятак", "пЯтка", "тяпка", "тяпка", "Листок", "слиток", "столик", "кофта"}

	m := SearchAnagram(&strs)
	for k, v := range *m {
		fmt.Println(k, *v)
	}

}

func SearchAnagram(strs *[]string) *map[string]*[]string {
	m := make(map[string][]string)
	ret := make(map[string]*[]string)

	for _, str := range *strs {
		str = strings.ToLower(str)
		m[sortStr(str)] = append(m[sortStr(str)], str)
	}
	for _, v := range m {
		if len(v) == 1 {
			continue
		}
		val := editSlise(&v)
		ret[v[0]] = val
	}
	return &ret
}

func sortStr(s string) string {
	ret := []rune(s)

	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})

	return string(ret)
}

func editSlise(strs *[]string) *[]string {
	m := make(map[string]struct{})
	ret := []string{}

	for _, s := range *strs {
		if _, ok := m[s]; !ok {
			m[s] = struct{}{}
			ret = append(ret, s)
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})

	return &ret
}
