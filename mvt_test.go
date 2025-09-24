package tile

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-mapbox/mvt"
)

// 创建测试图层
func createTestLayer(name string, featureCount int) *Layer {
	layer := &Layer{
		Name:     name,
		Features: make([]*geom.Feature, 0, featureCount),
		SRID:     3857, // Web墨卡托投影
	}

	// 添加点要素
	for i := 0; i < featureCount; i++ {
		// 创建点要素
		feature := geom.NewPointFeature([]float64{float64(i), float64(i)})
		feature.Properties = map[string]interface{}{
			"id":    i,
			"name":  "测试点" + string(rune(i+65)), // A, B, C...
			"value": float64(i * 10),
		}

		layer.Features = append(layer.Features, feature)
	}

	return layer
}

// 创建测试多边形图层
func createTestPolygonLayer(name string) *Layer {
	layer := &Layer{
		Name:     name,
		Features: make([]*geom.Feature, 0),
		SRID:     3857,
	}

	// 创建一个简单的多边形
	coords := [][][]float64{
		{
			{0, 0},
			{0, 100},
			{100, 100},
			{100, 0},
			{0, 0},
		},
	}
	// 创建多边形要素
	feature := geom.NewPolygonFeature(coords)
	feature.Properties = map[string]interface{}{
		"id":   1,
		"name": "测试多边形",
		"area": 10000.0,
	}

	layer.Features = append(layer.Features, feature)
	return layer
}

// 创建测试线图层
func createTestLineLayer(name string) *Layer {
	layer := &Layer{
		Name:     name,
		Features: make([]*geom.Feature, 0),
		SRID:     3857,
	}

	// 创建一条简单的线
	coords := [][]float64{
		{0, 0},
		{100, 100},
		{200, 150},
	}
	// 创建线要素
	feature := geom.NewLineStringFeature(coords)
	feature.Properties = map[string]interface{}{
		"id":     1,
		"name":   "测试线",
		"length": 300.0,
	}

	layer.Features = append(layer.Features, feature)
	return layer
}

// TestNewMVTTileExporter 测试创建新的MVT导出器
func TestNewMVTTileExporter(t *testing.T) {
	exporter := NewMVTTileExporter()

	if exporter == nil {
		t.Fatal("NewMVTTileExporter 返回了 nil")
	}

	// 检查默认选项
	if exporter.Options.FileMode != DefaultMVTOptions.FileMode {
		t.Errorf("默认 FileMode 不正确，期望 %v，得到 %v",
			DefaultMVTOptions.FileMode, exporter.Options.FileMode)
	}

	if exporter.Options.Proto != mvt.PROTO_MAPBOX {
		t.Errorf("默认 Proto 不正确，期望 %v，得到 %v",
			mvt.PROTO_MAPBOX, exporter.Options.Proto)
	}
}

// TestNewMVTTileExporterWithOptions 测试使用自定义选项创建新的MVT导出器
func TestNewMVTTileExporterWithOptions(t *testing.T) {
	customOptions := MVTOptions{
		FileMode:     0600,
		DirMode:      0700,
		Proto:        mvt.PROTO_MAPBOX,
		UseEmptyTile: false,
		BufferSize:   1024 * 32,
	}

	exporter := NewMVTTileExporterWithOptions(customOptions)

	if exporter == nil {
		t.Fatal("NewMVTTileExporterWithOptions 返回了 nil")
	}

	// 检查自定义选项
	if exporter.Options.FileMode != customOptions.FileMode {
		t.Errorf("自定义 FileMode 不正确，期望 %v，得到 %v",
			customOptions.FileMode, exporter.Options.FileMode)
	}

	if exporter.Options.BufferSize != customOptions.BufferSize {
		t.Errorf("自定义 BufferSize 不正确，期望 %v，得到 %v",
			customOptions.BufferSize, exporter.Options.BufferSize)
	}
}

// TestGenerateMVT 测试生成MVT数据
func TestGenerateMVT(t *testing.T) {
	exporter := NewMVTTileExporter()

	// 创建测试图层
	layers := []*Layer{
		createTestLayer("points", 5),
		createTestPolygonLayer("polygons"),
		createTestLineLayer("lines"),
	}

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 生成MVT数据
	mvtData, err := exporter.GenerateMVT(layers, tile)

	if err != nil {
		t.Fatalf("GenerateMVT 返回错误: %v", err)
	}

	if len(mvtData) == 0 {
		t.Error("GenerateMVT 返回了空数据")
	}
}

