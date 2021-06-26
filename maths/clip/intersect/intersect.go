package intersect

import (
	"github.com/flywave/go-vector-tiler/container/singlelist/point/list"
)

type Intersect struct {
	list.List
}

func New() *Intersect {
	l := new(Intersect)

	return l
}

func (i *Intersect) FirstInboundPtWalker() *Inbound {
	if i == nil || i.Len() == 0 {
		return nil
	}
	var ok bool
	var pt *Point
	for p := i.Front(); ; p = p.Next() {
		if pt, ok = p.(*Point); ok && pt.Inward {
			break
		}
		if p == i.Back() {
			return nil
		}
	}
	return NewInbound(pt)
}

func (i *Intersect) ForEach(fn func(*Point) bool) {
	i.List.ForEach(func(e list.ElementerPointer) bool {
		pt, ok := e.(*Point)
		if !ok {
			return true
		}
		return fn(pt)
	})
}
