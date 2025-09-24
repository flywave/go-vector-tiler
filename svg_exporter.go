package tile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	svg "github.com/ajstarks/svgo"
	svgdraw "github.com/flywave/go-vector-tiler/draw/svg"

	"github.com/flywave/go-geom"
)

// SVGOptions SVG导出选项
type SVGOptions struct {
	// FileMode 文件权限
	FileMode os.FileMode
	// DirMode 目录权限
	DirMode os.FileMode
	// Width 画布宽度
	Width int
	// Height 画布高度
	Height int
	// Grid 是否显示网格
	Grid bool
	// BufferSize 缓冲区大小
	BufferSize int
	// TileSize 瓦片大小
	TileSize int
	// UseEmptyTile 是否使用空瓦片
	UseEmptyTile bool
	// DefaultStyle 默认样式
	DefaultStyle string
	// DefaultPointStyle 默认点样式
	DefaultPointStyle string
}

// DefaultSVGOptions 默认SVG选项
var DefaultSVGOptions = SVGOptions{
	FileMode:          0644,
	DirMode:           0755,
	Width:             512,
	Height:            512,
	Grid:              false,
	BufferSize:        1024 * 16, // 16KB
	TileSize:          256,
	UseEmptyTile:      true,
	DefaultStyle:      "fill:none;stroke:blue;stroke-width:1",
	DefaultPointStyle: "red",
}

// SVGExporter SVG格式导出器
type SVGExporter struct {
	// Options 导出选项
	Options SVGOptions
	// 互斥锁，用于并发安全
	mu sync.Mutex
}

// NewSVGExporter 创建新的SVG导出器
func NewSVGExporter() *SVGExporter {
	return &SVGExporter{
		Options: DefaultSVGOptions,
	}
}

// NewSVGExporterWithOptions 使用自定义选项创建新的SVG导出器
func NewSVGExporterWithOptions(options SVGOptions) *SVGExporter {
	return &SVGExporter{
		Options: options,
	}
}

// SaveTile 保存瓦片为SVG格式
func (s *SVGExporter) SaveTile(res []*Layer, tile *Tile, path string) error {
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

	// 生成SVG数据
	svgData, err := s.GenerateSVG(res, tile)
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(path, svgData, s.Options.FileMode)
}

// GenerateSVG 生成SVG数据
func (s *SVGExporter) GenerateSVG(layers []*Layer, tile *Tile) ([]byte, error) {
	if tile == nil {
		return nil, ErrInvalidTile
	}

	// 如果没有图层且配置为使用空瓦片
	if len(layers) == 0 {
		if s.Options.UseEmptyTile {
			return s.generateEmptySVG(tile), nil
		}
		return nil, ErrEmptyLayers
	}

	// 创建缓冲区
	var buf bytes.Buffer

	// 创建SVG画布
	canvas := &svgdraw.Canvas{
		SVG: svg.New(&buf),
		Board: svgdraw.MinMax{
			MinX: 0,
			MinY: 0,
			MaxX: int64(s.Options.Width),
			MaxY: int64(s.Options.Height),
		},
		Region: svgdraw.MinMax{
			MinX: 0,
			MinY: 0,
			MaxX: int64(s.Options.TileSize),
			MaxY: int64(s.Options.TileSize),
		},
	}

	// 初始化SVG
	canvas.Start(s.Options.Width, s.Options.Height)

	// 添加瓦片信息注释
	canvas.Comment(fmt.Sprintf("Tile: Z=%d, X=%d, Y=%d", tile.Z, tile.X, tile.Y))

	// 如果启用网格，绘制网格
	if s.Options.Grid {
		canvas.DrawGrid(10, false, "stroke:lightgray;stroke-width:0.5")
		canvas.DrawGrid(100, true, "stroke:gray;stroke-width:1")
	}

	// 绘制各个图层
	for i, layer := range layers {
		if layer == nil || len(layer.Features) == 0 {
			continue // 跳过空图层
		}

		// 为每个图层创建组
		canvas.Group(fmt.Sprintf(`id="layer_%s"`, layer.Name), `style="opacity:1"`)
		canvas.Comment(fmt.Sprintf("Layer: %s, Features: %d", layer.Name, len(layer.Features)))

		// 绘制图层中的要素
		for j, feature := range layer.Features {
			if feature == nil || feature.Geometry == nil {
				continue // 跳过无效要素
			}

			featureId := fmt.Sprintf("layer_%d_feature_%d", i, j)
			style := s.getFeatureStyle(feature, s.Options.DefaultStyle)
			pointStyle := s.getFeaturePointStyle(feature, s.Options.DefaultPointStyle)

			// 根据几何类型绘制
			canvas.DrawGeometry(feature.Geometry, featureId, style, pointStyle, false)
		}

		canvas.Gend() // 结束图层组
	}

	canvas.End() // 结束SVG

	return buf.Bytes(), nil
}

