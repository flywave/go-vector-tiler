package tile

import (
	"math"

	"github.com/flywave/go-geom/general"
	"github.com/flywave/go-vector-tiler/maths/webmercator"
)

type Grid interface {
	Iterator(z uint32) []*Tile
}

type MercGrid struct {
	Bounds *general.Extent
}

func NewMercGrid(bound *general.Extent) *MercGrid {
	if bound == nil {
		bound = &general.Extent{webmercator.MinXExtent, webmercator.MinYExtent, webmercator.MaxXExtent, webmercator.MaxYExtent}
	}
	return &MercGrid{Bounds: bound}
}

func (g *MercGrid) Iterator(z uint32) []*Tile {
	lvlCount := math.Pow(2, float64(z))
	span := webmercator.MaxXExtent * 2 / lvlCount
	minx := math.Floor(g.Bounds[0] / span)
	miny := math.Floor(g.Bounds[1] / span)
	maxx := math.Ceil(g.Bounds[2] / span)
	maxy := math.Ceil(g.Bounds[3] / span)
	ts := []*Tile{}
	for x := minx; x < maxx; x++ {
		for y := miny; y < maxy; y++ {
			ts = append(ts, NewTile(uint(z), uint(x), uint(y)))
		}
	}
	return ts
}
