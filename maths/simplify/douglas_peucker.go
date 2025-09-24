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
	for i := 1; i < len(points)-1; i++ {
		d := l.DistanceFromPoint(points[i])
		if d > dmax {
			dmax = d
			idx = i
		}
	}

	if dmax > epsilon {
		rec1 := DouglasPeucker(points[0:idx+1], tolerance)
		rec2 := DouglasPeucker(points[idx:], tolerance)

		// 移除重复的中间点
		newpts := make([]maths.Pt, 0, len(rec1)+len(rec2)-1)
		newpts = append(newpts, rec1...)
		newpts = append(newpts, rec2[1:]...)

		return newpts
	}

	return []maths.Pt{points[0], points[len(points)-1]}
}
