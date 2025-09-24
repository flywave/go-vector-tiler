package maths

import (
	"sort"
	"testing"
)

// TestTriangle_FindEdge 测试三角形边查找功能
func TestTriangle_FindEdge(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0}, // 顶点0
		Pt{X: 1, Y: 0}, // 顶点1
		Pt{X: 0, Y: 1}, // 顶点2
	}

	testCases := []struct {
		name        string
		edge        Line
		expectedIdx int
		expectError bool
	}{
		{
			name:        "边0-1",
			edge:        Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}},
			expectedIdx: 0,
			expectError: false,
		},
		{
			name:        "边1-0 (反向)",
			edge:        Line{Pt{X: 1, Y: 0}, Pt{X: 0, Y: 0}},
			expectedIdx: 0,
			expectError: false,
		},
		{
			name:        "边1-2",
			edge:        Line{Pt{X: 1, Y: 0}, Pt{X: 0, Y: 1}},
			expectedIdx: 1,
			expectError: false,
		},
		{
			name:        "边2-0",
			edge:        Line{Pt{X: 0, Y: 1}, Pt{X: 0, Y: 0}},
			expectedIdx: 2,
			expectError: false,
		},
		{
			name:        "不存在的边",
			edge:        Line{Pt{X: 2, Y: 2}, Pt{X: 3, Y: 3}},
			expectedIdx: -1,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			idx, err := triangle.FindEdge(tc.edge)

			if tc.expectError && err == nil {
				t.Errorf("期望错误但未返回错误")
			}
			if !tc.expectError && err != nil {
				t.Errorf("不期望错误但返回错误: %v", err)
			}
			if idx != tc.expectedIdx {
				t.Errorf("FindEdge() = %v, want %v", idx, tc.expectedIdx)
			}
		})
	}
}

// TestTriangle_Edge 测试获取三角形边功能
func TestTriangle_Edge(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0}, // 顶点0
		Pt{X: 1, Y: 0}, // 顶点1
		Pt{X: 0, Y: 1}, // 顶点2
	}

	testCases := []struct {
		name     string
		n        int
		expected Line
	}{
		{
			name:     "边0",
			n:        0,
			expected: Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}},
		},
		{
			name:     "边1",
			n:        1,
			expected: Line{Pt{X: 1, Y: 0}, Pt{X: 0, Y: 1}},
		},
		{
			name:     "边2",
			n:        2,
			expected: Line{Pt{X: 0, Y: 1}, Pt{X: 0, Y: 0}},
		},
		{
			name:     "负数索引",
			n:        -1,
			expected: Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}},
		},
		{
			name:     "超出范围索引",
			n:        5,
			expected: Line{Pt{X: 0, Y: 1}, Pt{X: 0, Y: 0}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := triangle.Edge(tc.n)
			if !result[0].IsEqual(tc.expected[0]) || !result[1].IsEqual(tc.expected[1]) {
				t.Errorf("Edge(%d) = %v, want %v", tc.n, result, tc.expected)
			}
		})
	}
}

// TestTriangle_LREdge 测试获取三角形LR边功能
func TestTriangle_LREdge(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0}, // 顶点0
		Pt{X: 1, Y: 0}, // 顶点1
		Pt{X: 0, Y: 1}, // 顶点2
	}

	testCases := []struct {
		name     string
		n        int
		expected Line
	}{
		{
			name:     "LR边0",
			n:        0,
			expected: Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}},
		},
		{
			name:     "LR边1",
			n:        1,
			expected: Line{Pt{X: 1, Y: 0}, Pt{X: 0, Y: 1}},
		},
		{
			name:     "LR边2",
			n:        2,
			expected: Line{Pt{X: 0, Y: 0}, Pt{X: 0, Y: 1}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := triangle.LREdge(tc.n)
			if !result[0].IsEqual(tc.expected[0]) || !result[1].IsEqual(tc.expected[1]) {
				t.Errorf("LREdge(%d) = %v, want %v", tc.n, result, tc.expected)
			}
		})
	}
}

