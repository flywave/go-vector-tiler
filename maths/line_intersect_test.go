package maths

import (
	"testing"
)

// TestEvent_Point 测试事件点获取功能
func TestEvent_Point(t *testing.T) {
	pt := Pt{X: 1.0, Y: 2.0}
	e := event{
		edge:     0,
		edgeType: LEFT,
		ev:       &pt,
	}

	result := e.Point()
	if result != &pt {
		t.Errorf("Point() = %v, want %v", result, &pt)
	}
}

// TestEvent_Edge 测试事件边获取功能
func TestEvent_Edge(t *testing.T) {
	pt := Pt{X: 1.0, Y: 2.0}
	e := event{
		edge:     5,
		edgeType: RIGHT,
		ev:       &pt,
	}

	result := e.Edge()
	if result != 5 {
		t.Errorf("Edge() = %v, want %v", result, 5)
	}
}

// TestXYOrderedEventPtr 测试事件排序功能
func TestXYOrderedEventPtr(t *testing.T) {
	events := []event{
		{ev: &Pt{X: 3.0, Y: 4.0}},
		{ev: &Pt{X: 1.0, Y: 2.0}},
		{ev: &Pt{X: 2.0, Y: 3.0}},
	}

	orderedEvents := XYOrderedEventPtr(events)

	// 测试长度
	if orderedEvents.Len() != 3 {
		t.Errorf("Len() = %v, want %v", orderedEvents.Len(), 3)
	}

	// 测试Less函数
	if !orderedEvents.Less(1, 0) { // (1,2) < (3,4) 应该为真
		t.Error("Less(1, 0) should be true")
	}

	// 测试交换功能
	orderedEvents.Swap(0, 1)
	if orderedEvents[0].ev.X != 1.0 || orderedEvents[1].ev.X != 3.0 {
		t.Error("Swap did not work correctly")
	}
}

// TestNewEventQueue 测试事件队列创建
func TestNewEventQueue(t *testing.T) {
	testCases := []struct {
		name     string
		segments []Line
		expected int // 期望的事件数量
	}{
		{
			name: "单线段",
			segments: []Line{
				{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 1}},
			},
			expected: 2,
		},
		{
			name: "多线段",
			segments: []Line{
				{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 1}},
				{Pt{X: 2, Y: 2}, Pt{X: 3, Y: 3}},
			},
			expected: 4,
		},
		{
			name:     "空线段",
			segments: []Line{},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			eq := NewEventQueue(tc.segments)
			if len(eq) != tc.expected {
				t.Errorf("NewEventQueue() 事件数量 = %v, want %v", len(eq), tc.expected)
			}

			// 验证事件队列是否按XY顺序排序
			for i := 1; i < len(eq); i++ {
				if XYOrder(*eq[i-1].ev, *eq[i].ev) == 1 {
					t.Errorf("事件队列未正确排序: 位置 %d", i)
				}
			}
		})
	}
}

// TestDoesIntersect 测试线段相交检测
func TestDoesIntersect(t *testing.T) {
	testCases := []struct {
		name     string
		s1       Line
		s2       Line
		expected bool
	}{
		{
			name:     "相交线段 - X形状",
			s1:       Line{Pt{X: 0, Y: 0}, Pt{X: 2, Y: 2}},
			s2:       Line{Pt{X: 0, Y: 2}, Pt{X: 2, Y: 0}},
			expected: true,
		},
		{
			name:     "平行线段 - 不相交",
			s1:       Line{Pt{X: 0, Y: 0}, Pt{X: 2, Y: 0}},
			s2:       Line{Pt{X: 0, Y: 1}, Pt{X: 2, Y: 1}},
			expected: false,
		},
		{
			name:     "相同线段",
			s1:       Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 1}},
			s2:       Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 1}},
			expected: true,
		},
		{
			name:     "端点相交",
			s1:       Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}},
			s2:       Line{Pt{X: 1, Y: 0}, Pt{X: 2, Y: 0}},
			expected: true,
		},
		{
			name:     "不相交的线段",
			s1:       Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}},
			s2:       Line{Pt{X: 2, Y: 2}, Pt{X: 3, Y: 3}},
			expected: false,
		},
		{
			name:     "垂直线段相交",
			s1:       Line{Pt{X: 1, Y: 0}, Pt{X: 1, Y: 2}},
			s2:       Line{Pt{X: 0, Y: 1}, Pt{X: 2, Y: 1}},
			expected: true,
		},
		{
			name:     "T形相交",
			s1:       Line{Pt{X: 0, Y: 1}, Pt{X: 2, Y: 1}},
			s2:       Line{Pt{X: 1, Y: 0}, Pt{X: 1, Y: 1}},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DoesIntersect(tc.s1, tc.s2)
			if result != tc.expected {
				t.Errorf("DoesIntersect(%v, %v) = %v, want %v", tc.s1, tc.s2, result, tc.expected)
			}
		})
	}
}

