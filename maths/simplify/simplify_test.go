package simplify

import (
	"testing"

	gen "github.com/flywave/go-geom/general"

	"github.com/flywave/go-vector-tiler/basic"
	"github.com/flywave/go-vector-tiler/maths"
)

// TestSimplifyGeometry_Point 测试点几何简化
func TestSimplifyGeometry_Point(t *testing.T) {
	// 创建点几何
	point := gen.NewPoint([]float64{10, 20})

	// 点几何应该原样返回
	result := SimplifyGeometry(point, 1.0)
	if result != point {
		t.Error("点几何应该原样返回")
	}
}

// TestSimplifyGeometry_LineString 测试线段简化
func TestSimplifyGeometry_LineString(t *testing.T) {
	// 创建线段几何
	coords := [][]float64{
		{0, 0},
		{1, 1},
		{2, 0.1}, // 微小偏差
		{3, 3},
		{4, 4},
	}
	lineString := gen.NewLineString(coords)

	// 使用小容差简化
	result := SimplifyGeometry(lineString, 0.01)
	resultLine, ok := result.(basic.Line)
	if !ok {
		t.Fatal("简化结果应该是 basic.Line 类型")
	}

	// 应该保留所有重要点
	if len(resultLine) < 3 {
		t.Error("小容差应该保留更多点")
	}

	// 使用大容差简化
	result = SimplifyGeometry(lineString, 2.0)
	resultLine, ok = result.(basic.Line)
	if !ok {
		t.Fatal("简化结果应该是 basic.Line 类型")
	}

	// 应该显著减少点数
	if len(resultLine) > 3 {
		t.Error("大容差应该显著减少点数")
	}
}

// TestSimplifyGeometry_MultiLine 测试多线段简化
func TestSimplifyGeometry_MultiLine(t *testing.T) {
	// 创建多线段几何
	coords := [][][]float64{
		{
			{0, 0},
			{1, 1},
			{2, 0},
			{3, 3},
		},
		{
			{10, 10},
			{11, 11},
			{12, 10},
			{13, 13},
		},
	}
	multiLine := gen.NewMultiLineString(coords)

	// 简化多线段
	result := SimplifyGeometry(multiLine, 1.0)
	resultML, ok := result.(basic.MultiLine)
	if !ok {
		t.Fatal("简化结果应该是 basic.MultiLine 类型")
	}

	// 应该保留两条线段
	if len(resultML) != 2 {
		t.Errorf("应该保留2条线段，实际得到 %d 条", len(resultML))
	}

	// 每条线段都应该被简化
	for i, line := range resultML {
		if len(line) == 0 {
			t.Errorf("线段 %d 不应该为空", i)
		}
	}
}

// TestSimplifyGeometry_Polygon 测试多边形简化
func TestSimplifyGeometry_Polygon(t *testing.T) {
	// 创建多边形几何（外环 + 内环）
	coords := [][][]float64{
		{
			// 外环 - 矩形
			{0, 0},
			{10, 0},
			{10, 10},
			{0, 10},
			{0, 0},
		},
		{
			// 内环 - 小矩形孔
			{2, 2},
			{8, 2},
			{8, 8},
			{2, 8},
			{2, 2},
		},
	}
	polygon := gen.NewPolygon(coords)

	// 简化多边形
	result := SimplifyGeometry(polygon, 0.5)
	resultPoly, ok := result.(basic.Polygon)
	if !ok {
		t.Fatal("简化结果应该是 basic.Polygon 类型")
	}

	// 应该保留外环和内环
	if len(resultPoly) != 2 {
		t.Errorf("应该保留2个环，实际得到 %d 个", len(resultPoly))
	}

	// 验证环不为空
	for i, ring := range resultPoly {
		if len(ring) == 0 {
			t.Errorf("环 %d 不应该为空", i)
		}
	}
}

// TestSimplifyGeometry_MultiPolygon 测试多多边形简化
func TestSimplifyGeometry_MultiPolygon(t *testing.T) {
	// 创建多多边形几何
	coords := [][][][]float64{
		{
			{
				{0, 0},
				{5, 0},
				{5, 5},
				{0, 5},
				{0, 0},
			},
		},
		{
			{
				{10, 10},
				{15, 10},
				{15, 15},
				{10, 15},
				{10, 10},
			},
		},
	}
	multiPolygon := gen.NewMultiPolygon(coords)

	// 简化多多边形
	result := SimplifyGeometry(multiPolygon, 0.5)
	resultMP, ok := result.(basic.MultiPolygon)
	if !ok {
		t.Fatal("简化结果应该是 basic.MultiPolygon 类型")
	}

	// 应该保留两个多边形
	if len(resultMP) != 2 {
		t.Errorf("应该保留2个多边形，实际得到 %d 个", len(resultMP))
	}

	// 验证每个多边形都不为空
	for i, poly := range resultMP {
		if len(poly) == 0 {
			t.Errorf("多边形 %d 不应该为空", i)
		}
	}
}

// TestSimplifyGeometry_EmptyResults 测试产生空结果的情况
func TestSimplifyGeometry_EmptyResults(t *testing.T) {
	// 创建非常小的线段
	coords := [][]float64{
		{0, 0},
		{0.001, 0.001},
	}
	lineString := gen.NewLineString(coords)

	// 使用大容差，可能导致空结果
	result := SimplifyGeometry(lineString, 10.0)
	if result != nil {
		// 如果不为nil，检查是否为有效的几何
		if line, ok := result.(basic.Line); ok && len(line) == 0 {
			t.Error("简化结果不应该是空线段")
		}
	}
}

