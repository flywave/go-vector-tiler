package tile

import (
	"context"

	gen "github.com/flywave/go-geom/general"
	"github.com/flywave/go-vector-tiler/basic"
	"github.com/flywave/go-vector-tiler/maths/simplify"
	"github.com/flywave/go-vector-tiler/maths/validate"
)

var (
	simplifyGeometries    bool = true
	simplificationMaxZoom uint = 10
)

func NewWebMercatorTiler(bound *gen.Extent, pro Provider) *Tiler {
	return &Tiler{
		Grid:       NewMercGrid(bound),
		TileExtent: 32768,
		TileBuffer: 64,
		Provider:   pro,
	}
}

type Tiler struct {
	TileExtent uint64
	TileBuffer uint64
	Grid       Grid
	Provider   Provider
}

func (m *Tiler) Tiler(zoom []uint32, cb func(*Tile, []*Layer)) {
	ctx := context.Background()
	for _, z := range zoom {
		ts := m.Grid.Iterator(z)
		for _, t := range ts {
			ls := m.Provider.GetDataByTile(t)
			var res []*Layer
			for _, l := range ls {
				newLayer := &Layer{Name: l.Name}
				for _, f := range l.Features {
					g := f.Geometry
					if m.Provider.GetSrid() != WebMercator {
						g, _ = basic.ToWebMercator(m.Provider.GetSrid(), f.Geometry)
					}
					if z < uint32(simplificationMaxZoom) && simplifyGeometries {
						g = simplify.SimplifyGeometry(g, t.ZEpislon())
					}
					g = PrepareGeo(g, t.extent, float64(m.TileExtent))
					pbb, _ := t.PixelBufferedBounds()
					clipRegion := gen.NewExtent([]float64{pbb[0], pbb[1]}, []float64{pbb[2], pbb[3]})
					g, _ = validate.CleanGeometry(ctx, g, clipRegion)
					f.Geometry = g
					newLayer.Features = append(newLayer.Features, f)
				}
				res = append(res, newLayer)
			}
			cb(t, res)
		}
	}
}