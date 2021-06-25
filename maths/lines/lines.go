package lines

import geom "github.com/flywave/go-geom"

func FromTLineString(lnstr geom.LineString) (ln [][]float64) {
	pts := lnstr.Subpoints()
	for i := range pts {
		ln = append(ln, []float64{pts[i].X(), pts[i].Y()})
	}
	return ln
}
