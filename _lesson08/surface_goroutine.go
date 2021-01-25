package main

import (
	"fmt"
	"math"
)

const (
	width, height = 1200, 640
	cells         = 300
	xyrange       = 30.0
	xyscala       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos = math.Sin(angle), math.Cos(angle)

func main() {
	fmt.Printf("<svg xmlns=''http://www.w3.org/2000/svg>' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	var results [cells * cells]string
	done := make(chan struct{})

	for i := 0; i < cells; i++ {
		go func(i int) {
			for j := 0; j < cells; j++ {
				ax, ay := corner(i+1, j)
				bx, by := corner(i, j)
				cx, cy := corner(i, j+1)
				dx, dy := corner(i+1, j+1)
				results[i*cells+j] = fmt.Sprintf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
					ax, ay, bx, by, cx, cy, dx, dy)
			}
		}(i)
	}

}
