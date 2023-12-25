Что выведет программа? Объяснить вывод программы.

package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

Данная программа выведет error.
При вызове функции test() возвращается структура customError, удовлетворяющая
интерфейсу error.
Интерфейсы реализованы как структура из двух элементов: тип и значение.
Значение интерфейса будет nil только в том случае, если и тип и значение будут nil,
следовательно err != nil.

