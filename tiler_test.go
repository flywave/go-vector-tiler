package tile

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	geom "github.com/flywave/go-geom"
	"github.com/flywave/go-vector-tiler/basic"
)

// Mock实现用于测试

// MockProvider 模拟数据提供者
type MockProvider struct {
	layers []*Layer
	srid   uint64
}

func (p *MockProvider) GetDataByTile(t *Tile) []*Layer {
	return p.layers
}

func (p *MockProvider) GetSrid() uint64 {
	return p.srid
}

// MockProgress 模拟进度跟踪器
type MockProgress struct {
	Total       int
	Current     int
	Logs        []string
	Warnings    []string
	IsCompleted bool
	mu          sync.Mutex
}

func (p *MockProgress) Init(total int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Total = total
}

func (p *MockProgress) Update(current, total int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Current = current
}

func (p *MockProgress) Complete() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.IsCompleted = true
}

func (p *MockProgress) Log(format string, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Logs = append(p.Logs, fmt.Sprintf(format, args...))
}

func (p *MockProgress) Warn(format string, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Warnings = append(p.Warnings, fmt.Sprintf(format, args...))
}

// MockExporter 模拟导出器
type MockExporter struct {
	SavedTiles []TileData
	mu         sync.Mutex
}

type TileData struct {
	Layers []*Layer
	Tile   *Tile
	Path   string
}

func (e *MockExporter) SaveTile(layers []*Layer, tile *Tile, path string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.SavedTiles = append(e.SavedTiles, TileData{
		Layers: layers,
		Tile:   tile,
		Path:   path,
	})
	return nil
}

func (e *MockExporter) Extension() string {
	return ".test"
}

func (e *MockExporter) RelativeTilePath(zoom, x, y int) string {
	return fmt.Sprintf("%d/%d/%d.test", zoom, x, y)
}

func (e *MockExporter) GetSavedTiles() []TileData {
	e.mu.Lock()
	defer e.mu.Unlock()
	result := make([]TileData, len(e.SavedTiles))
	copy(result, e.SavedTiles)
	return result
}

// TestNewTiler 测试Tiler创建功能
func TestNewTiler(t *testing.T) {
	testCases := []struct {
		name     string
		config   *Config
		expected func(*Tiler) bool
	}{
		{
			name:   "nil配置使用默认值",
			config: nil,
			expected: func(tiler *Tiler) bool {
				return tiler.config.TileExtent == DefaultConfig.TileExtent &&
					tiler.config.TileBuffer == DefaultConfig.TileBuffer &&
					tiler.config.Concurrency == DefaultConfig.Concurrency
			},
		},
		{
			name: "部分配置填充默认值",
			config: &Config{
				MinZoom: 5,
				MaxZoom: 10,
			},
			expected: func(tiler *Tiler) bool {
				return tiler.config.MinZoom == 5 &&
					tiler.config.MaxZoom == 10 &&
					tiler.config.TileExtent == DefaultConfig.TileExtent &&
					tiler.config.Concurrency == DefaultConfig.Concurrency
			},
		},
		{
			name: "完整配置",
			config: &Config{
				TileExtent:  4096,
				TileBuffer:  128,
				Concurrency: 8,
				MinZoom:     0,
				MaxZoom:     15,
				SRS:         WGS84_PROJ4,
				Bound:       &[4]float64{-180, -90, 180, 90},
				OutputDir:   "./test_tiles",
			},
			expected: func(tiler *Tiler) bool {
				return tiler.config.TileExtent == 4096 &&
					tiler.config.TileBuffer == 128 &&
					tiler.config.Concurrency == 8 &&
					tiler.config.OutputDir == "./test_tiles"
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tiler := NewTiler(tc.config)

			// 基本字段检查
			if tiler.config == nil {
				t.Error("config不应该为nil")
			}
			if tiler.ctx == nil {
				t.Error("context不应该为nil")
			}
			if tiler.cancel == nil {
				t.Error("cancel函数不应该为nil")
			}
			if tiler.taskQueue == nil {
				t.Error("taskQueue不应该为nil")
			}
			if tiler.errChan == nil {
				t.Error("errChan不应该为nil")
			}
			if tiler.grid == nil {
				t.Error("grid不应该为nil")
			}
			if tiler.bbox == nil {
				t.Error("bbox不应该为nil")
			}

			// 自定义检查
			if !tc.expected(tiler) {
				t.Error("配置检查失败")
			}

			// 清理
			tiler.Stop()
		})
	}
}

