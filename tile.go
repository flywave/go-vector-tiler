package tile

import (
	dvec3 "github.com/flywave/go3d/float64/vec3"
)

type Tile struct {
	Z      uint
	X      uint
	Y      uint
	Lat    float64
	Long   float64
	Extent float64

	extent  *dvec3.Box
	bufpext *dvec3.Box
}
