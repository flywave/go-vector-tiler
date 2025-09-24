package svg

import (
	"reflect"
	"testing"

	geom "github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"
	"github.com/gdey/tbltest"
)

func TestMinMaxBasicMethods(t *testing.T) {
	type testcase struct {
		desc           string
		minmax         MinMax
		expectedMin    [2]int64
		expectedMax    [2]int64
		expectedWidth  int64
		expectedHeight int64
	}

	tests := tbltest.Cases(
		testcase{
			desc:           "Empty MinMax",
			minmax:         MinMax{},
			expectedMin:    [2]int64{0, 0},
			expectedMax:    [2]int64{0, 0},
			expectedWidth:  0,
			expectedHeight: 0,
		},
		testcase{
			desc:           "Simple MinMax",
			minmax:         MinMax{10, 20, 30, 40, true},
			expectedMin:    [2]int64{10, 20},
			expectedMax:    [2]int64{30, 40},
			expectedWidth:  20,
			expectedHeight: 20,
		},
		testcase{
			desc:           "Negative coordinates MinMax",
			minmax:         MinMax{-5, -10, 5, 10, true},
			expectedMin:    [2]int64{-5, -10},
			expectedMax:    [2]int64{5, 10},
			expectedWidth:  10,
			expectedHeight: 20,
		},
	)

	tests.Run(func(idx int, test testcase) {
		minX, minY := test.minmax.Min()
		if minX != test.expectedMin[0] || minY != test.expectedMin[1] {
			t.Errorf("Test %v (%v): Min() expected %v, got (%v, %v)", idx, test.desc, test.expectedMin, minX, minY)
		}

		maxX, maxY := test.minmax.Max()
		if maxX != test.expectedMax[0] || maxY != test.expectedMax[1] {
			t.Errorf("Test %v (%v): Max() expected %v, got (%v, %v)", idx, test.desc, test.expectedMax, maxX, maxY)
		}

		width := test.minmax.Width()
		if width != test.expectedWidth {
			t.Errorf("Test %v (%v): Width() expected %v, got %v", idx, test.desc, test.expectedWidth, width)
		}

		height := test.minmax.Height()
		if height != test.expectedHeight {
			t.Errorf("Test %v (%v): Height() expected %v, got %v", idx, test.desc, test.expectedHeight, height)
		}
	})
}

func TestMinMaxSentinalPts(t *testing.T) {
	type testcase struct {
		desc     string
		minmax   MinMax
		expected [][]int64
	}

	tests := tbltest.Cases(
		testcase{
			desc:     "Simple MinMax",
			minmax:   MinMax{10, 20, 30, 40, true},
			expected: [][]int64{{10, 20}, {30, 20}, {30, 40}, {10, 40}},
		},
	)

	tests.Run(func(idx int, test testcase) {
		pts := test.minmax.SentinalPts()
		if !reflect.DeepEqual(pts, test.expected) {
			t.Errorf("Test %v (%v): SentinalPts() expected %v, got %v", idx, test.desc, test.expected, pts)
		}
	})
}

func TestMinMaxMerge(t *testing.T) {
	type testcase struct {
		desc     string
		mm1      *MinMax
		mm2      *MinMax
		expected *MinMax
	}

	tests := tbltest.Cases(
		testcase{
			desc:     "Merge two MinMax",
			mm1:      &MinMax{10, 20, 30, 40, true},
			mm2:      &MinMax{5, 15, 35, 45, true},
			expected: &MinMax{5, 15, 35, 45, true},
		},
		testcase{
			desc:     "Merge with nil mm1",
			mm1:      nil,
			mm2:      &MinMax{5, 15, 35, 45, true},
			expected: &MinMax{5, 15, 35, 45, true},
		},
		testcase{
			desc:     "Merge with nil mm2",
			mm1:      &MinMax{10, 20, 30, 40, true},
			mm2:      nil,
			expected: &MinMax{10, 20, 30, 40, true},
		},
		testcase{
			desc:     "Merge both nil",
			mm1:      nil,
			mm2:      nil,
			expected: &MinMax{0, 0, 0, 0, false},
		},
	)

	tests.Run(func(idx int, test testcase) {
		result := test.mm1.MinMax(test.mm2)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %v (%v): MinMax() expected %v, got %v", idx, test.desc, test.expected, result)
		}
	})
}

