package svg

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	geom "github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"
	"github.com/flywave/go-vector-tiler/maths"
	"github.com/gdey/tbltest"
)

// TestCanvasInit 测试Canvas初始化功能
func TestCanvasInit(t *testing.T) {
	type testcase struct {
		desc        string
		canvas      *Canvas
		w, h        int
		grid        bool
		shouldPanic bool
	}

	tests := tbltest.Cases(
		testcase{
			desc:        "正常初始化",
			canvas:      &Canvas{Board: MinMax{0, 0, 100, 100, true}},
			w:           400,
			h:           300,
			grid:        false,
			shouldPanic: false,
		},
		testcase{
			desc:        "带网格初始化",
			canvas:      &Canvas{Board: MinMax{10, 10, 90, 90, true}},
			w:           500,
			h:           400,
			grid:        true,
			shouldPanic: false,
		},
		testcase{
			desc:        "nil canvas应该panic",
			canvas:      nil,
			w:           400,
			h:           300,
			grid:        false,
			shouldPanic: true,
		},
	)

	tests.Run(func(idx int, test testcase) {
		var buf bytes.Buffer

		defer func() {
			if r := recover(); r != nil {
				if !test.shouldPanic {
					t.Errorf("Test %v (%v): 未预期的panic: %v", idx, test.desc, r)
				}
			} else if test.shouldPanic {
				t.Errorf("Test %v (%v): 期望panic但没有发生", idx, test.desc)
			}
		}()

		if test.canvas != nil {
			result := test.canvas.Init(&buf, test.w, test.h, test.grid)

			if !test.shouldPanic {
				if result == nil {
					t.Errorf("Test %v (%v): Init() 返回nil", idx, test.desc)
				}
				if result.SVG == nil {
					t.Errorf("Test %v (%v): SVG未初始化", idx, test.desc)
				}

				output := buf.String()
				if !strings.Contains(output, "<svg") {
					t.Errorf("Test %v (%v): 输出不包含SVG标签", idx, test.desc)
				}

				if test.grid && !strings.Contains(output, "grid") {
					t.Errorf("Test %v (%v): 期望包含网格但未找到", idx, test.desc)
				}
			}
		} else {
			// 这里会触发panic
			test.canvas.Init(&buf, test.w, test.h, test.grid)
		}
	})
}

// TestCanvasDrawPoint 测试点绘制功能
func TestCanvasDrawPoint(t *testing.T) {
	type testcase struct {
		desc string
		x, y int
		fill string
	}

	tests := tbltest.Cases(
		testcase{
			desc: "绘制红色点",
			x:    50,
			y:    75,
			fill: "red",
		},
		testcase{
			desc: "绘制蓝色点",
			x:    0,
			y:    0,
			fill: "blue",
		},
		testcase{
			desc: "绘制空填充点",
			x:    10,
			y:    20,
			fill: "",
		},
	)

	tests.Run(func(idx int, test testcase) {
		var buf bytes.Buffer
		canvas := &Canvas{Board: MinMax{0, 0, 100, 100, true}}
		canvas.Init(&buf, 400, 300, false)

		canvas.DrawPoint(test.x, test.y, test.fill)

		output := buf.String()

		// 检查是否包含圆形元素
		if !strings.Contains(output, "<circle") {
			t.Errorf("Test %v (%v): 输出不包含圆形元素", idx, test.desc)
		}

		// 检查坐标
		expectedCoord := fmt.Sprintf(`cx="%d" cy="%d"`, test.x, test.y)
		if !strings.Contains(output, expectedCoord) {
			t.Errorf("Test %v (%v): 输出不包含预期坐标 %s", idx, test.desc, expectedCoord)
		}

		// 检查填充颜色
		if test.fill != "" {
			expectedFill := "fill:" + test.fill
			if !strings.Contains(output, expectedFill) {
				t.Errorf("Test %v (%v): 输出不包含预期填充颜色 %s", idx, test.desc, expectedFill)
			}
		}
	})
}

