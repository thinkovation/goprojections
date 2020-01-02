package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	svg "github.com/ajstarks/svgo"
)

type GeoJSONPolygon struct {
	Coordinates [][][]float64 `json:"coordinates"`
}

type GeoJSON struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Crs  struct {
		Type       string            `json:"type"`
		Properties map[string]string `json:"properties"`
	} `json:"crs"`
	Features []struct {
		Type       string            `json:"type"`
		Properties map[string]string `json:"properties"`
		Geometry   struct {
			Type        string          `json:"type"`
			Coordinates json.RawMessage `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}
type GeoXML interface{}

func circle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(500, 500)
	s.Circle(250, 250, 125, "fill:none;stroke:black")
	s.End()
}
func world(w http.ResponseWriter, req *http.Request) {
	keys, ok := req.URL.Query()["projection"]
	projection := ""
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		projection = ""
	} else {
		projection = keys[0]
	}

	log.Println("World Called")
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(3000, 1500)
	backpts := BlueBack()

	addPolygon(s, backpts, 5, "fill:#bbbbFF", projection)

	addCountries(s, projection)
	s.End()
}

type PolygonGeometry [][][]float64
type MultiPolygonGeometry [][][][]float64

type GeoPoly struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}
type Point struct {
	X float64
	Y float64
}

func BlueBack() []Point {
	var pts []Point
	lng := -180.0

	for lng <= 180 {
		var pt Point
		pt.Y = lng
		pt.X = -90
		lng = lng + 2
		pts = append(pts, pt)
	}
	lat := -90.0
	for lat <= 90 {
		var pt Point
		pt.Y = 180
		pt.X = lat
		lat = lat + 2
		pts = append(pts, pt)
	}
	lng = 180
	for lng >= -180 {
		var pt Point
		pt.Y = lng
		pt.X = 90
		lng = lng - 2
		pts = append(pts, pt)
	}
	lat = 90
	for lat >= -90 {
		var pt Point
		pt.Y = -180
		pt.X = lat
		lat = lat - 2
		pts = append(pts, pt)
	}

	return pts
}
func Sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	return math.Sin(x) / x
}

// NoTransform takes a lat lng and returns it as a point
func NoTransform(lat float64, lng float64) Point {
	var pt Point
	pt.X = lng
	pt.Y = lat
	return pt
}

func addPolygon(s *svg.SVG, pts []Point, multfactor float64, style string, projection string) {
	if style == "" {
		style = "fill:#dd3333;stroke:black"
	}
	var xcoll []int
	var ycoll []int
	for _, sp := range pts {

		//Latitude is the Y axis, longitude is the X axis.
		switch projection {
		case "winkeltripel":
			sp = Winkeltripel(sp.X, sp.Y)
		case "aitoff":
			sp = Aitoff(sp.X, sp.Y)
		case "robinson":
			sp = Robinson(sp.X, sp.Y, 300, 200)
		default:
			sp = NoTransform(sp.X, sp.Y)
		}

		//
		tx := (180 + sp.X) * multfactor
		ty := (90 - sp.Y) * multfactor
		xcoll = append(xcoll, int(tx))
		ycoll = append(ycoll, int(ty))

	}
	//	fmt.Println(xcoll)
	//	fmt.Println(ycoll)

	s.Polygon(xcoll, ycoll, style)
}

func addCountries(s *svg.SVG, projection string) {

	//	data, err := ioutil.ReadFile("dat.geojson")
	//	data, err := ioutil.ReadFile("dat full.geojson")
	data, err := ioutil.ReadFile("countries.geojson.txt")

	if err == nil {
		//		fmt.Println(data[:20])
	}
	var j GeoJSON
	JSONerr := json.Unmarshal(data, &j)
	if JSONerr != nil {
		fmt.Println(JSONerr)
	}
	for _, v := range j.Features {

		polygonstyle := ""
		if v.Properties["ADMIN"] == "United Kingdom" {
			polygonstyle = "fill:#3333ff;stroke:black"

		}
		if v.Geometry.Type == "Polygon" {
			var candpoly PolygonGeometry
			candpolyerr := json.Unmarshal(v.Geometry.Coordinates, &candpoly)
			if candpolyerr != nil {
				fmt.Println(candpolyerr)
			} else {
				for _, cv := range candpoly {
					var polypts []Point

					for _, ccv := range cv {
						var polypt Point
						//Latitude is the Y axis, longitude is the X axis.
						polypt.X = ccv[1]
						polypt.Y = ccv[0]
						polypts = append(polypts, polypt)
					}
					addPolygon(s, polypts, 5, polygonstyle, projection)

				}
			}

		}
		if v.Geometry.Type == "MultiPolygon" {
			var candMultiPoly MultiPolygonGeometry
			candMultiPolyErr := json.Unmarshal(v.Geometry.Coordinates, &candMultiPoly)
			if candMultiPolyErr != nil {
				fmt.Println(candMultiPolyErr)

			} else {

				for _, cmp := range candMultiPoly {
					//fmt.Println(cmp)
					for _, ccmp := range cmp {
						//	fmt.Println(ccmp)
						var polypts []Point
						for _, cccmp := range ccmp {
							var polypt Point
							//Latitude is the Y axis, longitude is the X axis.
							polypt.X = cccmp[1] // Lng
							polypt.Y = cccmp[0] // Lat
							polypts = append(polypts, polypt)
						}
						addPolygon(s, polypts, 5, polygonstyle, projection)
					}
				}

			}
		}

	}

}
