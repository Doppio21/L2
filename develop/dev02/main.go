package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)
var ErrInvalidLine = errors.New("invalid line")

func Unpacked(str string) (string, error) {
	ret := strings.Builder{}

	for i := 0; i < len(str); i++ {
		// проверка на escape-последовательность и ее обработка
		if str[i] == 92 {
			ret.WriteByte(str[i+1])
			i = i + 1
			continue
		}
		// Проверка является ли символ числом от 0 до 9
		if str[i] >= 48 && str[i] <= 57 {
			// если 1 символ является числом это некорректная строка
			if i == 0 {
				return "", ErrInvalidLine
			}

			num := strings.Builder{}
			// запись в num числа, в том числе, если оно является многозначным
			for j := i; j < len(str) && str[j] >= 48 && str[j] <= 57; j++ {
				num.WriteByte(str[j])
			}

			n, err := strconv.Atoi(num.String())
			if err != nil {
				return "", err
			}
			// если символ должен повториться 1 раз, переходим к следующей итерации
			if n == 1 {
				continue
			}

			// переменная для предыдущего символа
			prev := str[i-1]
			// повтор символа и его запись в ret
			rep := strings.Repeat(string(prev), n-1)
			ret.WriteString(rep)

			// если достигли конеца строки, возвращаем результат
			if i == len(str)-len(num.String()) {
				return ret.String(), nil
			}
			// перемещаем индекс на колличество элементов в числе
			i += len(num.String())-1
			continue
		}
		
		ret.WriteByte(str[i])
	}
	return ret.String(), nil
}

func main() {
	var str string
	fmt.Scan(&str)

	res, err := Unpacked(str)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%s => %s\n", str, res)
}
