package tile

import (
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

// TilerConfig 配置结构体
type TilerConfig struct {
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
}

// DefaultTilerConfig 默认配置
var DefaultTilerConfig = TilerConfig{
	TileExtent:            32768,
	TileBuffer:            64,
	SimplifyGeometries:    true,
	SimplificationMaxZoom: 10,
	Concurrency:           4,
	MinZoom:               0,
	MaxZoom:               14,
	SRS:                   WGS84_PROJ4,
	Bound:                 &[4]float64{-180, -90, 180, 90},
}

// NewTiler 创建Tiler实例
func NewTiler(config *TilerConfig) *Tiler {
	// 使用默认配置填充缺失字段
	if config == nil {
		config = &DefaultTilerConfig
	} else {
		if config.TileExtent == 0 {
			config.TileExtent = DefaultTilerConfig.TileExtent
		}
		if config.TileBuffer == 0 {
			config.TileBuffer = DefaultTilerConfig.TileBuffer
		}
		if config.Concurrency <= 0 {
			config.Concurrency = DefaultTilerConfig.Concurrency
		}
		if config.SimplificationMaxZoom == 0 {
			config.SimplificationMaxZoom = DefaultTilerConfig.SimplificationMaxZoom
		}
		// 设置默认SRS和Bound
		if config.SRS == "" {
			config.SRS = WGS84_PROJ4
		}
		if config.Bound == nil {
			config.Bound = &[4]float64{-180, -90, 180, 90} // 默认全球范围
		}
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
	config     *TilerConfig
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
func (m *Tiler) Tiler(cb func(*Tile, []*Layer) error) error {
	defer m.cancel()
	defer close(m.errChan)
	defer func() {
		if m.config.Progress != nil {
			m.config.Progress.Complete()
		}
	}()

	// 计算总任务数
	zooms := m.getZoomLevels()
	var totalTasks int64
	for _, zoom := range zooms {
		ts := m.Iterator(uint32(zoom))
		totalTasks += int64(len(ts))
	}
	m.totalTasks = totalTasks

	// 初始化进度
	if m.config.Progress != nil {
		m.config.Progress.Init(int(m.totalTasks))
	}

	// 启动工作池
	m.startWorkers(cb)

	// 生成任务
	m.generateTasks(zooms)

	// 等待所有任务完成
	m.wg.Wait()

	// 检查是否有错误发生
	return m.firstError
}

// startWorkers 启动工作池
func (m *Tiler) startWorkers(cb func(*Tile, []*Layer) error) {
	for i := 0; i < m.config.Concurrency; i++ {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			m.worker(cb)
		}()
	}
}

// worker 工作函数
func (m *Tiler) worker(cb func(*Tile, []*Layer) error) {
	for {
		select {
		case <-m.ctx.Done():
			return
		case task, ok := <-m.taskQueue:
			if !ok {
				return
			}
			m.processTile(task, cb)
		}
	}
}

// generateTasks 生成任务
func (m *Tiler) generateTasks(zooms []int) {
	defer close(m.taskQueue)

	for _, zoom := range zooms {
		ts := m.Iterator(uint32(zoom))

		for _, t := range ts {
			select {
			case <-m.ctx.Done():
				return
			case m.taskQueue <- &tileTask{z: t.Z, x: t.X, y: t.Y}:
			}
		}
	}
}

// processTile 处理单个瓦片
func (m *Tiler) processTile(task *tileTask, cb func(*Tile, []*Layer) error) {
	// 更新进度
	processed := atomic.AddInt64(&m.processed, 1)
	if m.config.Progress != nil {
		m.config.Progress.Update(int(processed), int(m.totalTasks))
	}

	// 获取瓦片
	var t *Tile
	tiles := m.Iterator(task.z)
	for _, tile := range tiles {
		if tile.X == task.x && tile.Y == task.y {
			t = tile
			break
		}
	}
	if t == nil {
		m.reportError(fmt.Errorf("未找到瓦片 z=%d x=%d y=%d", task.z, task.x, task.y))
		return
	}

	// 详细日志
	if m.config.Progress != nil {
		m.config.Progress.Log(fmt.Sprintf("处理瓦片 z=%d x=%d y=%d (%d/%d)",
			task.z, task.x, task.y, processed, m.totalTasks))
	}

	// 获取数据
	ls := m.config.Provider.GetDataByTile(t)
	if len(ls) == 0 {
		return
	}

	var res []*Layer
	for _, l := range ls {
		newLayer := &Layer{Name: l.Name}
		for _, f := range l.Features {
			g := f.Geometry
			if m.config.Provider.GetSrid() != util.WebMercator {
				g, _ = basic.ToWebMercator(m.config.Provider.GetSrid(), f.Geometry)
			}
			if task.z < uint32(m.config.SimplificationMaxZoom) && m.config.SimplifyGeometries {
				g = simplify.SimplifyGeometry(g, t.ZEpislon())
			}
			g = PrepareGeo(g, t.extent, float64(m.config.TileExtent))
			pbb, _ := t.PixelBufferedBounds()
			clipRegion := gen.NewExtent([]float64{pbb[0], pbb[1]}, []float64{pbb[2], pbb[3]})
			g, _ = validate.CleanGeometry(m.ctx, g, clipRegion)
			f.Geometry = g
			newLayer.Features = append(newLayer.Features, f)
		}
		res = append(res, newLayer)
	}

	// 调用回调函数
	err := cb(t, res)
	if err != nil {
		m.reportError(fmt.Errorf("处理瓦片失败: %w", err))
		return
	}
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
