package maths

import (
	"sort"
)

type eventType uint8

const (
	LEFT eventType = iota
	RIGHT
)

type event struct {
	edge     int
	edgeType eventType
	ev       *Pt
}

func (e *event) Point() *Pt {
	return e.ev
}

func (e *event) Edge() int {
	return e.edge
}

type Eventer interface {
	Point() *Pt
	Edge() int
}

type XYOrderedEventPtr []event

func (a XYOrderedEventPtr) Len() int           { return len(a) }
func (a XYOrderedEventPtr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a XYOrderedEventPtr) Less(i, j int) bool { return XYOrder(*(a[i].ev), *(a[j].ev)) == -1 }

func NewEventQueue(segments []Line) []event {
	ne := len(segments) * 2

	eq := make([]event, ne)

	for i := range segments {
		idx := 2 * i
		eq[idx].edge = i
		eq[idx+1].edge = i
		eq[idx].ev = &(segments[i][0])
		eq[idx+1].ev = &(segments[i][1])
		if XYOrder(segments[i][0], segments[i][1]) < 0 {
			eq[idx].edgeType = LEFT
			eq[idx+1].edgeType = RIGHT
		} else {
			eq[idx].edgeType = RIGHT
			eq[idx+1].edgeType = LEFT
		}
	}
	sort.Sort(XYOrderedEventPtr(eq))
	return eq
}

func findinter_doesNotIntersect(s1x0, s1y0, s1x1, s1y1, s2x0, s2y0, s2x1, s2y1 float64) bool {
	var swap float64

	if s1x0 > s1x1 {
		swap = s1x0
		s1x0 = s1x1
		s1x1 = swap

		swap = s1y0
		s1y0 = s1y1
		s1y1 = swap
	} else {
		if s1x0 == s1x1 && s1y0 > s1y1 {
			swap = s1x0
			s1x0 = s1x1
			s1x1 = swap

			swap = s1y0
			s1y0 = s1y1
			s1y1 = swap
		}
	}
	if s2x0 > s2x1 {
		swap = s2x0
		s2x0 = s2x1
		s2x1 = swap

		swap = s2y0
		s2y0 = s2y1
		s2y1 = swap
	} else {
		if s2x0 == s2x1 && s2y0 > s2y1 {
			swap = s2x0
			s2x0 = s2x1
			s2x1 = swap

			swap = s2y0
			s2y0 = s2y1
			s2y1 = swap
		}
	}

	if ((((s1x1 - s1x0) * (s2y0 - s1y0)) - ((s1y1 - s1y0) * (s2x0 - s1x0))) * (((s1x1 - s1x0) * (s2y1 - s1y0)) - ((s1y1 - s1y0) * (s2x1 - s1x0)))) > 0 {
		return true
	}
	if ((((s2x1 - s2x0) * (s1y0 - s2y0)) - ((s2y1 - s2y0) * (s1x0 - s2x0))) * (((s2x1 - s2x0) * (s1y1 - s2y0)) - ((s2y1 - s2y0) * (s1x1 - s2x0)))) > 0 {
		return true
	}

	return false

}

func DoesIntersect(s1, s2 Line) bool {
	switch {
	case s1[0].X > s1[1].X:
		s1[0].X, s1[0].Y, s1[1].X, s1[1].Y = s1[1].X, s1[1].Y, s1[0].X, s1[0].Y
	case s1[0].X < s1[1].X:
	case s1[0].Y > s1[1].Y:
		s1[0].X, s1[0].Y, s1[1].X, s1[1].Y = s1[1].X, s1[1].Y, s1[0].X, s1[0].Y
	}
	switch {
	case s2[0].X > s2[1].X:
		s2[0].X, s2[0].Y, s2[1].X, s2[1].Y = s2[1].X, s2[1].Y, s2[0].X, s2[0].Y
	case s2[0].X < s2[1].X:
	case s2[0].Y > s2[1].Y:
		s2[0].X, s2[0].Y, s2[1].X, s2[1].Y = s2[1].X, s2[1].Y, s2[0].X, s2[0].Y
	}

	s1sign := ((((s1[1].X - s1[0].X) * (s2[0].Y - s1[0].Y)) - ((s1[1].Y - s1[0].Y) * (s2[0].X - s1[0].X))) * (((s1[1].X - s1[0].X) * (s2[1].Y - s1[0].Y)) - ((s1[1].Y - s1[0].Y) * (s2[1].X - s1[0].X)))) > 0
	s2sign := ((((s2[1].X - s2[0].X) * (s1[0].Y - s2[0].Y)) - ((s2[1].Y - s2[0].Y) * (s1[0].X - s2[0].X))) * (((s2[1].X - s2[0].X) * (s1[1].Y - s2[0].Y)) - ((s2[1].Y - s2[0].Y) * (s1[1].X - s2[0].X)))) > 0

	return !(s1sign || s2sign)

}

