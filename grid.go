package tile

import (
	geo "github.com/flywave/go-geo"
	vec2d "github.com/flywave/go3d/float64/vec2"
)

const WGS84_PROJ4 = "+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs"
const GMERC_PROJ4 = "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0.0 +k=1.0 +units=m +nadgrids=@null +wktext +no_defs +over"

type Grid struct {
	grid     *geo.TileGrid
	currentX *uint32
	currentY *uint32
	currentZ *uint32
}

func NewGrid(bound *[4]float64, srs string) *Grid {
	conf := geo.DefaultTileGridOptions()
	conf[geo.TILEGRID_BBOX_SRS] = srs
	conf[geo.TILEGRID_SRS] = GMERC_PROJ4
	conf[geo.TILEGRID_BBOX] = &vec2d.Rect{Min: vec2d.T{bound[0], bound[1]}, Max: vec2d.T{bound[2], bound[3]}}
	grid := geo.NewTileGrid(conf)
	return &Grid{grid: grid}
}

func (g *Grid) Count(zs []uint32) int {
	c := 0
	bbx := g.grid.BBox
	for _, z := range zs {
		_, rc, _, _ := g.grid.GetAffectedLevelTiles(*bbx, int(z))
		c += rc[0] * rc[1]
	}
	return c
}

func (g *Grid) TileBounds(z uint32) (uint32, uint32, uint32, uint32) {
	bbx := g.grid.BBox
	rt, _, _, _ := g.grid.GetAffectedLevelTiles(*bbx, int(z))
	return uint32(rt.Min[0]), uint32(rt.Min[1]), uint32(rt.Max[0]), uint32(rt.Max[1])
}

func (g *Grid) Iterator(z uint32) []*Tile {
	ts := []*Tile{}
	bbx := g.grid.BBox
	rt, _, _, _ := g.grid.GetAffectedLevelTiles(*bbx, int(z))
	minx := uint32(rt.Min[0])
	miny := uint32(rt.Min[1])
	maxx := uint32(rt.Max[0])
	maxy := uint32(rt.Max[1])
	if g.currentZ != nil && *g.currentZ == z {
		minx = *g.currentX
		miny = *g.currentY
	}
	for ; miny <= maxy; miny++ {
		for ; minx <= maxx; minx++ {
			ts = append(ts, NewTile(uint32(z), uint32(minx), uint32(miny)))
		}
	}
	return ts
}

func (g *Grid) SkipTile(x, y, z uint32) {
	g.currentX = &x
	g.currentY = &y
	g.currentZ = &z
}
