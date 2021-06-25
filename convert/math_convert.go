package convert

import "github.com/flywave/go-vector-tiler/maths"

// TODO (gdey): Remove once conversion is over.
func FromMathPoint(pt ...maths.Pt) (gpts [][]float64) {
	gpts = make([][]float64, len(pt))
	for i := range pt {
		gpts[i] = []float64{pt[i].X, pt[i].Y}
	}
	return gpts
}