// TestCanvasDrawGrid 测试网格绘制功能
func TestCanvasDrawGrid(t *testing.T) {
	type testcase struct {
		desc  string
		n     int
		label bool
		style string
		board MinMax
	}

	tests := tbltest.Cases(
		testcase{
			desc:  "基本网格",
			n:     10,
			label: false,
			style: "stroke:gray",
			board: MinMax{0, 0, 100, 100, true},
		},
		testcase{
			desc:  "带标签网格",
			n:     20,
			label: true,
			style: "stroke:black",
			board: MinMax{10, 10, 90, 90, true},
		},
		testcase{
			desc:  "大间距网格",
			n:     50,
			label: false,
			style: "stroke:red",
			board: MinMax{-50, -50, 50, 50, true},
		},
	)

	tests.Run(func(idx int, test testcase) {
		var buf bytes.Buffer
		canvas := &Canvas{Board: test.board}
		canvas.Init(&buf, 400, 300, false)

		canvas.DrawGrid(test.n, test.label, test.style)

		output := buf.String()

		// 检查是否包含线条元素
		if !strings.Contains(output, "<line") {
			t.Errorf("Test %v (%v): 输出不包含线条元素", idx, test.desc)
		}

		// 检查网格ID
		expectedID := fmt.Sprintf("board_%d", test.n)
		if !strings.Contains(output, expectedID) {
			t.Errorf("Test %v (%v): 输出不包含预期网格ID %s", idx, test.desc, expectedID)
		}

		// 检查样式
		if !strings.Contains(output, test.style) {
			t.Errorf("Test %v (%v): 输出不包含预期样式 %s", idx, test.desc, test.style)
		}

		// 如果有标签，检查文本元素
		if test.label {
			if !strings.Contains(output, "<text") {
				t.Errorf("Test %v (%v): 期望包含文本元素但未找到", idx, test.desc)
			}
		}
	})
}

// TestCanvasDrawGeometry 测试几何形状绘制功能
func TestCanvasDrawGeometry(t *testing.T) {
	type testcase struct {
		desc         string
		geometry     geom.Geometry
		id           string
		style        string
		pointStyle   string
		drawPoints   bool
		expectedType string
	}

	// 创建测试几何对象
	point := gen.NewPoint([]float64{10, 20})
	multiPoint := gen.NewMultiPoint([][]float64{{5, 15}, {15, 25}})
	lineString := gen.NewLineString([][]float64{{0, 0}, {30, 30}})
	multiLine := gen.NewMultiLineString([][][]float64{
		{{10, 10}, {20, 20}},
		{{30, 30}, {40, 40}},
	})
	polygon := gen.NewPolygon([][][]float64{
		{{10, 10}, {10, 30}, {30, 30}, {30, 10}, {10, 10}},
	})
	multiPolygon := gen.NewMultiPolygon([][][][]float64{
		{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}},
		{{{20, 20}, {20, 30}, {30, 30}, {30, 20}, {20, 20}}},
	})

	tests := tbltest.Cases(
		testcase{
			desc:         "绘制点",
			geometry:     point,
			id:           "test_point",
			style:        "fill:red",
			pointStyle:   "red",
			drawPoints:   false,
			expectedType: "point",
		},
		testcase{
			desc:         "绘制多点",
			geometry:     multiPoint,
			id:           "test_multipoint",
			style:        "fill:blue",
			pointStyle:   "blue",
			drawPoints:   false,
			expectedType: "multipoint",
		},
		testcase{
			desc:         "绘制线串",
			geometry:     lineString,
			id:           "test_line",
			style:        "stroke:green",
			pointStyle:   "green",
			drawPoints:   true,
			expectedType: "line",
		},
		testcase{
			desc:         "绘制多线",
			geometry:     multiLine,
			id:           "test_multiline",
			style:        "stroke:orange",
			pointStyle:   "orange",
			drawPoints:   false,
			expectedType: "multiline",
		},
		testcase{
			desc:         "绘制多边形",
			geometry:     polygon,
			id:           "test_polygon",
			style:        "fill:yellow;stroke:black",
			pointStyle:   "black",
			drawPoints:   false,
			expectedType: "polygon",
		},
		testcase{
			desc:         "绘制多多边形",
			geometry:     multiPolygon,
			id:           "test_multipolygon",
			style:        "fill:purple;stroke:black",
			pointStyle:   "black",
			drawPoints:   false,
			expectedType: "multipolygon",
		},
	)

	tests.Run(func(idx int, test testcase) {
		var buf bytes.Buffer
		canvas := &Canvas{Board: MinMax{-10, -10, 50, 50, true}}
		canvas.Init(&buf, 400, 300, false)

		count := canvas.DrawGeometry(test.geometry, test.id, test.style, test.pointStyle, test.drawPoints)

		output := buf.String()

		// 检查是否包含预期的几何类型ID
		expectedID := test.expectedType + "_" + test.id
		if !strings.Contains(output, expectedID) {
			t.Errorf("Test %v (%v): 输出不包含预期ID %s", idx, test.desc, expectedID)
		}

		// 对于多边形类型，检查返回的点数
		switch test.geometry.(type) {
		case geom.Polygon, geom.MultiPolygon:
			if count <= 0 {
				t.Errorf("Test %v (%v): 多边形应该返回正数点数，得到 %d", idx, test.desc, count)
			}
		}

		// 检查样式
		if test.style != "" {
			if !strings.Contains(output, test.style) {
				t.Errorf("Test %v (%v): 输出不包含预期样式 %s", idx, test.desc, test.style)
			}
		}

		// 如果绘制点，检查是否包含点元素
		if test.drawPoints {
			if !strings.Contains(output, "points") {
				t.Errorf("Test %v (%v): 期望包含点元素但未找到", idx, test.desc)
			}
		}
	})
}

