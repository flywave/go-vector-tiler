package list

import (
	"fmt"
	"log"
)

type Elementer interface {
	Prev() Elementer
	Next() Elementer
	SetNext(n Elementer) Elementer
	SetPrev(n Elementer) Elementer
	List() *List
	SetList(l *List) *List
}

type Sentinel struct {
	next Elementer
	prev Elementer
	list *List
}

func (s *Sentinel) SetNext(e Elementer) (oldElement Elementer) {
	oldElement = s.next
	s.next = e
	return oldElement
}

func (s *Sentinel) SetPrev(e Elementer) (oldElement Elementer) {
	oldElement = s.prev
	s.prev = e
	return oldElement
}
func (s *Sentinel) Next() Elementer {
	if p := s.next; s.list != nil && p != s.list.root {
		return p
	}
	return nil
}
func (s *Sentinel) Prev() Elementer {
	if p := s.prev; s.list != nil && p != s.list.root {
		return p
	}
	return nil
}

func (s *Sentinel) String() string {
	return fmt.Sprintf("Sentinel(%p)[list: %p][n:%p,p:%p]", s, s.list, s.next, s.prev)
}

func (s *Sentinel) List() *List {
	return s.list
}
func (s *Sentinel) SetList(l *List) (oldList *List) {
	oldList = s.list
	s.list = l
	return oldList
}

type Element struct {
	Sentinel
	Value interface{}
}

func (e Element) String() string {
	return fmt.Sprintf("%v", e.Value)
}

func NewElement(v interface{}) *Element { return &Element{Value: v} }

func SliceOfElements(vals ...interface{}) []*Element {
	els := make([]*Element, 0, len(vals))
	for _, v := range vals {
		els = append(els, NewElement(v))
	}
	return els
}

type List struct {
	root Elementer
	len  int
}

func (l *List) Init() *List {
	s := &Sentinel{}
	s.SetList(l)
	l.root = s
	l.len = 0
	return l
}

func (l *List) lazyInit() {
	if l.root == nil {
		l.Init()
	}
}

func New() *List { return new(List).Init() }

func (l *List) Len() int { return l.len }

func (l *List) Front() Elementer {
	if l.len == 0 {
		return nil
	}
	return l.root.Next()
}

func (l *List) Back() Elementer {
	if l.len == 0 {
		return nil
	}
	return l.root.Prev()
}

func (l *List) insert(e Elementer, at Elementer) Elementer {
	if e == nil {
		return e
	}
	root := l.root
	n := at.Next()
	if n == nil {
		n = root
	}

	at.SetNext(e)

	e.SetPrev(at)
	e.SetNext(n)

	n.SetPrev(e)
	e.SetList(l)

	l.len++
	return e
}

func (l *List) remove(e Elementer) Elementer {
	p := e.Prev()
	if p == nil && l.root.Next() == e {
		p = l.root
	}
	n := e.Next()
	if n == nil && l.root.Prev() == e {
		n = l.root
	}

	if p != nil {
		p.SetNext(n)
	}
	if n != nil {
		n.SetPrev(p)
	}
	e.SetNext(nil)
	e.SetPrev(nil)
	e.SetList(nil)
	l.len--
	return e
}

func (l *List) Remove(e Elementer) Elementer {
	if e.List() == l {
		l.remove(e)
	}
	return e
}

func (l *List) PushFront(e Elementer) Elementer {
	l.lazyInit()
	return l.insert(e, l.root)
}

func (l *List) PushBack(e Elementer) Elementer {
	l.lazyInit()
	p := l.root.Prev()
	if p == nil {
		p = l.root
	}
	return l.insert(e, p)
}

func (l *List) InsertBefore(e Elementer, mark Elementer) Elementer {
	if mark.List() != l {
		log.Println("List don't match.")
		return nil
	}
	p := mark.Prev()
	if p == nil {
		log.Println("Using root for previous.")
		p = l.root
	}
	return l.insert(e, p)
}

func (l *List) InsertAfter(e Elementer, mark Elementer) Elementer {
	if mark.List() != l {
		return nil
	}
	return l.insert(e, mark)
}

func (l *List) MoveToFront(e Elementer) {
	if e.List() != l || l.root.Next() == e {
		return
	}
	l.insert(l.remove(e), l.root)
}

func (l *List) MoveToBack(e Elementer) {
	if e.List() != l || l.root.Prev() == e {
		return
	}
	p := l.root.Prev()
	if p == nil {
		p = l.root
	}
	l.insert(l.remove(e), p)
}

func (l *List) MoveBefore(e, mark Elementer) {
	if e.List() != l || e == mark || mark.List() != l {
		return
	}
	p := mark.Prev()
	if p == nil {
		p = l.root
	}
	if p == e {
		return
	}
	l.insert(l.remove(e), p)
}

func (l *List) MoveAfter(e, mark Elementer) {
	if e.List() != l || e == mark || mark.List() != l {
		return
	}
	l.insert(l.remove(e), mark)
}

func (l *List) Replace(e, mark Elementer) Elementer {
	if mark.List() != l {
		return nil
	}
	n := mark.Next()
	l.Remove(mark)
	l.InsertBefore(e, n)
	return mark
}

func (l *List) FindElementForward(start, end Elementer, finder func(e Elementer) (didFind bool)) (found Elementer) {
	if l == nil || l.len == 0 {
		return nil
	}
	if start == nil {
		start = l.Front()
	}
	if end == nil {
		end = l.Back()
	}
	if start.List() != l || end.List() != l {
		return nil
	}
	sawNil := false
	for e := start; ; e = e.Next() {
		if e == nil {
			if sawNil {
				break
			}
			sawNil = true
			e = l.Front()
		}
		if finder(e) {
			return e
		}
		if e == end {
			break
		}
	}
	return nil
}

func (l *List) FindElementBackward(start, end Elementer, finder func(e Elementer) (didFind bool)) (found Elementer) {
	if l == nil || l.len == 0 {
		return nil
	}
	if start == nil {
		start = l.Back()
	}
	if end == nil {
		end = l.Front()
	}
	if start.List() != l || end.List() != l {
		return nil
	}

	for e := start; e != end.Prev(); e = e.Prev() {
		if finder(e) {
			return e
		}
	}
	return nil
}

func (l *List) IsSentinel(e Elementer) bool { return e.List() == l && e == l.root }
