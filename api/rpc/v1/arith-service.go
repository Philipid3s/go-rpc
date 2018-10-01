package v1

import (
	"errors"
)

// Arithmetic arguments
type Args struct {
	A, B int
}

// Quotient and Remaining
type Quotient struct {
	Quo, Rem int
}

// Type Arith: int
type Arith int

// Multiply 2 integers
func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

// Divide 2 integers
func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
