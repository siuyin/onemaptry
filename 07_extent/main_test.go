package main

import (
	"fmt"
	"math"
	"testing"
)

type point struct {
	x float64
	y float64
}

func (p point) String() string {
	return fmt.Sprintf("(%v,%v)", p.x, p.y)
}

func testExtent(t *testing.T) {
	dat := []struct {
		p  point
		r  float64
		sw point
		ne point
	}{
		{point{20000, 30000}, 200, point{20000 - 100, 30000 - 100}, point{20000 + 100, 30000 + 100}},
	}
	eps := 0.01
	for i, d := range dat {
		if sw, ne := extent(d.p, d.r); diff(sw, d.sw) > eps || diff(ne, d.ne) > eps {
			t.Errorf("case %d: expected %s and %s, got: %s and %s", i, d.sw, d.ne, sw, ne)
		}
	}
}

func diff(p, q point) float64 {
	return math.Max(math.Abs(p.x-q.x),
		math.Abs(p.y-q.y))
}

func extent(p point, r float64) (point, point) {
	sw := point{p.x - r/2, p.y - r/2}
	ne := point{p.x + r/2, p.y + r/2}
	return sw, ne
}
