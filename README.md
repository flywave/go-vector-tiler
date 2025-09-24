# go-vector-tiler

[![Go Report Card](https://goreportcard.com/badge/github.com/flywave/go-vector-tiler)](https://goreportcard.com/report/github.com/flywave/go-vector-tiler)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

go-vector-tiler 是一个高性能的矢量瓦片生成库，用于将地理空间数据转换为矢量瓦片(Mapbox Vector Tiles)格式。

## 功能特性

- 支持多种矢量数据格式输入
- 高效的瓦片生成算法
- 支持多线程并发处理
- 可配置的几何简化
- 支持Web墨卡托和WGS84坐标系
- 提供GeoJSON和MVT两种输出格式
- 完善的错误处理和进度监控

## 安装

```bash
go get github.com/flywave/go-vector-tiler
```

## 快速开始

```go
package main

import (
	"github.com/flywave/go-vector-tiler/tile"
)

func main() {
	// 创建配置
	config := &tile.TilerConfig{
		Provider:           &MyDataProvider{}, // 实现Provider接口
		MinZoom:           0,
		MaxZoom:           14,
		Concurrency:       4,
		OutputDir:         "./tiles",
		Exporter:          tile.NewMVTTileExporter(),
	}

	// 创建Tiler实例
	tiler := tile.NewTiler(config)

	// 生成瓦片
	if err := tiler.Tiler(); err != nil {
		panic(err)
	}
}
```

## API文档

### TilerConfig

```go
type TilerConfig struct {
	Provider              Provider      // 数据提供者
	Progress              Progress      // 进度监控
	TileExtent            uint64        // 瓦片范围(默认32768)
	TileBuffer            uint64        // 瓦片缓冲区(默认64)
	SimplifyGeometries    bool          // 是否简化几何(默认true)
	SimplificationMaxZoom uint          // 最大简化级别(默认10)
	Concurrency           int           // 并发数(默认4)
	MinZoom               int           // 最小级别(默认0)
	MaxZoom               int           // 最大级别(默认14)
	SpecificZooms         []int         // 指定级别
	Bound                 *[4]float64   // 边界范围
	SRS                   string        // 空间参考系统
	Exporter              TileExporter  // 导出器
	OutputDir             string        // 输出目录
}
```

### TileExporter接口

```go
type TileExporter interface {
	SaveTile(res []*Layer, tile *Tile, path string) error
	Extension() string
	RelativeTilePath(zoom, x, y int) string
}
```

### 内置导出器

1. **GeoJSONTileExporter** - 导出为GeoJSON格式
2. **MVTTileExporter** - 导出为MVT(Mapbox Vector Tiles)格式

#### MVTTileExporter配置

```go
type MVTOptions struct {
	FileMode     os.FileMode  // 文件权限(默认0644)
	DirMode      os.FileMode  // 目录权限(默认0755)
	Proto        mvt.ProtoType // 协议版本(默认PROTO_MAPBOX)
	UseEmptyTile bool         // 使用空瓦片(默认true)
	BufferSize   int          // 缓冲区大小(默认16KB)
}
```

## 示例

### 自定义导出器

```go
// 创建自定义MVT导出器
options := tile.MVTOptions{
	FileMode:     0644,
	DirMode:      0755,
	Proto:        mvt.PROTO_MAPBOX,
	UseEmptyTile: true,
	BufferSize:   1024 * 32, // 32KB
}
exporter := tile.NewMVTTileExporterWithOptions(options)

// 创建配置
config := &tile.TilerConfig{
	Provider:   &MyDataProvider{},
	MinZoom:    0,
	MaxZoom:    14,
	Concurrency: 8,
	Exporter:   exporter,
	OutputDir:  "./mvt_tiles",
}
```

### 进度监控

```go
type MyProgress struct{}

func (p *MyProgress) Init(total int) {
	fmt.Printf("开始处理，总任务数: %d\n", total)
}

func (p *MyProgress) Update(current, total int) {
	fmt.Printf("进度: %d/%d\n", current, total)
}

func (p *MyProgress) Log(msg string) {
	fmt.Println(msg)
}

func (p *MyProgress) Warn(msg string, args ...interface{}) {
	fmt.Printf("警告: "+msg+"\n", args...)
}

// 使用自定义进度监控
config := &tile.TilerConfig{
	Progress: &MyProgress{},
	// 其他配置...
}
```

## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork仓库
2. 创建特性分支 (`git checkout -b feature/your-feature`)
3. 提交更改 (`git commit -am 'Add some feature'`)
4. 推送到分支 (`git push origin feature/your-feature`)
5. 创建Pull Request

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件