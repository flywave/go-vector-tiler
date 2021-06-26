package intersect

import (
	"testing"

	"github.com/flywave/go-vector-tiler/maths"
	"github.com/flywave/go-vector-tiler/maths/clip/region"
	"github.com/flywave/go-vector-tiler/maths/clip/subject"
)

func TestNewIntersect(t *testing.T) {

	sl, err := subject.New([]float64{-5, -5, -5, 5, 5, 5, 5, -5})
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
		return
	}
	rl := region.New(sl.Winding(), maths.Pt{0, 0}, maths.Pt{10, 10})

	l := New()
	inwardPt := NewPt(maths.Pt{0, 5}, true)
	l.PushBack(inwardPt)
	rl.Axis(0).PushInBetween(inwardPt.AsRegionPoint())
	if slp := sl.GetPair(2); slp != nil {
		slp.PushInBetween(inwardPt.AsSubjectPoint())
	}

	outwardPt := NewPt(maths.Pt{5, 0}, false)
	l.PushBack(outwardPt)
	rl.Axis(3).PushInBetween(outwardPt.AsRegionPoint())
	if slp := sl.GetPair(3); slp != nil {
		slp.PushInBetween(outwardPt.AsSubjectPoint())
	}

	expectedWalk := [][]maths.Pt{
		{
			{0, 5}, {5, 5}, {5, 0}, {0, 0},
		},
	}
	current := 0
	for ib := l.FirstInboundPtWalker(); ib != nil; ib = ib.Next() {
		if len(expectedWalk) <= current {
			t.Fatalf("Too many paths: expected: %v got: %v", len(expectedWalk), current)
		}
		ib.Walk(func(idx int, pt maths.Pt) bool {
			if len(expectedWalk[current]) <= idx {
				t.Fatalf("Too many points for (%v): expected: %v got: %v", current, len(expectedWalk[current]), idx)
			}
			if !expectedWalk[current][idx].IsEqual(pt) {
				t.Errorf("Point(%v) not correct of line %v: Expected: %v got %v", idx, current, expectedWalk[current][idx], pt)

			}
			if idx == 10 {
				t.Error("More then 10 paths returned!!")
				return false
			}
			return true
		})
		current++
	}

}