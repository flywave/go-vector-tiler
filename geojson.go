package tile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// GeoJSONOptions GeoJSON导出选项
type GeoJSONOptions struct {
	FileMode os.FileMode // 文件权限(默认0644)
	DirMode  os.FileMode // 目录权限(默认0755)
	Indent   bool        // 是否缩进JSON(默认true)
}

// DefaultGeoJSONOptions 默认GeoJSON选项
var DefaultGeoJSONOptions = GeoJSONOptions{
	FileMode: 0644,
	DirMode:  0755,
	Indent:   true,
}

// GeoJSONExporter GeoJSON格式导出器
type GeoJSONExporter struct {
	Options GeoJSONOptions
	mu      sync.Mutex // 互斥锁，用于并发安全
}

// NewGeoJSONExporter 创建新的GeoJSON导出器
func NewGeoJSONExporter() *GeoJSONExporter {
	return &GeoJSONExporter{
		Options: DefaultGeoJSONOptions,
	}
}

// NewGeoJSONExporterWithOptions 使用自定义选项创建GeoJSON导出器
func NewGeoJSONExporterWithOptions(options GeoJSONOptions) *GeoJSONExporter {
	return &GeoJSONExporter{
		Options: options,
	}
}

func (s *GeoJSONExporter) SaveTile(res []*Layer, tile *Tile, path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if tile == nil {
		return ErrInvalidTile
	}

	if path == "" {
		return ErrInvalidPath
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, s.Options.DirMode); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 构建GeoJSON FeatureCollection
	geoJSONData := map[string]interface{}{
		"type":     "FeatureCollection",
		"features": []interface{}{},
		"properties": map[string]interface{}{
			"zoom": tile.Z,
			"x":    tile.X,
			"y":    tile.Y,
		},
	}

	// 遍历所有图层和要素，添加到GeoJSON
	for _, layer := range res {
		for _, feature := range layer.Features {
			featureData := map[string]interface{}{
				"type":       "Feature",
				"geometry":   feature.Geometry,
				"properties": feature.Properties,
			}
			geoJSONData["features"] = append(geoJSONData["features"].([]interface{}), featureData)
		}
	}

	// 将GeoJSON数据写入文件
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if s.Options.Indent {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(geoJSONData); err != nil {
		return fmt.Errorf("编码GeoJSON失败: %w", err)
	}

	return nil
}

func (s *GeoJSONExporter) Extension() string {
	return "geojson"
}

func (s *GeoJSONExporter) RelativeTilePath(zoom, x, y int) string {
	return filepath.Join(fmt.Sprintf("%d", zoom), fmt.Sprintf("%d", x), fmt.Sprintf("%d.%s", y, s.Extension()))
}

func (s *GeoJSONExporter) SaveTileToWriter(res []*Layer, tile *Tile, writer io.Writer) error {
	if writer == nil {
		return errors.New("invalid writer")
	}

	geoJSONData := map[string]interface{}{
		"type":     "FeatureCollection",
		"features": []interface{}{},
	}

	for _, layer := range res {
		for _, feature := range layer.Features {
			featureData := map[string]interface{}{
				"type":       "Feature",
				"geometry":   feature.Geometry,
				"properties": feature.Properties,
			}
			geoJSONData["features"] = append(geoJSONData["features"].([]interface{}), featureData)
		}
	}

	encoder := json.NewEncoder(writer)
	if s.Options.Indent {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(geoJSONData)
}

// DefaultExporter 默认瓦片导出器
var DefaultExporter Exporter = NewGeoJSONExporter()