// TestTiler_getZoomLevels 测试缩放级别获取
func TestTiler_getZoomLevels(t *testing.T) {
	testCases := []struct {
		name     string
		config   *Config
		expected []int
	}{
		{
			name: "使用MinZoom和MaxZoom",
			config: &Config{
				MinZoom: 5,
				MaxZoom: 8,
			},
			expected: []int{5, 6, 7, 8},
		},
		{
			name: "使用SpecificZooms",
			config: &Config{
				MinZoom:       0,
				MaxZoom:       10,
				SpecificZooms: []int{2, 5, 8},
			},
			expected: []int{2, 5, 8},
		},
		{
			name: "单一缩放级别",
			config: &Config{
				MinZoom: 3,
				MaxZoom: 3,
			},
			expected: []int{3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tiler := NewTiler(tc.config)
			defer tiler.Stop()

			result := tiler.getZoomLevels()

			if len(result) != len(tc.expected) {
				t.Errorf("getZoomLevels() 长度 = %v, want %v", len(result), len(tc.expected))
				return
			}

			for i, zoom := range result {
				if zoom != tc.expected[i] {
					t.Errorf("getZoomLevels()[%d] = %v, want %v", i, zoom, tc.expected[i])
				}
			}
		})
	}
}

// TestTiler_Count 测试瓦片数量计算
func TestTiler_Count(t *testing.T) {
	config := &Config{
		SRS:   WGS84_PROJ4,
		Bound: &[4]float64{-1, -1, 1, 1}, // 小范围测试
	}
	tiler := NewTiler(config)
	defer tiler.Stop()

	testCases := []struct {
		name  string
		zooms []uint32
	}{
		{
			name:  "单一缩放级别",
			zooms: []uint32{0},
		},
		{
			name:  "多个缩放级别",
			zooms: []uint32{0, 1, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count := tiler.Count(tc.zooms)
			if count <= 0 {
				t.Errorf("Count() = %v, want > 0", count)
			}
		})
	}
}

// TestTiler_TileBounds 测试瓦片边界计算
func TestTiler_TileBounds(t *testing.T) {
	config := &Config{
		SRS:   WGS84_PROJ4,
		Bound: &[4]float64{-180, -90, 180, 90},
	}
	tiler := NewTiler(config)
	defer tiler.Stop()

	testCases := []struct {
		name string
		zoom uint32
	}{
		{"缩放级别0", 0},
		{"缩放级别1", 1},
		{"缩放级别5", 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			minx, miny, maxx, maxy := tiler.TileBounds(tc.zoom)

			if minx > maxx {
				t.Errorf("TileBounds() minx(%v) > maxx(%v)", minx, maxx)
			}
			if miny > maxy {
				t.Errorf("TileBounds() miny(%v) > maxy(%v)", miny, maxy)
			}
		})
	}
}

// TestTiler_Iterator 测试瓦片迭代器
func TestTiler_Iterator(t *testing.T) {
	config := &Config{
		SRS:   WGS84_PROJ4,
		Bound: &[4]float64{-1, -1, 1, 1}, // 小范围以减少瓦片数量
	}
	tiler := NewTiler(config)
	defer tiler.Stop()

	testCases := []struct {
		name string
		zoom uint32
	}{
		{"缩放级别0", 0},
		{"缩放级别1", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tiles := tiler.Iterator(tc.zoom)

			if len(tiles) == 0 {
				t.Error("Iterator() 应该返回至少一个瓦片")
			}

			// 验证所有瓦片的缩放级别正确
			for _, tile := range tiles {
				if tile.Z != tc.zoom {
					t.Errorf("瓦片缩放级别 = %v, want %v", tile.Z, tc.zoom)
				}
			}

			// 验证瓦片坐标在边界内
			minx, miny, maxx, maxy := tiler.TileBounds(tc.zoom)
			for _, tile := range tiles {
				if tile.X < minx || tile.X > maxx {
					t.Errorf("瓦片X坐标 %v 超出边界 [%v, %v]", tile.X, minx, maxx)
				}
				if tile.Y < miny || tile.Y > maxy {
					t.Errorf("瓦片Y坐标 %v 超出边界 [%v, %v]", tile.Y, miny, maxy)
				}
			}
		})
	}
}

// TestTiler_Stop 测试停止功能
func TestTiler_Stop(t *testing.T) {
	tiler := NewTiler(&Config{})

	// 测试context在stop前是未取消的
	select {
	case <-tiler.ctx.Done():
		t.Error("context在Stop()前不应该被取消")
	default:
		// 正常情况
	}

	tiler.Stop()

	// 测试context在stop后被取消
	select {
	case <-tiler.ctx.Done():
		// 正常情况，context应该被取消
	case <-time.After(100 * time.Millisecond):
		t.Error("context在Stop()后应该被取消")
	}
}

