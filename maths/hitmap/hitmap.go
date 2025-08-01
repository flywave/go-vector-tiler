package hitmap

import (
	"sort"

	geom "github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"

	"github.com/flywave/go-vector-tiler/convert"
	"github.com/flywave/go-vector-tiler/maths"
)

type Interface interface {
	LabelFor(pt maths.Pt) maths.Label
}

type allwaysInside struct{}

func (ai allwaysInside) LabelFor(_ maths.Pt) maths.Label { return maths.Inside }

var AllwaysInside allwaysInside

type bbox struct {
	f    [4]float64
	init bool
}

func (bb *bbox) Contains(pt maths.Pt) bool {
	return pt.X >= bb.f[0] && pt.Y >= bb.f[1] && pt.X <= bb.f[2] && pt.Y <= bb.f[3]
}

func (bb *bbox) Add(pts ...maths.Pt) {
	if bb == nil {
		return
	}
	for _, pt := range pts {
		if !bb.init {
			bb.f = [4]float64{pt.X, pt.Y, pt.X, pt.Y}
			bb.init = true
			return
		}
		if bb.f[0] > pt.X {
			bb.f[0] = pt.X
		}
		if bb.f[1] > pt.Y {
			bb.f[1] = pt.Y
		}
		if bb.f[2] < pt.X {
			bb.f[2] = pt.X
		}
		if bb.f[3] < pt.Y {
			bb.f[3] = pt.Y
		}
	}
}

func (bb *bbox) Coords() [4]float64 {
	if bb == nil {
		return [4]float64{}
	}
	return bb.f
}

type segEvent struct {
	x1         float64
	y1         int64
	x2         float64
	y2         int64
	m          float64
	b          float64
	isMDefined bool
}

type segEvents []segEvent

func (se segEvents) Len() int { return len(se) }
func (se segEvents) Less(i, j int) bool {
	if se[i].x1 == se[j].x1 {
		return se[i].y1 < se[j].y1
	}
	return se[i].x1 < se[j].x1
}
func (se segEvents) Swap(i, j int) { se[i], se[j] = se[j], se[i] }
func (se *segEvents) Add(l maths.Line) {
	if se == nil {
		return
	}
	var ev segEvent
	if l[0].IsEqual(l[1]) {
		return
	}
	switch {
	case l[0].X == l[1].X && l[0].Y > l[1].Y:
		fallthrough
	case l[0].X < l[1].X:
		ev.x1 = l[0].X
		ev.y1 = int64(l[0].Y * 100)
		ev.x2 = l[1].X
		ev.y2 = int64(l[1].Y * 100)
	default:
		ev.x1 = l[1].X
		ev.y1 = int64(l[1].Y * 100)
		ev.x2 = l[0].X
		ev.y2 = int64(l[0].Y * 100)
	}
	ev.m, ev.b, ev.isMDefined = l.SlopeIntercept()
	*se = append(*se, ev)
}

func (se segEvents) Contains(pt maths.Pt) (ok bool) {
	var i, count int
	var y, uy, ly int64
	var y1100, y2100 int64
	var y100 = int64(pt.Y * 100)
	for i = 0; i < len(se) && se[i].x1 <= pt.X; i++ {
		if se[i].y1 <= se[i].y2 {
			uy, ly = se[i].y1, se[i].y2
		} else {
			uy, ly = se[i].y2, se[i].y1
		}

		if y100 < uy || y100 > ly {
			continue
		}

		y1100, y2100 = se[i].y1, se[i].y2

		if y1100 == y2100 {
			if y100 == y1100 {
				if se[i].x1 <= pt.X &&
					pt.X <= se[i].x2 {
					return true
				}
				continue
			}
		}

		if y100 == y1100 && se[i].x1 < pt.X {
			if y2100 <= y100 {
				count++
			}
			continue
		}
		if y100 == y2100 && se[i].x2 < pt.X {
			if y1100 <= y100 {
				count++
			}
			continue
		}

		if !se[i].isMDefined && pt.X == se[i].x1 {
			return true
		}

		if pt.X > se[i].x2 {
			count++
			continue
		}

		y = int64((se[i].m*pt.X + se[i].b) * 100)
		if y == y100 {
			return true
		}

		if (se[i].m < 0 && y < y100) ||
			(se[i].m > 0 && y > y100) {
			count++
			continue
		}
	}
	return count%2 != 0

}

