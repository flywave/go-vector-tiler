package list

import (
	"fmt"
	"log"
	"strings"

	"github.com/flywave/go-vector-tiler/container/list"
	"github.com/flywave/go-vector-tiler/maths"
)

type Elementer interface {
	list.Elementer
}

type ElementerPointer interface {
	list.Elementer
	maths.Pointer
}

type Pt struct {
	maths.Pt
	list.Sentinel
}

func (p *Pt) Point() (pt maths.Pt) { return p.Pt }
func (p *Pt) String() string {
	if p == nil {
		return "(nil)"
	}
	return p.Pt.String()
}
func (p *Pt) GoString() string {
	if p == nil {
		return "(nil)"
	}
	return fmt.Sprintf("[%v,%v]", p.Pt.X, p.Pt.Y)
}

func NewPt(pt maths.Pt) *Pt {
	return &Pt{Pt: pt}
}
func NewPoint(x, y float64) *Pt {
	return &Pt{
		Pt: maths.Pt{
			X: x,
			Y: y,
		},
	}
}

func NewPointSlice(pts ...maths.Pt) (rpts []*Pt) {
	for _, pt := range pts {
		rpts = append(rpts, &Pt{Pt: pt})
	}
	return rpts
}

type List struct {
	list.List
}

func (l *List) ForEachPt(fn func(idx int, pt maths.Pt) (cont bool)) {
	for i, p := 0, l.Front(); p != nil; i, p = i+1, p.Next() {
		pt := p.(maths.Pointer).Point()
		if !fn(i, pt) {
			break
		}
	}
}

func (l *List) PushInBetween(start, end ElementerPointer, element ElementerPointer) (r bool) {
	spt := start.Point().Truncate()
	ept := end.Point().Truncate()
	var mark Elementer

	defer func() {
		if r && (element.Prev() == nil || element.Next() == nil) {
			log.Println("nil!")
			log.Printf("\tstart: %v[%[1]p] %v[%[2]p] %v[%[3]p]", start.Prev(), start, start.Next())
			log.Printf("\t   pt: %v[%[1]p] %v[%[2]p] %v[%[3]p]", element.Prev(), element, element.Next())
			log.Printf("\t  end: %v[%[1]p] %v[%[2]p] %v[%[3]p]", end.Prev(), end, end.Next())
			log.Printf("\t mark: %v[%[1]p] %v[%[2]p] %v[%[3]p]", mark.Prev(), mark, mark.Next())
			panic("Stop!")
		}
	}()

	mpt := element.Point().Truncate()
	{
		line := maths.Line{spt, ept}
		if !line.InBetween(mpt) {
			return false
		}
	}

	deltaX := ept.X - spt.X
	deltaY := ept.Y - spt.Y
	xIncreasing := deltaX > 0
	yIncreasing := deltaY > 0

	if ept.IsEqual(mpt) {
		l.InsertBefore(element, end)
		return true
	}

	if spt.IsEqual(mpt) {
		l.InsertAfter(element, start)
		return true
	}

	mark = l.FindElementForward(start.Next(), end, func(e list.Elementer) bool {
		var goodX, goodY = true, true
		if ele, ok := e.(maths.Pointer); ok {
			pt := ele.Point()

			if deltaX != 0 {
				if xIncreasing {
					goodX = int64(mpt.X) < int64(pt.X)
				} else {
					goodX = int64(mpt.X) > int64(pt.X)
				}
			}
			if deltaY != 0 {
				if yIncreasing {
					goodY = int64(mpt.Y) < int64(pt.Y)
				} else {
					goodY = int64(mpt.Y) > int64(pt.Y)
				}
			}
			return goodX && goodY
		}

		return false
	})
	if mark == nil {
		l.InsertBefore(element, end)
		return true
	}

	l.InsertBefore(element, mark)

	return true
}

func (l *List) GoString() string {
	if l == nil || l.Len() == 0 {
		return "List{}"
	}
	strs := []string{"List{"}
	for p := l.Front(); p != nil; p = p.Next() {
		pt := p.(maths.Pointer)
		strs = append(strs, fmt.Sprintf("%v(%p:%[2]T)", pt.Point(), p))
	}
	strs = append(strs, "}")
	return strings.Join(strs, "")
}

func New() *List {
	return &List{
		List: *list.New(),
	}
}
