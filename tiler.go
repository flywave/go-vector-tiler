package tile

import (
	"os"
	"path/filepath"

	"context"
	"fmt"
	"sync"
	"sync/atomic"

	geo "github.com/flywave/go-geo"
	gen "github.com/flywave/go-geom/general"
	vec2d "github.com/flywave/go3d/float64/vec2"

	"github.com/flywave/go-vector-tiler/basic"
	"github.com/flywave/go-vector-tiler/maths/simplify"
	"github.com/flywave/go-vector-tiler/maths/validate"
	"github.com/flywave/go-vector-tiler/util"
)

const WGS84_PROJ4 = "+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs"
const GMERC_PROJ4 = "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0.0 +k=1.0 +units=m +nadgrids=@null +wktext +no_defs +over"

// NewTiler 创建Tiler实例
func NewTiler(config *Config) *Tiler {
	// 使用默认配置填充缺失字段
	if config == nil {
		config = &DefaultConfig
	} else {
		if config.TileExtent == 0 {
			config.TileExtent = DefaultConfig.TileExtent
		}
		if config.TileBuffer == 0 {
			config.TileBuffer = DefaultConfig.TileBuffer
		}
		if config.Concurrency <= 0 {
			config.Concurrency = DefaultConfig.Concurrency
		}
		if config.SimplificationMaxZoom == 0 {
			config.SimplificationMaxZoom = DefaultConfig.SimplificationMaxZoom
		}
		// 设置默认SRS和Bound
		if config.SRS == "" {
			config.SRS = WGS84_PROJ4
		}
		if config.Bound == nil {
			config.Bound = &[4]float64{-180, -90, 180, 90} // 默认全球范围
		}
		// 设置默认输出目录
		if config.OutputDir == "" {
			config.OutputDir = DefaultConfig.OutputDir
		}
		// Exporter默认为nil，将使用DefaultExporter
	}

	// 创建网格配置
	conf := geo.DefaultTileGridOptions()
	conf[geo.TILEGRID_BBOX_SRS] = config.SRS
	conf[geo.TILEGRID_SRS] = GMERC_PROJ4
	grid := geo.NewTileGrid(conf)
	bbx := &vec2d.Rect{Min: vec2d.T{config.Bound[0], config.Bound[1]}, Max: vec2d.T{config.Bound[2], config.Bound[3]}}

	ctx, cancel := context.WithCancel(context.Background())

	return &Tiler{
		config:     config,
		ctx:        ctx,
		cancel:     cancel,
		taskQueue:  make(chan *tileTask, config.Concurrency*2),
		wg:         sync.WaitGroup{},
		errChan:    make(chan error, config.Concurrency),
		firstError: nil,
		errOnce:    sync.Once{},
		processed:  0,
		totalTasks: 0,
		grid:       grid,
		bbox:       bbx,
	}
}

// tileTask 表示单个瓦片处理任务
type tileTask struct {
	z uint32
	x uint32
	y uint32
}

// Tiler 瓦片生成器
type Tiler struct {
	config     *Config
	ctx        context.Context
	cancel     context.CancelFunc
	taskQueue  chan *tileTask
	wg         sync.WaitGroup
	errChan    chan error
	firstError error
	errOnce    sync.Once
	processed  int64
	totalTasks int64

	// Grid相关字段
	grid *geo.TileGrid
	bbox *vec2d.Rect
}

// getZoomLevels 获取需要处理的缩放级别列表
func (m *Tiler) getZoomLevels() []int {
	if len(m.config.SpecificZooms) > 0 {
		return m.config.SpecificZooms
	}

	zooms := make([]int, 0, m.config.MaxZoom-m.config.MinZoom+1)
	for z := m.config.MinZoom; z <= m.config.MaxZoom; z++ {
		zooms = append(zooms, z)
	}
	return zooms
}

// Count 计算指定缩放级别的瓦片数量
func (m *Tiler) Count(zs []uint32) int {
	c := 0
	bbx := m.bbox
	for _, z := range zs {
		_, rc, _, _ := m.grid.GetAffectedLevelTiles(*bbx, int(z))
		c += rc[0] * rc[1]
	}
	return c
}

// TileBounds 获取指定缩放级别的瓦片边界
func (m *Tiler) TileBounds(z uint32) (uint32, uint32, uint32, uint32) {
	bbx := m.bbox
	_, _, iter, _ := m.grid.GetAffectedLevelTiles(*bbx, int(z))
	bd := iter.GetTileBound()
	return bd[0], bd[1], bd[2], bd[3]
}

// Iterator 生成指定缩放级别的瓦片迭代器
func (m *Tiler) Iterator(z uint32) []*Tile {
	ts := []*Tile{}
	minx, miny, maxx, maxy := m.TileBounds(z)
	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			ts = append(ts, NewTile(uint32(z), uint32(x), uint32(y)))
		}
	}
	return ts
}

// Stop 停止瓦片生成过程
func (m *Tiler) Stop() {
	m.cancel()
}