type Segment struct {
	bbox   bbox
	label  maths.Label
	events segEvents
}

func (seg Segment) Contains(pt maths.Pt) bool {
	if !seg.bbox.Contains(pt) {
		return false
	}
	return seg.events.Contains(pt)
}

func NewSegment(label maths.Label, linestring geom.LineString) (seg Segment) {

	subpts := linestring.Subpoints()

	seg.label = label
	seg.events = make(segEvents, 0, len(subpts))

	j := len(subpts) - 1
	for i := range subpts {
		l := maths.Line{
			maths.Pt{X: subpts[j].X(), Y: subpts[j].Y()},
			maths.Pt{X: subpts[i].X(), Y: subpts[i].Y()},
		}
		seg.bbox.Add(l[:]...)
		seg.events.Add(l)
		j = i
	}
	sort.Sort(seg.events)
	return seg
}

func NewSegmentFromRing(label maths.Label, ring []maths.Pt) (seg Segment) {
	seg.label = label
	seg.events = make(segEvents, 0, len(ring))

	j := len(ring) - 1
	pts := convert.FromMathPoint(ring...)
	seg.bbox.f = gen.NewExtent(pts...).Extent()
	seg.bbox.init = true
	for i := range ring {
		l := maths.Line{ring[j], ring[i]}
		seg.events.Add(l)
		j = i
	}
	sort.Sort(seg.events)
	return seg
}

func NewSegmentFromLines(label maths.Label, lines []maths.Line) (seg Segment) {
	seg.label = label
	seg.events = make(segEvents, 0, len(lines))
	for i := range lines {
		seg.bbox.Add(lines[i][:]...)
		seg.events.Add(lines[i])
	}
	sort.Sort(seg.events)
	return seg
}

type M struct {
	s      []Segment
	DoClip bool
	Clip   maths.Rectangle
}

func (hm *M) AppendSegment(seg ...Segment) *M {
	hm.s = append(hm.s, seg...)
	return hm
}

func (hm *M) LabelFor(pt maths.Pt) maths.Label {
	if hm == nil {
		return maths.Outside
	}
	if hm.DoClip {
		if !hm.Clip.Contains(pt) {
			return maths.Outside
		}
	}
	if len(hm.s) == 0 {
		return maths.Outside
	}
	for i := len(hm.s) - 1; i >= 0; i-- {
		if hm.s[i].Contains(pt) {
			return hm.s[i].label
		}
	}
	return maths.Outside
}

func NewFromPolygon(p geom.Polygon) (hm M) {
	sl := p.Sublines()
	if len(sl) == 0 {
		return hm
	}
	hm.s = make([]Segment, len(sl))
	hm.s[0] = NewSegment(maths.Inside, sl[0])
	for i := range sl[1:] {
		hm.s[i+1] = NewSegment(maths.Outside, sl[i+1])
	}
	return hm
}

func NewFromMultiPolygon(mp geom.MultiPolygon) (hm M) {
	plgs := mp.Polygons()
	for i := range plgs {
		hm.s = append(hm.s, NewFromPolygon(plgs[i]).s...)
	}
	return hm
}

func NewFromGeometry(g geom.Geometry) (hm M) {
	switch gg := g.(type) {
	case geom.Polygon:
		return NewFromPolygon(gg)
	case geom.MultiPolygon:
		return NewFromMultiPolygon(gg)
	default:
		return hm
	}
}

func NewFromLines(ln []maths.MultiLine) (hm M) {
	hm.s = make([]Segment, len(ln))
	hm.s[0] = NewSegmentFromLines(maths.Inside, ln[0])
	for i := range ln[1:] {
		hm.s[i+1] = NewSegmentFromLines(maths.Outside, ln[i+1])
	}
	return hm
}
