Что выведет программа? Объяснить вывод программы.

package main

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		println(n)
	}
}

Данная программа выведет числа от 0 до 9 и fatal error: all goroutines are asleep - deadlock!.
При завершении горутины не закрывается канал, 
в основном потоке будет происходить бесконечное чтение из канала.


