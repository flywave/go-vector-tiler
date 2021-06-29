package basic

import (
	"log"

	"github.com/flywave/go-vector-tiler/maths"
	"github.com/flywave/go-vector-tiler/util"
)

func (l Line) IsValid() bool {
	var seen map[string]struct{}
	for _, pt := range l {
		if _, ok := seen[pt.String()]; ok {
			return false
		}
	}
	pt0 := l[len(l)-1]
	endj := len(l) - 1
	for i, pt1 := range l[:len(l)-2] {
		inpt0 := l[i+1]
		for _, inpt1 := range l[i+2 : endj] {
			if util.IsPointEqual(pt0, inpt0) && util.IsPointEqual(pt1, inpt1) {
				continue
			}
			l1, l2 := maths.Line{pt0.AsPt(), pt1.AsPt()}, maths.Line{inpt0.AsPt(), inpt1.AsPt()}
			if ipt, ok := maths.Intersect(l1, l2); ok {
				if l1.InBetween(ipt) && l2.InBetween(ipt) {
					return false
				}
			}
			inpt0 = inpt1
		}
		pt0 = pt1
		endj = len(l)
	}
	return true
}

func (p Polygon) IsValid() bool {
	if len(p) == 0 {
		return false
	}

	if !(p[0].IsValid() && p[0].Direction() == maths.Clockwise) {
		log.Println("Line 0", p[0].IsValid(), p[0].Direction())
		return false
	}

	for _, l := range p[1:] {
		if !(l.IsValid() && l.Direction() == maths.CounterClockwise && p[0].ContainsLine(l)) {
			return false
		}
	}
	return true
}