// generateEmptySVG 生成空SVG
func (s *SVGExporter) generateEmptySVG(tile *Tile) []byte {
	var buf bytes.Buffer

	svgCanvas := svg.New(&buf)
	svgCanvas.Start(s.Options.Width, s.Options.Height)
	svgCanvas.Rect(0, 0, s.Options.Width, s.Options.Height, "fill:white;stroke:none")
	svgCanvas.Text(s.Options.Width/2, s.Options.Height/2,
		fmt.Sprintf("Empty Tile: %d/%d/%d", tile.Z, tile.X, tile.Y),
		"text-anchor:middle;font-family:Arial;font-size:14;fill:gray")
	svgCanvas.End()

	return buf.Bytes()
}

// getFeatureStyle 获取要素样式
func (s *SVGExporter) getFeatureStyle(feature *geom.Feature, defaultStyle string) string {
	if feature.Properties != nil {
		if style, ok := feature.Properties["style"]; ok {
			if styleStr, ok := style.(string); ok {
				return styleStr
			}
		}

		// 根据要素属性生成样式
		if color, ok := feature.Properties["color"]; ok {
			if colorStr, ok := color.(string); ok {
				return fmt.Sprintf("fill:none;stroke:%s;stroke-width:1", colorStr)
			}
		}
	}

	return defaultStyle
}

// getFeaturePointStyle 获取要素点样式
func (s *SVGExporter) getFeaturePointStyle(feature *geom.Feature, defaultStyle string) string {
	if feature.Properties != nil {
		if style, ok := feature.Properties["pointStyle"]; ok {
			if styleStr, ok := style.(string); ok {
				return styleStr
			}
		}

		if color, ok := feature.Properties["color"]; ok {
			if colorStr, ok := color.(string); ok {
				return colorStr
			}
		}
	}

	return defaultStyle
}

// SaveTileToWriter 将瓦片数据写入io.Writer
func (s *SVGExporter) SaveTileToWriter(res []*Layer, tile *Tile, writer io.Writer) error {
	if writer == nil {
		return errors.New("无效的writer")
	}

	// 生成SVG数据
	svgData, err := s.GenerateSVG(res, tile)
	if err != nil {
		return err
	}

	// 写入writer
	_, err = writer.Write(svgData)
	return err
}

// Extension 返回文件扩展名
func (s *SVGExporter) Extension() string {
	return "svg"
}

// RelativeTilePath 返回相对路径
func (s *SVGExporter) RelativeTilePath(zoom, x, y int) string {
	return filepath.Join(fmt.Sprintf("%d", zoom), fmt.Sprintf("%d", x), fmt.Sprintf("%d.%s", y, s.Extension()))
}

// AbsoluteTilePath 返回绝对路径
func (s *SVGExporter) AbsoluteTilePath(baseDir string, zoom, x, y int) string {
	return filepath.Join(baseDir, s.RelativeTilePath(zoom, x, y))
}

// BatchSaveTiles 批量保存瓦片
func (s *SVGExporter) BatchSaveTiles(layersByTile map[*Tile][]*Layer, baseDir string) error {
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
