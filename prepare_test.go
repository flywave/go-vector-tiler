package tile

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"
)

// TestPrepareGeo 测试PrepareGeo函数对不同几何类型的处理
func TestPrepareGeo(t *testing.T) {
	// 创建测试用的瓦片范围
	tile := &gen.Extent{
		116.0,
		39.0,
		117.0,
		40.0,
	}
	// 验证瓦片范围是否有效
	if tile.XSpan() <= 0 || tile.YSpan() <= 0 {
		panic("无效的瓦片范围: 宽度或高度必须大于0")
	}
	pixelExtent := float64(4096)

	// 测试点
	pointTests := []struct {
		name     string
		input    geom.Point
		expected geom.Point
	}{{
		name:     "普通点",
		input:    gen.NewPoint([]float64{116.5, 39.5, 10.0}),
		expected: gen.NewPoint([]float64{2048.0, 2048.0, 10.0}),
	}, {
		name:     "边界点",
		input:    gen.NewPoint([]float64{116.0, 39.0, 5.0}),
		expected: gen.NewPoint([]float64{0.0, 4096.0, 5.0}),
	}, {
		name:     "瓦片外的点",
		input:    gen.NewPoint([]float64{115.0, 38.0, 0.0}),
		expected: nil,
	}}

	// 测试多点
	multiPointTests := []struct {
		name     string
		input    geom.MultiPoint
		expected geom.MultiPoint
	}{{
		name:     "普通多点",
		input:    gen.NewMultiPoint([][]float64{{116.25, 39.25}, {116.75, 39.75}}),
		expected: gen.NewMultiPoint([][]float64{{1024.0, 3072.0, 0.0}, {3072.0, 1024.0, 0.0}}),
	}, {
		name:     "包含瓦片外点的多点",
		input:    gen.NewMultiPoint([][]float64{{116.25, 39.25}, {115.0, 38.0}}),
		expected: gen.NewMultiPoint([][]float64{{1024.0, 3072.0, 0.0}}),
	}, {
		name:     "空多点",
		input:    gen.NewMultiPoint([][]float64{}),
		expected: nil,
	}}

	// 测试线字符串
	lineStringTests := []struct {
		name     string
		input    geom.LineString
		expected geom.LineString
	}{{
		name:     "普通线",
		input:    gen.NewLineString([][]float64{{116.0, 39.0}, {117.0, 40.0}}),
		expected: gen.NewLineString([][]float64{{0.0, 4096.0, 0.0}, {4096.0, 0.0, 0.0}}),
	}, {
		name:     "部分在瓦片外的线",
		input:    gen.NewLineString([][]float64{{115.0, 39.5}, {116.5, 39.5}}),
		expected: gen.NewLineString([][]float64{{0.0, 2048.0, 0.0}, {2048.0, 2048.0, 0.0}}),
	}}

	// 测试多边形
	polygonTests := []struct {
		name     string
		input    geom.Polygon
		expected geom.Polygon
	}{{
		name: "普通多边形",
		input: gen.NewPolygon([][][]float64{{
			{116.0, 39.0},
			{117.0, 39.0},
			{117.0, 40.0},
			{116.0, 40.0},
			{116.0, 39.0},
		}}),
		expected: gen.NewPolygon([][][]float64{{
			{0.0, 4096.0, 0.0},
			{4096.0, 4096.0, 0.0},
			{4096.0, 0.0, 0.0},
			{0.0, 0.0, 0.0},
			{0.0, 4096.0, 0.0},
		}}),
	}}

	// 测试几何集合
	collectionTests := []struct {
		name     string
		input    geom.Collection
		expected geom.Collection
	}{{
		name: "包含点和线的集合",
		input: gen.NewGeometryCollection(
			gen.NewPoint([]float64{116.5, 39.5}),
			gen.NewLineString([][]float64{{116.0, 39.0}, {117.0, 40.0}}),
		),
		expected: gen.NewGeometryCollection(
			gen.NewPoint([]float64{2048.0, 2048.0, 0.0}),
			gen.NewLineString([][]float64{{0.0, 4096.0, 0.0}, {4096.0, 0.0, 0.0}}),
		),
	}}

	// 运行所有测试
	for _, tc := range pointTests {
		t.Run("Point/"+tc.name, func(t *testing.T) {
			result := PrepareGeo(tc.input, tile, pixelExtent)
			if !isGeometryEqual(result, tc.expected) {
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, result)
			}
		})
	}

	for _, tc := range multiPointTests {
		t.Run("MultiPoint/"+tc.name, func(t *testing.T) {
			result := PrepareGeo(tc.input, tile, pixelExtent)
			if !isGeometryEqual(result, tc.expected) {
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, result)
			}
		})
	}

	for _, tc := range lineStringTests {
		t.Run("LineString/"+tc.name, func(t *testing.T) {
			// 打印测试用例信息
			fmt.Printf("测试用例: %s\n", tc.name)
			fmt.Printf("输入线: %v\n", tc.input)
			fmt.Printf("预期线: %v\n", tc.expected)

			result := PrepareGeo(tc.input, tile, pixelExtent)
			fmt.Printf("实际结果: %v\n", result)

			if !isGeometryEqual(result, tc.expected) {
				// 添加调试信息
				if lsResult, ok := result.(geom.LineString); ok {
					lsExpected := tc.expected
					resultData := lsResult.Data()
					expectedData := lsExpected.Data()
					fmt.Printf("实际点数: %d, 预期点数: %d\n", len(resultData), len(expectedData))
					if len(resultData) != len(expectedData) {
						t.Errorf("%s: 点数量不同: 预期 %d, 实际 %d", tc.name, len(expectedData), len(resultData))
					} else {
						for i := range resultData {
							fmt.Printf("点 %d: 实际=%v, 预期=%v\n", i, resultData[i], expectedData[i])
							if len(resultData[i]) != len(expectedData[i]) {
								t.Errorf("%s: 点 %d 维度不同: 预期 %d, 实际 %d", tc.name, i, len(expectedData[i]), len(resultData[i]))
							} else {
								for j := range resultData[i] {
									if math.Abs(resultData[i][j]-expectedData[i][j]) > 1e-6 {
										t.Errorf("%s: 点 %d 坐标 %d 不同: 预期 %.6f, 实际 %.6f", tc.name, i, j, expectedData[i][j], resultData[i][j])
									}
								}
							}
						}
					}
				}
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, result)
			}
		})
	}

	for _, tc := range polygonTests {
		t.Run("Polygon/"+tc.name, func(t *testing.T) {
			result := PrepareGeo(tc.input, tile, pixelExtent)
			if !isGeometryEqual(result, tc.expected) {
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, result)
			}
		})
	}

	for _, tc := range collectionTests {
		t.Run("Collection/"+tc.name, func(t *testing.T) {
			result := PrepareGeo(tc.input, tile, pixelExtent)
			if !isGeometryEqual(result, tc.expected) {
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, result)
			}
		})
	}
}