// TestTiler_processTile 测试瓦片处理功能
func TestTiler_processTile(t *testing.T) {
	// 创建测试数据
	feature := &geom.Feature{
		Geometry: basic.Point{0, 0},
	}

	layer := &Layer{
		Name:     "test_layer",
		Features: []*geom.Feature{feature},
	}

	provider := &MockProvider{
		layers: []*Layer{layer},
		srid:   4326,
	}

	progress := &MockProgress{}
	exporter := &MockExporter{}

	config := &Config{
		Provider:           provider,
		Progress:           progress,
		Exporter:           exporter,
		TileExtent:         4096,
		TileBuffer:         64,
		SimplifyGeometries: false,
		OutputDir:          "./test_output",
	}

	tiler := NewTiler(config)
	defer tiler.Stop()

	task := &tileTask{z: 1, x: 0, y: 0}

	// 处理瓦片
	tiler.processTile(task)

	// 验证进度更新
	if progress.Current != 1 {
		t.Errorf("进度更新错误: current = %v, want 1", progress.Current)
	}

	// 验证导出结果
	savedTiles := exporter.GetSavedTiles()
	if len(savedTiles) != 1 {
		t.Errorf("导出瓦片数量 = %v, want 1", len(savedTiles))
	}

	if len(savedTiles) > 0 {
		savedTile := savedTiles[0]
		if savedTile.Tile.Z != 1 || savedTile.Tile.X != 0 || savedTile.Tile.Y != 0 {
			t.Errorf("导出瓦片坐标错误: z=%v, x=%v, y=%v",
				savedTile.Tile.Z, savedTile.Tile.X, savedTile.Tile.Y)
		}
	}
}

// TestTiler_exportTile 测试瓦片导出功能
func TestTiler_exportTile(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	exporter := &MockExporter{}
	config := &Config{
		Exporter:  exporter,
		OutputDir: tempDir,
	}

	tiler := NewTiler(config)
	defer tiler.Stop()

	// 创建测试数据
	feature := &geom.Feature{
		Geometry: basic.Point{0, 0},
	}

	layer := &Layer{
		Name:     "test_layer",
		Features: []*geom.Feature{feature},
	}

	tile := NewTile(1, 2, 3)

	// 测试导出
	err := tiler.exportTile([]*Layer{layer}, tile)
	if err != nil {
		t.Errorf("exportTile() 错误 = %v", err)
	}

	// 验证导出结果
	savedTiles := exporter.GetSavedTiles()
	if len(savedTiles) != 1 {
		t.Errorf("导出瓦片数量 = %v, want 1", len(savedTiles))
	}

	if len(savedTiles) > 0 {
		savedTile := savedTiles[0]
		expectedPath := filepath.Join(tempDir, "1/2/3.test")
		if savedTile.Path != expectedPath {
			t.Errorf("导出路径 = %v, want %v", savedTile.Path, expectedPath)
		}
	}
}

// TestTiler_exportTile_NoExporter 测试无导出器情况
func TestTiler_exportTile_NoExporter(t *testing.T) {
	config := &Config{
		Exporter: nil, // 无导出器
	}

	tiler := NewTiler(config)
	defer tiler.Stop()

	layer := &Layer{Name: "test"}
	tile := NewTile(0, 0, 0)

	// 无导出器时应该正常返回而不报错
	err := tiler.exportTile([]*Layer{layer}, tile)
	if err != nil {
		t.Errorf("无导出器时exportTile() 不应该报错，但得到: %v", err)
	}
}

// TestTiler_count 测试总任务数计算
func TestTiler_count(t *testing.T) {
	config := &Config{
		SRS:   WGS84_PROJ4,
		Bound: &[4]float64{-1, -1, 1, 1},
	}
	tiler := NewTiler(config)
	defer tiler.Stop()

	testCases := []struct {
		name  string
		zooms []int
	}{
		{
			name:  "单一缩放级别",
			zooms: []int{0},
		},
		{
			name:  "多个缩放级别",
			zooms: []int{0, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count := tiler.count(tc.zooms)
			if count <= 0 {
				t.Errorf("count() = %v, want > 0", count)
			}

			// 验证计算正确性：手动计算总数
			var expected int64
			for _, z := range tc.zooms {
				minx, miny, maxx, maxy := tiler.TileBounds(uint32(z))
				expected += int64((maxx - minx + 1) * (maxy - miny + 1))
			}

			if count != expected {
				t.Errorf("count() = %v, want %v", count, expected)
			}
		})
	}
}

// BenchmarkTiler_TileBounds 基准测试瓦片边界计算性能
func BenchmarkTiler_TileBounds(b *testing.B) {
	config := &Config{
		SRS:   WGS84_PROJ4,
		Bound: &[4]float64{-180, -90, 180, 90},
	}
	tiler := NewTiler(config)
	defer tiler.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tiler.TileBounds(uint32(i % 15)) // 测试不同缩放级别
	}
}

// BenchmarkTiler_Count 基准测试瓦片数量计算性能
func BenchmarkTiler_Count(b *testing.B) {
	config := &Config{
		SRS:   WGS84_PROJ4,
		Bound: &[4]float64{-10, -10, 10, 10},
	}
	tiler := NewTiler(config)
	defer tiler.Stop()

	zooms := []uint32{0, 1, 2, 3, 4}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tiler.Count(zooms)
	}
}

// BenchmarkNewTiler 基准测试Tiler创建性能
func BenchmarkNewTiler(b *testing.B) {
	config := &Config{
		TileExtent:  4096,
		TileBuffer:  64,
		Concurrency: 4,
		MinZoom:     0,
		MaxZoom:     10,
		SRS:         WGS84_PROJ4,
		Bound:       &[4]float64{-180, -90, 180, 90},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tiler := NewTiler(config)
		tiler.Stop()
	}
}
