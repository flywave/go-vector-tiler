package tile

import (
	"math"

	"github.com/flywave/go-geom/general"
	"github.com/flywave/go-vector-tiler/maths/webmercator"
)

type Grid interface {
	Iterator(z uint32) []*Tile
	SkipTile(x, y, z uint32)
}

type MercGrid struct {
	Bounds   *general.Extent
	currentX *uint32
	currentY *uint32
	currentZ *uint32
}

func NewMercGrid(bound *general.Extent) *MercGrid {
	if bound == nil {
		bound = &general.Extent{webmercator.MinXExtent, webmercator.MinYExtent, webmercator.MaxXExtent, webmercator.MaxYExtent}
	}
	return &MercGrid{Bounds: bound}
}

func (g *MercGrid) Count(zs []uint32) int {
	c := 0
	for _, z := range zs {
		c += len(g.Iterator(z))
	}
	return c
}

func (g *MercGrid) Iterator(z uint32) []*Tile {
	if g.currentZ != nil && z < *g.currentZ {
		return nil
	}
	lvlCount := math.Pow(2, float64(z))
	span := webmercator.MaxXExtent * 2 / lvlCount
	minx := math.Floor((g.Bounds[0] + webmercator.MaxXExtent) / span)
	miny := math.Floor((webmercator.MaxXExtent - g.Bounds[1]) / span)
	maxx := math.Ceil((g.Bounds[2] + webmercator.MaxXExtent) / span)
	maxy := math.Ceil((webmercator.MaxXExtent - g.Bounds[3]) / span)
	ts := []*Tile{}

	if g.currentZ != nil && z == *g.currentZ {
		if g.currentX != nil {
			minx = float64(*g.currentX)
		}
		if g.currentY != nil {
			miny = float64(*g.currentY)
		}
	}
	for x := minx; x < maxx; x++ {
		for y := miny; y < maxy; y++ {
			ts = append(ts, NewTile(uint(z), uint(x), uint(y)))
		}
	}
	return ts
}

func (g *MercGrid) SkipTile(x, y, z uint32) {
	g.currentX = &x
	g.currentY = &y
	g.currentZ = &z
}
