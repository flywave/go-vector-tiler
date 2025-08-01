package util

import (
	"math"

	geom "github.com/flywave/go-geom"
)

func IsPointEqual(p1, p2 geom.Point) bool {
	if p1 == nil || p2 == nil {
		return p1 == p2
	}
	// 使用容差比较浮点数
	deltaX := math.Abs(p1.X() - p2.X())
	deltaY := math.Abs(p1.Y() - p2.Y())
	return deltaX < 1e-6 && deltaY < 1e-6
}

func IsPoint3Equal(p1, p2 geom.Point3) bool {
	if p1 == nil || p2 == nil {
		return p1 == p2
	}
	// 使用容差比较浮点数
	deltaX := math.Abs(p1.X() - p2.X())
	deltaY := math.Abs(p1.Y() - p2.Y())
	deltaZ := math.Abs(p1.Z() - p2.Z())
	return deltaX < 1e-6 && deltaY < 1e-6 && deltaZ < 1e-6
}

func IsMultiPointEqual(mp1, mp2 geom.MultiPoint) bool {
	pts1, pts2 := mp1.Points(), mp2.Points()
	if len(pts1) != len(pts2) {
		return false
	}
	for i, pt := range pts1 {
		if !IsPointEqual(pt, pts2[i]) {
			return false
		}
	}
	return true
}

func IsMultiPoint3Equal(mp1, mp2 geom.MultiPoint3) bool {
	pts1, pts2 := mp1.Points(), mp2.Points()
	if len(pts1) != len(pts2) {
		return false
	}
	for i, pt := range pts1 {
		if !IsPoint3Equal(pt, pts2[i]) {
			return false
		}
	}
	return true
}

func IsLineStringEqual(l1, l2 geom.LineString) bool {
	pts1, pts2 := l1.Subpoints(), l2.Subpoints()
	if len(pts1) != len(pts2) {
		return false
	}
	for i, pt := range pts1 {
		if !IsPointEqual(pt, pts2[i]) {
			return false
		}
	}
	return true
}

func IsLineString3Equal(l1, l2 geom.LineString3) bool {
	pts1, pts2 := l1.Subpoints(), l2.Subpoints()
	if len(pts1) != len(pts2) {
		return false
	}
	for i, pt := range pts1 {
		if !IsPoint3Equal(pt, pts2[i]) {
			return false
		}
	}
	return true
}

func IsMultiLineEqual(ml1, ml2 geom.MultiLine) bool {
	lns1, lns2 := ml1.Lines(), ml2.Lines()
	if len(lns1) != len(lns2) {
		return false
	}
	for i, ln := range lns1 {
		if !IsLineStringEqual(ln, lns2[i]) {
			return false
		}
	}
	return true
}

func IsMultiLine3Equal(ml1, ml2 geom.MultiLine3) bool {
	lns1, lns2 := ml1.Lines(), ml2.Lines()
	if len(lns1) != len(lns2) {
		return false
	}
	for i, ln := range lns1 {
		if !IsLineString3Equal(ln, lns2[i]) {
			return false
		}
	}
	return true
}

func IsPolygonEqual(p1, p2 geom.Polygon) bool {
	lns1, lns2 := p1.Sublines(), p2.Sublines()
	if len(lns1) != len(lns2) {
		return false
	}
	for i, ln := range lns1 {
		if !IsLineStringEqual(ln, lns2[i]) {
			return false
		}
	}
	return true
}

func IsPolygon3Equal(p1, p2 geom.Polygon3) bool {
	lns1, lns2 := p1.Sublines(), p2.Sublines()
	if len(lns1) != len(lns2) {
		return false
	}
	for i, ln := range lns1 {
		if !IsLineString3Equal(ln, lns2[i]) {
			return false
		}
	}
	return true
}

func IsMultiPolygonEqual(mp1, mp2 geom.MultiPolygon) bool {
	pgs1, pgs2 := mp1.Polygons(), mp2.Polygons()
	if len(pgs1) != len(pgs2) {
		return false
	}
	for i, pg := range pgs1 {
		if !IsPolygonEqual(pg, pgs2[i]) {
			return false
		}
	}
	return true
}

func IsMultiPolygon3Equal(mp1, mp2 geom.MultiPolygon3) bool {
	pgs1, pgs2 := mp1.Polygons(), mp2.Polygons()
	if len(pgs1) != len(pgs2) {
		return false
	}
	for i, pg := range pgs1 {
		if !IsPolygon3Equal(pg, pgs2[i]) {
			return false
		}
	}
	return true
}

func IsGeometryEqual(g1, g2 geom.Geometry) bool {
	// 处理 nil 情况
	if g1 == nil || g2 == nil {
		return g1 == g2
	}

	// 检查类型是否相同
	if g1.GetType() != g2.GetType() {
		return false
	}

	switch geo1 := g1.(type) {
	case geom.Point3:
		geo2, ok := g2.(geom.Point3)
		if !ok {
			return false
		}
		return IsPoint3Equal(geo1, geo2)
	case geom.Point:
		geo2, ok := g2.(geom.Point)
		if !ok {
			return false
		}
		return IsPointEqual(geo1, geo2)
	case geom.MultiPoint:
		geo2, ok := g2.(geom.MultiPoint)
		if !ok {
			return false
		}
		return IsMultiPointEqual(geo1, geo2)
	case geom.MultiPoint3:
		geo2, ok := g2.(geom.MultiPoint3)
		if !ok {
			return false
		}
		return IsMultiPoint3Equal(geo1, geo2)
	case geom.LineString:
		geo2, ok := g2.(geom.LineString)
		if !ok {
			return false
		}
		return IsLineStringEqual(geo1, geo2)
	case geom.LineString3:
		geo2, ok := g2.(geom.LineString3)
		if !ok {
			return false
		}
		return IsLineString3Equal(geo1, geo2)
	case geom.MultiLine:
		geo2, ok := g2.(geom.MultiLine)
		if !ok {
			return false
		}
		return IsMultiLineEqual(geo1, geo2)
	case geom.MultiLine3:
		geo2, ok := g2.(geom.MultiLine3)
		if !ok {
			return false
		}
		return IsMultiLine3Equal(geo1, geo2)
	case geom.Polygon:
		geo2, ok := g2.(geom.Polygon)
		if !ok {
			return false
		}
		return IsPolygonEqual(geo1, geo2)
	case geom.Polygon3:
		geo2, ok := g2.(geom.Polygon3)
		if !ok {
			return false
		}
		return IsPolygon3Equal(geo1, geo2)
	case geom.MultiPolygon:
		geo2, ok := g2.(geom.MultiPolygon)
		if !ok {
			return false
		}
		return IsMultiPolygonEqual(geo1, geo2)
	case geom.MultiPolygon3:
		geo2, ok := g2.(geom.MultiPolygon3)
		if !ok {
			return false
		}
		return IsMultiPolygon3Equal(geo1, geo2)
	case geom.Collection:
		geo2, ok := g2.(geom.Collection)
		if !ok {
			return false
		}
		return IsCollectionEqual(geo1, geo2)
	default:
		// 处理未知类型
		return false
	}
}

func IsCollectionEqual(c1, c2 geom.Collection) bool {
	geos1, geos2 := c1.Geometries(), c2.Geometries()
	if len(geos1) != len(geos2) {
		return false
	}
	for i, geo := range geos1 {
		if !IsGeometryEqual(geo, geos2[i]) {
			return false
		}
	}
	return true
}
