package tile

import (
	"math"

	"github.com/flywave/go-geom/general"
	"github.com/flywave/go-proj"
	proj "github.com/flywave/go-proj"
	"github.com/flywave/go-vector-tiler/maths/webmercator"
)

const WGS84_PROJ4 = "+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs"
const GMERC_PROJ4 = "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0.0 +k=1.0 +units=m +nadgrids=@null +wktext +no_defs +over"

type Grid struct {
	Bounds   *general.Extent
	currentX *uint32
	currentY *uint32
	currentZ *uint32
}

func NewMercGrid(bound *general.Extent) *Grid {
	if bound == nil {
		bound = &general.Extent{webmercator.MinXExtent, webmercator.MinYExtent, webmercator.MaxXExtent, webmercator.MaxYExtent}
	}
	return &Grid{Bounds: bound}
}

func NewGrid(bound *general.Extent, srs string) *Grid {
	if bound == nil {
		bound = &general.Extent{webmercator.MinXExtent, webmercator.MinYExtent, webmercator.MaxXExtent, webmercator.MaxYExtent}
	} else {
		p1, _ := proj.NewProj(srs)
		p2, _ := proj.NewProj(GMERC_PROJ4)
		tran, _ := proj.NewTransformation(p1, p2)
		x := []float64{bound[0], bound[2]}
		y := []float64{bound[1], bound[3]}
		x, y, _, _ = tran.Transform(x, y, nil)
		bound = &general.Extent{x[0], y[0], x[1], y[1]}
	}
	return &Grid{Bounds: bound}
}

func (g *Grid) Count(zs []uint32) int {
	c := 0
	for _, z := range zs {
		c += len(g.Iterator(z))
	}
	return c
}

func (g *Grid) Iterator(z uint32) []*Tile {
	if g.currentZ != nil && z < *g.currentZ {
		return nil
	}
	lvlCount := math.Pow(2, float64(z))
	span := webmercator.MaxXExtent * 2 / lvlCount
	minx := math.Floor(g.Bounds[0] / span)
	miny := math.Floor(g.Bounds[1] / span)
	maxx := math.Ceil(g.Bounds[2] / span)
	maxy := math.Ceil(g.Bounds[3] / span)
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
			ts = append(ts, NewTile(uint32(z), uint32(x), uint32(y)))
		}
	}
	return ts
}

func (g *Grid) SkipTile(x, y, z uint32) {
	g.currentX = &x
	g.currentY = &y
	g.currentZ = &z
}
