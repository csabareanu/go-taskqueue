package main

import (
	"fmt"
)

type Counter struct {
	Value int
}

// --- Value Receiver --- 
func (c Counter) IncrementValue() {
	fmt.Printf("[Value Receiver] Before increment: %d\n", c.Value)
	c.Value++
	fmt.Printf("[Value Receiver] After increment: %d (inside method)\n", c.Value)
}

// --- Pointer Receiver --- 
func (c *Counter) IncrementPointer() {
	fmt.Printf("[Pointer Receiver] Before increment: %d\n", c.Value)
	c.Value++
	fmt.Printf("[Pointer Receiver] After increment: %d (inside method)\n", c.Value)
}

func main()	{
	c := Counter{Value: 10}
	fmt.Println("Initial counter value: ", c.Value)
	fmt.Println("Calling IncrementValue (by Value) -----")
	c.IncrementValue()
	fmt.Println("After calling IncrementValue (outside):", c.Value)

	fmt.Println("Calling IncrementPointer (by Pointer) -----")
	c.IncrementPointer()
	fmt.Println("After calling IncrementPointer (outside):", c.Value)
}