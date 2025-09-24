package tile

// Config 配置结构体
type Config struct {
	Provider              Provider
	Progress              Progress
	TileExtent            uint64
	TileBuffer            uint64
	SimplifyGeometries    bool
	SimplificationMaxZoom uint
	Concurrency           int
	MinZoom               int
	MaxZoom               int
	SpecificZooms         []int
	Bound                 *[4]float64
	SRS                   string
	Exporter              Exporter
	OutputDir             string
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	TileExtent:            32768,
	TileBuffer:            64,
	SimplifyGeometries:    true,
	SimplificationMaxZoom: 10,
	Concurrency:           4,
	MinZoom:               0,
	MaxZoom:               14,
	SRS:                   WGS84_PROJ4,
	Bound:                 &[4]float64{-180, -90, 180, 90},
	Exporter:              nil,
	OutputDir:             "./tiles",
}
