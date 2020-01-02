package main

import (
	"fmt"
	"math"
)

// Robinson converts an incoming lat lng to a Point (x,y) using the Robinson formula
// as described here - https://en.wikipedia.org/wiki/Robinson_projection
func Robinson(lat float64, lng float64, mapWidth float64, mapHeight float64) Point {
	robinsonAA := []float64{
		0.84870000,
		0.84751182,
		0.84479598,
		0.84021300,
		0.83359314,
		0.82578510,
		0.81475200,
		0.80006949,
		0.78216192,
		0.76060494,
		0.73658673,
		0.70866450,
		0.67777182,
		0.64475739,
		0.60987582,
		0.57134484,
		0.52729731,
		0.48562614,
		0.45167814}
	robinsonBB := []float64{
		0.00000000,
		0.08384260,
		0.16768520,
		0.25152780,
		0.33537040,
		0.41921300,
		0.50305560,
		0.58689820,
		0.67047034,
		0.75336633,
		0.83518048,
		0.91537187,
		0.99339958,
		1.06872269,
		1.14066505,
		1.20841528,
		1.27035062,
		1.31998003,
		1.35230000}

	var mapOffsetX = 0.0
	var mapOffsetY = 0.0
	var heightFactor = 1.0

	// Robinson's latitude interpolation points are in 5-degree-steps
	fmt.Println(lat)
	var latitudeAbs = math.Abs(lat)
	fmt.Println(latitudeAbs)
	latitudeStepFloor := math.Floor(latitudeAbs / 5)
	var latitudeStepCeil = int(math.Ceil(latitudeAbs / 5))
	// calc interpolation factor (>=0 to <1) between two steps
	var latitudeInterpolation = (latitudeAbs - latitudeStepFloor*5) / 5
	// interpolate robinson table values

	fmt.Println(int(latitudeStepFloor))

	var AA = robinsonAA[int(latitudeStepFloor)] + (robinsonAA[latitudeStepCeil]-robinsonAA[int(latitudeStepFloor)])*latitudeInterpolation
	var BB = robinsonBB[int(latitudeStepFloor)] + (robinsonBB[latitudeStepCeil]-robinsonBB[int(latitudeStepFloor)])*latitudeInterpolation

	var robinsonWidth = 2 * math.Pi * robinsonAA[0]
	var widthFactor = mapWidth / robinsonWidth
	var latitudeSign = 1.0
	if lat < 0 {
		latitudeSign = -1.0
	}
	// Lat is Y , Long is X
	Y := widthFactor*BB*latitudeSign*heightFactor + mapOffsetY
	X := (widthFactor*AA*lng*math.Pi)/180 + mapOffsetX
	var pt Point
	pt.X = X
	pt.Y = Y
	return pt

}