// TestGenerateMVTEmptyLayers 测试生成空图层的MVT数据
func TestGenerateMVTEmptyLayers(t *testing.T) {
	// 创建不使用空瓦片的导出器
	customOptions := MVTOptions{
		UseEmptyTile: false,
	}
	exporter := NewMVTTileExporterWithOptions(customOptions)

	// 创建空图层
	var emptyLayers []*Layer

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 生成MVT数据
	_, err := exporter.GenerateMVT(emptyLayers, tile)

	// 应该返回空图层错误
	if err != ErrEmptyLayers {
		t.Errorf("期望错误 %v，得到 %v", ErrEmptyLayers, err)
	}

	// 使用默认导出器（使用空瓦片）
	defaultExporter := NewMVTTileExporter()

	// 生成MVT数据
	mvtData, err := defaultExporter.GenerateMVT(emptyLayers, tile)

	if err != nil {
		t.Fatalf("使用空瓦片时 GenerateMVT 返回错误: %v", err)
	}

	// 应该返回空瓦片数据
	if !bytes.Equal(mvtData, mvtEmpty) {
		t.Error("空图层应该返回空瓦片数据")
	}
}

// TestSaveTileToWriter 测试将瓦片数据写入io.Writer
func TestSaveTileToWriter(t *testing.T) {
	exporter := NewMVTTileExporter()

	// 创建测试图层
	layers := []*Layer{
		createTestLayer("points", 3),
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

	// 测试无效的writer
	err = exporter.SaveTileToWriter(layers, tile, nil)
	if err == nil {
		t.Error("对于无效的writer，SaveTileToWriter 应该返回错误")
	}
}

// TestRelativeTilePath 测试相对路径生成
func TestRelativeTilePath(t *testing.T) {
	exporter := NewMVTTileExporter()

	// 测试路径生成
	path := exporter.RelativeTilePath(10, 512, 512)
	expected := filepath.Join("10", "512", "512.mvt")

	if path != expected {
		t.Errorf("RelativeTilePath 返回的路径不正确，期望 %s，得到 %s",
			expected, path)
	}
}

// TestAbsoluteTilePath 测试绝对路径生成
func TestAbsoluteTilePath(t *testing.T) {
	exporter := NewMVTTileExporter()

	// 测试路径生成
	baseDir := "/tmp/tiles"
	path := exporter.AbsoluteTilePath(baseDir, 10, 512, 512)
	expected := filepath.Join(baseDir, "10", "512", "512.mvt")

	if path != expected {
		t.Errorf("AbsoluteTilePath 返回的路径不正确，期望 %s，得到 %s",
			expected, path)
	}
}

// TestSaveTile 测试保存瓦片到文件
func TestSaveTile(t *testing.T) {
	exporter := NewMVTTileExporter()

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "mvt_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试图层
	layers := []*Layer{
		createTestLayer("points", 3),
	}

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 保存瓦片
	path := filepath.Join(tempDir, "test.mvt")
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

// TestBatchSaveTiles 测试批量保存瓦片
func TestBatchSaveTiles(t *testing.T) {
	exporter := NewMVTTileExporter()

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "mvt_batch_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试图层和瓦片
	tile1 := NewTile(10, 512, 512)
	tile2 := NewTile(10, 513, 512)

	layers1 := []*Layer{createTestLayer("points1", 3)}
	layers2 := []*Layer{createTestPolygonLayer("polygons1")}

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
	path1 := filepath.Join(tempDir, "10", "512", "512.mvt")
	path2 := filepath.Join(tempDir, "10", "513", "512.mvt")

	if _, err := os.Stat(path1); os.IsNotExist(err) {
		t.Errorf("BatchSaveTiles 没有创建文件 %s", path1)
	}

	if _, err := os.Stat(path2); os.IsNotExist(err) {
		t.Errorf("BatchSaveTiles 没有创建文件 %s", path2)
	}
}

// 模拟文件系统错误的Writer
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	exporter := NewMVTTileExporter()

	// 创建测试图层
	layers := []*Layer{
		createTestLayer("points", 3),
	}

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 测试写入错误
	err := exporter.SaveTileToWriter(layers, tile, &errorWriter{})
	if err == nil {
		t.Error("对于写入错误，SaveTileToWriter 应该返回错误")
	}
}
