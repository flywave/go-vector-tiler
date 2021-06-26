package list

type Sentinel struct {
	next Elementer
	list *List
}

func (s *Sentinel) Next() Elementer {
	if s == nil {
		return nil
	}
	return s.next
}

func (s *Sentinel) SetNext(e Elementer) Elementer {
	if s == nil {
		return nil
	}
	n := s.next
	s.next = e
	return n
}

func (s *Sentinel) List() *List {
	if s == nil {
		return nil
	}
	return s.list
}

func (s *Sentinel) SetList(l *List) *List {
	if s == nil {
		return nil
	}
	ol := s.list
	s.list = l
	return ol
}
