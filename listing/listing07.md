Что выведет программа? Объяснить вывод программы.

package main
import (
	"fmt"
	"math/rand"
	"time"
)
func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}
func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}
func main() {
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}

Данная программа выведет числа от 1 до 8, а затем бесконечные 0.
После записи в каналы переданных значений в функции asChan() и чтения этих значений в канал с,
происходит закрытие каналов а и b, но чтение из них бесконечно продолжается,
и в канал с записываютя дефолтные значения = 0.