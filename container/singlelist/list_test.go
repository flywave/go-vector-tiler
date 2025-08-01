package list

import (
	"fmt"
	"testing"
)

func checkListLen(t *testing.T, desc string, l *List, len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("%v: l.Len() = %d, want %d", desc, n, len)
		return false
	}
	return true
}

func checkListPointers(t *testing.T, desc string, l *List, es []*Element) {

	desc = fmt.Sprintf("%v:%v", MyCallerFileLine(), desc)
	if !checkListLen(t, desc, l, len(es)) {
		return
	}

	if len(es) == 0 {
		if l.root != nil && l.end != nil {
			t.Errorf("%v: l.root = %p, l.end = %p; both should be nil or %[2]p", desc, l.root, l.end)
		}
		return
	}

	if l.Len() > 0 {
		if l.root != l.end.Next() {
			t.Errorf("%v: l.root.Next() != %p, l.end = %p; both should  or %[2]p", desc, l.root, l.end.Next())
		}
		return
	}

	current := l.root
	for i := range es {
		if es[i] != current {
			t.Errorf("%s: es[%v] is not equal to the pos[%[2]v] of the list. expected %p == %p  ", desc, i, es[i], current)
		}

		current = current.Next()
		if current == nil {
			t.Errorf("%s: pos[%[2]v] of list is nil. Unexpected nil element", desc, i)
			return
		}
	}
}

func TestList(t *testing.T) {
	e := SliceOfElements("a", 1, 2, 3, "banana")
	l := New()
	checkListPointers(t, "Zero Element test on New List", l, []*Element{})

	l.PushFront(e[0])
	checkListPointers(t, "One element test", l, e[0:1])

	l.Remove(e[0])
	checkListPointers(t, "zero element after Remove", l, []*Element{})

	l.PushFront(e[2])
	l.PushFront(e[1])
	l.PushBack(e[3])
	l.PushBack(e[4])

	checkListPointers(t, "4 element list", l, e[1:])

	l.Remove(e[2])
	checkListPointers(t, "3 element list after removing e2", l, []*Element{e[1], e[3], e[4]})

	l.InsertBefore(e[2], e[1])
	checkListPointers(t, "4 element inserted e2 before e1", l, []*Element{e[2], e[1], e[4], e[3]})
	l.Remove(e[2])
	l.InsertBefore(e[2], e[4])
	checkListPointers(t, "4 element inserted e2 before e4", l, []*Element{e[1], e[2], e[4], e[3]})
	l.Remove(e[2])

	l.InsertAfter(e[2], e[1])
	checkListPointers(t, "4 element inserted e2 after e1", l, []*Element{e[1], e[2], e[4], e[3]})
	l.Remove(e[2])
	l.InsertAfter(e[2], e[4])
	checkListPointers(t, "4 element inserted e2 after e4", l, []*Element{e[1], e[4], e[2], e[3]})
	l.Remove(e[2])
	l.InsertAfter(e[2], e[3])
	checkListPointers(t, "4 element inserted e2 after e3", l, []*Element{e[1], e[4], e[3], e[2]})
	l.Remove(e[2])

	sum := 0
	for e := l.Front(); ; e = e.Next() {
		if elem, ok := e.(*Element); ok {
			if i, ok := elem.Value.(int); ok {
				sum += i
			}
		}
		if e == l.Back() {
			break
		}
	}
	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	l.Clear()
	checkListPointers(t, "Cleared list", l, []*Element{})

}

func checkList(t *testing.T, desc string, l *List, es []int) {
	if !checkListLen(t, desc, l, len(es)) {
		return
	}
	l.ForEachIdx(func(i int, el Elementer) bool {
		e, ok := el.(*Element)
		if !ok {
			t.Errorf("%s:elt[%d] is not of type Element.", desc, i)
			return false
		}

		le := e.Value.(int)
		if le != es[i] {
			t.Errorf("%s:elt[%d].Value = %v, want %v", desc, i, le, es[i])
		}
		return true

	})
}

func TestRemove(t *testing.T) {
	l := New()
	e := SliceOfElements(1, 2)
	l.PushBack(e[0])
	l.PushBack(e[1])
	checkListPointers(t, "List with two items", l, e)
	ef := l.Front()
	l.Remove(ef)
	checkListPointers(t, "List with only e1", l, []*Element{e[1]})
	l.Remove(ef)
	checkListPointers(t, "Noop remove", l, []*Element{e[1]})
}

func TestIssue4102(t *testing.T) {
	e1 := SliceOfElements(1, 2, 8)
	l1 := New()
	l1.PushBack(e1[0])
	l1.PushBack(e1[1])

	e2 := SliceOfElements(3, 4)
	l2 := New()
	l2.PushBack(e2[0])
	l2.PushBack(e2[1])

	ef1 := l1.Front()
	l2.Remove(ef1) // l2 should not change because ef1 is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}
	l1.InsertBefore(e1[2], ef1)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestIssue6349(t *testing.T) {
	l := New()
	l.PushBack(NewElement(1))
	l.PushBack(NewElement(2))
	e := l.Front()
	l.Remove(e)
	el := e.(*Element)
	i := el.Value.(int)
	if i != 1 {
		t.Errorf("e.value = %d, want 1", i)
	}
	if e.Next() != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.List() != nil {
		t.Errorf("e.List() != nil")
	}

}

func TestZeroList(t *testing.T) {
	var l1 = new(List)
	l1.PushFront(NewElement(1))
	checkList(t, "PushFront uninit list", l1, []int{1})

	l1.PushFront(nil)
	checkList(t, "PushFront nil value", l1, []int{1})

	var l2 = new(List)
	l2.PushBack(NewElement(1))
	checkList(t, "PushBack uninit list", l2, []int{1})

	l2.PushBack(nil)
	checkList(t, "PushBack nil", l2, []int{1})

}

func TestPushBack(t *testing.T) {
	var l List
	l.PushBack(NewElement(1))
	l.PushBack(NewElement(2))
	l.PushBack(NewElement(3))
	checkList(t, "Check insert before unknown mark", &l, []int{1, 2, 3})
}

func TestPushFront(t *testing.T) {
	var l List
	l.PushFront(NewElement(1))
	l.PushFront(NewElement(2))
	l.PushFront(NewElement(3))
	checkList(t, "Check insert before unknown mark", &l, []int{3, 2, 1})
}

func TestInsertBeforeUnknownMark(t *testing.T) {
	var l List
	l.PushBack(NewElement(1))
	l.PushBack(NewElement(2))
	l.PushBack(NewElement(3))
	l.InsertBefore(NewElement(4), new(Element))
	checkList(t, "Check insert before unknown mark", &l, []int{1, 2, 3})
}

func TestInsertAfterUnknownMark(t *testing.T) {
	var l List
	l.PushBack(NewElement(1))
	l.PushBack(NewElement(2))
	l.PushBack(NewElement(3))
	l.InsertAfter(NewElement(4), new(Element))
	checkList(t, "Check insert after unknown mark", &l, []int{1, 2, 3})
}