// TestLine_DoesIntersect 测试Line类型的DoesIntersect方法
func TestLine_DoesIntersect(t *testing.T) {
	s1 := Line{Pt{X: 0, Y: 0}, Pt{X: 2, Y: 2}}
	s2 := Line{Pt{X: 0, Y: 2}, Pt{X: 2, Y: 0}}

	result := s1.DoesIntersect(s2)
	if !result {
		t.Errorf("Line.DoesIntersect() = %v, want %v", result, true)
	}
}

// TestIntersectfn_PtFn 测试相交点计算函数
func TestIntersectfn_PtFn(t *testing.T) {
	l1 := Line{Pt{X: 0, Y: 0}, Pt{X: 2, Y: 2}}
	l2 := Line{Pt{X: 0, Y: 2}, Pt{X: 2, Y: 0}}
	ifn := intersectfn{l1, l2}

	pt := ifn.PtFn()
	expected := Pt{X: 1, Y: 1}

	if pt.X != expected.X || pt.Y != expected.Y {
		t.Errorf("PtFn() = %v, want %v", pt, expected)
	}
}

// TestFindinter_doesNotIntersect 测试内部相交检测函数
func TestFindinter_doesNotIntersect(t *testing.T) {
	testCases := []struct {
		name       string
		s1x0, s1y0 float64
		s1x1, s1y1 float64
		s2x0, s2y0 float64
		s2x1, s2y1 float64
		expected   bool
	}{
		{
			name: "相交线段应返回false",
			s1x0: 0, s1y0: 0, s1x1: 2, s1y1: 2,
			s2x0: 0, s2y0: 2, s2x1: 2, s2y1: 0,
			expected: false,
		},
		{
			name: "平行线段应返回true",
			s1x0: 0, s1y0: 0, s1x1: 2, s1y1: 0,
			s2x0: 0, s2y0: 1, s2x1: 2, s2y1: 1,
			expected: true,
		},
		{
			name: "不相交线段应返回true",
			s1x0: 0, s1y0: 0, s1x1: 1, s1y1: 0,
			s2x0: 2, s2y0: 2, s2x1: 3, s2y1: 3,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findinter_doesNotIntersect(tc.s1x0, tc.s1y0, tc.s1x1, tc.s1y1,
				tc.s2x0, tc.s2y0, tc.s2x1, tc.s2y1)
			if result != tc.expected {
				t.Errorf("findinter_doesNotIntersect() = %v, want %v", result, tc.expected)
			}
		})
	}
}

// TestFindIntersectsWithEventQueue 测试事件队列相交查找
func TestFindIntersectsWithEventQueue(t *testing.T) {
	segments := []Line{
		{Pt{X: 0, Y: 0}, Pt{X: 2, Y: 2}},
		{Pt{X: 0, Y: 2}, Pt{X: 2, Y: 0}},
	}
	eq := NewEventQueue(segments)

	intersections := make([]struct {
		src, dest int
		pt        Pt
	}, 0)

	fn := func(srcIdx, destIdx int, ptfn func() Pt) bool {
		intersections = append(intersections, struct {
			src, dest int
			pt        Pt
		}{srcIdx, destIdx, ptfn()})
		return true
	}

	FindIntersectsWithEventQueue(false, eq, segments, fn)

	if len(intersections) != 1 {
		t.Errorf("FindIntersectsWithEventQueue() 发现相交点数量 = %v, want %v", len(intersections), 1)
	}

	if len(intersections) > 0 {
		expected := Pt{X: 1, Y: 1}
		if intersections[0].pt.X != expected.X || intersections[0].pt.Y != expected.Y {
			t.Errorf("相交点 = %v, want %v", intersections[0].pt, expected)
		}
	}
}