// isGeometryEqual 比较两个几何对象是否相等
func isGeometryEqual(g1, g2 geom.Geometry) bool {
	// 处理nil情况
	if g1 == nil || g2 == nil {
		return g1 == g2
	}

	// 检查类型是否相同
	if g1.GetType() != g2.GetType() {
		return false
	}

	// 根据类型进行具体比较
	switch geo1 := g1.(type) {
	case geom.Point:
		geo2, ok := g2.(geom.Point)
		if !ok {
			return false
		}
		return isPointEqual(geo1, geo2)
	case geom.MultiPoint:
		geo2, ok := g2.(geom.MultiPoint)
		if !ok {
			return false
		}
		return isMultiPointEqual(geo1, geo2)
	case geom.LineString:
		geo2, ok := g2.(geom.LineString)
		if !ok {
			return false
		}
		return isLineStringEqual(geo1, geo2)
	case geom.MultiLine:
		geo2, ok := g2.(geom.MultiLine)
		if !ok {
			return false
		}
		return isMultiLineEqual(geo1, geo2)
	case geom.Polygon:
		geo2, ok := g2.(geom.Polygon)
		if !ok {
			return false
		}
		return isPolygonEqual(geo1, geo2)
	case geom.MultiPolygon:
		geo2, ok := g2.(geom.MultiPolygon)
		if !ok {
			return false
		}
		return isMultiPolygonEqual(geo1, geo2)
	case geom.Collection:
		geo2, ok := g2.(geom.Collection)
		if !ok {
			return false
		}
		return isCollectionEqual(geo1, geo2)
	default:
		// 处理未知类型
		return false
	}
}

// isPointEqual 比较两个点是否相等
func isPointEqual(p1, p2 geom.Point) bool {
	if p1 == nil || p2 == nil {
		return p1 == p2
	}
	return reflect.DeepEqual(p1.Data(), p2.Data())
}

// isMultiPointEqual 比较两个多点是否相等
func isMultiPointEqual(mp1, mp2 geom.MultiPoint) bool {
	pts1, pts2 := mp1.Points(), mp2.Points()
	if len(pts1) != len(pts2) {
		return false
	}
	for i, pt := range pts1 {
		if !isPointEqual(pt, pts2[i]) {
			return false
		}
	}
	return true
}

// isLineStringEqual 比较两个线字符串是否相等
func isLineStringEqual(l1, l2 geom.LineString) bool {
	if l1 == nil || l2 == nil {
		return l1 == l2
	}
	if len(l1.Data()) != len(l2.Data()) {
		return false
	}
	for i := range l1.Data() {
		if len(l1.Data()[i]) != len(l2.Data()[i]) {
			return false
		}
		// 比较坐标值，允许一定的浮点精度误差
		for j := range l1.Data()[i] {
			if math.Abs(l1.Data()[i][j]-l2.Data()[i][j]) > 1e-6 {
				return false
			}
		}
	}
	return true
}

// isMultiLineEqual 比较两个多线是否相等
func isMultiLineEqual(ml1, ml2 geom.MultiLine) bool {
	lns1, lns2 := ml1.Lines(), ml2.Lines()
	if len(lns1) != len(lns2) {
		return false
	}
	for i, ln := range lns1 {
		if !isLineStringEqual(ln, lns2[i]) {
			return false
		}
	}
	return true
}

// isPolygonEqual 比较两个多边形是否相等
func isPolygonEqual(p1, p2 geom.Polygon) bool {
	lns1, lns2 := p1.Sublines(), p2.Sublines()
	if len(lns1) != len(lns2) {
		return false
	}
	for i, ln := range lns1 {
		if !isLineStringEqual(ln, lns2[i]) {
			return false
		}
	}
	return true
}

// isMultiPolygonEqual 比较两个多多边形是否相等
func isMultiPolygonEqual(mp1, mp2 geom.MultiPolygon) bool {
	pgs1, pgs2 := mp1.Polygons(), mp2.Polygons()
	if len(pgs1) != len(pgs2) {
		return false
	}
	for i, pg := range pgs1 {
		if !isPolygonEqual(pg, pgs2[i]) {
			return false
		}
	}
	return true
}

// isCollectionEqual 比较两个几何集合是否相等
func isCollectionEqual(c1, c2 geom.Collection) bool {
	geos1, geos2 := c1.Geometries(), c2.Geometries()
	if len(geos1) != len(geos2) {
		return false
	}
	for i, geo := range geos1 {
		if !isGeometryEqual(geo, geos2[i]) {
			return false
		}
	}
	return true
}
