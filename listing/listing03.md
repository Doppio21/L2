Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}

Данная программа выведет nil, false.
PathError удовлетворяет интерфейсу error, поэтому его можно вернуть как error.
Интерфейсы реализованы как структура из двух элементов: тип и значение
type iface struct {
    tab  *itab
    data unsafe.Pointer
}

Значение интерфейса будет nil только в том случае, если и тип и значение будут nil,
следовательно при проверке err == nil будет выведен false.

Пустой интерфейс это интерфейс, не имеющий методов, все типы удовлетворяют пустому интерфейсу.
type eface struct {
	_type *_type
	data  unsafe.Pointer
}