// TestFindIntersectsWithoutIntersect 测试不需要相交点的相交查找
func TestFindIntersectsWithoutIntersect(t *testing.T) {
	segments := []Line{
		{Pt{X: 0, Y: 0}, Pt{X: 2, Y: 2}},
		{Pt{X: 0, Y: 2}, Pt{X: 2, Y: 0}},
		{Pt{X: 3, Y: 3}, Pt{X: 4, Y: 4}}, // 不与其他线段相交
	}

	intersectionCount := 0
	fn := func(srcIdx, destIdx int) bool {
		intersectionCount++
		return true
	}

	FindIntersectsWithoutIntersect(segments, fn)

	if intersectionCount != 1 {
		t.Errorf("FindIntersectsWithoutIntersect() 相交对数量 = %v, want %v", intersectionCount, 1)
	}
}

// TestFindIntersects 测试查找相交点函数
func TestFindIntersects(t *testing.T) {
	segments := []Line{
		{Pt{X: 0, Y: 0}, Pt{X: 2, Y: 2}},
		{Pt{X: 0, Y: 2}, Pt{X: 2, Y: 0}},
	}

	intersectionCount := 0
	fn := func(srcIdx, destIdx int, ptfn func() Pt) bool {
		intersectionCount++
		pt := ptfn()
		expected := Pt{X: 1, Y: 1}
		if pt.X != expected.X || pt.Y != expected.Y {
			t.Errorf("相交点 = %v, want %v", pt, expected)
		}
		return true
	}

	FindIntersects(segments, fn)

	if intersectionCount != 1 {
		t.Errorf("FindIntersects() 相交点数量 = %v, want %v", intersectionCount, 1)
	}
}

// TestFindPolygonIntersects 测试多边形相交查找
func TestFindPolygonIntersects(t *testing.T) {
	// 测试少于3个线段的情况
	segments := []Line{
		{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}},
		{Pt{X: 1, Y: 0}, Pt{X: 1, Y: 1}},
	}

	intersectionCount := 0
	fn := func(srcIdx, destIdx int, ptfn func() Pt) bool {
		intersectionCount++
		return true
	}

	FindPolygonIntersects(segments, fn)

	// 少于3个线段应该直接返回，不处理
	if intersectionCount != 0 {
		t.Errorf("FindPolygonIntersects() 对于少于3个线段应该不处理，但得到 %v 个相交点", intersectionCount)
	}

	// 测试有效的多边形（3个或更多线段）
	segments = []Line{
		{Pt{X: 0, Y: 0}, Pt{X: 2, Y: 0}},
		{Pt{X: 2, Y: 0}, Pt{X: 2, Y: 2}},
		{Pt{X: 2, Y: 2}, Pt{X: 0, Y: 2}},
		{Pt{X: 0, Y: 2}, Pt{X: 0, Y: 0}},
		{Pt{X: 0, Y: 1}, Pt{X: 2, Y: 1}}, // 与多边形边相交的线段
	}

	intersectionCount = 0
	FindPolygonIntersects(segments, fn)

	// 应该找到一些相交点
	if intersectionCount == 0 {
		t.Error("FindPolygonIntersects() 应该找到相交点")
	}
}

