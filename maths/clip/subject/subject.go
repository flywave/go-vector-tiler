package subject

import (
	"errors"
	"log"

	"github.com/flywave/go-vector-tiler/container/singlelist/point/list"
	"github.com/flywave/go-vector-tiler/maths"
)

var ErrInvalidCoordsNumber = errors.New("Event number of coords expected.")

type Subject struct {
	winding maths.WindingOrder
	list.List
}

func New(coords []float64) (*Subject, error) {
	return new(Subject).Init(coords)
}

func (s *Subject) Init(coords []float64) (*Subject, error) {
	if len(coords)%2 != 0 {
		return nil, ErrInvalidCoordsNumber
	}
	s.winding = maths.WindingOrderOf(coords)

	for x, y := 0, 1; x < len(coords); x, y = x+2, y+2 {
		s.PushBack(list.NewPoint(coords[x], coords[y]))
	}
	return s, nil
}

func (s *Subject) Winding() maths.WindingOrder { return s.winding }

func (s *Subject) FirstPair() *Pair {
	if s == nil {
		return nil
	}
	var first, last *list.Pt
	var ok bool
	l, f := s.Back(), s.Front()
	if last, ok = l.(*list.Pt); !ok {
		return nil
	}
	if first, ok = f.(*list.Pt); !ok {
		return nil
	}
	return &Pair{
		l:   &(s.List),
		pts: [2]*list.Pt{last, first},
	}
}

func (s *Subject) GetPair(idx int) *Pair {
	p := s.FirstPair()
	log.Println(p)
	for i := 0; i < idx; i++ {
		p = p.Next()
		if p == nil {
			return p
		}
	}
	return p
}

func (s *Subject) Contains(pt maths.Pt) bool {

	line := maths.Line{pt, maths.Pt{X: pt.X - 1, Y: pt.Y}}
	count := 0
	var lpt maths.Pt
	var haveLpt bool
	for p := s.FirstPair(); p != nil; p = p.Next() {
		pline := p.AsLine()
		if ipt, ok := maths.Intersect(line, pline); ok {
			ipt = ipt.Truncate()
			if haveLpt {
				if lpt.IsEqual(ipt) {
					continue
				}
			}

			if pline.InBetween(ipt) && ipt.X < pt.X {
				count++
			}
			lpt = ipt
			haveLpt = true

		}
	}

	log.Println("Contains Count:", count)

	return count%2 != 0
}
