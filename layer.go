package tile

import (
	geom "github.com/flywave/go-geom"
)

type Layer struct {
	Name     string
	Features []*geom.Feature
	SRID     int
}

func (l *Layer) GetName() string {
	return l.Name
}
