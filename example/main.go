package main

import (
	"context"
	"fmt"
	"time"

	waiter "github.com/budimanlai/go-waiter"
)

func main() {
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