// TestTriangle_Edges 测试获取三角形所有边
func TestTriangle_Edges(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0}, // 顶点0
		Pt{X: 1, Y: 0}, // 顶点1
		Pt{X: 0, Y: 1}, // 顶点2
	}

	edges := triangle.Edges()

	expectedEdges := [3]Line{
		{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}}, // 边0-1
		{Pt{X: 1, Y: 0}, Pt{X: 0, Y: 1}}, // 边1-2
		{Pt{X: 0, Y: 1}, Pt{X: 0, Y: 0}}, // 边2-0
	}

	for i, edge := range edges {
		if !edge[0].IsEqual(expectedEdges[i][0]) || !edge[1].IsEqual(expectedEdges[i][1]) {
			t.Errorf("Edges()[%d] = %v, want %v", i, edge, expectedEdges[i])
		}
	}
}

// TestTriangle_LREdges 测试获取三角形所有LR边
func TestTriangle_LREdges(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0}, // 顶点0
		Pt{X: 1, Y: 0}, // 顶点1
		Pt{X: 0, Y: 1}, // 顶点2
	}

	edges := triangle.LREdges()

	expectedEdges := [3]Line{
		{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}}, // 边0-1
		{Pt{X: 1, Y: 0}, Pt{X: 0, Y: 1}}, // 边1-2
		{Pt{X: 0, Y: 0}, Pt{X: 0, Y: 1}}, // 边0-2
	}

	for i, edge := range edges {
		if !edge[0].IsEqual(expectedEdges[i][0]) || !edge[1].IsEqual(expectedEdges[i][1]) {
			t.Errorf("LREdges()[%d] = %v, want %v", i, edge, expectedEdges[i])
		}
	}
}

// TestTriangle_EdgeIdx 测试边索引查找
func TestTriangle_EdgeIdx(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0}, // 顶点0
		Pt{X: 1, Y: 0}, // 顶点1
		Pt{X: 0, Y: 1}, // 顶点2
	}

	testCases := []struct {
		name     string
		pt1, pt2 Pt
		expected int
	}{
		{
			name:     "边0-1",
			pt1:      Pt{X: 0, Y: 0},
			pt2:      Pt{X: 1, Y: 0},
			expected: 0,
		},
		{
			name:     "边1-2",
			pt1:      Pt{X: 1, Y: 0},
			pt2:      Pt{X: 0, Y: 1},
			expected: 1,
		},
		{
			name:     "边2-0",
			pt1:      Pt{X: 0, Y: 1},
			pt2:      Pt{X: 0, Y: 0},
			expected: 2,
		},
		{
			name:     "不存在的边",
			pt1:      Pt{X: 2, Y: 2},
			pt2:      Pt{X: 3, Y: 3},
			expected: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := triangle.EdgeIdx(tc.pt1, tc.pt2)
			if result != tc.expected {
				t.Errorf("EdgeIdx(%v, %v) = %v, want %v", tc.pt1, tc.pt2, result, tc.expected)
			}
		})
	}

	// 测试nil三角形
	var nilTriangle *Triangle
	result := nilTriangle.EdgeIdx(Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0})
	if result != -1 {
		t.Errorf("nil三角形EdgeIdx应该返回-1，但得到: %v", result)
	}
}

// TestTriangle_Key 测试三角形键生成
func TestTriangle_Key(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 1, Y: 0},
		Pt{X: 0, Y: 0},
		Pt{X: 0, Y: 1},
	}

	key := triangle.Key()

	// 键应该基于排序后的顶点
	if key == "" {
		t.Error("Key()不应该返回空字符串")
	}

	// 测试nil三角形
	var nilTriangle *Triangle
	nilKey := nilTriangle.Key()
	if nilKey != "" {
		t.Errorf("nil三角形Key应该返回空字符串，但得到: %v", nilKey)
	}

	// 测试相同三角形生成相同的键
	triangle2 := &Triangle{
		Pt{X: 0, Y: 1},
		Pt{X: 1, Y: 0},
		Pt{X: 0, Y: 0},
	}
	key2 := triangle2.Key()
	if key != key2 {
		t.Errorf("相同的三角形应该生成相同的键: %v != %v", key, key2)
	}
}