func TestMinMaxMinMaxPt(t *testing.T) {
	type testcase struct {
		desc     string
		mm       *MinMax
		ptX, ptY int64
		expected *MinMax
	}

	tests := tbltest.Cases(
		testcase{
			desc:     "Add point to MinMax",
			mm:       &MinMax{10, 20, 30, 40, true},
			ptX:      5,
			ptY:      15,
			expected: &MinMax{5, 15, 30, 40, true},
		},
		testcase{
			desc:     "Add point outside existing MinMax",
			mm:       &MinMax{10, 20, 30, 40, true},
			ptX:      40,
			ptY:      50,
			expected: &MinMax{10, 20, 40, 50, true},
		},
		testcase{
			desc:     "Add point inside existing MinMax",
			mm:       &MinMax{10, 20, 30, 40, true},
			ptX:      15,
			ptY:      25,
			expected: &MinMax{10, 20, 30, 40, true},
		},
		testcase{
			desc:     "Add point to nil MinMax",
			mm:       nil,
			ptX:      5,
			ptY:      15,
			expected: &MinMax{5, 15, 5, 15, true},
		},
	)

	tests.Run(func(idx int, test testcase) {
		result := test.mm.MinMaxPt(test.ptX, test.ptY)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %v (%v): MinMaxPt() expected %v, got %v", idx, test.desc, test.expected, result)
		}
	})
}

func TestMinMaxOfGeometry(t *testing.T) {
	type testcase struct {
		desc     string
		geoms    []geom.Geometry
		expected *MinMax
	}

	point := gen.NewPoint([]float64{10, 20})
	multiPoint := gen.NewMultiPoint([][]float64{{5, 15}, {15, 25}})
	lineString := gen.NewLineString([][]float64{{0, 0}, {30, 30}})
	multiLString := gen.NewMultiLineString([][][]float64{
		{{-10, -10}, {5, 5}},
		{{25, 25}, {40, 40}},
	})
	polygon := gen.NewPolygon([][][]float64{
		{{10, 10}, {10, 30}, {30, 30}, {30, 10}, {10, 10}},
	})
	multiPolygon := gen.NewMultiPolygon([][][][]float64{
		{{{10, 10}, {10, 30}, {30, 30}, {30, 10}, {10, 10}}},
	})

	tests := tbltest.Cases(
		testcase{
			desc:     "Single point",
			geoms:    []geom.Geometry{point},
			expected: &MinMax{10, 20, 10, 20, true},
		},
		testcase{
			desc:     "MultiPoint",
			geoms:    []geom.Geometry{multiPoint},
			expected: &MinMax{5, 15, 15, 25, true},
		},
		testcase{
			desc:     "LineString",
			geoms:    []geom.Geometry{lineString},
			expected: &MinMax{0, 0, 30, 30, true},
		},
		testcase{
			desc:     "MultiLine",
			geoms:    []geom.Geometry{multiLString},
			expected: &MinMax{-10, -10, 40, 40, true},
		},
		testcase{
			desc:     "Polygon",
			geoms:    []geom.Geometry{polygon},
			expected: &MinMax{10, 10, 30, 30, true},
		},
		testcase{
			desc:     "MultiPolygon",
			geoms:    []geom.Geometry{multiPolygon},
			expected: &MinMax{10, 10, 30, 30, true},
		},
		testcase{
			desc:     "Multiple geometries",
			geoms:    []geom.Geometry{point, lineString, polygon},
			expected: &MinMax{0, 0, 30, 30, true},
		},
		testcase{
			desc:     "Nil geometries",
			geoms:    nil,
			expected: &MinMax{},
		},
	)

	tests.Run(func(idx int, test testcase) {
		mm := &MinMax{}
		result := mm.OfGeometry(test.geoms...)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %v (%v): OfGeometry() expected %v, got %v", idx, test.desc, test.expected, result)
		}
	})
}

