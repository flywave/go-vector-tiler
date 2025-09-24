package tile

import (
	"fmt"

	gen "github.com/flywave/go-geom/general"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-vector-tiler/maths/clip"
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
	case geom.Point3:
		// 处理3D点
		px := int64((g.X() - tile.MinX()) / tile.XSpan() * pixelExtent)
		py := int64((tile.MaxY() - g.Y()) / tile.YSpan() * pixelExtent)
		z := g.Z()
		return gen.NewPoint([]float64{float64(px), float64(py), z})

	case geom.Point:
		return preparept(g, tile, pixelExtent)

	case geom.MultiPoint3:
		// 处理3D多点
		pts := g.Points()
		if len(pts) == 0 {
			return nil
		}
		mp := make([][]float64, len(pts))
		for i, pt := range pts {
			px := int64((pt.X() - tile.MinX()) / tile.XSpan() * pixelExtent)
			py := int64((tile.MaxY() - pt.Y()) / tile.YSpan() * pixelExtent)
			z := pt.Z()
			mp[i] = []float64{float64(px), float64(py), z}
		}
		return gen.NewMultiPoint(mp)

	case geom.MultiPoint:
		pts := g.Points()
		if len(pts) == 0 {
			return nil
		}
		mp := make([][]float64, 0, len(pts))
		for _, pt := range g.Points() {
			preparedPt := preparept(pt, tile, pixelExtent)
			if preparedPt != nil {
				mp = append(mp, preparedPt.Data())
			}
		}
		if len(mp) == 0 {
			return nil
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

	case geom.Collection:
		// 处理几何集合
		var geoms []geom.Geometry
		for _, geo := range g.Geometries() {
			ng := PrepareGeo(geo, tile, pixelExtent)
			if ng != nil {
				geoms = append(geoms, ng)
			}
		}
		if len(geoms) == 0 {
			return nil
		}
		return gen.NewGeometryCollection(geoms...)
	}

	return nil
}

func preparept(g geom.Point, tile *gen.Extent, pixelExtent float64) geom.Point {
	// 参数验证
	if g == nil || tile == nil || pixelExtent <= 0 {
		fmt.Printf("参数无效: g=%v, tile=%v, pixelExtent=%v\n", g, tile, pixelExtent)
		return nil
	}

	// 检查点是否在瓦片范围内
	minX, minY := tile.MinX(), tile.MinY()
	maxX, maxY := tile.MaxX(), tile.MaxY()
	gx, gy := g.X(), g.Y()
	if gx < minX || gx > maxX || gy < minY || gy > maxY {
		fmt.Printf("点在瓦片范围外: g=(%v,%v), 瓦片范围=(%v,%v,%v,%v)\n", gx, gy, minX, minY, maxX, maxY)
		return nil
	}

	// 计算跨度，避免除以零
	xSpan := tile.XSpan()
	ySpan := tile.YSpan()
	if xSpan <= 0 || ySpan <= 0 {
		fmt.Printf("瓦片跨度无效: xSpan=%v, ySpan=%v\n", xSpan, ySpan)
		return nil
	}

	// 计算像素坐标
	px := (gx - minX) / xSpan * pixelExtent
	py := (maxY - gy) / ySpan * pixelExtent
	fmt.Printf("坐标转换: g=(%v,%v), 瓦片范围=(%v,%v,%v,%v), 跨度=(%v,%v), 像素范围=%v, 转换后=(%v,%v)\n",
		gx, gy, minX, minY, maxX, maxY, xSpan, ySpan, pixelExtent, px, py)
	var z float64 = 0
	if len(g.Data()) > 2 {
		z = g.Data()[2]
	}
	return gen.NewPoint([]float64{px, py, z})
}