// TestTriangle_Points 测试获取三角形顶点
func TestTriangle_Points(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0},
		Pt{X: 1, Y: 0},
		Pt{X: 0, Y: 1},
	}

	points := triangle.Points()
	expectedPoints := []Pt{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
	}

	if len(points) != 3 {
		t.Errorf("Points()应该返回3个点，但得到: %v", len(points))
	}

	for i, pt := range points {
		if !pt.IsEqual(expectedPoints[i]) {
			t.Errorf("Points()[%d] = %v, want %v", i, pt, expectedPoints[i])
		}
	}
}

// TestTriangle_Point 测试通过索引获取顶点
func TestTriangle_Point(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0},
		Pt{X: 1, Y: 0},
		Pt{X: 0, Y: 1},
	}

	testCases := []struct {
		name     string
		index    int
		expected Pt
	}{
		{
			name:     "索引0",
			index:    0,
			expected: Pt{X: 0, Y: 0},
		},
		{
			name:     "索引1",
			index:    1,
			expected: Pt{X: 1, Y: 0},
		},
		{
			name:     "索引2",
			index:    2,
			expected: Pt{X: 0, Y: 1},
		},
		{
			name:     "负索引",
			index:    -1,
			expected: Pt{X: 0, Y: 0},
		},
		{
			name:     "超出范围索引",
			index:    5,
			expected: Pt{X: 0, Y: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := triangle.Point(tc.index)
			if !result.IsEqual(tc.expected) {
				t.Errorf("Point(%d) = %v, want %v", tc.index, result, tc.expected)
			}
		})
	}
}

// TestTriangle_Len 测试三角形长度
func TestTriangle_Len(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0},
		Pt{X: 1, Y: 0},
		Pt{X: 0, Y: 1},
	}

	length := triangle.Len()
	if length != 3 {
		t.Errorf("Len() = %v, want %v", length, 3)
	}

	// 测试nil三角形
	var nilTriangle *Triangle
	nilLength := nilTriangle.Len()
	if nilLength != 0 {
		t.Errorf("nil三角形Len应该返回0，但得到: %v", nilLength)
	}
}

// TestTriangle_Sort 测试三角形排序功能
func TestTriangle_Sort(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 1, Y: 1},
		Pt{X: 0, Y: 0},
		Pt{X: 2, Y: 0},
	}

	// 排序前检查
	originalTriangle := *triangle

	// 执行排序
	sort.Sort(triangle)

	// 验证排序结果
	for i := 0; i < 2; i++ {
		if XYOrder(triangle[i], triangle[i+1]) == 1 {
			t.Errorf("排序后三角形顶点应该按XY顺序排列，位置%d不正确", i)
		}
	}

	// 验证所有原始点仍然存在
	found := make([]bool, 3)
	for _, originalPt := range originalTriangle {
		for j, sortedPt := range triangle {
			if originalPt.IsEqual(sortedPt) {
				found[j] = true
				break
			}
		}
	}
	for i, isFound := range found {
		if !isFound {
			t.Errorf("排序后丢失了原始顶点%d", i)
		}
	}
}

// TestTriangle_Equal 测试三角形相等性比较
func TestTriangle_Equal(t *testing.T) {
	triangle1 := &Triangle{
		Pt{X: 0, Y: 0},
		Pt{X: 1, Y: 0},
		Pt{X: 0, Y: 1},
	}

	triangle2 := &Triangle{
		Pt{X: 1, Y: 0}, // 不同顺序，但相同的点
		Pt{X: 0, Y: 1},
		Pt{X: 0, Y: 0},
	}

	triangle3 := &Triangle{
		Pt{X: 2, Y: 2}, // 不同的点
		Pt{X: 3, Y: 3},
		Pt{X: 4, Y: 4},
	}

	testCases := []struct {
		name     string
		t1       *Triangle
		t2       *Triangle
		expected bool
	}{
		{
			name:     "相同三角形",
			t1:       triangle1,
			t2:       triangle1,
			expected: true,
		},
		{
			name:     "相同点不同顺序",
			t1:       triangle1,
			t2:       triangle2,
			expected: true,
		},
		{
			name:     "不同三角形",
			t1:       triangle1,
			t2:       triangle3,
			expected: false,
		},
		{
			name:     "nil与非nil",
			t1:       nil,
			t2:       triangle1,
			expected: false,
		},
		{
			name:     "非nil与nil",
			t1:       triangle1,
			t2:       nil,
			expected: false,
		},
		{
			name:     "两个nil",
			t1:       nil,
			t2:       nil,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.t1.Equal(tc.t2)
			if result != tc.expected {
				t.Errorf("Equal() = %v, want %v", result, tc.expected)
			}
		})
	}
}