func TestMinMaxString(t *testing.T) {
	type testcase struct {
		desc     string
		mm       *MinMax
		expected string
	}

	tests := tbltest.Cases(
		testcase{
			desc:     "Nil MinMax",
			mm:       nil,
			expected: "(nil)[0 0 , 0 0]",
		},
		testcase{
			desc:     "Empty MinMax",
			mm:       &MinMax{},
			expected: "[0 0 , 0 0]",
		},
		testcase{
			desc:     "Simple MinMax",
			mm:       &MinMax{10, 20, 30, 40, true},
			expected: "[10 20 , 30 40]",
		},
		testcase{
			desc:     "Negative coordinates MinMax",
			mm:       &MinMax{-5, -10, 5, 10, true},
			expected: "[-5 -10 , 5 10]",
		},
	)

	tests.Run(func(idx int, test testcase) {
		result := test.mm.String()
		if result != test.expected {
			t.Errorf("Test %v (%v): String() expected %q, got %q", idx, test.desc, test.expected, result)
		}
	})
}

func TestMinMaxIsZero(t *testing.T) {
	type testcase struct {
		desc     string
		mm       *MinMax
		expected bool
	}

	tests := tbltest.Cases(
		testcase{
			desc:     "Nil MinMax",
			mm:       nil,
			expected: true,
		},
		testcase{
			desc:     "Empty MinMax",
			mm:       &MinMax{},
			expected: true,
		},
		testcase{
			desc:     "Non-zero MinMax",
			mm:       &MinMax{10, 20, 30, 40, true},
			expected: false,
		},
		testcase{
			desc:     "Partially zero MinMax",
			mm:       &MinMax{0, 0, 30, 40, true},
			expected: false,
		},
	)

	tests.Run(func(idx int, test testcase) {
		result := test.mm.IsZero()
		if result != test.expected {
			t.Errorf("Test %v (%v): IsZero() expected %v, got %v", idx, test.desc, test.expected, result)
		}
	})
}

func TestMinMaxExpandBy(t *testing.T) {
	type testcase struct {
		desc     string
		mm       *MinMax
		n        int64
		expected *MinMax
	}

	tests := tbltest.Cases(
		testcase{
			desc:     "Expand simple MinMax",
			mm:       &MinMax{10, 20, 30, 40, true},
			n:        5,
			expected: &MinMax{5, 15, 35, 45, true},
		},
		testcase{
			desc:     "Expand with zero",
			mm:       &MinMax{10, 20, 30, 40, true},
			n:        0,
			expected: &MinMax{10, 20, 30, 40, true},
		},
		testcase{
			desc:     "Expand with negative value",
			mm:       &MinMax{10, 20, 30, 40, true},
			n:        -5,
			expected: &MinMax{15, 25, 25, 35, true},
		},
		testcase{
			desc:     "Expand empty MinMax",
			mm:       &MinMax{},
			n:        5,
			expected: &MinMax{-5, -5, 5, 5, true},
		},
	)

	tests.Run(func(idx int, test testcase) {
		result := test.mm.ExpandBy(test.n)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %v (%v): ExpandBy() expected %v, got %v", idx, test.desc, test.expected, result)
		}
	})
}

func TestMinMaxFn(t *testing.T) {
	type testcase struct {
		desc     string
		mm       *MinMax
		expected *MinMax
	}

	tests := tbltest.Cases(
		testcase{
			desc:     "Simple MinMax",
			mm:       &MinMax{10, 20, 30, 40, true},
			expected: &MinMax{10, 20, 30, 40, true},
		},
		testcase{
			desc:     "Nil MinMax",
			mm:       nil,
			expected: nil,
		},
	)

	tests.Run(func(idx int, test testcase) {
		result := test.mm.Fn()
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %v (%v): Fn() expected %v, got %v", idx, test.desc, test.expected, result)
		}
	})
}

func TestMinMaxMinMaxFn(t *testing.T) {
	type testcase struct {
		desc     string
		mm       *MinMax
		fn       func() *MinMax
		expected *MinMax
	}

	tests := tbltest.Cases(
		testcase{
			desc:     "Merge with function result",
			mm:       &MinMax{10, 20, 30, 40, true},
			fn:       func() *MinMax { return &MinMax{5, 15, 35, 45, true} },
			expected: &MinMax{5, 15, 35, 45, true},
		},
		testcase{
			desc:     "Merge with nil function result",
			mm:       &MinMax{10, 20, 30, 40, true},
			fn:       func() *MinMax { return nil },
			expected: &MinMax{10, 20, 30, 40, true},
		},
	)

	tests.Run(func(idx int, test testcase) {
		result := test.mm.MinMaxFn(test.fn)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %v (%v): MinMaxFn() expected %v, got %v", idx, test.desc, test.expected, result)
		}
	})
}
