package intersect

import (
	"github.com/flywave/go-vector-tiler/container/singlelist/point/list"
	"github.com/flywave/go-vector-tiler/maths"
)

type Inbound struct {
	pt    *Point
	seen  map[list.Elementer]bool
	iseen map[*Point]bool
}

func NewInbound(pt *Point) *Inbound {
	if pt == nil {
		return nil
	}
	seen := make(map[list.Elementer]bool)
	iseen := make(map[*Point]bool)
	return &Inbound{pt: pt, seen: seen, iseen: iseen}
}

func (ib *Inbound) Next() (nib *Inbound) {
	var pt *Point
	var ok bool

	for p := ib.pt.Next(); ib.pt != p; p = p.Next() {
		if p == nil {
			return
		}
		if pt, ok = p.(*Point); !ok {
			continue
		}
		ipp := asIntersect(p)
		if pt.Inward && !ib.iseen[ipp] {
			nib := NewInbound(pt)
			nib.seen = ib.seen
			nib.iseen = ib.iseen
			return nib
		}
	}
	return nil
}

func next(p list.Elementer) list.Elementer {
	switch ppt := p.(type) {
	case *Point:
		return ppt.NextWalk()
	case *SubjectPoint:
		ipt := ppt.AsIntersectPoint()
		return ipt.NextWalk()
	case *RegionPoint:
		ipt := ppt.AsIntersectPoint()
		return ipt.NextWalk()
	case list.ElementerPointer:
		return ppt.Next()
	default:
		return nil
	}
}

func asIntersect(p list.Elementer) *Point {
	switch ppt := p.(type) {
	case *Point:
		return ppt
	case *SubjectPoint:
		return ppt.AsIntersectPoint()
	case *RegionPoint:
		return ppt.AsIntersectPoint()

	default:
		return nil
	}
}

func (ib *Inbound) Walk(fn func(idx int, pt maths.Pt) bool) {

	firstInboundPoint := ib.pt
	if ib.iseen[firstInboundPoint] {
		return
	}

	if !fn(0, firstInboundPoint.Point()) {
		return
	}

	ib.seen[firstInboundPoint] = true
	ib.iseen[firstInboundPoint] = true

	for i, p := 1, next(firstInboundPoint); ; i, p = i+1, next(p) {
		ipp := asIntersect(p)
		if ipp == firstInboundPoint {
			return
		}

		if ib.seen[p] {
			return
		}

		ib.seen[p] = true
		if ipp != nil {
			ib.iseen[ipp] = true
		}

		pter, ok := p.(list.ElementerPointer)
		if !ok {
			continue
		}

		if !fn(i, pter.Point()) {
			return
		}
	}
}