// preparelinestr 将线字符串转换为瓦片像素坐标并裁剪到瓦片范围内
// g: 输入的线字符串
// tile: 瓦片范围
// pixelExtent: 瓦片像素尺寸
// 返回转换后的线字符串，如果线完全在瓦片外则返回nil
func preparelinestr(g geom.LineString, tile *gen.Extent, pixelExtent float64) geom.LineString {
	// 参数验证
	if g == nil || tile == nil || pixelExtent <= 0 {
		fmt.Printf("preparelinestr: 参数无效: g=%v, tile=%v, pixelExtent=%v\n", g, tile, pixelExtent)
		return nil
	}

	// 检查线是否有足够的点
	if len(g.Data()) < 2 {
		fmt.Printf("preparelinestr: 线点数不足: %d\n", len(g.Data()))
		return nil
	}

	// 裁剪线
	clippedLines, err := clip.LineString(g, tile)
	if err != nil {
		fmt.Printf("preparelinestr: 裁剪线错误: %v\n", err)
		return nil
	}
	if len(clippedLines) == 0 {
		fmt.Printf("preparelinestr: 没有裁剪结果\n")
		return nil
	}

	// 打印裁剪后的线数量
	fmt.Printf("preparelinestr: 裁剪后得到 %d 条线\n", len(clippedLines))

	// 转换裁剪后的线到像素坐标
	clippedLine := clippedLines[0]
	fmt.Printf("preparelinestr: 第一条裁剪线有 %d 个点\n", len(clippedLine.Data()))

	points := make([][]float64, len(clippedLine.Data()))

	for i, pt := range clippedLine.Data() {
		fmt.Printf("preparelinestr: 处理裁剪线的点 %d: %v\n", i, pt)
		ptGeom := gen.NewPoint([]float64{pt[0], pt[1]})
		preparedPt := preparept(ptGeom, tile, pixelExtent)
		if preparedPt == nil {
			fmt.Printf("preparelinestr: 点 %d 转换后为nil\n", i)
			continue
		}
		points[i] = preparedPt.Data()
		fmt.Printf("preparelinestr: 点 %d 转换后: %v\n", i, points[i])
	}

	if len(points) < 2 {
		fmt.Printf("preparelinestr: 转换后点数不足: %d\n", len(points))
		return nil
	}

	result := gen.NewLineString(points)
	fmt.Printf("preparelinestr: 最终线: %v\n", result)
	return result
}

// preparePolygon 将多边形转换为瓦片像素坐标并确保符合MVT规范
func preparePolygon(g geom.Polygon, tile *gen.Extent, pixelExtent float64) geom.Polygon {
	// 裁剪多边形
	clippedPolys, err := clip.Polygon(g, tile)
	if err != nil || len(clippedPolys) == 0 {
		return nil
	}

	clippedPoly := clippedPolys[0]
	lines := gen.NewMultiLineString(clippedPoly.Data())
	p := make([][][]float64, 0, len(lines.Data()))

	if len(lines.Data()) == 0 {
		return gen.NewPolygon(p)
	}

	for _, line := range lines.Lines() {
		// 转换线到像素坐标
		ln := preparelinestr(line, tile, pixelExtent)
		if ln == nil || len(ln.Data()) < 2 {
			continue
		}

		// 检查首尾点是否相同，确保环是闭合的（MVT规范要求）
		coords := ln.Data()
		if len(coords) < 2 {
			continue
		}

		first := coords[0]
		last := coords[len(coords)-1]

		// 确保至少有2个不同的点
		if isPointsEqual(first, last) && len(coords) == 2 {
			continue
		}

		// 如果首尾点不同，则闭合多边形环
		if !isPointsEqual(first, last) {
			// 添加第一个点到末尾以闭合环
			coords = append(coords, first)
			ln = gen.NewLineString(coords)
		}

		p = append(p, ln.Data())
	}

	return gen.NewPolygon(p)
}

// isPointsEqual 检查两个点是否相等
func isPointsEqual(a, b []float64) bool {
	if len(a) < 2 || len(b) < 2 {
		return false
	}
	return a[0] == b[0] && a[1] == b[1]
}