type intersectfn [2]Line

func (ifn intersectfn) PtFn() Pt {
	pt, _ := Intersect(ifn[0], ifn[1])
	return pt
}

func FindIntersectsWithEventQueue(polygonCheck bool, eq []event, segments []Line, fn func(srcIdx, destIdx int, ptfn func() Pt) bool) {
	ns := len(segments)
	var val struct{}

	isegmap := make(map[int]struct{})
	for _, ev := range eq {
		_, ok := isegmap[ev.edge]

		if !ok {
			isegmap[ev.edge] = val
			continue
		}

		delete(isegmap, ev.edge)
		if len(isegmap) == 0 {
			continue
		}
		edge := segments[ev.edge]
		sslice := make([]int, 0, len(isegmap))
		for s := range isegmap {
			sslice = append(sslice, s)
		}

		for _, s := range sslice {

			src, dest := (s+1)%ns, (ev.edge+1)%ns

			if ev.edge == s {
				continue
			}
			if polygonCheck && (src == ev.edge || dest == s) {
				continue
			}

			sedge := segments[s]
			if !DoesIntersect(edge, sedge) {
				continue
			}

			src, dest = ev.edge, s
			if src > dest {
				src, dest = dest, src
			}
			ptfn := intersectfn{edge, sedge}
			if !fn(src, dest, ptfn.PtFn) {
				return
			}
		}
	}
}

func FindIntersectsWithEventQueueWithoutIntersectNew(polygonCheck bool, eq []event, segments []Line, fn func(srcIdx, destIdx int) bool) {
	ns := len(segments)
	isegmap := make([]uint8, ns)
	var haveSeenAll uint

	for i, ev := range eq {
		edgeidx := eq[i].edge

		if isegmap[edgeidx] < 2 {
			isegmap[edgeidx] = 1
			haveSeenAll++
			continue
		}

		isegmap[edgeidx] = 2
		haveSeenAll--
		if haveSeenAll == 0 {
			continue
		}

		for s := range isegmap {
			if isegmap[s] != 1 || edgeidx == s {
				continue
			}

			src, dest := s+1, ev.edge+1
			if dest >= ns {
				dest = 0
			}
			if src >= ns {
				src = 0
			}

			if polygonCheck && (src == edgeidx || dest == s) {
				continue
			}

			if !DoesIntersect(segments[edgeidx], segments[s]) {
				continue
			}

			src, dest = edgeidx, s
			if src > dest {
				src, dest = dest, src
			}
			if !fn(src, dest) {
				return
			}
		}
	}
}

func FindIntersectsWithEventQueueWithoutIntersect(polygonCheck bool, eq []event, segments []Line, fn func(srcIdx, destIdx int) bool) {
	ns := len(segments)

	isegmap := make(map[int]bool, ns)
	seenEdgeCount := 0

	for i := range eq {
		edgeidx := eq[i].edge

		if !isegmap[edgeidx] {
			isegmap[edgeidx] = true
			seenEdgeCount++
			continue
		}

		isegmap[edgeidx] = false
		seenEdgeCount--
		if seenEdgeCount <= 0 {
			seenEdgeCount = 0
			continue
		}

		for s, sv := range isegmap {

			if edgeidx == s || !sv {
				continue
			}
			src, dest := s+1, edgeidx+1
			if dest >= ns {
				dest = 0
			}
			if src >= ns {
				src = 0
			}

			if polygonCheck && (src == edgeidx || dest == s) {
				continue
			}

			if !DoesIntersect(segments[edgeidx], segments[s]) {
				continue
			}

			src, dest = edgeidx, s
			if src > dest {
				src, dest = dest, src
			}
			if !fn(src, dest) {
				return
			}
		}
	}
}

