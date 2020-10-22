package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"time"
)

var maxIter int = 3

func foo(channel chan string) {
	// TODO: Write an infinite loop of sending "pings" and receiving "pongs"
	for {
		// ? Pre-send code
		fmt.Println("Foo is sending: ping")
		channel <- "ping"

		// ? Post-receive code
		fmt.Println("Foo has received:", <-channel)
	}
}

func bar(channel chan string) {
	// TODO: Write an infinite loop of receiving "pings" and sending "pongs"
	for maxIter > 0 {
		// ? Post-receive code
		fmt.Println("Bar has received:", <-channel)

		// ? Pre-send code
		fmt.Println("Bar is sending: pong")
		channel <- "pong"

		maxIter--
	}
}

func pingPong() {
	// TODO: make channel of type string and pass it to foo and bar
	pingPongChan := make(chan string)
	go foo(pingPongChan) // Nil is similar to null. Sending or receiving from a nil chan blocks forever.
	go bar(pingPongChan)
	time.Sleep(500 * time.Millisecond)
}

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	pingPong()
}
