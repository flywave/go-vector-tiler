package tile

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
)

func TestNewGeoJSONExporter(t *testing.T) {
	exporter := NewGeoJSONExporter()

	if exporter == nil {
		t.Fatal("NewGeoJSONExporter 返回了 nil")
	}

	// 检查默认选项
	if exporter.Options.FileMode != DefaultGeoJSONOptions.FileMode {
		t.Errorf("默认 FileMode 不正确，期望 %v，得到 %v",
			DefaultGeoJSONOptions.FileMode, exporter.Options.FileMode)
	}

	if exporter.Options.Indent != DefaultGeoJSONOptions.Indent {
		t.Errorf("默认 Indent 不正确，期望 %v，得到 %v",
			DefaultGeoJSONOptions.Indent, exporter.Options.Indent)
	}
}

func TestNewGeoJSONExporterWithOptions(t *testing.T) {
	customOptions := GeoJSONOptions{
		FileMode: 0600,
		DirMode:  0700,
		Indent:   false,
	}

	exporter := NewGeoJSONExporterWithOptions(customOptions)

	if exporter == nil {
		t.Fatal("NewGeoJSONExporterWithOptions 返回了 nil")
	}

	// 检查自定义选项
	if exporter.Options.FileMode != customOptions.FileMode {
		t.Errorf("自定义 FileMode 不正确，期望 %v，得到 %v",
			customOptions.FileMode, exporter.Options.FileMode)
	}

	if exporter.Options.Indent != customOptions.Indent {
		t.Errorf("自定义 Indent 不正确，期望 %v，得到 %v",
			customOptions.Indent, exporter.Options.Indent)
	}
}

func TestGeoJSONExporter_SaveTile(t *testing.T) {
	exporter := NewGeoJSONExporter()

	// 创建测试图层
	layer := &Layer{
		Name: "test",
		Features: []*geom.Feature{
			{
				Geometry: general.NewPoint([]float64{0, 0}),
				Properties: map[string]interface{}{
					"name": "test point",
				},
			},
		},
	}

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "geojson_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 测试保存瓦片
	path := filepath.Join(tempDir, "test.geojson")
	err = exporter.SaveTile([]*Layer{layer}, tile, path)

	if err != nil {
		t.Fatalf("SaveTile 返回错误: %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("SaveTile 没有创建文件")
	}

	// 测试无效瓦片
	err = exporter.SaveTile([]*Layer{layer}, nil, path)
	if err != ErrInvalidTile {
		t.Errorf("对于无效瓦片，期望错误 %v，得到 %v", ErrInvalidTile, err)
	}

	// 测试无效路径
	err = exporter.SaveTile([]*Layer{layer}, tile, "")
	if err != ErrInvalidPath {
		t.Errorf("对于无效路径，期望错误 %v，得到 %v", ErrInvalidPath, err)
	}
}

func TestGeoJSONExporter_Extension(t *testing.T) {
	exporter := NewGeoJSONExporter()
	ext := exporter.Extension()

	if ext != "geojson" {
		t.Errorf("Extension() 返回了错误的扩展名，期望 'geojson'，得到 '%s'", ext)
	}
}

func TestGeoJSONExporter_RelativeTilePath(t *testing.T) {
	exporter := NewGeoJSONExporter()
	path := exporter.RelativeTilePath(10, 512, 512)
	expected := filepath.Join("10", "512", "512.geojson")

	if path != expected {
		t.Errorf("RelativeTilePath 返回了错误的路径，期望 '%s'，得到 '%s'", expected, path)
	}
}

func TestGeoJSONExporter_SaveTileToWriter(t *testing.T) {
	exporter := NewGeoJSONExporter()

	// 创建测试图层
	layer := &Layer{
		Name: "test",
		Features: []*geom.Feature{
			{
				Geometry: general.NewPoint([]float64{0, 0}),
				Properties: map[string]interface{}{
					"name": "test point",
				},
			},
		},
	}

	// 创建测试瓦片
	tile := NewTile(10, 512, 512)

	// 创建缓冲区作为writer
	var buf bytes.Buffer

	// 测试保存到writer
	err := exporter.SaveTileToWriter([]*Layer{layer}, tile, &buf)

	if err != nil {
		t.Fatalf("SaveTileToWriter 返回错误: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("SaveTileToWriter 没有写入数据")
	}

	// 测试无效writer
	err = exporter.SaveTileToWriter([]*Layer{layer}, tile, nil)
	if err == nil {
		t.Error("对于无效writer，SaveTileToWriter 应该返回错误")
	}
}
