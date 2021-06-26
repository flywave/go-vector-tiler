package basic

type Geometry interface {
	basicType()
	String() string
}

type Collection []Geometry

func (c Collection) Geometeries() (geometeries []G) {
	geometeries = make([]G, 0, len(c))
	for i := range c {
		geometeries = append(geometeries, G{c[i]})
	}
	return geometeries
}

func (Collection) String() string {
	return "Collection"
}

func (Collection) basicType() {}
