package main

import (
	"fmt"
	"time"
)

func main() {
	chanel := make(chan int, 10)
	go produce(chanel)
	consume(chanel)
}

func produce(chanel chan<- int) {
	for i := 0; i < 10; i++ {
		fmt.Println("produce", i)
		chanel <- i
		time.Sleep(1 * time.Second)
	}
	close(chanel)
}

func consume(chanel <-chan int) {
	for i := range chanel {
		fmt.Println("consume", i)
		time.Sleep(1 * time.Second)
	}
}