// TestCanvasDrawRegion 测试区域绘制功能
func TestCanvasDrawRegion(t *testing.T) {
	type testcase struct {
		desc     string
		region   MinMax
		withGrid bool
	}

	tests := tbltest.Cases(
		testcase{
			desc:     "简单区域",
			region:   MinMax{10, 10, 50, 50, true},
			withGrid: false,
		},
		testcase{
			desc:     "带网格区域",
			region:   MinMax{0, 0, 100, 100, true},
			withGrid: true,
		},
		testcase{
			desc:     "负坐标区域",
			region:   MinMax{-20, -20, 20, 20, true},
			withGrid: false,
		},
	)

	tests.Run(func(idx int, test testcase) {
		var buf bytes.Buffer
		canvas := &Canvas{
			Board:  MinMax{-50, -50, 150, 150, true},
			Region: test.region,
		}
		canvas.Init(&buf, 400, 300, false)

		canvas.DrawRegion(test.withGrid)

		output := buf.String()

		// 检查是否包含区域ID
		if !strings.Contains(output, `id="region"`) {
			t.Errorf("Test %v (%v): 输出不包含区域ID", idx, test.desc)
		}

		// 检查是否包含矩形元素
		if !strings.Contains(output, "<rect") {
			t.Errorf("Test %v (%v): 输出不包含矩形元素", idx, test.desc)
		}

		// 检查矩形参数
		expectedRect := fmt.Sprintf(`x="%d" y="%d" width="%d" height="%d"`,
			int(test.region.MinX), int(test.region.MinY),
			int(test.region.Width()), int(test.region.Height()))
		if !strings.Contains(output, expectedRect) {
			t.Errorf("Test %v (%v): 输出不包含预期矩形参数 %s", idx, test.desc, expectedRect)
		}

		// 如果有网格，检查网格元素
		if test.withGrid {
			if !strings.Contains(output, "region_10") || !strings.Contains(output, "region_100") {
				t.Errorf("Test %v (%v): 期望包含网格但未找到", idx, test.desc)
			}
		}
	})
}

// TestCanvasComment 测试注释功能
func TestCanvasComment(t *testing.T) {
	type testcase struct {
		desc    string
		comment string
	}

	tests := tbltest.Cases(
		testcase{
			desc:    "简单注释",
			comment: "这是一个测试注释",
		},
		testcase{
			desc:    "包含特殊字符的注释",
			comment: "测试<>&\"'字符",
		},
		testcase{
			desc:    "空注释",
			comment: "",
		},
	)

	tests.Run(func(idx int, test testcase) {
		var buf bytes.Buffer
		canvas := &Canvas{Board: MinMax{0, 0, 100, 100, true}}
		canvas.Init(&buf, 400, 300, false)

		result := canvas.Comment(test.comment)

		// 检查返回值
		if result != canvas {
			t.Errorf("Test %v (%v): Comment() 应该返回self", idx, test.desc)
		}

		output := buf.String()

		// 检查是否包含注释
		if !strings.Contains(output, "<!--") || !strings.Contains(output, "-->") {
			t.Errorf("Test %v (%v): 输出不包含注释标签", idx, test.desc)
		}

		// 对于非空注释，检查内容（特殊字符应该被转义）
		if test.comment != "" {
			if !strings.Contains(output, test.comment) && !strings.Contains(output, "&#") {
				t.Errorf("Test %v (%v): 输出不包含注释内容", idx, test.desc)
			}
		}
	})
}

