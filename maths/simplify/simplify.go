package simplify

import (
	"github.com/flywave/go-geom"
	"github.com/flywave/go-vector-tiler/basic"
	"github.com/flywave/go-vector-tiler/maths"
	"github.com/flywave/go-vector-tiler/maths/points"
)

func SimplifyGeometry(g geom.Geometry, tolerance float64) geom.Geometry {
	switch gg := g.(type) {
	case geom.Polygon:
		return simplifyPolygon(gg, tolerance)

	case geom.MultiPolygon:
		var newMP basic.MultiPolygon

		for _, p := range gg.Polygons() {
			sp := simplifyPolygon(p, tolerance)
			if sp == nil {
				continue
			}
			newMP = append(newMP, sp)
		}

		if len(newMP) == 0 {
			return nil
		}

		return newMP

	case geom.LineString:
		return simplifyLineString(gg, tolerance)

	case geom.MultiLine:
		var newML basic.MultiLine

		for _, l := range gg.Lines() {
			sl := simplifyLineString(l, tolerance)
			if sl == nil {
				continue
			}
			newML = append(newML, sl)
		}

		if len(newML) == 0 {
			return nil
		}

		return newML
	}

	return g
}

func simplifyLineString(g geom.LineString, tolerance float64) basic.Line {
	line := basic.CloneLine(g)
	if len(line) <= 4 || maths.DistOfLine(g) < tolerance {
		return line
	}

	pts := line.AsPts()
	pts = DouglasPeucker(pts, tolerance)
	if len(pts) == 0 {
		return nil
	}

	return basic.NewLineTruncatedFromPt(pts...)
}

func simplifyPolygon(g geom.Polygon, tolerance float64) basic.Polygon {
	lines := g.Sublines()
	if len(lines) <= 0 {
		return nil
	}

	var poly basic.Polygon
	sqTolerance := tolerance * tolerance
	for i := range lines {
		area := maths.AreaOfPolygonLineString(lines[i])
		l := basic.CloneLine(lines[i])

		if area < sqTolerance {
			if i == 0 {
				return basic.ClonePolygon(g)
			}
			poly = append(poly, l)
			continue
		}

		pts := l.AsPts()
		if len(pts) <= 2 {
			if i == 0 {
				return nil
			}
			continue
		}

		pts = normalizePoints(pts)
		if len(pts) <= 4 {
			if i == 0 {
				return basic.ClonePolygon(g)
			}
			poly = append(poly, l)
			continue
		}

		pts = DouglasPeucker(pts, sqTolerance)
		if len(pts) <= 2 {
			if i == 0 {
				return nil
			}
			continue
		}

		poly = append(poly, basic.NewLineTruncatedFromPt(pts...))
	}

	if len(poly) == 0 {
		return nil
	}

	return poly
}

func normalizePoints(pts []maths.Pt) (pnts []maths.Pt) {
	if pts[0] == pts[len(pts)-1] {
		pts = pts[1:]
	}

	if len(pts) <= 4 {
		return pts
	}

	lpt := 0
	pnts = append(pnts, pts[0])

	for i := 1; i < len(pts); i++ {
		ni := i + 1
		if ni >= len(pts) {
			ni = 0
		}
		m1, _, sdef1 := points.SlopeIntercept(pts[lpt], pts[i])
		m2, _, sdef2 := points.SlopeIntercept(pts[lpt], pts[ni])
		if m1 != m2 || sdef1 != sdef2 {
			pnts = append(pnts, pts[i])
		}
	}

	return pnts
}
