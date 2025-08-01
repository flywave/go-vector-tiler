package validate

import (
	"context"

	geom "github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
	"github.com/flywave/go-vector-tiler/basic"
	"github.com/flywave/go-vector-tiler/maths"
	"github.com/flywave/go-vector-tiler/maths/clip"
	"github.com/flywave/go-vector-tiler/maths/hitmap"
	"github.com/flywave/go-vector-tiler/maths/makevalid"
)

func CleanLinestring(g []float64) (l []float64, err error) {

	var ptsMap = make(map[maths.Pt][]int)
	var pts []maths.Pt
	i := 0
	for x, y := 0, 1; y < len(g); x, y = x+2, y+2 {

		p := maths.Pt{X: g[x], Y: g[y]}
		ptsMap[p] = append(ptsMap[p], i)
		pts = append(pts, p)
		i++
	}

	for i := 0; i < len(pts); i++ {
		pt := pts[i]
		fpts := ptsMap[pt]
		l = append(l, pt.X, pt.Y)
		if len(fpts) > 1 {
			i = fpts[len(fpts)-1]
		}
	}
	return l, nil
}

func LineAsPointPairs(l geom.LineString) (pp []float64) {
	spts := l.Subpoints()
	pp = make([]float64, 0, len(spts)*2)
	for _, pt := range spts {
		pp = append(pp, pt.X(), pt.Y())
	}
	return pp
}

func LineStringToSegments(l geom.LineString) ([]maths.Line, error) {
	ppln := LineAsPointPairs(l)
	return maths.NewSegments(ppln)
}

func makePolygonValid(ctx context.Context, hm *hitmap.M, extent *general.Extent, gs ...geom.Polygon) (mp basic.MultiPolygon, err error) {
	var plygLines []maths.MultiLine
	for _, g := range gs {
		for _, l := range g.Sublines() {
			segs, cerr := LineStringToSegments(l)
			if cerr != nil {
				return mp, cerr
			}
			plygLines = append(plygLines, segs)
			if cerr := ctx.Err(); cerr != nil {
				return mp, cerr
			}
		}
	}
	plyPoints, err := makevalid.MakeValid(ctx, hm, extent, plygLines...)
	if err != nil {
		return mp, err
	}
	for i := range plyPoints {
		var p basic.Polygon
		for j := range plyPoints[i] {
			nl := basic.NewLineFromPt(plyPoints[i][j]...)
			p = append(p, nl)
			if cerr := ctx.Err(); cerr != nil {
				return mp, cerr
			}
		}
		mp = append(mp, p)
	}
	return mp, err
}

func scalePolygon(p geom.Polygon, factor float64) (bp basic.Polygon) {
	lines := p.Sublines()
	bp = make(basic.Polygon, len(lines))
	for i := range lines {
		pts := lines[i].Subpoints()
		bp[i] = make(basic.Line, len(pts))
		for j := range pts {
			bp[i][j] = basic.Point{pts[j].X() * factor, pts[j].Y() * factor}
		}
	}
	return bp
}

func scaleMultiPolygon(p geom.MultiPolygon, factor float64) (bmp basic.MultiPolygon) {
	polygons := p.Polygons()
	bmp = make(basic.MultiPolygon, len(polygons))
	for i := range polygons {
		bmp[i] = scalePolygon(polygons[i], factor)
	}
	return bmp
}

func CleanGeometry(ctx context.Context, g geom.Geometry, extent *general.Extent) (geo geom.Geometry, err error) {
	if g == nil {
		return nil, nil
	}
	switch gg := g.(type) {
	case geom.Polygon:
		expp := scalePolygon(gg, 10.0)
		ext := extent.ScaleBy(10.0)
		hm := hitmap.NewFromGeometry(expp)
		mp, err := makePolygonValid(ctx, &hm, ext, expp)
		if err != nil {
			return nil, err
		}
		return scaleMultiPolygon(mp, 0.10), nil

	case geom.MultiPolygon:
		expp := scaleMultiPolygon(gg, 10.0)
		ext := extent.ScaleBy(10.0)
		hm := hitmap.NewFromGeometry(expp)
		mp, err := makePolygonValid(ctx, &hm, ext, expp.Polygons()...)
		if err != nil {
			return nil, err
		}
		return scaleMultiPolygon(mp, 0.10), nil

	case geom.MultiLine:
		var ml basic.MultiLine
		lns := gg.Lines()
		for i := range lns {
			nls, err := clip.LineString(lns[i], extent)
			if err != nil {
				return ml, err
			}
			ml = append(ml, nls...)
		}
		return ml, nil
	case geom.LineString:
		nls, err := clip.LineString(gg, extent)
		return basic.MultiLine(nls), err
	}
	return g, nil
}
