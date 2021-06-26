package tile

import (
	gen "github.com/flywave/go-geom/general"
)

type Tile struct {
	Z       uint
	X       uint
	Y       uint
	Lat     float64
	Long    float64
	Extent  float64
	extent  *gen.Extent
	bufpext *gen.Extent
}
