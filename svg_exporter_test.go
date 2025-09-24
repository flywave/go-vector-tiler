package tile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"
)

// createTestSVGLayer 创建测试SVG图层
func createTestSVGLayer(name string, featureCount int) *Layer {
	layer := &Layer{
		Name:     name,
		Features: make([]*geom.Feature, 0, featureCount),
		SRID:     4326,
	}

	for i := 0; i < featureCount; i++ {
		// 创建点几何
		point := gen.NewPoint([]float64{float64(i * 10), float64(i * 10)})

		// 创建要素
		feature := &geom.Feature{
			Geometry: point,
			Properties: map[string]interface{}{
				"id":   i,
				"name": fmt.Sprintf("point_%d", i),
			},
		}

		layer.Features = append(layer.Features, feature)
	}

	return layer
}

// createTestSVGPolygonLayer 创建测试多边形图层
func createTestSVGPolygonLayer(name string) *Layer {
	layer := &Layer{
		Name:     name,
		Features: make([]*geom.Feature, 0, 1),
		SRID:     4326,
	}

	// 创建一个简单的矩形多边形
	coords := [][][]float64{
		{
			{0, 0},
			{100, 0},
			{100, 100},
			{0, 100},
			{0, 0}, // 闭合
		},
	}

	polygon := gen.NewPolygon(coords)

	feature := &geom.Feature{
		Geometry: polygon,
		Properties: map[string]interface{}{
			"id":    1,
			"name":  "test_polygon",
			"color": "blue",
		},
	}

	layer.Features = append(layer.Features, feature)
	return layer
}

// TestNewSVGExporter 测试创建SVG导出器
func TestNewSVGExporter(t *testing.T) {
	exporter := NewSVGExporter()

	if exporter == nil {
		t.Fatal("NewSVGExporter 返回了 nil")
	}

	// 检查默认选项
	if exporter.Options.Width != DefaultSVGOptions.Width {
		t.Errorf("默认宽度不正确，期望 %d，得到 %d", DefaultSVGOptions.Width, exporter.Options.Width)
	}

	if exporter.Options.Height != DefaultSVGOptions.Height {
		t.Errorf("默认高度不正确，期望 %d，得到 %d", DefaultSVGOptions.Height, exporter.Options.Height)
	}

	if exporter.Options.FileMode != DefaultSVGOptions.FileMode {
		t.Errorf("默认文件权限不正确，期望 %o，得到 %o", DefaultSVGOptions.FileMode, exporter.Options.FileMode)
	}
}

// TestNewSVGExporterWithOptions 测试使用自定义选项创建SVG导出器
func TestNewSVGExporterWithOptions(t *testing.T) {
	customOptions := SVGOptions{
		FileMode: 0755,
		DirMode:  0766,
		Width:    1024,
		Height:   768,
		Grid:     true,
	}

	exporter := NewSVGExporterWithOptions(customOptions)

	if exporter == nil {
		t.Fatal("NewSVGExporterWithOptions 返回了 nil")
	}

	if exporter.Options.Width != customOptions.Width {
		t.Errorf("自定义宽度不正确，期望 %d，得到 %d", customOptions.Width, exporter.Options.Width)
	}

	if exporter.Options.Grid != customOptions.Grid {
		t.Errorf("自定义网格选项不正确，期望 %t，得到 %t", customOptions.Grid, exporter.Options.Grid)
	}
}

// TestGenerateSVG 测试SVG生成
func TestGenerateSVG(t *testing.T) {
	exporter := NewSVGExporter()

	// 创建测试图层
	layers := []*Layer{
		createTestSVGLayer("points", 3),
		createTestSVGPolygonLayer("polygons"),
	}

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 生成SVG数据
	svgData, err := exporter.GenerateSVG(layers, tile)

	if err != nil {
		t.Fatalf("GenerateSVG 返回错误: %v", err)
	}

	if len(svgData) == 0 {
		t.Error("GenerateSVG 返回空数据")
	}

	// 检查SVG内容
	svgContent := string(svgData)
	if !strings.Contains(svgContent, "<svg") {
		t.Error("生成的数据不包含SVG标签")
	}

	if !strings.Contains(svgContent, "Tile: Z=10, X=512, Y=512") {
		t.Error("生成的SVG不包含瓦片信息")
	}

	if !strings.Contains(svgContent, "Layer: points") {
		t.Error("生成的SVG不包含图层信息")
	}
}

// TestGenerateSVGEmptyLayers 测试空图层处理
func TestGenerateSVGEmptyLayers(t *testing.T) {
	// 测试不使用空瓦片的导出器
	options := DefaultSVGOptions
	options.UseEmptyTile = false
	exporter := NewSVGExporterWithOptions(options)

	// 创建空图层
	var emptyLayers []*Layer

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 生成SVG数据
	_, err := exporter.GenerateSVG(emptyLayers, tile)

	if err != ErrEmptyLayers {
		t.Errorf("期望错误 %v，得到 %v", ErrEmptyLayers, err)
	}

	// 使用默认导出器（使用空瓦片）
	defaultExporter := NewSVGExporter()

	// 生成SVG数据
	svgData, err := defaultExporter.GenerateSVG(emptyLayers, tile)

	if err != nil {
		t.Fatalf("使用空瓦片时 GenerateSVG 返回错误: %v", err)
	}

	// 应该返回空瓦片数据
	svgContent := string(svgData)
	if !strings.Contains(svgContent, "Empty Tile") {
		t.Error("空图层应该返回空瓦片SVG")
	}
}

