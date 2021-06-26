package region

import (
	"github.com/flywave/go-vector-tiler/container/singlelist/point/list"
	"github.com/flywave/go-vector-tiler/maths"
)

type Region struct {
	list.List
	sentinelPoints [4]*list.Pt
	winding        maths.WindingOrder
	aDownOrRight   [4]bool
	max, min       maths.Pt
}

func New(winding maths.WindingOrder, Min, Max maths.Pt) *Region {
	return new(Region).Init(winding, Min, Max)
}

func (r *Region) Init(winding maths.WindingOrder, Min, Max maths.Pt) *Region {
	r.winding = winding
	r.max = Max
	r.min = Min
	var pts [4][2]float64
	if winding == maths.Clockwise {
		pts = [4][2]float64{
			{Min.X, Max.Y},
			{Min.X, Min.Y},
			{Max.X, Min.Y},
			{Max.X, Max.Y},
		}
		r.aDownOrRight = [4]bool{false, true, true, false}
	} else {
		pts = [4][2]float64{{Min.X, Min.Y}, {Min.X, Max.Y}, {Max.X, Max.Y}, {Max.X, Min.Y}}
		r.aDownOrRight = [4]bool{true, true, false, false}
	}
	for i, pt := range pts {
		point := list.NewPoint(pt[0], pt[1])
		r.sentinelPoints[i] = point
		r.PushBack(point)
	}
	return r
}

func (r *Region) Axis(idx int) *Axis {
	s, e := idx%4, (idx+1)%4
	return &Axis{
		region:      r,
		idx:         s,
		pt0:         r.sentinelPoints[s],
		pt1:         r.sentinelPoints[e],
		downOrRight: r.aDownOrRight[s],
		winding:     r.winding,
	}
}
func (r *Region) FirstAxis() *Axis { return r.Axis(0) }
func (r *Region) LineString() []float64 {
	return []float64{
		r.sentinelPoints[0].Pt.X, r.sentinelPoints[0].Pt.Y,
		r.sentinelPoints[1].Pt.X, r.sentinelPoints[1].Pt.Y,
		r.sentinelPoints[2].Pt.X, r.sentinelPoints[2].Pt.Y,
		r.sentinelPoints[3].Pt.X, r.sentinelPoints[3].Pt.Y,
	}
}

func (r *Region) Max() maths.Pt                    { return r.max }
func (r *Region) Min() maths.Pt                    { return r.min }
func (r *Region) WindingOrder() maths.WindingOrder { return r.winding }
func (r *Region) Contains(pt maths.Pt) bool {
	return r.max.X > pt.X && pt.X > r.min.X &&
		r.max.Y > pt.Y && pt.Y > r.min.Y
}
func (r *Region) SentinalPoints() (pts []maths.Pt) {
	for _, p := range r.sentinelPoints {
		pts = append(pts, p.Point())
	}
	return pts
}

type Intersect struct {
	Pt        maths.Pt
	Inward    bool
	Idx       int
	isNotZero bool
}

func (r *Region) Intersections(l maths.Line) (out []Intersect, Pt1Placement, Pt2Placement PlacementCode) {
	pt1, pt2 := l[0], l[1]

	if r.Contains(pt1) && r.Contains(pt2) {
		return out, Pt1Placement, Pt2Placement
	}
	var ai [4]Intersect
	for i := 0; i < len(ai); i++ {

		a := r.Axis(i)

		Pt1Placement |= a.Placement(pt1)
		Pt2Placement |= a.Placement(pt2)

		pt, doesIntersect := a.Intersect(l)
		if !doesIntersect {
			continue
		}
		inward, err := a.IsInward(l)
		if err != nil {
			continue
		}

		out = append(out, Intersect{
			Pt:        pt,
			Inward:    inward,
			Idx:       i,
			isNotZero: true,
		})

	}
	return out, Pt1Placement, Pt2Placement
}
