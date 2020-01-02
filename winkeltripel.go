package main

import (
	"math"
)

// Winkeltripel converts an incoming lat lng to a Point (x,y) using the WinkelTripel formula
// as described here - https://en.wikipedia.org/wiki/Winkel_tripel_projection
func Winkeltripel(lat float64, lng float64) Point {
	Lam := lng * (math.Pi / 180)
	Phi := lat * (math.Pi / 180)
	Alpha := math.Acos(math.Cos(Phi) * math.Cos(Lam/2))
	Phi1 := math.Acos(2 / math.Pi)
	var SincAlpha float64
	if Alpha == 0 {
		SincAlpha = 1

	} else {
		SincAlpha = math.Sin(Alpha) / Alpha
	}

	X := .5 * (Lam*math.Cos(Phi1) + ((2 * math.Cos(Phi) * math.Sin(Lam/2)) / SincAlpha))
	Y := .5 * (Phi + (math.Sin(Phi) / SincAlpha))
	var pt Point
	pt.X = X / (math.Pi / 180)
	pt.Y = Y / (math.Pi / 180)

	return pt
}
