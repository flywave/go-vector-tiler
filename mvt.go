package tile

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-mapbox/mvt"
	m "github.com/flywave/go-mapbox/tileid"
)

//go:embed empty.mvt
var mvtEmpty []byte

var (
	// ErrInvalidTile 表示无效的瓦片
	ErrInvalidTile = errors.New("无效的瓦片")
	// ErrEmptyLayers 表示空图层
	ErrEmptyLayers = errors.New("空图层")
	// ErrInvalidPath 表示无效的路径
	ErrInvalidPath = errors.New("无效的路径")
)

// MVTOptions MVT导出选项
type MVTOptions struct {
	// FileMode 文件权限
	FileMode os.FileMode
	// DirMode 目录权限
	DirMode os.FileMode
	// Proto 协议版本
	Proto mvt.ProtoType
	// UseEmptyTile 是否使用空瓦片
	UseEmptyTile bool
	// BufferSize 缓冲区大小
	BufferSize int
}

// DefaultMVTOptions 默认MVT选项
var DefaultMVTOptions = MVTOptions{
	FileMode:     0644,
	DirMode:      0755,
	Proto:        mvt.PROTO_MAPBOX,
	UseEmptyTile: true,
	BufferSize:   1024 * 16, // 16KB
}

// MVTTileExporter MVT格式导出器
type MVTTileExporter struct {
	// Options 导出选项
	Options MVTOptions
	// 互斥锁，用于并发安全
	mu sync.Mutex
}

// NewMVTTileExporter 创建新的MVT导出器
func NewMVTTileExporter() *MVTTileExporter {
	return &MVTTileExporter{
		Options: DefaultMVTOptions,
	}
}

// NewMVTTileExporterWithOptions 使用自定义选项创建新的MVT导出器
func NewMVTTileExporterWithOptions(options MVTOptions) *MVTTileExporter {
	return &MVTTileExporter{
		Options: options,
	}
}

// SaveTile 保存瓦片为MVT格式
func (s *MVTTileExporter) SaveTile(res []*Layer, tile *Tile, path string) error {
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

	// 生成MVT数据
	mvtData, err := s.GenerateMVT(res, tile)
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(path, mvtData, s.Options.FileMode)
}

// GenerateMVT 生成MVT数据
func (s *MVTTileExporter) GenerateMVT(layers []*Layer, tile *Tile) ([]byte, error) {
	if tile == nil {
		return nil, ErrInvalidTile
	}

	// 创建MVT瓦片
	tileID := m.TileID{Z: uint64(tile.Z), X: int64(tile.X), Y: int64(tile.Y)}

	// 创建MVT图层
	var mvtLayers [][]byte

	for _, layer := range layers {
		if layer == nil || len(layer.Features) == 0 {
			continue // 跳过空图层
		}

		config := mvt.Config{
			TileID: tileID,
			Name:   layer.Name,
			Proto:  s.Options.Proto,
		}

		mvtLayer := mvt.NewLayerConfig(config)

		// 添加要素到图层
		for _, feature := range layer.Features {
			if feature == nil || feature.Geometry == nil {
				continue // 跳过无效要素
			}

			// 转换几何对象为geom.Feature
			geomFeature := &geom.Feature{
				Geometry:   feature.Geometry,
				Properties: feature.Properties,
			}

			mvtLayer.AddFeature(geomFeature)
		}

		// 刷新图层并添加到瓦片
		layerBytes := mvtLayer.Flush()
		if len(layerBytes) > 0 {
			mvtLayers = append(mvtLayers, layerBytes)
		}
	}

	// 如果没有有效图层数据且配置为使用空瓦片
	if len(mvtLayers) == 0 {
		if s.Options.UseEmptyTile {
			return mvtEmpty, nil
		}
		return nil, ErrEmptyLayers
	}

	// 合并所有图层，预分配缓冲区以提高性能
	totalSize := 0
	for _, layer := range mvtLayers {
		totalSize += len(layer)
	}

	mvtData := make([]byte, 0, totalSize)
	for _, layer := range mvtLayers {
		mvtData = append(mvtData, layer...)
	}

	return mvtData, nil
}

// SaveTileToWriter 将瓦片数据写入io.Writer
func (s *MVTTileExporter) SaveTileToWriter(res []*Layer, tile *Tile, writer io.Writer) error {
	if writer == nil {
		return errors.New("无效的writer")
	}

	// 生成MVT数据
	mvtData, err := s.GenerateMVT(res, tile)
	if err != nil {
		return err
	}

	// 写入writer
	_, err = writer.Write(mvtData)
	return err
}

// Extension 返回文件扩展名
func (s *MVTTileExporter) Extension() string {
	return "mvt"
}

// RelativeTilePath 返回相对路径
func (s *MVTTileExporter) RelativeTilePath(zoom, x, y int) string {
	return filepath.Join(fmt.Sprintf("%d", zoom), fmt.Sprintf("%d", x), fmt.Sprintf("%d.%s", y, s.Extension()))
}

// AbsoluteTilePath 返回绝对路径
func (s *MVTTileExporter) AbsoluteTilePath(baseDir string, zoom, x, y int) string {
	return filepath.Join(baseDir, s.RelativeTilePath(zoom, x, y))
}

// BatchSaveTiles 批量保存瓦片
func (s *MVTTileExporter) BatchSaveTiles(layersByTile map[*Tile][]*Layer, baseDir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for tile, layers := range layersByTile {
		if tile == nil {
			continue
		}

		path := s.AbsoluteTilePath(baseDir, int(tile.Z), int(tile.X), int(tile.Y))
		if err := s.SaveTile(layers, tile, path); err != nil {
			return fmt.Errorf("保存瓦片 %d/%d/%d 失败: %w", tile.Z, tile.X, tile.Y, err)
		}
	}

	return nil
}
