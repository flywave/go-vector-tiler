package list

import (
	"fmt"
	"log"
	"strings"

	list "github.com/flywave/go-vector-tiler/container/singlelist"
	"github.com/flywave/go-vector-tiler/maths"
)

type List struct {
	list.List
}

func (l *List) ForEachIdx(fn func(int, ElementerPointer) bool) {
	l.List.ForEachIdx(func(idx int, e list.Elementer) bool {
		el, ok := e.(ElementerPointer)
		if !ok {
			return true
		}
		return fn(idx, el)
	})
}

func (l *List) ForEach(fn func(ElementerPointer) bool) {
	l.List.ForEach(func(e list.Elementer) bool {
		el, ok := e.(ElementerPointer)
		if !ok {
			return true
		}
		return fn(el)
	})
}

func (l *List) ForEachPt(fn func(int, maths.Pt) bool) {
	l.ForEachIdx(func(idx int, pt ElementerPointer) bool {
		return fn(idx, pt.Point())
	})
}

func (l *List) ForEachPtBetween(start, end ElementerPointer, fn func(int, maths.Pt) bool) {
	count := 0
	l.FindElementsBetween(start, end, func(e list.Elementer) bool {
		pt, ok := e.(maths.Pointer)
		count++
		if !ok {
			return false
		}

		return !fn(count-1, pt.Point())
	})
}

func (l *List) PushInBetween(start, end ElementerPointer, element ElementerPointer) (r bool) {
	spt := start.Point().Truncate()
	ept := end.Point().Truncate()
	var mark ElementerPointer

	defer func() {
		if r && element.Next() == nil {
			log.Println("nil!")
			log.Printf("\tstart: %v[%[1]p] %v[%[2]p]", start, start.Next())
			log.Printf("\t   pt: %v[%[1]p] %v[%[2]p]", element, element.Next())
			log.Printf("\t  end: %v[%[1]p] %v[%[2]p]", end, end.Next())
			log.Printf("\t mark: %v[%[1]p] ]", mark)
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

	fmark := l.FindElementsBetween(start, end, func(e list.Elementer) bool {
		var goodX, goodY = true, true
		if ele, ok := e.(maths.Pointer); ok {
			pt := ele.Point()

			if deltaX != 0 {
				if xIncreasing {
					goodX = int64(mpt.X) <= int64(pt.X)
				} else {
					goodX = int64(mpt.X) >= int64(pt.X)
				}
			}
			if deltaY != 0 {
				if yIncreasing {
					goodY = int64(mpt.Y) <= int64(pt.Y)
				} else {
					goodY = int64(mpt.Y) >= int64(pt.Y)
				}
			}
			return goodX && goodY
		}

		return false
	})
	mark, _ = fmark.(ElementerPointer)

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
	l.ForEach(func(pt ElementerPointer) bool {
		strs = append(strs, fmt.Sprintf("%v(%p:%[2]T)", pt.Point(), pt))
		return true
	})
	strs = append(strs, "}")
	return strings.Join(strs, "")
}

func New() *List {
	return &List{
		List: *list.New(),
	}
}
