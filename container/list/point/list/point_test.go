package list

import (
	"fmt"
	"testing"

	"github.com/flywave/go-vector-tiler/maths"
	"github.com/gdey/tbltest"
)

func checkListLen(t *testing.T, desc string, l *List, es []*Pt) bool {
	eslen := len(es)
	if n := l.Len(); n != eslen {
		t.Errorf("%v: got l.Len() = %d, want %d: %#v,%#v", desc, n, eslen, l, es)

		return false
	}
	return true
}

func checkListPointers(t *testing.T, desc string, l *List, es []*Pt) {

	if !checkListLen(t, desc, l, es) {
		return
	}
	if len(es) == 0 {
		return
	}
	var didError bool

	for i, e := range es {
		var Next, Prev *Pt
		if i > 0 {
			Prev = es[i-1]
		}

		if p := e.Prev(); (p != nil || Prev != nil) && p != Prev {
			t.Errorf("%s: elt[%d](%p).Prev() = %p, want %p", desc, i, e, p, Prev)
			didError = true
		}
		if i < len(es)-1 {
			Next = es[i+1]
		}
		if n := e.Next(); (n != nil || Next != nil) && n != Next {
			t.Errorf("%s: elt[%d](%p).Next() = %p, want %p", desc, i, e, n, Next)
			didError = true
		}
	}
	if didError {
		t.Errorf("list:%#v", l)
	}
}

func checkListInBetween(t *testing.T, desc string, i maths.Pt, loc int, mpts ...maths.Pt) {

	l := New()
	pts := NewPointSlice(mpts...)
	insert := NewPoint(i.X, i.Y)
	offset := 1
	if loc < 0 {
		offset = 0
	}

	cpts := make([]*Pt, len(pts)+offset)
	for i, p := range pts {
		l.PushBack(p)
		switch {
		case i < loc:
			cpts[i] = p
		case i == loc:
			cpts[i] = insert
			fallthrough
		case i > loc:
			cpts[i+offset] = p
		}

	}
	if loc >= len(mpts) {
		cpts[loc] = insert
	}

	checkListPointers(t, fmt.Sprintf("list check: %v", desc), l, pts)
	l.PushInBetween(pts[0], pts[len(pts)-1], insert)
	checkListPointers(t, desc, l, cpts)
}

func TestPushInBetween(t *testing.T) {

	type testcase struct {
		desc      string
		insertPt  maths.Pt
		pointList []maths.Pt
		pos       int
	}
	tests := tbltest.Cases(
		testcase{
			"Simple two point(3,1), after 1,1.",
			maths.Pt{X: 3, Y: 1},
			[]maths.Pt{
				{X: 1, Y: 1},
				{X: 4, Y: 1},
			},
			1,
		},
		testcase{
			"Simple three point(3,1), after 2,1.",
			maths.Pt{X: 3, Y: 1},
			[]maths.Pt{
				{X: 1, Y: 1},
				{X: 2, Y: 1},
				{X: 4, Y: 1},
			},
			2,
		},
		testcase{
			"Simple three point(-1,1), Not included.",
			maths.Pt{X: -1, Y: 1},
			[]maths.Pt{
				{X: 1, Y: 1},
				{X: 2, Y: 1},
				{X: 4, Y: 1},
			},
			-1,
		},
		testcase{
			"Dup three point(1,1), after 1,1.",
			maths.Pt{X: 1, Y: 1},
			[]maths.Pt{
				{X: 1, Y: 1},
				{X: 2, Y: 1},
				{X: 4, Y: 1},
			},
			1,
		},
		testcase{
			"Dup three point(2,1), after 2,1.",
			maths.Pt{X: 2, Y: 1},
			[]maths.Pt{
				{X: 1, Y: 1},
				{X: 2, Y: 1},
				{X: 4, Y: 1},
			},
			2,
		},
	)
	tests.Run(func(idx int, test testcase) {
		checkListInBetween(t, test.desc, test.insertPt, test.pos, test.pointList...)
	})
}