// TestCanvasCommentf 测试格式化注释功能
func TestCanvasCommentf(t *testing.T) {
	type testcase struct {
		desc   string
		format string
		args   []interface{}
	}

	tests := tbltest.Cases(
		testcase{
			desc:   "格式化注释",
			format: "点数: %d, 名称: %s",
			args:   []interface{}{5, "测试"},
		},
		testcase{
			desc:   "无参数格式化",
			format: "简单文本",
			args:   []interface{}{},
		},
		testcase{
			desc:   "多参数格式化",
			format: "坐标: (%d, %d), 半径: %.2f",
			args:   []interface{}{10, 20, 5.5},
		},
	)

	tests.Run(func(idx int, test testcase) {
		var buf bytes.Buffer
		canvas := &Canvas{Board: MinMax{0, 0, 100, 100, true}}
		canvas.Init(&buf, 400, 300, false)

		result := canvas.Commentf(test.format, test.args...)

		// 检查返回值
		if result != canvas {
			t.Errorf("Test %v (%v): Commentf() 应该返回self", idx, test.desc)
		}

		output := buf.String()

		// 检查是否包含注释
		if !strings.Contains(output, "<!--") || !strings.Contains(output, "-->") {
			t.Errorf("Test %v (%v): 输出不包含注释标签", idx, test.desc)
		}

		// 验证格式化结果
		expectedContent := fmt.Sprintf(test.format, test.args...)
		if !strings.Contains(output, expectedContent) && !strings.Contains(output, "&#") {
			t.Errorf("Test %v (%v): 输出不包含预期格式化内容", idx, test.desc)
		}
	})
}

// TestCanvasDrawMathPoints 测试数学点绘制功能
func TestCanvasDrawMathPoints(t *testing.T) {
	type testcase struct {
		desc   string
		points []maths.Pt
		styles []string
	}

	tests := tbltest.Cases(
		testcase{
			desc: "单个点",
			points: []maths.Pt{
				{X: 10, Y: 20},
			},
			styles: []string{"stroke:red"},
		},
		testcase{
			desc: "多个点",
			points: []maths.Pt{
				{X: 0, Y: 0},
				{X: 10, Y: 10},
				{X: 20, Y: 0},
			},
			styles: []string{"stroke:blue", "fill:none"},
		},
		testcase{
			desc:   "空点列表",
			points: []maths.Pt{},
			styles: []string{"stroke:black"},
		},
	)

	tests.Run(func(idx int, test testcase) {
		var buf bytes.Buffer
		canvas := &Canvas{Board: MinMax{-10, -10, 30, 30, true}}
		canvas.Init(&buf, 400, 300, false)

		canvas.DrawMathPoints(test.points, test.styles...)

		output := buf.String()

		if len(test.points) > 0 {
			// 检查是否包含路径元素
			if !strings.Contains(output, "<path") {
				t.Errorf("Test %v (%v): 输出不包含路径元素", idx, test.desc)
			}

			// 检查是否包含M指令（路径开始）
			if !strings.Contains(output, "M") {
				t.Errorf("Test %v (%v): 路径不包含M指令", idx, test.desc)
			}

			// 如果有多个点，检查L指令
			if len(test.points) > 1 && !strings.Contains(output, "L") {
				t.Errorf("Test %v (%v): 多点路径不包含L指令", idx, test.desc)
			}
		}

		// 检查样式
		for _, style := range test.styles {
			if !strings.Contains(output, style) {
				t.Errorf("Test %v (%v): 输出不包含预期样式 %s", idx, test.desc, style)
			}
		}
	})
}