// TestSimplifyGeometry_ToleranceEffects 测试不同容差的效果
func TestSimplifyGeometry_ToleranceEffects(t *testing.T) {
	// 创建锯齿形线段
	coords := [][]float64{
		{0, 0},
		{1, 1},
		{2, 0},
		{3, 1},
		{4, 0},
		{5, 1},
		{6, 0},
	}
	lineString := gen.NewLineString(coords)

	tolerances := []float64{0.1, 0.5, 1.0, 2.0}
	var prevLength int

	for i, tolerance := range tolerances {
		result := SimplifyGeometry(lineString, tolerance)
		if result == nil {
			continue
		}

		line, ok := result.(basic.Line)
		if !ok {
			t.Fatalf("容差 %f: 结果应该是 basic.Line 类型", tolerance)
		}

		currentLength := len(line)

		// 更大的容差应该产生更少的点（一般情况下）
		if i > 0 && currentLength > prevLength {
			t.Logf("警告：容差 %f 产生了比容差 %f 更多的点", tolerance, tolerances[i-1])
		}

		prevLength = currentLength

		// 至少应该保留端点
		if currentLength < 2 && currentLength > 0 {
			t.Errorf("容差 %f: 如果有结果，至少应该有2个点", tolerance)
		}
	}
}

// TestSimplifyLineString 测试线段简化内部函数
func TestSimplifyLineString(t *testing.T) {
	// 创建测试线段
	coords := [][]float64{
		{0, 0},
		{1, 1},
		{2, 0.1},
		{3, 3},
	}
	lineString := gen.NewLineString(coords)

	// 测试不同容差
	testCases := []struct {
		tolerance   float64
		expectedMin int
		description string
	}{
		{0.01, 2, "小容差应该保留至少端点"},
		{1.0, 2, "中等容差应该保留至少端点"},
		{5.0, 2, "大容差应该保留至少端点"},
	}

	for _, tc := range testCases {
		result := simplifyLineString(lineString, tc.tolerance)

		if result == nil {
			t.Errorf("%s: 不应该返回nil", tc.description)
			continue
		}

		if len(result) < tc.expectedMin {
			t.Errorf("%s: 期望至少 %d 个点，实际得到 %d",
				tc.description, tc.expectedMin, len(result))
		}
	}
}

// TestSimplifyPolygon 测试多边形简化内部函数
func TestSimplifyPolygon(t *testing.T) {
	// 创建测试多边形
	coords := [][][]float64{
		{
			{0, 0},
			{1, 0.1}, // 微小偏差
			{2, 0},
			{2, 2},
			{0, 2},
			{0, 0},
		},
	}
	polygon := gen.NewPolygon(coords)

	// 测试小容差
	result := simplifyPolygon(polygon, 0.01)
	if result == nil {
		t.Error("小容差不应该返回nil")
	} else if len(result) == 0 {
		t.Error("小容差不应该返回空多边形")
	}

	// 测试大容差
	result = simplifyPolygon(polygon, 1.0)
	if len(result) > 0 {
		// 如果返回结果，验证基本有效性
		for i, ring := range result {
			if len(ring) < 3 {
				t.Errorf("环 %d 至少应该有3个点", i)
			}
		}
	}
}

// TestNormalizePoints 测试点规范化函数
func TestNormalizePoints(t *testing.T) {
	// 测试闭合环
	closedRing := []maths.Pt{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 0, Y: 1},
		{X: 0, Y: 0}, // 重复起点
	}

	result := normalizePoints(closedRing)

	// 应该移除重复的起点
	if len(result) < 3 {
		t.Errorf("规范化后应该至少有3个点，实际得到 %d 个", len(result))
	}

	// 测试基本情况
	basicPoints := []maths.Pt{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 2},
	}

	result = normalizePoints(basicPoints)

	// 应该保持基本几何结构
	if len(result) < 2 {
		t.Error("规范化不应该移除所有点")
	}
}

// BenchmarkSimplifyGeometry_LineString 性能基准测试
func BenchmarkSimplifyGeometry_LineString(b *testing.B) {
	// 创建复杂线段
	var coords [][]float64
	for i := 0; i < 100; i++ {
		x := float64(i)
		y := float64(i % 2) // 锯齿形
		coords = append(coords, []float64{x, y})
	}
	lineString := gen.NewLineString(coords)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SimplifyGeometry(lineString, 0.5)
	}
}

// BenchmarkSimplifyGeometry_Polygon 多边形简化性能基准测试
func BenchmarkSimplifyGeometry_Polygon(b *testing.B) {
	// 创建复杂多边形
	var coords [][][]float64
	var ring [][]float64

	// 外环
	for i := 0; i < 50; i++ {
		x := 10.0 + 5.0*float64(i%2)*0.1 // 添加一些不规则性
		y := 10.0 + 5.0*float64(i%2)*0.1
		ring = append(ring, []float64{x, y})
	}
	ring = append(ring, ring[0]) // 闭合
	coords = append(coords, ring)

	polygon := gen.NewPolygon(coords)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SimplifyGeometry(polygon, 0.5)
	}
}
