package main

import (
	"context"
	"fmt"
	"time"

	waiter "github.com/budimanlai/go-waiter"
)

func main() {
	// runExample()
	// runWithTimeout()
	runRequest()
}

func runRequest() {
	l := waiter.NewWaiter[string]()

	resp, err := l.Request(context.Background(), func(key string) error {
		fmt.Println("Simulating work... " + key)
		time.Sleep(3 * time.Second)
		l.Resolve(key, "Hello, World!")
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Received response:", resp)
}

func runWithTimeout() {
	l := waiter.NewWaiter[string]()
	key := "exampleKey"
	l.Register(key)

	// Simulate resolving the listener in a separate goroutine
	go func() {
		fmt.Println("Simulating work...")
		time.Sleep(6 * time.Second)
		l.Resolve(key, "Hello, World!")
	}()

	fmt.Println("Waiting for response...")
	// Listen for the response
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := l.Listen(ctx, key)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Received response:", response)
}

func runExample() {
	l := waiter.NewWaiter[string]()
	key := "exampleKey"
	l.Register(key)

	// Simulate resolving the listener in a separate goroutine
	go func() {
		fmt.Println("Simulating work...")
		time.Sleep(3 * time.Second)
		l.Resolve(key, "Hello, World!")
	}()

	fmt.Println("Waiting for response...")
	// Listen for the response
	ctx := context.Background()
	response, err := l.Listen(ctx, key)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Received response:", response)
}