// TestCanvasDrawMathSegments 测试数学线段绘制功能
func TestCanvasDrawMathSegments(t *testing.T) {
	type testcase struct {
		desc     string
		segments []maths.Line
		styles   []string
	}

	tests := tbltest.Cases(
		testcase{
			desc: "单条线段",
			segments: []maths.Line{
				{maths.Pt{X: 0, Y: 0}, maths.Pt{X: 10, Y: 10}},
			},
			styles: []string{"stroke:red"},
		},
		testcase{
			desc: "多条线段",
			segments: []maths.Line{
				{maths.Pt{X: 0, Y: 0}, maths.Pt{X: 10, Y: 0}},
				{maths.Pt{X: 10, Y: 0}, maths.Pt{X: 10, Y: 10}},
				{maths.Pt{X: 10, Y: 10}, maths.Pt{X: 0, Y: 10}},
			},
			styles: []string{"stroke:blue", "stroke-width:2"},
		},
		testcase{
			desc:     "空线段列表",
			segments: []maths.Line{},
			styles:   []string{"stroke:black"},
		},
	)

	tests.Run(func(idx int, test testcase) {
		var buf bytes.Buffer
		canvas := &Canvas{Board: MinMax{-5, -5, 15, 15, true}}
		canvas.Init(&buf, 400, 300, false)

		canvas.DrawMathSegments(test.segments, test.styles...)

		output := buf.String()

		if len(test.segments) > 0 {
			// 检查是否包含线条元素
			if !strings.Contains(output, "<line") {
				t.Errorf("Test %v (%v): 输出不包含线条元素", idx, test.desc)
			}

			// 检查线段数量（大致验证）
			lineCount := strings.Count(output, "<line")
			if lineCount != len(test.segments) {
				t.Errorf("Test %v (%v): 线段数量不匹配，期望 %d，得到 %d",
					idx, test.desc, len(test.segments), lineCount)
			}

			// 检查样式（只在有线段时检查）
			for _, style := range test.styles {
				if !strings.Contains(output, style) {
					t.Errorf("Test %v (%v): 输出不包含预期样式 %s", idx, test.desc, style)
				}
			}
		}
	})
}

// TestCanvasGroupFn 测试组功能
func TestCanvasGroupFn(t *testing.T) {
	type testcase struct {
		desc       string
		attributes []string
		drawAction func(*Canvas)
	}

	tests := tbltest.Cases(
		testcase{
			desc:       "简单组",
			attributes: []string{`id="test_group"`},
			drawAction: func(c *Canvas) {
				c.DrawPoint(10, 10, "red")
			},
		},
		testcase{
			desc:       "带样式组",
			attributes: []string{`id="styled_group"`, `style="opacity:0.5"`},
			drawAction: func(c *Canvas) {
				c.DrawPoint(20, 20, "blue")
				c.DrawPoint(30, 30, "green")
			},
		},
		testcase{
			desc:       "空组",
			attributes: []string{`id="empty_group"`},
			drawAction: func(c *Canvas) {
				// 不绘制任何内容
			},
		},
	)

	tests.Run(func(idx int, test testcase) {
		var buf bytes.Buffer
		canvas := &Canvas{Board: MinMax{0, 0, 50, 50, true}}
		canvas.Init(&buf, 400, 300, false)

		canvas.GroupFn(test.attributes, test.drawAction)

		output := buf.String()

		// 检查是否包含组标签
		if !strings.Contains(output, "<g") {
			t.Errorf("Test %v (%v): 输出不包含组开始标签", idx, test.desc)
		}

		if !strings.Contains(output, "</g>") {
			t.Errorf("Test %v (%v): 输出不包含组结束标签", idx, test.desc)
		}

		// 检查属性
		for _, attr := range test.attributes {
			if !strings.Contains(output, attr) {
				t.Errorf("Test %v (%v): 输出不包含预期属性 %s", idx, test.desc, attr)
			}
		}
	})
}

// BenchmarkCanvasDrawPoint 基准测试点绘制性能
func BenchmarkCanvasDrawPoint(b *testing.B) {
	var buf bytes.Buffer
	canvas := &Canvas{Board: MinMax{0, 0, 100, 100, true}}
	canvas.Init(&buf, 400, 300, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		canvas.DrawPoint(i%100, (i*2)%100, "red")
	}
}

// BenchmarkCanvasDrawGeometry 基准测试几何形状绘制性能
func BenchmarkCanvasDrawGeometry(b *testing.B) {
	var buf bytes.Buffer
	canvas := &Canvas{Board: MinMax{0, 0, 100, 100, true}}
	canvas.Init(&buf, 400, 300, false)

	// 创建测试几何对象
	geometries := []geom.Geometry{
		gen.NewPoint([]float64{10, 20}),
		gen.NewLineString([][]float64{{0, 0}, {10, 10}, {20, 0}}),
		gen.NewPolygon([][][]float64{
			{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
		}),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		geom := geometries[i%len(geometries)]
		canvas.DrawGeometry(geom, fmt.Sprintf("bench_%d", i), "stroke:black", "fill:red", false)
	}
}

// BenchmarkCanvasInit 基准测试Canvas初始化性能
func BenchmarkCanvasInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		canvas := &Canvas{Board: MinMax{0, 0, 100, 100, true}}
		canvas.Init(&buf, 400, 300, true)
	}
}
