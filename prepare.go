package tile

import (
	"github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"
)

// PrepareGeo converts the geometry's coordinates to tile pixel coordinates. tile should be the
// extent of the tile, in the same projection as geo. pixelExtent is the dimension of the
// (square) tile in pixels usually 4096, see DefaultExtent.
// This function treats the tile extent elements as left, top, right, bottom. This is fine
// when working with a north-positive projection such as lat/long (epsg:4326)
// and web mercator (epsg:3857), but a south-positive projection (ie. epsg:2054) or west-postive
// projection would then flip the geomtery. To properly render these coordinate systems, simply
// swap the X's or Y's in the tile extent.
func PrepareGeo(geo geom.Geometry, tile *gen.Extent, pixelExtent float64) geom.Geometry {
	switch g := geo.(type) {
	case geom.Point:
		return preparept(g, tile, pixelExtent)

	case geom.MultiPoint:
		pts := g.Points()
		if len(pts) == 0 {
			return nil
		}
		mp := make([][]float64, len(pts))
		for i, pt := range g.Points() {
			mp[i] = preparept(pt, tile, pixelExtent).Data()
		}

		return gen.NewMultiPoint(mp)

	case geom.LineString:
		return preparelinestr(g, tile, pixelExtent)

	case geom.MultiLine:
		var ml [][][]float64
		for _, l := range g.Lines() {
			nl := preparelinestr(l, tile, pixelExtent).Data()
			if len(nl) > 0 {
				ml = append(ml, nl)
			}
		}
		return gen.NewMultiLineString(ml)

	case geom.Polygon:
		return preparePolygon(g, tile, pixelExtent)

	case geom.MultiPolygon:
		var mp [][][][]float64
		for _, p := range g.Polygons() {
			np := preparePolygon(p, tile, pixelExtent).Data()
			if len(np) > 0 {
				mp = append(mp, np)
			}
		}
		return gen.NewMultiPolygon(mp)
	}

	return nil
}

func preparept(g geom.Point, tile *gen.Extent, pixelExtent float64) geom.Point {
	px := int64((g.X() - tile.MinX()) / tile.XSpan() * pixelExtent)
	py := int64((tile.MaxY() - g.Y()) / tile.YSpan() * pixelExtent)
	var z float64 = 0
	if len(g.Data()) > 2 {
		z = g.Data()[2]
	}
	return gen.NewPoint([]float64{float64(px), float64(py), z})
}

func preparelinestr(g geom.LineString, tile *gen.Extent, pixelExtent float64) geom.LineString {
	pts := g
	// If the linestring
	if len(pts.Data()) < 2 {
		// Not enought points to make a line.
		return nil
	}

	ls := make([][]float64, len(pts.Data()))
	for i := 0; i < len(pts.Subpoints()); i++ {
		ls[i] = preparept(pts.Subpoints()[i], tile, pixelExtent).Data()
	}

	return gen.NewLineString(ls)
}

func preparePolygon(g geom.Polygon, tile *gen.Extent, pixelExtent float64) geom.Polygon {
	lines := gen.NewMultiLineString(g.Data())
	p := make([][][]float64, 0, len(g.Data()))

	if len(lines.Data()) == 0 {
		return gen.NewPolygon(p)
	}

	for _, line := range lines.Lines() {
		ln := preparelinestr(line, tile, pixelExtent)
		if len(ln.Data()) < 2 {
			continue
		}
		// TODO: check the last and first point to make sure
		// they are not the same, per the mvt spec
		p = append(p, ln.Data())
	}
	return gen.NewPolygon(p)
}