// TestTriangle_EqualAnyPt 测试三角形是否包含任意指定点
func TestTriangle_EqualAnyPt(t *testing.T) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0},
		Pt{X: 1, Y: 0},
		Pt{X: 0, Y: 1},
	}

	testCases := []struct {
		name     string
		pts      []Pt
		expected bool
	}{
		{
			name:     "包含一个顶点",
			pts:      []Pt{{X: 0, Y: 0}},
			expected: true,
		},
		{
			name:     "包含多个顶点",
			pts:      []Pt{{X: 0, Y: 0}, {X: 1, Y: 0}},
			expected: true,
		},
		{
			name:     "不包含任何顶点",
			pts:      []Pt{{X: 2, Y: 2}, {X: 3, Y: 3}},
			expected: false,
		},
		{
			name:     "空点列表",
			pts:      []Pt{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := triangle.EqualAnyPt(tc.pts...)
			if result != tc.expected {
				t.Errorf("EqualAnyPt() = %v, want %v", result, tc.expected)
			}
		})
	}

	// 测试nil三角形
	var nilTriangle *Triangle
	result := nilTriangle.EqualAnyPt(Pt{X: 0, Y: 0})
	if result != false {
		t.Errorf("nil三角形EqualAnyPt应该返回false，但得到: %v", result)
	}
}

// TestAreaOfTriangle 测试三角形面积计算
func TestAreaOfTriangle(t *testing.T) {
	testCases := []struct {
		name       string
		v0, v1, v2 Pt
		expected   float64
	}{
		{
			name:     "标准直角三角形",
			v0:       Pt{X: 0, Y: 0},
			v1:       Pt{X: 1, Y: 0},
			v2:       Pt{X: 0, Y: 1},
			expected: 0.5,
		},
		{
			name:     "正方形对角三角形",
			v0:       Pt{X: 0, Y: 0},
			v1:       Pt{X: 2, Y: 0},
			v2:       Pt{X: 0, Y: 2},
			expected: 2.0,
		},
		{
			name:     "零面积（共线点）",
			v0:       Pt{X: 0, Y: 0},
			v1:       Pt{X: 1, Y: 0},
			v2:       Pt{X: 2, Y: 0},
			expected: 0.0,
		},
		{
			name:     "负面积（逆时针）",
			v0:       Pt{X: 0, Y: 0},
			v1:       Pt{X: 0, Y: 1},
			v2:       Pt{X: 1, Y: 0},
			expected: -0.5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := AreaOfTriangle(tc.v0, tc.v1, tc.v2)
			if result != tc.expected {
				t.Errorf("AreaOfTriangle() = %v, want %v", result, tc.expected)
			}
		})
	}
}

// TestTriangle_Area 测试三角形面积方法
func TestTriangle_Area(t *testing.T) {
	testCases := []struct {
		name     string
		triangle *Triangle
		expected float64
	}{
		{
			name: "标准直角三角形",
			triangle: &Triangle{
				Pt{X: 0, Y: 0},
				Pt{X: 1, Y: 0},
				Pt{X: 0, Y: 1},
			},
			expected: 0.5,
		},
		{
			name: "逆时针三角形（绝对值）",
			triangle: &Triangle{
				Pt{X: 0, Y: 0},
				Pt{X: 0, Y: 1},
				Pt{X: 1, Y: 0},
			},
			expected: 0.5,
		},
		{
			name:     "nil三角形",
			triangle: nil,
			expected: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.triangle.Area()
			if result != tc.expected {
				t.Errorf("Area() = %v, want %v", result, tc.expected)
			}
		})
	}
}

