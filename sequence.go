package main

import (
	"github.com/go-gl/mathgl/mgl32"
)


func squareSpiral( n int ) []mgl32.Vec3 {
	out := make( []mgl32.Vec3, n )
	up := mgl32.Vec3{ 0, 1, 0 }
	left := mgl32.Vec3{ -1, 0, 0 }
	down := mgl32.Vec3{ 0, -1, 0 }
	right := mgl32.Vec3{ 1, 0, 0 }
	orbit := 0;
        i:= 0
	position:= mgl32.Vec3{ 0, 0, 0 }
	base :=-1
	for i< n {
		out[i]= position;
		i++;
		// compute the next position
		switch {
		case i == (2*orbit+1)*(2*orbit+1):
			// go to the next orbit
			orbit++;
			base=i-1 // last index in the previous orbit
			position= position.Add(left) // the first in the left arm
		case i-base <= 2*orbit:
			// still in the left arm
			position= position.Add(up)
		case i-base <= 4*orbit:
			// still in the upper arm
			position= position.Add(right)
		case i-base <= 6*orbit:
			// still in the right arm
			position= position.Add(down)
		case i-base <= 8*orbit:
			// still in the lower arm
			position= position.Add(left)			
		}
	}

	
	return out
}
