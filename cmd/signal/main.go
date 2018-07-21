package main

import (
	"os"
	"fmt"
	"os/signal"
	"time"
)
//buf:queue destroy
// kill

var (
	running = true
)

func main() {
	c := make(chan os.Signal)
	println(os.Getpid())
	//signal.Notify(c, syscall.SIGTERM,os.Kill, os.Interrupt)
	signal.Notify(c)
	fmt.Println("boot...")
	go func(){
		for running{
			fmt.Println("work..")
			time.Sleep(time.Second *3)
		}
		fmt.Println("clean...")
		os.Exit(0)
		//gc operations
	}()

	for{
		select{
		case s := <- c:
			fmt.Println(s)
			println(s.String())
			switch s.String() {
			case "terminated":
				println("term ...")
				running = false
			case "interrupt":
				println("interrupt ...")
				//interruptHandle()

			}
		}
	}
}
