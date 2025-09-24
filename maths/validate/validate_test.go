package validate

import (
	"context"
	"testing"

	"github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"
)

// TestCleanLinestring 测试线段清理功能
func TestCleanLinestring(t *testing.T) {
	testCases := []struct {
		name     string
		input    []float64
		expected []float64
	}{
		{
			name:     "简单线段",
			input:    []float64{0, 0, 1, 1, 2, 2},
			expected: []float64{0, 0, 1, 1, 2, 2},
		},
		{
			name:     "有重复点的线段",
			input:    []float64{0, 0, 1, 1, 1, 1, 2, 2},
			expected: []float64{0, 0, 1, 1, 2, 2},
		},
		{
			name:     "空线段",
			input:    []float64{},
			expected: []float64{},
		},
		{
			name:     "单点",
			input:    []float64{1, 1},
			expected: []float64{1, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CleanLinestring(tc.input)
			if err != nil {
				t.Errorf("CleanLinestring返回错误: %v", err)
				return
			}

			if len(result) != len(tc.expected) {
				t.Errorf("结果长度不匹配，期望 %d，得到 %d", len(tc.expected), len(result))
				return
			}

			for i, v := range tc.expected {
				if result[i] != v {
					t.Errorf("索引 %d 处值不匹配，期望 %f，得到 %f", i, v, result[i])
				}
			}
		})
	}
}

// TestLineAsPointPairs 测试线段转点对功能
func TestLineAsPointPairs(t *testing.T) {
	// 创建测试线段
	coords := [][]float64{
		{0, 0},
		{1, 1},
		{2, 2},
	}
	lineString := gen.NewLineString(coords)

	result := LineAsPointPairs(lineString)

	expected := []float64{0, 0, 1, 1, 2, 2}
	if len(result) != len(expected) {
		t.Errorf("结果长度不匹配，期望 %d，得到 %d", len(expected), len(result))
		return
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("索引 %d 处值不匹配，期望 %f，得到 %f", i, v, result[i])
		}
	}
}

// TestLineStringToSegments 测试线段转分段功能
func TestLineStringToSegments(t *testing.T) {
	// 创建测试线段
	coords := [][]float64{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
	}
	lineString := gen.NewLineString(coords)

	segments, err := LineStringToSegments(lineString)
	if err != nil {
		t.Errorf("LineStringToSegments返回错误: %v", err)
		return
	}

	// maths.NewSegments会创建闭环，所以4个点会产生4个线段
	expectedSegmentCount := 4
	if len(segments) != expectedSegmentCount {
		t.Errorf("分段数量不正确，期望 %d，得到 %d", expectedSegmentCount, len(segments))
	}

	// 验证第一个分段
	if len(segments) > 0 {
		firstSegment := segments[0]
		if len(firstSegment) != 2 {
			t.Error("线段应该包含两个点")
		}
	}
}

// TestCleanGeometry 测试几何清理功能
func TestCleanGeometry(t *testing.T) {
	ctx := context.Background()
	extent := gen.NewExtent([]float64{-10, -10, 10, 10})

	testCases := []struct {
		name     string
		geometry geom.Geometry
		hasError bool
	}{
		{
			name:     "nil几何",
			geometry: nil,
			hasError: false,
		},
		{
			name: "简单线段",
			geometry: gen.NewLineString([][]float64{
				{0, 0},
				{1, 1},
				{2, 2},
			}),
			hasError: false,
		},
		{
			name: "多线段",
			geometry: gen.NewMultiLineString([][][]float64{
				{
					{0, 0},
					{1, 1},
				},
				{
					{2, 2},
					{3, 3},
				},
			}),
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CleanGeometry(ctx, tc.geometry, extent)

			if tc.hasError && err == nil {
				t.Error("期望有错误，但没有返回错误")
				return
			}

			if !tc.hasError && err != nil {
				t.Errorf("不期望有错误，但返回了错误: %v", err)
				return
			}

			if tc.geometry == nil {
				if result != nil {
					t.Error("nil几何应该返回nil")
				}
			} else {
				if result == nil && tc.geometry != nil {
					t.Error("非nil几何不应该返回nil结果")
				}
			}
		})
	}
}

// TestScalePolygon 测试多边形缩放功能
func TestScalePolygon(t *testing.T) {
	// 创建测试多边形
	coords := [][][]float64{
		{
			{0, 0},
			{2, 0},
			{2, 2},
			{0, 2},
			{0, 0},
		},
	}
	polygon := gen.NewPolygon(coords)

	// 缩放因子为2
	scaleFactor := 2.0
	result := scalePolygon(polygon, scaleFactor)

	if len(result) != 1 {
		t.Errorf("结果应该包含1个环，但得到 %d 个", len(result))
		return
	}

	ring := result[0]
	if len(ring) != 5 {
		t.Errorf("环应该包含5个点，但得到 %d 个", len(ring))
		return
	}

	// 验证第一个点是否正确缩放
	if ring[0].X() != 0 || ring[0].Y() != 0 {
		t.Errorf("第一个点缩放错误，期望 (0,0)，得到 (%v,%v)", ring[0].X(), ring[0].Y())
	}

	// 验证第二个点是否正确缩放
	if ring[1].X() != 4 || ring[1].Y() != 0 { // 2 * 2 = 4
		t.Errorf("第二个点缩放错误，期望 (4,0)，得到 (%v,%v)", ring[1].X(), ring[1].Y())
	}
}

// BenchmarkCleanLinestring 性能基准测试
func BenchmarkCleanLinestring(b *testing.B) {
	// 创建包含重复点的大线段
	var coords []float64
	for i := 0; i < 1000; i++ {
		coords = append(coords, float64(i), float64(i))
		if i%10 == 0 { // 每10个点添加一个重复点
			coords = append(coords, float64(i), float64(i))
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CleanLinestring(coords)
	}
}
