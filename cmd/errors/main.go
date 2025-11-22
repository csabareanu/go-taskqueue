package main

import (
	"fmt"
	"errors"
)

type DivideError struct {
	A, B float64
	Msg string
}

func (e *DivideError) Error() string {
    return fmt.Sprintf("DivideError: %s (A=%.2f, B=%.2f)", e.Msg, e.A, e.B)
}

func divide(a float64, b float64) (float64, error)	{
	if b==0 {
		// return 0, errors.New("Cannot divide by 0")
		return 0, &DivideError{A: a, B:b, Msg: "division by zero"}
	}
	return a/b, nil
}

func main()	{
	res, err := divide(10, 0)
	if err != nil {
		var divErr *DivideError
		if errors.As(err, &divErr)	{
			fmt.Println("---Caught divider error---")
			fmt.Println(err)
			return
		}
		
	}
	fmt.Println("Result: ", res)
}