package svg

import (
	"fmt"

	"github.com/flywave/go-geom"
)

type MinMax struct {
	MinX, MinY, MaxX, MaxY int64
	initialized            bool
}

func (mm MinMax) Min() (int64, int64) {
	return mm.MinX, mm.MinY
}

func (mm MinMax) Max() (int64, int64) {
	return mm.MaxX, mm.MaxY
}

func (mm MinMax) Width() int64 {
	return mm.MaxX - mm.MinX
}

func (mm MinMax) Height() int64 {
	return mm.MaxY - mm.MinY
}

func (mm MinMax) SentinalPts() [][]int64 {
	return [][]int64{
		{mm.MinX, mm.MinY},
		{mm.MaxX, mm.MinY},
		{mm.MaxX, mm.MaxY},
		{mm.MinX, mm.MaxY},
	}
}

// Remove the special case from MinMax method

// Remove special zero handling from MinMax method
func (mm *MinMax) MinMax(m1 *MinMax) *MinMax {
	if m1 == nil || !m1.initialized {
		return mm
	}
	if !mm.initialized {
		mm.MinX = m1.MinX
		mm.MinY = m1.MinY
		mm.MaxX = m1.MaxX
		mm.MaxY = m1.MaxY
		mm.initialized = true
		return mm
	}

	if m1.MinX < mm.MinX {
		mm.MinX = m1.MinX
	}
	if m1.MinY < mm.MinY {
		mm.MinY = m1.MinY
	}
	if m1.MaxX > mm.MaxX {
		mm.MaxX = m1.MaxX
	}
	if m1.MaxY > mm.MaxY {
		mm.MaxY = m1.MaxY
	}
	return mm
}

func (mm *MinMax) Fn() *MinMax                        { return mm }
func (mm *MinMax) MinMaxFn(fn func() *MinMax) *MinMax { return mm.MinMax(fn()) }

// Update MinMaxPt to handle initialization properly
func (mm *MinMax) MinMaxPt(x, y int64) *MinMax {
	if !mm.initialized {
		mm.MinX = x
		mm.MinY = y
		mm.MaxX = x
		mm.MaxY = y
		mm.initialized = true
		return mm
	}

	if x < mm.MinX {
		mm.MinX = x
	}
	if y < mm.MinY {
		mm.MinY = y
	}
	if x > mm.MaxX {
		mm.MaxX = x
	}
	if y > mm.MaxY {
		mm.MaxY = y
	}
	return mm
}
func (mm *MinMax) OfGeometry(gs ...geom.Geometry) *MinMax {
	for _, g := range gs {
		switch geo := g.(type) {
		case geom.Point:
			mm.MinMaxPt(int64(geo.X()), int64(geo.Y()))
		case geom.MultiPoint:
			for _, pt := range geo.Points() {
				mm.MinMaxPt(int64(pt.X()), int64(pt.Y()))
			}
		case geom.LineString:
			for _, pt := range geo.Subpoints() {
				mm.MinMaxPt(int64(pt.X()), int64(pt.Y()))
			}
		case geom.MultiLine:
			for _, ln := range geo.Lines() {
				for _, pt := range ln.Subpoints() {
					mm.MinMaxPt(int64(pt.X()), int64(pt.Y()))
				}
			}
		case geom.Polygon:
			for _, ln := range geo.Sublines() {
				for _, pt := range ln.Subpoints() {
					mm.MinMaxPt(int64(pt.X()), int64(pt.Y()))
				}
			}
		case geom.MultiPolygon:
			for _, p := range geo.Polygons() {
				for _, ln := range p.Sublines() {
					for _, pt := range ln.Subpoints() {
						mm.MinMaxPt(int64(pt.X()), int64(pt.Y()))
					}
				}
			}
		}
	}
	return mm
}

func (mm *MinMax) String() string {
	if mm == nil {
		return "(nil)[0 0 , 0 0]"
	}
	return fmt.Sprintf("[%v %v , %v %v]", mm.MinX, mm.MinY, mm.MaxX, mm.MaxY)
}

func (mm *MinMax) IsZero() bool {
	return mm == nil || !mm.initialized
}
func (mm *MinMax) ExpandBy(n int64) *MinMax {
	mm.MinX -= n
	mm.MinY -= n
	mm.MaxX += n
	mm.MaxY += n
	return mm
}
