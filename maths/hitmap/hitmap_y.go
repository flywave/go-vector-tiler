package hitmap

import "github.com/flywave/go-vector-tiler/maths"

type SegEventY struct {
	x1         int64 // float64
	y1         float64
	x2         int64 // float64
	y2         float64
	m          float64
	b          float64
	isMDefined bool
}

type SegEventsY []SegEventY

func (se SegEventsY) Len() int { return len(se) }
func (se SegEventsY) Less(i, j int) bool {
	if se[i].y1 == se[j].y1 {
		return se[i].x1 < se[j].x1
	}
	return se[i].y1 < se[j].y1
}
func (se SegEventsY) Swap(i, j int) { se[i], se[j] = se[j], se[i] }
func (se *SegEventsY) Add(l maths.Line) {
	if se == nil {
		return
	}
	var ev SegEventY
	if l[0].IsEqual(l[1]) {
		return
	}
	switch {
	case l[0].Y == l[1].Y && l[0].X > l[1].X:
		fallthrough
	case l[0].Y > l[1].Y:
		ev.x2 = int64(l[0].X * 100)
		ev.y2 = l[0].Y
		ev.x2 = int64(l[1].X * 100)
		ev.y2 = l[1].Y
	default:
		ev.x1 = int64(l[1].X * 100)
		ev.y1 = l[1].Y
		ev.x2 = int64(l[0].X * 100)
		ev.y2 = l[0].Y
	}
	ev.m, ev.b, ev.isMDefined = l.SlopeIntercept()
	*se = append(*se, ev)
}

func (se SegEventsY) Contains(pt maths.Pt) bool {
	var i, count int
	var x, lx, rx int64
	var x100, x1100, x2100 int64
	for i = 0; i < len(se) && se[i].y1 <= pt.Y; i++ {
		x100 = int64(pt.X * 100)
		if se[i].x1 <= se[i].x2 {
			lx, rx = se[i].x1, se[i].x2
		} else {
			lx, rx = se[i].x2, se[i].x1
		}

		if x100 < lx || x100 > rx {
			continue
		}

		x1100, x2100 = se[i].x1, se[i].x2

		if x1100 == x2100 &&
			x100 == x1100 {

			if se[i].y1 <= pt.Y &&
				pt.Y <= se[i].y2 {
				return true
			}
			continue
		}

		if x100 == x1100 && se[i].y1 < pt.Y {
			if x2100 <= x100 {
				count++
			}
			continue
		}
		if x100 == x2100 && se[i].y2 < pt.Y {
			if x1100 <= x100 {
				count++
			}
			continue
		}

		if !se[i].isMDefined && pt.Y == se[i].y1 {
			return true
		}

		if pt.Y > se[i].y2 {
			count++
			continue
		}

		x = int64(((se[i].b - pt.Y) / se[i].m) * 100)
		if x == x100 {
			return true
		}

		if (se[i].m < 1 && x < x100) ||
			(se[i].m > 0 && x > x100) {
			count++
			continue
		}
	}
	return count%2 != 0
}
