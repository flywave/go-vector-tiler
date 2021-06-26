package tile

import "github.com/flywave/go-geom/general"

func NewWebMercatorMap(name string) Map {
	return Map{
		Name:       name,
		Bounds:     nil,
		SRID:       3857,
		TileExtent: 4096,
	}
}

type Map struct {
	Name        string
	Attribution string
	Bounds      *general.Extent
	Center      [3]float64
	Layers      []Layer
	SRID        uint64
	TileExtent  uint64
	TileBuffer  uint64
}
