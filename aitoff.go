package main

import (
	"math"
)

// Aitoff converts an incoming lat lng to a Point (x,y) using the Aitoff formula
// as described here - https://en.wikipedia.org/wiki/Aitoff_projection
func Aitoff(lat float64, lng float64) Point {
	Lam := lng * (math.Pi / 180)
	Phi := lat * (math.Pi / 180)
	Alpha := math.Acos(math.Cos(Phi) * math.Cos(Lam/2))
	var SincAlpha float64
	if Alpha == 0 {
		SincAlpha = 1

	} else {
		SincAlpha = math.Sin(Alpha) / Alpha
	}
	X := (2 * math.Cos(Phi) * math.Sin(Lam/2)) / SincAlpha
	Y := math.Sin(Phi) / SincAlpha

	var pt Point
	pt.X = X / (math.Pi / 180)
	pt.Y = Y / (math.Pi / 180)
	return pt
}
