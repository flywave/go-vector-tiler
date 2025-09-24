package tile

// Exporter 定义瓦片导出接口
type Exporter interface {
	SaveTile(res []*Layer, tile *Tile, path string) error
	Extension() string
	RelativeTilePath(zoom, x, y int) string
}
