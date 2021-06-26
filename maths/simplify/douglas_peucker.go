package simplify

import (
	"github.com/flywave/go-vector-tiler/maths"
)

func DouglasPeucker(points []maths.Pt, tolerance float64) []maths.Pt {
	if tolerance <= 0 || len(points) <= 2 {
		return points
	}

	epsilon := tolerance * tolerance

	l := maths.Line{points[0], points[len(points)-1]}
	dmax := 0.0
	idx := 0
	for i := 1; i < len(points)-2; i++ {
		d := l.DistanceFromPoint(points[i])
		if d > dmax {
			dmax = d
			idx = i
		}
	}

	if dmax > epsilon {
		rec1 := DouglasPeucker(points[0:idx], epsilon)
		rec2 := DouglasPeucker(points[idx:], epsilon)

		newpts := append(rec1, rec2...)

		return newpts
	}

	return []maths.Pt{points[0], points[len(points)-1]}
}
