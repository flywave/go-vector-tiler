package tile

import dvec3 "github.com/flywave/go3d/float64/vec3"

const (
	WebMercator = 3857
	WGS84       = 4326
)

var (
	WebMercatorBounds = &dvec3.Box{Min: dvec3.T{-20026376.39, -20048966.10, 0}, Max: dvec3.T{20026376.39, 20048966.10, 0}}
	WGS84Bounds       = &dvec3.Box{Min: dvec3.T{-180.0, -85.0511, 0}, Max: dvec3.T{180.0, 85.0511, 0}}
)
