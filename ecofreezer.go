package main

import(
	"fmt"
	"time"
)

const EcoFreezerDelay = 15 * time.Second // EcoFrezer delay in milliseconds

// EcoFreezer is a goroutine that sets g.PauseRequest if there was no action for a long period.
func (g *Mki3dGame) EcoFreezer() {
	fmt.Println("Starting EcoFreezer ...") // tests
	g.WasAction.Set() // set before the first testing

	for {
		time.Sleep( EcoFreezerDelay )
		if( ! g.Paused ) {
			wasAction := g.WasAction.TestAndCancel()
			// fmt.Println("EcoFreezer: Testing: wasAction =", wasAction )
			if !wasAction {
				fmt.Println("EcoFreezer: Setting PauseRequest ...")
				g.PauseRequest.Set()
			}
		}
	}
}
