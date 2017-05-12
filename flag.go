package main

import(
	// "fmt"
)

// Flag  can be concurrently accessed by gorutines.
// It is implemented as a channel with capacity one.
// The methods are: Set and TestAndCancel.
type Flag chan bool

// You have to use it to create a working Flag
func MakeFlag() Flag{
	return make(chan bool, 1)
}

// Set the flag
func (f Flag) Set() {
	select {
	case f <- true:
		// f was canceled, now it is set
		return
		
	default:
		// f was set before. Nothing to do
		return
	}
}


// Cancel the flag and return whether it was set before.
func (f Flag) TestAndCancel() bool {
	select {
	case  <- f:
		// f was set, return true. (Now f is canceled)
		return true
		
	default:
		// f was set canceled. Return false
		return false
	}
}
