package simplify

import (
	"math"
	"testing"

	"github.com/flywave/go-vector-tiler/maths"
)

// TestDouglasPeucker_BasicLine 测试基本线段简化
func TestDouglasPeucker_BasicLine(t *testing.T) {
	// 创建一条简单的线段
	points := []maths.Pt{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 0},
	}

	// 使用较大的容差，应该简化为两个端点
	result := DouglasPeucker(points, 2.0)
	if len(result) != 2 {
		t.Errorf("期望简化为2个点，实际得到 %d 个点", len(result))
	}

	// 验证端点保持不变
	if result[0] != points[0] || result[len(result)-1] != points[len(points)-1] {
		t.Error("端点应该保持不变")
	}
}

// TestDouglasPeucker_StraightLine 测试直线简化
func TestDouglasPeucker_StraightLine(t *testing.T) {
	// 创建一条直线（共线点）
	points := []maths.Pt{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 3, Y: 3},
		{X: 4, Y: 4},
	}

	// 对于直线，应该简化为两个端点
	result := DouglasPeucker(points, 0.1)
	if len(result) != 2 {
		t.Errorf("直线应该简化为2个点，实际得到 %d 个点", len(result))
	}

	// 验证端点
	expected := []maths.Pt{
		{X: 0, Y: 0},
		{X: 4, Y: 4},
	}

	for i, pt := range expected {
		if result[i] != pt {
			t.Errorf("端点 %d 不匹配：期望 %v，得到 %v", i, pt, result[i])
		}
	}
}

// TestDouglasPeucker_EdgeCases 测试边界情况
func TestDouglasPeucker_EdgeCases(t *testing.T) {
	// 测试空切片
	result := DouglasPeucker([]maths.Pt{}, 1.0)
	if len(result) != 0 {
		t.Error("空切片应该返回空切片")
	}

	// 测试单个点
	singlePoint := []maths.Pt{{X: 1, Y: 1}}
	result = DouglasPeucker(singlePoint, 1.0)
	if len(result) != 1 || result[0] != singlePoint[0] {
		t.Error("单个点应该原样返回")
	}

	// 测试两个点
	twoPoints := []maths.Pt{{X: 0, Y: 0}, {X: 1, Y: 1}}
	result = DouglasPeucker(twoPoints, 1.0)
	if len(result) != 2 {
		t.Error("两个点应该原样返回")
	}

	// 测试零容差
	points := []maths.Pt{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 0},
	}
	result = DouglasPeucker(points, 0)
	if len(result) != len(points) {
		t.Error("零容差应该保留所有点")
	}
}

// TestDouglasPeucker_AlgorithmCorrectness 测试算法正确性
func TestDouglasPeucker_AlgorithmCorrectness(t *testing.T) {
	// 创建一个具有明确偏差的线段
	points := []maths.Pt{
		{X: 0, Y: 0},
		{X: 2, Y: 2}, // 偏离较大
		{X: 4, Y: 0}, // 回到基线
	}

	// 小容差应该保留关键点
	result := DouglasPeucker(points, 0.1)
	if len(result) < 2 {
		t.Error("算法应该至少保留端点")
	}

	// 验证端点保持不变
	if result[0] != points[0] || result[len(result)-1] != points[len(points)-1] {
		t.Error("端点应该保持不变")
	}

	// 大容差可能简化
	result = DouglasPeucker(points, 5.0)
	if len(result) < 2 {
		t.Error("即使大容差也应该保留端点")
	}
}

// BenchmarkDouglasPeucker 性能基准测试
func BenchmarkDouglasPeucker(b *testing.B) {
	// 创建测试数据
	var points []maths.Pt
	for i := 0; i < 100; i++ {
		x := float64(i)
		y := math.Sin(float64(i) * 0.1) // 正弦波
		points = append(points, maths.Pt{X: x, Y: y})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DouglasPeucker(points, 0.1)
	}
}