// TestTriangle_Center 测试三角形中心点计算
func TestTriangle_Center(t *testing.T) {
	testCases := []struct {
		name     string
		triangle *Triangle
		expected Pt
	}{
		{
			name: "标准三角形",
			triangle: &Triangle{
				Pt{X: 0, Y: 0},
				Pt{X: 3, Y: 0},
				Pt{X: 0, Y: 3},
			},
			expected: Pt{X: 1, Y: 1},
		},
		{
			name: "等边三角形",
			triangle: &Triangle{
				Pt{X: 0, Y: 0},
				Pt{X: 1, Y: 0},
				Pt{X: 0.5, Y: 1},
			},
			expected: Pt{X: 0.5, Y: 1.0 / 3.0},
		},
		{
			name:     "nil三角形",
			triangle: nil,
			expected: Pt{X: 0, Y: 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.triangle.Center()
			if !result.IsEqual(tc.expected) {
				t.Errorf("Center() = %v, want %v", result, tc.expected)
			}
		})
	}
}

// TestNewTriangle 测试三角形构造函数
func TestNewTriangle(t *testing.T) {
	testCases := []struct {
		name          string
		pt1, pt2, pt3 Pt
		expectPanic   bool
	}{
		{
			name:        "正常三角形",
			pt1:         Pt{X: 0, Y: 0},
			pt2:         Pt{X: 1, Y: 0},
			pt3:         Pt{X: 0, Y: 1},
			expectPanic: false,
		},
		{
			name:        "两个相同点",
			pt1:         Pt{X: 0, Y: 0},
			pt2:         Pt{X: 0, Y: 0},
			pt3:         Pt{X: 1, Y: 1},
			expectPanic: true,
		},
		{
			name:        "三个相同点",
			pt1:         Pt{X: 1, Y: 1},
			pt2:         Pt{X: 1, Y: 1},
			pt3:         Pt{X: 1, Y: 1},
			expectPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewTriangle应该panic，但没有")
					}
				}()
				NewTriangle(tc.pt1, tc.pt2, tc.pt3)
			} else {
				triangle := NewTriangle(tc.pt1, tc.pt2, tc.pt3)

				// 验证三角形包含所有输入点
				containsPt1 := triangle[0].IsEqual(tc.pt1) || triangle[1].IsEqual(tc.pt1) || triangle[2].IsEqual(tc.pt1)
				containsPt2 := triangle[0].IsEqual(tc.pt2) || triangle[1].IsEqual(tc.pt2) || triangle[2].IsEqual(tc.pt2)
				containsPt3 := triangle[0].IsEqual(tc.pt3) || triangle[1].IsEqual(tc.pt3) || triangle[2].IsEqual(tc.pt3)

				if !containsPt1 || !containsPt2 || !containsPt3 {
					t.Errorf("NewTriangle未包含所有输入点")
				}

				// 验证排序
				for i := 0; i < 2; i++ {
					if XYOrder(triangle[i], triangle[i+1]) == 1 {
						t.Errorf("NewTriangle返回的三角形未正确排序")
					}
				}
			}
		})
	}
}

// TestPointPairs 测试点对生成函数
func TestPointPairs(t *testing.T) {
	testCases := []struct {
		name        string
		pts         []Pt
		expected    int
		expectError bool
	}{
		{
			name:        "一个点",
			pts:         []Pt{{X: 0, Y: 0}},
			expected:    0,
			expectError: true,
		},
		{
			name:        "两个点",
			pts:         []Pt{{X: 0, Y: 0}, {X: 1, Y: 1}},
			expected:    1,
			expectError: false,
		},
		{
			name:        "三个点",
			pts:         []Pt{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}},
			expected:    3,
			expectError: false,
		},
		{
			name:        "四个点",
			pts:         []Pt{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 3, Y: 3}},
			expected:    6,
			expectError: false,
		},
		{
			name:        "五个点",
			pts:         []Pt{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 3, Y: 3}, {X: 4, Y: 4}},
			expected:    10,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pairs, err := PointPairs(tc.pts)

			if tc.expectError && err == nil {
				t.Errorf("期望错误但未返回错误")
			}
			if !tc.expectError && err != nil {
				t.Errorf("不期望错误但返回错误: %v", err)
			}

			if !tc.expectError {
				if len(pairs) != tc.expected {
					t.Errorf("PointPairs() 生成对数 = %v, want %v", len(pairs), tc.expected)
				}

				// 验证所有对都不相同
				for i := 0; i < len(pairs); i++ {
					for j := i + 1; j < len(pairs); j++ {
						if (pairs[i][0].IsEqual(pairs[j][0]) && pairs[i][1].IsEqual(pairs[j][1])) ||
							(pairs[i][0].IsEqual(pairs[j][1]) && pairs[i][1].IsEqual(pairs[j][0])) {
							t.Errorf("发现重复的点对: %v 和 %v", pairs[i], pairs[j])
						}
					}
				}
			}
		})
	}
}

