package basic

import (
	"fmt"

	geom "github.com/flywave/go-geom"
	"github.com/flywave/go-vector-tiler/maths"
)

type Point [2]float64

func (Point) basicType() {}

func (p *Point) AsPt() maths.Pt {
	if p == nil {
		return maths.Pt{X: 0, Y: 0}
	}
	return maths.Pt{X: p[0], Y: p[1]}
}

func (bp Point) Data() []float64 {
	return bp[:]
}

func (bp Point) X() float64 {
	return bp[0]
}

func (bp Point) Y() float64 {
	return bp[1]
}

func (p Point) String() string {
	return fmt.Sprintf("Point(%v,%v)", p[0], p[1])
}

func (Point) GetType() string { return string(geom.GeometryPoint) }

type Point3 [3]float64

func (Point3) basicType() {}

func (bp Point3) X() float64 {
	return bp[0]
}

func (p Point3) String() string {
	return fmt.Sprintf("Point3(%v,%v,%v)", p[0], p[1], p[2])
}

func (Point3) GetType() string { return string(geom.GeometryPoint) }

func (bp Point3) Data() []float64 {
	return bp[:]
}

func (bp Point3) Y() float64 {
	return bp[1]
}

func (bp Point3) Z() float64 {
	return bp[2]
}

type MultiPoint []Point

func (MultiPoint) basicType() {}

func (MultiPoint) String() string { return "MultiPoint" }

func (v MultiPoint) Points() (points []geom.Point) {
	for i := range v {
		points = append(points, v[i])
	}
	return points
}

type MultiPoint3 []Point3

func (MultiPoint3) basicType() {}

func (v MultiPoint3) Points() (points []geom.Point3) {
	for i := range v {
		points = append(points, v[i])
	}
	return points
}

func (MultiPoint3) String() string {
	return "MultiPoint3"
}
