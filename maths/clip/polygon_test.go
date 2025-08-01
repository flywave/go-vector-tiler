package clip

import (
	"fmt"
	"testing"

	geom "github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"
)

func TestPolygon(t *testing.T) {
	type tcase struct {
		name       string
		clipExtent *gen.Extent
		poly       geom.Polygon
		expected   []geom.Polygon
		err        error
	}

	// 辅助函数：比较两个多边形是否相等（带容差）
	isPolygonEqual := func(p1, p2 geom.Polygon) bool {
		p1Data := p1.Data()
		p2Data := p2.Data()

		if len(p1Data) != len(p2Data) {
			return false
		}

		for i := range p1Data {
			ring1 := p1Data[i]
			ring2 := p2Data[i]

			// 检查环是否闭合
			if !pointsEqual(ring1[0], ring1[len(ring1)-1]) {
				ring1 = append(ring1, ring1[0])
			}
			if !pointsEqual(ring2[0], ring2[len(ring2)-1]) {
				ring2 = append(ring2, ring2[0])
			}

			// 环的长度可能在闭合后发生变化，重新检查
			if len(ring1) != len(ring2) {
				return false
			}

			// 检查点是否相同（考虑环的起点可能不同且可能顺时针/逆时针方向不同）
			foundStart := -1
			for j := range ring2 {
				if pointsEqual(ring1[0], ring2[j]) {
					foundStart = j
					break
				}
			}

			// 如果找不到起始点，尝试反转ring2后再找
			if foundStart == -1 {
				// 反转ring2
				reversedRing2 := make([][]float64, len(ring2))
				for j := 0; j < len(ring2); j++ {
					reversedRing2[j] = ring2[len(ring2)-1-j]
				}
				// 再次查找起始点
				for j := range reversedRing2 {
					if pointsEqual(ring1[0], reversedRing2[j]) {
						foundStart = j
						// 使用反转后的环
						ring2 = reversedRing2
						break
					}
				}
			}

			if foundStart == -1 {
				// 尝试使用ring1的第二个点作为起始点
				if len(ring1) > 1 {
					for j := range ring2 {
						if pointsEqual(ring1[1], ring2[j]) {
							foundStart = j
							break
						}
					}
				}
				if foundStart == -1 {
					return false
				}
			}

			// 按顺序比较所有点
			match := true
			for j := range ring1 {
				p1Idx := j
				p2Idx := (foundStart + j) % len(ring2)
				if !pointsEqual(ring1[p1Idx], ring2[p2Idx]) {
					match = false
					break
				}
			}

			// 如果不匹配，尝试反方向比较
			if !match {
				match = true
				for j := range ring1 {
					p1Idx := j
					p2Idx := (foundStart - j) % len(ring2)
					if p2Idx < 0 {
						p2Idx += len(ring2)
					}
					if !pointsEqual(ring1[p1Idx], ring2[p2Idx]) {
						match = false
						break
					}
				}
			}

			if !match {
				return false
			}
		}

		return true
	}

	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			got, gerr := Polygon(tc.poly, tc.clipExtent)
			switch {
			case tc.err != nil && gerr == nil:
				t.Errorf("%s: expected error, expected %v, got nil", tc.name, tc.err.Error())
				return
			case tc.err != nil && gerr != nil && tc.err.Error() != gerr.Error():
				t.Errorf("%s: unexpected error, expected %v, got %v", tc.name, tc.err.Error(), gerr.Error())
				return
			case tc.err == nil && gerr != nil:
				t.Errorf("%s: unexpected error, expected nil, got %v", tc.name, gerr.Error())
				return
			}
			if tc.err != nil {
				return
			}
			if len(tc.expected) != len(got) {
				t.Errorf("%s: number of polygons, expected %v got %v", tc.name, len(tc.expected), len(got))
				// 打印详细信息
				if len(got) == 0 {
					t.Logf("%s: no polygons returned", tc.name)
				} else {
					for i, p := range got {
						t.Logf("%s: returned polygon %d has %d rings", tc.name, i, len(p.Data()))
						for j, ring := range p.Data() {
							t.Logf("%s: ring %d has %d points", tc.name, j, len(ring))
							if len(ring) > 0 {
								t.Logf("First point: %v", ring[0])
								t.Logf("Last point: %v", ring[len(ring)-1])
							}
						}
					}
				}
				return
			}
			// 比较每个多边形
			for i := range tc.expected {
				if !isPolygonEqual(tc.expected[i], got[i]) {
					t.Errorf("%s: polygon %v not equal", tc.name, i)
					t.Errorf("%s: expected: %v", tc.name, tc.expected[i].Data())
					t.Errorf("%s: got: %v", tc.name, got[i].Data())
				}
			}
		}
	}

	// 创建测试范围
	extent00 := &gen.Extent{0, 0, 10, 10}
	extent01 := &gen.Extent{2, 2, 9, 9}
	extent03 := &gen.Extent{5, 1, 7, 3}
	// 打印范围信息
	fmt.Printf("extent00: %v\n", extent00)
	fmt.Printf("extent01: %v\n", extent01)
	fmt.Printf("extent03: %v\n", extent03)

	// 测试用例
	tests := map[string]tcase{
		"完全包含的多边形": {
			name:       "完全包含的多边形",
			clipExtent: extent00,
			poly: gen.NewPolygon([][][]float64{{{
				1, 1},
				{9, 1},
				{9, 9},
				{1, 9},
				{1, 1},
			}}),
			expected: []geom.Polygon{
				gen.NewPolygon([][][]float64{{{
					1, 1},
					{9, 1},
					{9, 9},
					{1, 9},
					{1, 1},
				}}),
			},
		},
		"部分重叠的多边形": {
			name:       "部分重叠的多边形",
			clipExtent: extent01,
			poly: gen.NewPolygon([][][]float64{{{
				0, 0},
				{10, 0},
				{10, 10},
				{0, 10},
				{0, 0},
			}}),
			expected: []geom.Polygon{
				gen.NewPolygon([][][]float64{{{
					2, 2},
					{9, 2},
					{9, 9},
					{2, 9},
					{2, 2},
				}}),
			},
		},
		"完全在裁剪区域外的多边形": {
			name:       "完全在裁剪区域外的多边形",
			clipExtent: extent03,
			poly: gen.NewPolygon([][][]float64{{{
				0, 0},
				{1, 0},
				{1, 1},
				{0, 1},
				{0, 0},
			}}),
			expected: nil,
		},
		"空多边形": {
			name:       "空多边形",
			clipExtent: extent00,
			poly:       gen.NewPolygon([][][]float64{}),
			expected:   nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}
}
