package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan bool, 0)

	go func() {
		for {
			ch <- true
			time.Sleep(time.Second * 1)
		}
	}()

	go func() {

		for {

			select {
			case c := <-ch:
				fmt.Println(c)

			}

			fmt.Println("break...")
		}
		fmt.Println("退出。。。")

	}()

	time.Sleep(time.Second * 60)

}
