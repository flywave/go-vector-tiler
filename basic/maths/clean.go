package maths

import (
	"fmt"

	geom "github.com/flywave/go-geom"

	"github.com/flywave/go-vector-tiler/basic"
	"github.com/flywave/go-vector-tiler/maths"
)

var ErrUnableToClean = fmt.Errorf("unable to clean MultiPolygon")

func cleanPolygon(p geom.Polygon) (polygons []basic.Polygon, invalids basic.Polygon) {
	if p == nil {
		return polygons, invalids
	}
	lines := p.Sublines()
	if len(lines) == 0 {
		return polygons, invalids
	}
	var currentPolygon basic.Polygon
	for _, l := range lines {
		bl := basic.CloneLine(l)
		if len(bl) == 0 {
			continue
		}
		switch bl.Direction() {
		case maths.Clockwise:
			if currentPolygon != nil {
				polygons = append(polygons, currentPolygon)
				currentPolygon = nil
			}
		case maths.CounterClockwise:
			if currentPolygon == nil {
				invalids = append(invalids, bl)
				continue
			}
		}
		currentPolygon = append(currentPolygon, bl)
	}
	if currentPolygon != nil {
		polygons = append(polygons, currentPolygon)
	}
	return polygons, invalids
}

func cleanMultiPolygon(mpolygon geom.MultiPolygon) (mp basic.MultiPolygon, err error) {
	for _, p := range mpolygon.Polygons() {
		poly, invalids := cleanPolygon(p)
		invalidLen := len(invalids)
		mpLen := len(mp)
		switch {
		case invalidLen != 0 && mpLen == 0:
			return mp, ErrUnableToClean
		case invalidLen != 0 && mpLen != 0:
			mp[len(mp)-1] = append(mp[len(mp)-1], invalids...)
			continue
		}
		mp = append(mp, poly...)
	}
	return mp, nil
}

func MakeValid(geo geom.Geometry) (basic.Geometry, error) {
	switch g := geo.(type) {
	case geom.MultiPolygon:
		return cleanMultiPolygon(g)
	}
	return basic.Clone(geo), nil
}
