package util

import (
	"testing"

	geom "github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"
)

func TestIsGeometryEqual(t *testing.T) {
	// 测试 nil 情况
	cases := []struct {
		name     string
		g1       geom.Geometry
		g2       geom.Geometry
		expected bool
	}{
		{
			name:     "both nil",
			g1:       nil,
			g2:       nil,
			expected: true,
		},
		{
			name:     "one nil",
			g1:       nil,
			g2:       gen.NewPoint([]float64{0, 0}),
			expected: false,
		},
	}

	// 测试点类型
	for i := 0; i < 10; i++ {
		p := gen.NewPoint([]float64{float64(i), float64(i)})
		cases = append(cases, struct {
			name     string
			g1       geom.Geometry
			g2       geom.Geometry
			expected bool
		}{
			name:     "same point",
			g1:       p,
			g2:       p,
			expected: true,
		})

		cases = append(cases, struct {
			name     string
			g1       geom.Geometry
			g2       geom.Geometry
			expected bool
		}{
			name:     "different points",
			g1:       p,
			g2:       gen.NewPoint([]float64{float64(i + 1), float64(i + 1)}),
			expected: false,
		})

		// 测试浮点数容差
		cases = append(cases, struct {
			name     string
			g1       geom.Geometry
			g2       geom.Geometry
			expected bool
		}{
			name:     "close points within tolerance",
			g1:       gen.NewPoint([]float64{0, 0}),
			g2:       gen.NewPoint([]float64{1e-7, 1e-7}),
			expected: true,
		})

		cases = append(cases, struct {
			name     string
			g1       geom.Geometry
			g2       geom.Geometry
			expected bool
		}{
			name:     "points outside tolerance",
			g1:       gen.NewPoint([]float64{0, 0}),
			g2:       gen.NewPoint([]float64{1e-5, 1e-5}),
			expected: false,
		})
	}

	// 测试 3D 点类型
	for i := 0; i < 5; i++ {
		p := gen.NewPoint3([]float64{float64(i), float64(i), float64(i)})
		cases = append(cases, struct {
			name     string
			g1       geom.Geometry
			g2       geom.Geometry
			expected bool
		}{
			name:     "same 3D point",
			g1:       p,
			g2:       p,
			expected: true,
		})

		cases = append(cases, struct {
			name     string
			g1       geom.Geometry
			g2       geom.Geometry
			expected bool
		}{
			name:     "different 3D points",
			g1:       p,
			g2:       gen.NewPoint3([]float64{float64(i + 1), float64(i + 1), float64(i + 1)}),
			expected: false,
		})
	}

	// 测试线串类型
	for i := 0; i < 5; i++ {
		ls := gen.NewLineString([][]float64{
			{float64(i), float64(i)},
			{float64(i + 1), float64(i + 1)},
		})

		cases = append(cases, struct {
			name     string
			g1       geom.Geometry
			g2       geom.Geometry
			expected bool
		}{
			name:     "same linestring",
			g1:       ls,
			g2:       ls,
			expected: true,
		})

		cases = append(cases, struct {
			name     string
			g1       geom.Geometry
			g2       geom.Geometry
			expected bool
		}{
			name: "different linestring",
			g1:   ls,
			g2: gen.NewLineString([][]float64{
				{float64(i), float64(i)},
				{float64(i + 2), float64(i + 2)},
			}),
			expected: false,
		})
	}

	// 测试多边形类型
	for i := 0; i < 3; i++ {
		poly := gen.NewPolygon([][][]float64{
			{
				{float64(i), float64(i)},
				{float64(i + 1), float64(i)},
				{float64(i + 1), float64(i + 1)},
				{float64(i), float64(i + 1)},
				{float64(i), float64(i)},
			},
		})

		cases = append(cases, struct {
			name     string
			g1       geom.Geometry
			g2       geom.Geometry
			expected bool
		}{
			name:     "same polygon",
			g1:       poly,
			g2:       poly,
			expected: true,
		})

		cases = append(cases, struct {
			name     string
			g1       geom.Geometry
			g2       geom.Geometry
			expected bool
		}{
			name: "different polygon",
			g1:   poly,
			g2: gen.NewPolygon([][][]float64{
				{
					{float64(i), float64(i)},
					{float64(i + 2), float64(i)},
					{float64(i + 2), float64(i + 2)},
					{float64(i), float64(i + 2)},
					{float64(i), float64(i)},
				},
			}),
			expected: false,
		})
	}

	// 测试不同类型
	cases = append(cases, struct {
		name     string
		g1       geom.Geometry
		g2       geom.Geometry
		expected bool
	}{
		name:     "different types",
		g1:       gen.NewPoint([]float64{0, 0}),
		g2:       gen.NewLineString([][]float64{{0, 0}, {1, 1}}),
		expected: false,
	})

	// 运行测试
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsGeometryEqual(tc.g1, tc.g2)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// 更多针对特定类型的测试
func TestIsPointEqual(t *testing.T) {
	cases := []struct {
		name     string
		p1       geom.Point
		p2       geom.Point
		expected bool
	}{
		{"both nil", nil, nil, true},
		{"one nil", nil, gen.NewPoint([]float64{0, 0}), false},
		{"same point", gen.NewPoint([]float64{1, 2}), gen.NewPoint([]float64{1, 2}), true},
		{"different points", gen.NewPoint([]float64{1, 2}), gen.NewPoint([]float64{3, 4}), false},
		{"close points", gen.NewPoint([]float64{0, 0}), gen.NewPoint([]float64{1e-7, 1e-7}), true},
		{"not close enough", gen.NewPoint([]float64{0, 0}), gen.NewPoint([]float64{1e-5, 1e-5}), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsPointEqual(tc.p1, tc.p2)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestIsPoint3Equal(t *testing.T) {
	cases := []struct {
		name     string
		p1       geom.Point3
		p2       geom.Point3
		expected bool
	}{
		{"both nil", nil, nil, true},
		{"one nil", nil, gen.NewPoint3([]float64{0, 0, 0}), false},
		{"same point", gen.NewPoint3([]float64{1, 2, 3}), gen.NewPoint3([]float64{1, 2, 3}), true},
		{"different points", gen.NewPoint3([]float64{1, 2, 3}), gen.NewPoint3([]float64{4, 5, 6}), false},
		{"close points", gen.NewPoint3([]float64{0, 0, 0}), gen.NewPoint3([]float64{1e-7, 1e-7, 1e-7}), true},
		{"not close enough", gen.NewPoint3([]float64{0, 0, 0}), gen.NewPoint3([]float64{1e-5, 1e-5, 1e-5}), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsPoint3Equal(tc.p1, tc.p2)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// 可以根据需要添加更多针对其他类型的测试函数
