package basic

import geom "github.com/flywave/go-geom"

func ClonePoint(pt geom.Point) Point {
	return Point{pt.X(), pt.Y()}
}

func ClonePoint3(pt geom.Point3) Point3 {
	return Point3{pt.X(), pt.Y(), pt.Z()}
}

func CloneMultiPoint(mpt geom.MultiPoint) MultiPoint {
	var bmpt MultiPoint
	for _, pt := range mpt.Points() {
		bmpt = append(bmpt, ClonePoint(pt))
	}
	return bmpt
}

func CloneLine(line geom.LineString) (l Line) {
	for _, pt := range line.Subpoints() {
		l = append(l, Point{pt.X(), pt.Y()})
	}
	return l
}

func CloneMultiLine(mline geom.MultiLine) (ml MultiLine) {
	for _, ln := range mline.Lines() {
		ml = append(ml, CloneLine(ln))
	}
	return ml
}

func ClonePolygon(polygon geom.Polygon) (ply Polygon) {
	for _, ln := range polygon.Sublines() {
		ply = append(ply, CloneLine(ln))
	}
	return ply
}

func CloneMultiPolygon(mpolygon geom.MultiPolygon) (mply MultiPolygon) {
	for _, ply := range mpolygon.Polygons() {
		mply = append(mply, ClonePolygon(ply))
	}
	return mply
}

func Clone(geo geom.Geometry) Geometry {
	switch g := geo.(type) {
	case geom.Point:
		return ClonePoint(g)
	case geom.MultiPoint:
		return CloneMultiPoint(g)
	case geom.LineString:
		return CloneLine(g)
	case geom.MultiLine:
		return CloneMultiLine(g)
	case geom.Polygon:
		return ClonePolygon(g)
	case geom.MultiPolygon:
		return CloneMultiPolygon(g)
	}
	return nil
}
