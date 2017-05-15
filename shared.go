package main

import (
// "fmt" // tests
)

// SharedBool getter function
type SharedBoolGet func() bool

// SharedBool setter function
type SharedBoolSet func(bool)

// Type for bool variables shared by goroutines
type SharedBool struct {
	get SharedBoolGet
	set SharedBoolSet
}

func MakeSharedBool() SharedBool {
	var value bool
	getChan := make(chan bool)
	setChan := make(chan bool)

	// start goroutine serving the variable
	go func() {
		for {
			select {
			case getChan <- value:
				// output value
			case value = <-setChan:
				// input value
				// value = v
			}
		}
	}()

	// getter
	getter := func() bool {
		return <-getChan
	}

	// setter
	setter := func(v bool) {
		setChan <- v
	}

	return SharedBool{get: getter, set: setter}

}
