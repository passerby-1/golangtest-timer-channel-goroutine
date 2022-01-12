package main

import (
	"fmt"
	"time"
)

func main() {

	countChan := time.NewTimer(time.Second)
	time.Sleep(time.Second)
	fizzChan := time.NewTimer(time.Second * 5)
	buzzChan := time.NewTimer(time.Second * 3)
	resetChan := make(chan bool)
	resetCompreteChan := make(chan bool)

	go sub(fizzChan, buzzChan, countChan, resetChan, resetCompreteChan)
	go Wait(resetChan, resetCompreteChan)

	for {
		time.Sleep(100)
	}
}

func sub(fizzChan *time.Timer, buzzChan *time.Timer, countChan *time.Timer, resetChan chan bool, resetCompreteChan chan bool) {

	var count int = 0

	for {

		select {
		case <-countChan.C:
			fmt.Printf("\ncount:%v: ", count+1)
			resetTimer(countChan, time.Second, resetCompreteChan)
			count++

		case <-fizzChan.C:
			fmt.Printf("Fizz")
			resetTimer(fizzChan, time.Second*5, resetCompreteChan)

		case <-buzzChan.C:
			fmt.Printf("Buzz")
			resetTimer(buzzChan, time.Second*3, resetCompreteChan)

		case reset := <-resetChan:
			if reset {

				resetTimer(fizzChan, time.Second*5, resetCompreteChan)
				resetTimer(buzzChan, time.Second*3, resetCompreteChan)
				resetTimer(countChan, time.Second, resetCompreteChan)
				fmt.Printf("COUNT RESET\n")

				count = 0

				resetCompreteChan <- true

			}
		}
	}
}

func Wait(resetChan chan bool, resetCompreteChan chan bool) {

	go func() {
		for {
			select {
			case <-resetCompreteChan:
				resetChan <- false
			default:
			}
		}
	}()

	for {
		time.Sleep(time.Second * 20)
		resetChan <- true
		time.Sleep(time.Second * 10)
		resetChan <- true
	}

}

func resetTimer(timer *time.Timer, d time.Duration, resetCompreteChan chan bool) {

	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
			resetCompreteChan <- false
		}
	}

	timer.Reset(d)
}
