package list

import "context"

type Elementer interface {
	Next() Elementer
	SetNext(n Elementer) Elementer
	List() *List
	SetList(l *List) *List
}

type List struct {
	root Elementer
	end  Elementer
	len  int
}

func New() *List { return new(List) }

func (l *List) Len() int { return l.len }

func (l *List) Front() Elementer {
	if l == nil || l.len == 0 {
		return nil
	}
	return l.root
}

func (l *List) Back() Elementer {
	if l == nil || l.len == 0 {
		return nil
	}
	return l.end
}

func (l *List) IsInList(e Elementer) bool {
	if e.List() != l {
		return false
	}
	return l.GetBefore(e) != nil
}

func (l *List) insert(e Elementer, at Elementer) Elementer {
	if e == nil || at == nil || !l.IsInList(at) {
		return e
	}

	if l.root == nil {
		e.SetNext(e)
		l.root = e
		l.end = e
		l.len++
		return e
	}

	n := at.Next()
	at.SetNext(e)
	e.SetNext(n)
	e.SetList(l)

	l.len++
	return e
}

func (l *List) GetBefore(m Elementer) Elementer {

	if m.List() != l {
		return nil
	}

	if l.Len() == 0 {
		return nil
	}

	if m == l.root {
		return l.end
	}

	last := l.root
	for e := l.root.Next(); e != l.root; e = e.Next() {
		if e == m {
			return last
		}
		last = e
	}
	return nil
}

func (l *List) remove(e Elementer) Elementer {

	if l.root == e {
		r := l.root
		if r.Next() == l.root {

			l.root = nil
			l.len = 0
			l.end = nil
			return e

		}
		l.root = r.Next()
		l.end.SetNext(l.root)

		l.len--
		r.SetList(nil)
		r.SetNext(nil)
		return e
	}

	p := l.GetBefore(e)

	if p == nil {
		return e
	}
	n := e.Next()
	p.SetNext(n)
	if e == l.end {
		l.end = p
	}
	l.len--
	e.SetNext(nil)
	e.SetList(nil)

	return e
}

func (l *List) Remove(e Elementer) Elementer {
	if e.List() == l {
		l.remove(e)
	}
	return e
}

func (l *List) PushFront(e Elementer) Elementer {

	if e == nil {
		return e
	}

	if l.Len() == 0 {
		e.SetList(l)
		e.SetNext(e)
		l.end = e
		l.root = e
		l.len++
		return e
	}

	if e.List() == l {
		if e == l.root {
			return e
		}
		l.remove(e)
	}
	e.SetNext(l.root)
	e.SetList(l)
	l.end.SetNext(e)
	l.root = e
	l.len++
	return e

}

func (l *List) PushBack(e Elementer) Elementer {
	if e == nil {
		return e
	}

	if l.Len() == 0 {
		e.SetList(l)
		e.SetNext(e)
		l.end = e
		l.root = e
		l.len++
		return e
	}

	if e.List() == l {
		if e == l.end {
			return e
		}
		l.remove(e)
	}
	e.SetNext(l.root)
	e.SetList(l)
	l.end.SetNext(e)
	l.end = e
	l.len++
	return e

}

func (l *List) InsertBefore(e Elementer, mark Elementer) Elementer {
	if !l.IsInList(mark) {
		return nil
	}

	if mark == l.root {
		return l.PushFront(e)

	}

	p := l.GetBefore(mark)
	if p == nil {
		return p
	}

	return l.insert(e, p)
}

func (l *List) InsertAfter(e Elementer, mark Elementer) Elementer {
	if !l.IsInList(mark) {
		return nil
	}
	if mark == l.end {
		return l.PushBack(e)
	}

	return l.insert(e, mark)
}

func (l *List) FindElementsBetween(start, end Elementer, finder func(e Elementer) (didFind bool)) (found Elementer) {
	if l == nil || l.len == 0 {
		return nil
	}
	if start == nil {
		start = l.root
	}
	if end == nil {
		end = l.end
	}

	if start.List() != l || end.List() != l {
		return nil
	}

	for e := start; ; e = e.Next() {
		if finder(e) {
			return e
		}
		if e == end {
			return nil
		}
	}
}

func (l *List) ForEach(fn func(e Elementer) bool) {
	if l == nil || l.len == 0 {
		return
	}

	for e := l.root; ; e = e.Next() {
		if !fn(e) || e == l.end {
			break
		}
	}
}

func (l *List) Range(ctx context.Context) <-chan Elementer {
	c := make(chan Elementer)

	go func() {
		var els []Elementer
		for e := l.root; ; e = e.Next() {
			els = append(els, e)
			if e == l.end {
				break
			}
		}
		for i := range els {
			select {
			case c <- els[i]:
			case <-ctx.Done():
				return
			}
		}
		close(c)
	}()
	return c
}

func (l *List) ForEachIdx(fn func(idx int, e Elementer) bool) {
	if l == nil || l.len == 0 {
		return
	}

	idx := 0
	for e := l.root; ; e = e.Next() {
		if !fn(idx, e) || e == l.end {
			break
		}
		idx++
	}
}

func (l *List) Clear() {
	if l == nil || l.len == 0 {
		return
	}

	l.end.SetNext(nil)
	for e := l.root; e != nil; {
		c := e
		e = e.Next()
		c.SetList(nil)
		c.SetNext(nil)
	}
	l.len = 0
	l.root = nil
	l.end = nil
}
