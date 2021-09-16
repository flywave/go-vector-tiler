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
	bbox     *vec2d.Rect
}

func NewGrid(bound *[4]float64, srs string) *Grid {
	conf := geo.DefaultTileGridOptions()
	conf[geo.TILEGRID_BBOX_SRS] = srs
	conf[geo.TILEGRID_SRS] = GMERC_PROJ4
	grid := geo.NewTileGrid(conf)
	bbx := &vec2d.Rect{Min: vec2d.T{bound[0], bound[1]}, Max: vec2d.T{bound[2], bound[3]}}
	return &Grid{grid: grid, bbox: bbx}
}

func (g *Grid) Count(zs []uint32) int {
	c := 0
	bbx := g.bbox
	for _, z := range zs {
		_, rc, _, _ := g.grid.GetAffectedLevelTiles(*bbx, int(z))
		c += rc[0] * rc[1]
	}
	return c
}

func (g *Grid) TileBounds(z uint32) (uint32, uint32, uint32, uint32) {
	bbx := g.bbox
	_, _, iter, _ := g.grid.GetAffectedLevelTiles(*bbx, int(z))
	bd := iter.GetTileBound()
	return bd[0], bd[1], bd[2], bd[3]
}

func (g *Grid) Iterator(z uint32) []*Tile {
	ts := []*Tile{}
	minx, miny, maxx, maxy := g.TileBounds(z)
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