// TestLine_IntersectsLines 测试线段与多条线段的相交检测
func TestLine_IntersectsLines(t *testing.T) {
	l := Line{Pt{X: 0, Y: 1}, Pt{X: 3, Y: 1}}

	testCases := []struct {
		name     string
		lines    []Line
		expected []int // 期望相交的线段索引
	}{
		{
			name:     "空线段数组",
			lines:    []Line{},
			expected: []int{},
		},
		{
			name: "单条相交线段",
			lines: []Line{
				{Pt{X: 1, Y: 0}, Pt{X: 1, Y: 2}},
			},
			expected: []int{0},
		},
		{
			name: "单条不相交线段",
			lines: []Line{
				{Pt{X: 4, Y: 0}, Pt{X: 4, Y: 2}},
			},
			expected: []int{},
		},
		{
			name: "多条线段混合",
			lines: []Line{
				{Pt{X: 1, Y: 0}, Pt{X: 1, Y: 2}}, // 相交
				{Pt{X: 4, Y: 0}, Pt{X: 4, Y: 2}}, // 不相交
				{Pt{X: 2, Y: 0}, Pt{X: 2, Y: 2}}, // 相交
			},
			expected: []int{0, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			intersected := []int{}
			l.IntersectsLines(tc.lines, func(idx int) bool {
				intersected = append(intersected, idx)
				return true
			})

			if len(intersected) != len(tc.expected) {
				t.Errorf("相交线段数量 = %v, want %v", len(intersected), len(tc.expected))
				return
			}

			for i, expectedIdx := range tc.expected {
				if intersected[i] != expectedIdx {
					t.Errorf("相交线段索引[%d] = %v, want %v", i, intersected[i], expectedIdx)
				}
			}
		})
	}
}

// TestLine_XYOrderedPtsIdx 测试线段点的XY排序索引
func TestLine_XYOrderedPtsIdx(t *testing.T) {
	testCases := []struct {
		name        string
		line        Line
		expectLeft  int
		expectRight int
	}{
		{
			name:        "正常顺序 - 左下到右上",
			line:        Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 1}},
			expectLeft:  0,
			expectRight: 1,
		},
		{
			name:        "反向顺序 - 右上到左下",
			line:        Line{Pt{X: 1, Y: 1}, Pt{X: 0, Y: 0}},
			expectLeft:  1,
			expectRight: 0,
		},
		{
			name:        "垂直线段 - 下到上",
			line:        Line{Pt{X: 1, Y: 0}, Pt{X: 1, Y: 2}},
			expectLeft:  0,
			expectRight: 1,
		},
		{
			name:        "垂直线段 - 上到下",
			line:        Line{Pt{X: 1, Y: 2}, Pt{X: 1, Y: 0}},
			expectLeft:  1,
			expectRight: 0,
		},
		{
			name:        "水平线段 - 左到右",
			line:        Line{Pt{X: 0, Y: 1}, Pt{X: 2, Y: 1}},
			expectLeft:  0,
			expectRight: 1,
		},
		{
			name:        "水平线段 - 右到左",
			line:        Line{Pt{X: 2, Y: 1}, Pt{X: 0, Y: 1}},
			expectLeft:  1,
			expectRight: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			left, right := tc.line.XYOrderedPtsIdx()
			if left != tc.expectLeft || right != tc.expectRight {
				t.Errorf("XYOrderedPtsIdx() = (%v, %v), want (%v, %v)",
					left, right, tc.expectLeft, tc.expectRight)
			}
		})
	}
}

// BenchmarkDoesIntersect 性能基准测试
func BenchmarkDoesIntersect(b *testing.B) {
	s1 := Line{Pt{X: 0, Y: 0}, Pt{X: 2, Y: 2}}
	s2 := Line{Pt{X: 0, Y: 2}, Pt{X: 2, Y: 0}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DoesIntersect(s1, s2)
	}
}

// BenchmarkNewEventQueue 性能基准测试
func BenchmarkNewEventQueue(b *testing.B) {
	segments := []Line{
		{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 1}},
		{Pt{X: 1, Y: 0}, Pt{X: 0, Y: 1}},
		{Pt{X: 2, Y: 2}, Pt{X: 3, Y: 3}},
		{Pt{X: 3, Y: 2}, Pt{X: 2, Y: 3}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewEventQueue(segments)
	}
}

// BenchmarkFindIntersects 性能基准测试
func BenchmarkFindIntersects(b *testing.B) {
	segments := []Line{
		{Pt{X: 0, Y: 0}, Pt{X: 2, Y: 2}},
		{Pt{X: 0, Y: 2}, Pt{X: 2, Y: 0}},
		{Pt{X: 1, Y: 0}, Pt{X: 1, Y: 2}},
		{Pt{X: 0, Y: 1}, Pt{X: 2, Y: 1}},
	}

	fn := func(srcIdx, destIdx int, ptfn func() Pt) bool {
		return true
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindIntersects(segments, fn)
	}
}
