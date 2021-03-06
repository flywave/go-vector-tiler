package basic

import (
	geom "github.com/flywave/go-geom"

	"github.com/flywave/go-vector-tiler/maths"
)

type Polygon []Line

func (Polygon) basicType() {}

func (p Polygon) Sublines() (slines []geom.LineString) {
	slines = make([]geom.LineString, 0, len(p))
	for i := range p {
		slines = append(slines, p[i])
	}
	return slines
}
func (Polygon) String() string {
	return "Polygon"
}

func (Polygon) GetType() string { return string(geom.GeometryPolygon) }

func (p Polygon) Data() [][][]float64 {
	ps := [][][]float64{}
	for _, ll := range p {
		ps = append(ps, ll.Data())
	}
	return ps
}

type MultiPolygon []Polygon

func (MultiPolygon) basicType() {}

func (mp MultiPolygon) Polygons() (polygons []geom.Polygon) {
	polygons = make([]geom.Polygon, 0, len(mp))
	for i := range mp {
		polygons = append(polygons, mp[i])
	}
	return polygons
}
func (MultiPolygon) String() string {
	return "MultiPolygon"
}

func (p MultiPolygon) Data() [][][][]float64 {
	ps := [][][][]float64{}
	for _, ll := range p {
		ps = append(ps, ll.Data())
	}
	return ps
}

func (MultiPolygon) GetType() string { return string(geom.GeometryMultiPolygon) }

func NewPolygon(main []maths.Pt, clines ...[]maths.Pt) Polygon {
	p := Polygon{NewLineFromPt(main...)}
	for _, l := range clines {
		p = append(p, NewLineFromPt(l...))
	}
	return p
}
func NewPolygonFromSubLines(lines ...geom.LineString) (p Polygon) {
	p = make(Polygon, 0, len(lines))
	for i := range lines {
		l := NewLineFromSubPoints(lines[i].Subpoints()...)
		p = append(p, l)
	}
	return p
}

func NewMultiPolygonFromPolygons(polygons ...geom.Polygon) (mp MultiPolygon) {
	mp = make(MultiPolygon, 0, len(polygons))
	for i := range polygons {
		p := NewPolygonFromSubLines(polygons[i].Sublines()...)
		mp = append(mp, p)
	}
	return mp
}