// TestByXY 测试点排序功能
func TestByXY(t *testing.T) {
	pts := []Pt{
		{X: 2, Y: 3},
		{X: 1, Y: 2},
		{X: 3, Y: 1},
		{X: 1, Y: 1},
	}

	sort.Sort(ByXY(pts))

	// 验证排序结果
	for i := 0; i < len(pts)-1; i++ {
		if XYOrder(pts[i], pts[i+1]) == 1 {
			t.Errorf("点排序不正确，位置%d: %v > %v", i, pts[i], pts[i+1])
		}
	}
}

// TestByXYLine 测试线段排序功能
func TestByXYLine(t *testing.T) {
	lines := []Line{
		{Pt{X: 2, Y: 2}, Pt{X: 3, Y: 3}},
		{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 1}},
		{Pt{X: 1, Y: 1}, Pt{X: 2, Y: 2}},
	}

	sort.Sort(ByXYLine(lines))

	// 验证排序结果
	for i := 0; i < len(lines)-1; i++ {
		li, lj := lines[i].LeftRightMostAsLine(), lines[i+1].LeftRightMostAsLine()
		ret := XYOrder(li[0], lj[0])
		if ret == 0 {
			ret = XYOrder(li[1], lj[1])
		}
		if ret == 1 {
			t.Errorf("线段排序不正确，位置%d", i)
		}
	}
}

// BenchmarkTriangle_Area 性能基准测试
func BenchmarkTriangle_Area(b *testing.B) {
	triangle := &Triangle{
		Pt{X: 0, Y: 0},
		Pt{X: 1, Y: 0},
		Pt{X: 0, Y: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		triangle.Area()
	}
}

// BenchmarkAreaOfTriangle 性能基准测试
func BenchmarkAreaOfTriangle(b *testing.B) {
	v0, v1, v2 := Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}, Pt{X: 0, Y: 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AreaOfTriangle(v0, v1, v2)
	}
}

// BenchmarkNewTriangle 性能基准测试
func BenchmarkNewTriangle(b *testing.B) {
	pt1, pt2, pt3 := Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}, Pt{X: 0, Y: 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewTriangle(pt1, pt2, pt3)
	}
}

// BenchmarkPointPairs 性能基准测试
func BenchmarkPointPairs(b *testing.B) {
	pts := []Pt{
		{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 3, Y: 3}, {X: 4, Y: 4},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PointPairs(pts)
	}
}

// TestTriangleNode_LabelAs 测试三角形节点标记功能
func TestTriangleNode_LabelAs(t *testing.T) {
	node := &TriangleNode{
		Triangle: Triangle{
			Pt{X: 0, Y: 0},
			Pt{X: 1, Y: 0},
			Pt{X: 0, Y: 1},
		},
		Label: Unknown,
	}

	// 测试第一次标记
	unlabeled := node.LabelAs(Inside, false)
	if node.Label != Inside {
		t.Errorf("LabelAs后标签应该是Inside，但得到: %v", node.Label)
	}

	// 测试强制标记
	unlabeled = node.LabelAs(Outside, true)
	if node.Label != Outside {
		t.Errorf("强制LabelAs后标签应该是Outside，但得到: %v", node.Label)
	}

	// 测试非强制标记已标记的节点
	unlabeled = node.LabelAs(Inside, false)
	if node.Label != Outside {
		t.Errorf("非强制LabelAs不应该改变已标记的节点，标签应该保持Outside，但得到: %v", node.Label)
	}
	if len(unlabeled) != 0 {
		t.Errorf("非强制标记已标记节点应该返回空切片，但得到: %v", len(unlabeled))
	}

	// 测试nil节点
	var nilNode *TriangleNode
	unlabeled = nilNode.LabelAs(Inside, false)
	if len(unlabeled) != 0 {
		t.Errorf("nil节点LabelAs应该返回空切片，但得到: %v", len(unlabeled))
	}
}