// TestSaveTile 测试保存瓦片
func TestSaveSVGTile(t *testing.T) {
	exporter := NewSVGExporter()

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "svg_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试图层
	layers := []*Layer{
		createTestSVGLayer("points", 3),
	}

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 保存瓦片
	path := filepath.Join(tempDir, "test.svg")
	err = exporter.SaveTile(layers, tile, path)

	if err != nil {
		t.Fatalf("SaveTile 返回错误: %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("SaveTile 没有创建文件")
	}

	// 检查文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读取保存的文件失败: %v", err)
	}

	if len(content) == 0 {
		t.Error("保存的文件内容为空")
	}

	// 检查是否包含SVG内容
	if !strings.Contains(string(content), "<svg") {
		t.Error("保存的文件不包含SVG内容")
	}

	// 测试无效的瓦片
	err = exporter.SaveTile(layers, nil, path)
	if err != ErrInvalidTile {
		t.Errorf("对于无效的瓦片，SaveTile 应该返回 ErrInvalidTile，但得到 %v", err)
	}

	// 测试无效的路径
	err = exporter.SaveTile(layers, tile, "")
	if err != ErrInvalidPath {
		t.Errorf("对于无效的路径，SaveTile 应该返回 ErrInvalidPath，但得到 %v", err)
	}
}

// TestSaveTileToWriter 测试将瓦片数据写入io.Writer
func TestSaveSVGTileToWriter(t *testing.T) {
	exporter := NewSVGExporter()

	// 创建测试图层
	layers := []*Layer{
		createTestSVGLayer("points", 3),
	}

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 创建缓冲区作为writer
	var buf bytes.Buffer

	// 保存瓦片到writer
	err := exporter.SaveTileToWriter(layers, tile, &buf)

	if err != nil {
		t.Fatalf("SaveTileToWriter 返回错误: %v", err)
	}

	// 检查写入的数据
	if buf.Len() == 0 {
		t.Error("SaveTileToWriter 没有写入数据")
	}

	// 检查内容
	content := buf.String()
	if !strings.Contains(content, "<svg") {
		t.Error("写入的数据不包含SVG内容")
	}

	// 测试无效的writer
	err = exporter.SaveTileToWriter(layers, tile, nil)
	if err == nil {
		t.Error("对于无效的writer，SaveTileToWriter 应该返回错误")
	}
}

// TestRelativeTilePath 测试相对路径生成
func TestSVGRelativeTilePath(t *testing.T) {
	exporter := NewSVGExporter()

	// 测试路径生成
	path := exporter.RelativeTilePath(10, 512, 512)
	expected := filepath.Join("10", "512", "512.svg")

	if path != expected {
		t.Errorf("RelativeTilePath 返回的路径不正确，期望 %s，得到 %s", expected, path)
	}
}

// TestAbsoluteTilePath 测试绝对路径生成
func TestSVGAbsoluteTilePath(t *testing.T) {
	exporter := NewSVGExporter()
	baseDir := "/tmp/tiles"

	// 测试绝对路径生成
	path := exporter.AbsoluteTilePath(baseDir, 10, 512, 512)
	expected := filepath.Join(baseDir, "10", "512", "512.svg")

	if path != expected {
		t.Errorf("AbsoluteTilePath 返回的路径不正确，期望 %s，得到 %s", expected, path)
	}
}

// TestExtension 测试文件扩展名
func TestSVGExtension(t *testing.T) {
	exporter := NewSVGExporter()

	extension := exporter.Extension()
	expected := "svg"

	if extension != expected {
		t.Errorf("Extension 返回的扩展名不正确，期望 %s，得到 %s", expected, extension)
	}
}

// TestBatchSaveTiles 测试批量保存瓦片
func TestSVGBatchSaveTiles(t *testing.T) {
	exporter := NewSVGExporter()

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "svg_batch_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试瓦片和图层
	tile1 := NewTile(10, 512, 512)
	tile2 := NewTile(10, 513, 512)

	layers1 := []*Layer{createTestSVGLayer("points1", 3)}
	layers2 := []*Layer{createTestSVGPolygonLayer("polygons1")}

	// 创建瓦片映射
	tileMap := map[*Tile][]*Layer{
		tile1: layers1,
		tile2: layers2,
	}

	// 批量保存瓦片
	err = exporter.BatchSaveTiles(tileMap, tempDir)

	if err != nil {
		t.Fatalf("BatchSaveTiles 返回错误: %v", err)
	}

	// 检查文件是否存在
	path1 := filepath.Join(tempDir, "10", "512", "512.svg")
	path2 := filepath.Join(tempDir, "10", "513", "512.svg")

	if _, err := os.Stat(path1); os.IsNotExist(err) {
		t.Errorf("BatchSaveTiles 没有创建文件 %s", path1)
	}

	if _, err := os.Stat(path2); os.IsNotExist(err) {
		t.Errorf("BatchSaveTiles 没有创建文件 %s", path2)
	}
}

// 模拟写入错误的Writer
type svgErrorWriter struct{}

func (w *svgErrorWriter) Write(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

// TestSVGErrorHandling 测试错误处理
func TestSVGErrorHandling(t *testing.T) {
	exporter := NewSVGExporter()

	// 创建测试图层
	layers := []*Layer{
		createTestSVGLayer("points", 3),
	}

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 测试写入错误
	err := exporter.SaveTileToWriter(layers, tile, &svgErrorWriter{})
	if err == nil {
		t.Error("对于写入错误，SaveTileToWriter 应该返回错误")
	}
}