// Tiler 生成指定缩放级别的瓦片
func (m *Tiler) Tiler() error {
	defer m.cancel()
	defer close(m.errChan)

	// 计算总任务数
	zooms := m.getZoomLevels()
	totalTasks := m.count(zooms)
	atomic.StoreInt64(&m.totalTasks, totalTasks)

	if m.config.Progress != nil {
		m.config.Progress.Init(int(totalTasks))
	}

	// 启动工作池
	m.startWorkers()

	// 生成任务
	go m.generateTasks(zooms)

	// 等待所有任务完成
	m.wg.Wait()

	// 检查是否有错误发生
	if m.firstError != nil {
		return fmt.Errorf("瓦片生成失败: %w", m.firstError)
	}
	return nil
}

// startWorkers 启动工作池
func (m *Tiler) startWorkers() {
	for i := 0; i < m.config.Concurrency; i++ {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			m.worker()
		}()
	}
}

// worker 工作函数
func (m *Tiler) worker() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case task, ok := <-m.taskQueue:
			if !ok {
				return
			}
			m.processTile(task)
		}
	}
}

// 计算总任务数
func (m *Tiler) count(zooms []int) int64 {
	var total int64
	for _, z := range zooms {
		minx, miny, maxx, maxy := m.TileBounds(uint32(z))
		count := int64((maxx - minx + 1) * (maxy - miny + 1))
		total += count
	}
	return total
}

// 生成任务
func (m *Tiler) generateTasks(zooms []int) {
	defer close(m.taskQueue)

	for _, zoom := range zooms {
		z := uint32(zoom)
		minx, miny, maxx, maxy := m.TileBounds(z)

		for y := miny; y <= maxy; y++ {
			for x := minx; x <= maxx; x++ {
				select {
				case <-m.ctx.Done():
					return
				case m.taskQueue <- &tileTask{z: z, x: x, y: y}:
				}
			}
		}
	}
}

// 处理单个瓦片
func (m *Tiler) processTile(task *tileTask) {
	// 创建瓦片对象
	t := NewTile(task.z, task.x, task.y)

	// 更新进度
	processed := atomic.AddInt64(&m.processed, 1)
	if m.config.Progress != nil {
		m.config.Progress.Update(int(processed), int(atomic.LoadInt64(&m.totalTasks)))
		m.config.Progress.Log(fmt.Sprintf("处理瓦片 z=%d x=%d y=%d (%d/%d)",
			task.z, task.x, task.y, processed, atomic.LoadInt64(&m.totalTasks)))
	}

	// 获取数据
	layers := m.config.Provider.GetDataByTile(t)
	if len(layers) == 0 {
		return
	}

	// 处理每个图层的要素
	var resultLayers []*Layer
	for _, layer := range layers {
		newLayer := &Layer{Name: layer.Name}

		for _, feature := range layer.Features {
			geom := feature.Geometry

			// 坐标转换
			if m.config.Provider.GetSrid() != util.WebMercator {
				var err error
				if geom, err = basic.ToWebMercator(m.config.Provider.GetSrid(), geom); err != nil {
					m.reportError(fmt.Errorf("坐标转换失败 (z=%d, x=%d, y=%d): %w",
						task.z, task.x, task.y, err))
					continue
				}
			}

			// 几何简化
			if task.z < uint32(m.config.SimplificationMaxZoom) && m.config.SimplifyGeometries {
				geom = simplify.SimplifyGeometry(geom, t.ZEpislon())
			}

			// 几何预处理
			geom = PrepareGeo(geom, t.extent, float64(m.config.TileExtent))

			// 几何裁剪
			pbb, _ := t.PixelBufferedBounds()
			clipRegion := gen.NewExtent([]float64{pbb[0], pbb[1]}, []float64{pbb[2], pbb[3]})
			if cleaned, err := validate.CleanGeometry(m.ctx, geom, clipRegion); err == nil {
				geom = cleaned
			}

			feature.Geometry = geom
			newLayer.Features = append(newLayer.Features, feature)
		}

		if len(newLayer.Features) > 0 {
			resultLayers = append(resultLayers, newLayer)
		}
	}

	// 导出瓦片
	if len(resultLayers) > 0 {
		if err := m.exportTile(resultLayers, t); err != nil {
			m.reportError(fmt.Errorf("导出瓦片失败 (z=%d, x=%d, y=%d): %w",
				task.z, task.x, task.y, err))
		}
	}
}

// 导出瓦片
func (m *Tiler) exportTile(layers []*Layer, t *Tile) error {
	exporter := m.config.Exporter
	if exporter == nil {
		exporter = DefaultExporter
	}

	if exporter == nil {
		return nil // 没有导出器配置
	}

	path := exporter.RelativeTilePath(int(t.Z), int(t.X), int(t.Y))
	fullPath := filepath.Join(m.config.OutputDir, path)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	return exporter.SaveTile(layers, t, fullPath)
}

// reportError 报告错误
func (m *Tiler) reportError(err error) {
	select {
	case m.errChan <- err:
	default:
	}

	select {
	case err := <-m.errChan:
		m.errOnce.Do(func() {
			m.firstError = err
			if m.config.Progress != nil {
				m.config.Progress.Warn("瓦片生成停止: %v", err)
			}
			m.cancel()
		})
	default:
	}
}
