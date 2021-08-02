package basic

import (
	"fmt"

	geom "github.com/flywave/go-geom"
	"github.com/flywave/go-vector-tiler/maths"
)

type Line []Point

func (Line) basicType()                      {}
func (Line) String() string                  { return "Line" }
func (l Line) Direction() maths.WindingOrder { return maths.WindingOrderOfLine(l) }
func (Line) GetType() string                 { return string(geom.GeometryLineString) }
func (l Line) AsPts() []maths.Pt {
	var line []maths.Pt
	for _, p := range l {
		line = append(line, p.AsPt())
	}
	return line
}

func (l Line) Data() [][]float64 {
	ps := [][]float64{}
	for _, ll := range l {
		p := make([]float64, len(ll))
		copy(p, ll[:])
		ps = append(ps, p)
	}
	return ps
}

func (l Line) AsGeomLineString() (ln [][]float64) {
	for i := range l {
		ln = append(ln, []float64{l[i].X(), l[i].Y()})
	}
	return ln
}

func (l Line) Contains(pt Point) bool {
	pt0 := l[len(l)-1]
	ptln := maths.Line{pt.AsPt(), maths.Pt{X: pt.X() + 1, Y: pt.Y()}}
	count := 0
	for _, pt1 := range l {
		ln := maths.Line{pt0.AsPt(), pt1.AsPt()}
		if ipt, ok := maths.Intersect(ln, ptln); ok {
			if ipt.IsEqual(pt.AsPt()) {
				return false
			}
			if ln.InBetween(ipt) && ipt.X < pt.X() {
				count++
			}
		}
		pt0 = pt1
	}
	return count%2 != 0
}

func (l Line) ContainsLine(ln Line) bool {
	for _, pt := range ln {
		if !l.Contains(pt) {
			return false
		}
	}
	return true
}

func NewLine(pointPairs ...float64) Line {
	var line Line
	if (len(pointPairs) % 2) != 0 {
		panic(fmt.Sprintf("NewLine requires pair of points. %v", len(pointPairs)%2))
	}
	for i := 0; i < len(pointPairs); i += 2 {
		line = append(line, Point{pointPairs[i], pointPairs[i+1]})
	}
	return line
}

func NewLineFromPt(points ...maths.Pt) Line {
	var line Line
	for _, p := range points {
		line = append(line, Point{p.X, p.Y})
	}
	return line
}
func NewLineTruncatedFromPt(points ...maths.Pt) Line {
	var line Line
	for _, p := range points {
		line = append(line, Point{float64(int64(p.X)), float64(int64(p.Y))})
	}
	return line
}

func NewLineFromSubPoints(points ...geom.Point) (l Line) {
	l = make(Line, 0, len(points))
	for i := range points {
		l = append(l, Point{points[i].X(), points[i].Y()})
	}
	return l
}

func NewLineFrom2Float64(points ...[]float64) (l Line) {
	l = make(Line, 0, len(points))
	for i := range points {
		l = append(l, Point{points[i][0], points[i][1]})
	}
	return l
}

func (l Line) Subpoints() (points []geom.Point) {
	points = make([]geom.Point, 0, len(l))
	for i := range l {
		points = append(points, geom.Point(l[i]))
	}
	return points
}

type MultiLine []Line

func NewMultiLine(pointPairLines ...[]float64) (ml MultiLine) {
	for _, pp := range pointPairLines {
		ml = append(ml, NewLine(pp...))
	}
	return ml
}

func (MultiLine) GetType() string { return "MultiLine" }

func (MultiLine) String() string { return "MultiLine" }

func (l MultiLine) Data() [][][]float64 {
	ps := [][][]float64{}
	for _, ll := range l {
		ps = append(ps, ll.Data())
	}
	return ps
}

func (MultiLine) basicType() {}

func (ml MultiLine) Lines() (lines []geom.LineString) {
	lines = make([]geom.LineString, 0, len(ml))
	for i := range ml {
		lines = append(lines, &ml[i])
	}
	return lines
}