func FindIntersectsWithEventQueueWithoutIntersectNotPolygon(eq []event, segments []Line, fn func(srcIdx, destIdx int) bool) {
	ns := len(segments)

	isegmap := make(map[int]bool, ns)
	seenEdgeCount := 0
	var shouldReturn bool

	for i := range eq {
		edgeidx := eq[i].edge

		if !isegmap[edgeidx] {
			isegmap[edgeidx] = true
			seenEdgeCount++
			continue
		}

		isegmap[edgeidx] = false
		seenEdgeCount--
		if seenEdgeCount <= 0 {
			seenEdgeCount = 0
			continue
		}

		for s, sv := range isegmap {
			if edgeidx == s || !sv {
				continue
			}
			if !DoesIntersect(segments[edgeidx], segments[s]) {
				continue
			}
			if edgeidx <= s {
				shouldReturn = !fn(edgeidx, s)
			} else {
				shouldReturn = !fn(s, edgeidx)
			}
			if shouldReturn {
				return
			}
		}
	}
}

func FindAllIntersectsWithEventQueueWithoutIntersectNotPolygon(eq []event, segments []Line, skipfn func(srcIdx, destIdx int) bool, fn func(srcIdx, destIdx int)) {
	ns := len(segments)
	isegmap := make([]bool, ns)
	seenEdgeCount := 0
	var edgeidx, s int
	var sv bool
	for i := range eq {
		edgeidx = eq[i].edge

		if !isegmap[edgeidx] {
			isegmap[edgeidx] = true
			seenEdgeCount++
			continue
		}

		isegmap[edgeidx] = false
		seenEdgeCount = seenEdgeCount - 1
		if seenEdgeCount <= 0 {
			seenEdgeCount = 0
			continue
		}

		for s, sv = range isegmap {
			if !sv {
				continue
			}
			if skipfn(edgeidx, s) {
				continue
			}
			if segments[edgeidx][0].X == segments[s][0].X && segments[edgeidx][0].Y == segments[s][0].Y {
				continue
			}
			if segments[edgeidx][0].X == segments[s][1].X && segments[edgeidx][0].Y == segments[s][1].Y {
				continue
			}
			if segments[edgeidx][1].X == segments[s][0].X && segments[edgeidx][1].Y == segments[s][0].Y {
				continue
			}
			if segments[edgeidx][1].X == segments[s][1].X && segments[edgeidx][1].Y == segments[s][1].Y {
				continue
			}

			if findinter_doesNotIntersect(segments[edgeidx][0].X, segments[edgeidx][0].Y, segments[edgeidx][1].X, segments[edgeidx][1].Y, segments[s][0].X, segments[s][0].Y, segments[s][1].X, segments[s][1].Y) {
				continue
			}
			fn(edgeidx, s)
		}
	}
}

func FindIntersectsWithoutIntersect(segments []Line, fn func(srcIdx, destIdx int) bool) {
	eq := NewEventQueue(segments)
	FindIntersectsWithEventQueueWithoutIntersectNotPolygon(eq, segments, fn)
}

func FindIntersects(segments []Line, fn func(srcIdx, destIdx int, ptfn func() Pt) bool) {
	eq := NewEventQueue(segments)
	FindIntersectsWithEventQueue(false, eq, segments, fn)
}

func FindPolygonIntersects(segments []Line, fn func(srcIdx, destIdx int, ptfn func() Pt) bool) {
	ns := len(segments)
	if ns < 3 {
		return
	}
	eq := NewEventQueue(segments)
	FindIntersectsWithEventQueue(true, eq, segments, fn)
}

func (s1 Line) DoesIntersect(s2 Line) bool {
	return DoesIntersect(s1, s2)
}

func (l Line) IntersectsLines(lines []Line, fn func(idx int) bool) {
	if len(lines) == 0 {
		return
	}
	if len(lines) == 1 {
		if l.DoesIntersect(lines[0]) {
			fn(0)
		}
		return
	}

	eq := NewEventQueue(append([]Line{l}, lines...))
	var val struct{}

	isegmap := make(map[int]struct{})

	for _, ev := range eq {

		if _, ok := isegmap[ev.edge]; !ok {
			isegmap[ev.edge] = val
			continue
		}

		delete(isegmap, ev.edge)
		if len(isegmap) == 0 {
			continue
		}

		if ev.edge == 0 {
			for s := range isegmap {
				if !l.DoesIntersect(lines[s-1]) {
					continue
				}
				if !fn(s - 1) {
					return
				}
			}
			return
		}

		if _, ok := isegmap[0]; !ok {
			continue
		}

		if !l.DoesIntersect(lines[ev.edge-1]) {
			continue
		}
		if !fn(ev.edge - 1) {
			return
		}
		continue
	}
}

func (l Line) XYOrderedPtsIdx() (left, right int) {
	if XYOrder(l[0], l[1]) == -1 {
		return 0, 1
	}
	return 1, 0
}
