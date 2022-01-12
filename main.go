package main

import (
	"fmt"
	"time"
)

func main() {

	timeChan := time.NewTimer(time.Second * 5)
	tickChan := time.NewTimer(time.Second * 1)
	resetChan := make(chan bool)
	resetCompreteChan := make(chan bool)

	go sub(timeChan, tickChan, resetChan, resetCompreteChan)
	go Wait(resetChan, resetCompreteChan)

	for {
		time.Sleep(100)
	}
}

func sub(timeChan *time.Timer, tickChan *time.Timer, resetChan chan bool, resetCompreteChan chan bool) {
	var count int = 0

	for {

		select {
		case <-timeChan.C:
			fmt.Printf("timer\n")
			resetTimer(timeChan, time.Second*5, resetCompreteChan)

		case <-tickChan.C:
			fmt.Printf("count:%v\n", count)
			fmt.Printf("tick\n")
			resetTimer(tickChan, time.Second, resetCompreteChan)

			count++

		case reset := <-resetChan:
			// fmt.Printf("resetchan\n")
			if reset {

				resetTimer(timeChan, time.Second*5, resetCompreteChan)
				resetTimer(tickChan, time.Second, resetCompreteChan)
				fmt.Printf("reset\n")

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
		time.Sleep(time.Second * 12)
		resetChan <- true
		time.Sleep(time.Second * 6)
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