// TestLabel_String 测试标签字符串表示
func TestLabel_String(t *testing.T) {
	testCases := []struct {
		label    Label
		expected string
	}{
		{Unknown, "unknown"},
		{Outside, "outside"},
		{Inside, "inside"},
		{Label(99), "unknown"}, // 未知标签
	}

	for _, tc := range testCases {
		result := tc.label.String()
		if result != tc.expected {
			t.Errorf("Label(%v).String() = %v, want %v", tc.label, result, tc.expected)
		}
	}
}

// TestNewPointList 测试点列表创建
func TestNewPointList(t *testing.T) {
	line := Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 1}}
	pl := NewPointList(line)

	if pl.Head == nil || pl.Tail == nil {
		t.Error("NewPointList应该创建非空的头尾节点")
	}

	if !pl.Head.Pt.IsEqual(line[0]) {
		t.Errorf("头节点应该是 %v，但得到: %v", line[0], pl.Head.Pt)
	}

	if !pl.Tail.Pt.IsEqual(line[1]) {
		t.Errorf("尾节点应该是 %v，但得到: %v", line[1], pl.Tail.Pt)
	}

	if pl.Head.Next != pl.Tail {
		t.Error("头节点应该指向尾节点")
	}

	if pl.isComplete {
		t.Error("新创建的点列表不应该是完整的")
	}
}

// TestPointList_TryAddLine 测试向点列表添加线段
func TestPointList_TryAddLine(t *testing.T) {
	initialLine := Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 1}}
	pl := NewPointList(initialLine)

	testCases := []struct {
		name     string
		line     Line
		expected bool
	}{
		{
			name:     "连接到头部",
			line:     Line{Pt{X: -1, Y: -1}, Pt{X: 0, Y: 0}},
			expected: true,
		},
		{
			name:     "连接到尾部",
			line:     Line{Pt{X: 1, Y: 1}, Pt{X: 2, Y: 2}},
			expected: true,
		},
		{
			name:     "不相关的线段",
			line:     Line{Pt{X: 5, Y: 5}, Pt{X: 6, Y: 6}},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 重新创建点列表以避免状态干扰
			pl := NewPointList(initialLine)
			result := pl.TryAddLine(tc.line)
			if result != tc.expected {
				t.Errorf("TryAddLine() = %v, want %v", result, tc.expected)
			}
		})
	}

	// 测试完成环路
	pl = NewPointList(Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 0}})
	pl.TryAddLine(Line{Pt{X: 1, Y: 0}, Pt{X: 1, Y: 1}})
	pl.TryAddLine(Line{Pt{X: 1, Y: 1}, Pt{X: 0, Y: 1}})

	// 闭合环路
	result := pl.TryAddLine(Line{Pt{X: 0, Y: 1}, Pt{X: 0, Y: 0}})
	if !result {
		t.Error("闭合环路应该成功")
	}
	if !pl.IsComplete() {
		t.Error("闭合后点列表应该是完整的")
	}

	// 测试完整列表拒绝新线段
	result = pl.TryAddLine(Line{Pt{X: 2, Y: 2}, Pt{X: 3, Y: 3}})
	if result {
		t.Error("完整的点列表不应该接受新线段")
	}
}

// TestPointList_AsRing 测试点列表转换为环
func TestPointList_AsRing(t *testing.T) {
	line := Line{Pt{X: 0, Y: 0}, Pt{X: 1, Y: 1}}
	pl := NewPointList(line)
	pl.TryAddLine(Line{Pt{X: 1, Y: 1}, Pt{X: 2, Y: 2}})

	ring := pl.AsRing()
	expectedPoints := []Pt{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 2},
	}

	if len(ring) != len(expectedPoints) {
		t.Errorf("环长度应该是 %v，但得到: %v", len(expectedPoints), len(ring))
	}

	for i, pt := range ring {
		if !pt.IsEqual(expectedPoints[i]) {
			t.Errorf("环点[%d]应该是 %v，但得到: %v", i, expectedPoints[i], pt)
		}
	}
}
