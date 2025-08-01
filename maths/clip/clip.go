package clip

import (
	"fmt"
	"math"
	"sort"

	geom "github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"

	"github.com/flywave/go-vector-tiler/basic"
	"github.com/flywave/go-vector-tiler/maths"
	"github.com/flywave/go-vector-tiler/maths/lines"
)

// byxy 用于对点数组进行排序
type byxy [][]float64

func (b byxy) Len() int      { return len(b) }
func (b byxy) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byxy) Less(i, j int) bool {
	if b[i][0] != b[j][0] {
		return b[i][0] < b[j][0]
	}
	return b[i][1] < b[j][1]
}

// intersectPt 计算线段与裁剪框的交点
func intersectPt(clipbox *gen.Extent, ln [2][]float64) (pts [][]float64, ok bool) {
	lln := maths.NewLineWith2Float64(ln)
loop:
	for _, edge := range clipbox.Edges(nil) {
		eln := maths.NewLineWith2Float64(edge)
		if pt, ok := maths.Intersect(eln, lln); ok {
			if !eln.InBetween(pt) || !lln.InBetween(pt) {
				continue loop
			}

			for i := range pts {
				if pts[i][0] == pt.X && pts[i][1] == pt.Y {
					continue loop
				}
			}
			pts = append(pts, []float64{pt.X, pt.Y})
		}
	}
	sort.Sort(byxy(pts))
	return pts, len(pts) > 0
}

// LineString 裁剪线字符串
func LineString(linestr geom.LineString, extent *gen.Extent) (ls []basic.Line, err error) {
	line := lines.FromTLineString(linestr)
	if len(line) == 0 {
		return ls, nil
	}

	// 打印线数据
	fmt.Printf("线数据: %v\n", line)

	// 快速路径：完全包含
	allIn := true
	for _, pt := range line {
		in := extent.ContainsPoint(pt)
		fmt.Printf("点 %v 是否在裁剪范围内: %v\n", pt, in)
		if !in {
			allIn = false
			break
		}
	}
	fmt.Printf("是否完全包含: %v\n", allIn)
	if allIn {
		fmt.Printf("完全包含，返回原始线\n")
		return []basic.Line{basic.NewLineFrom2Float64(line...)}, nil
	}

	var cpts [][]float64
	lptIsIn := extent.ContainsPoint(line[0])
	fmt.Printf("第一个点 %v 是否在裁剪范围内: %v\n", line[0], lptIsIn)
	if lptIsIn {
		cpts = append(cpts, line[0])
		fmt.Printf("添加第一个点到cpts: %v\n", cpts)
	}

	for i := 1; i < len(line); i++ {
		cptIsIn := extent.ContainsPoint(line[i])
		fmt.Printf("点 %d (%v) 是否在裁剪范围内: %v\n", i, line[i], cptIsIn)
		switch {
		case !lptIsIn && cptIsIn:
			fmt.Printf("情况1: 前一点不在内，当前点在内\n")
			if ipts, ok := intersectPt(extent, [2][]float64{line[i-1], line[i]}); ok && len(ipts) > 0 {
				fmt.Printf("找到交点: %v\n", ipts)
				if len(ipts) == 1 {
					cpts = append(cpts, ipts[0])
					fmt.Printf("添加交点到cpts: %v\n", cpts)
				} else {
					isLess := gen.PointLess(line[i-1], line[i])
					isCLess := gen.PointLess(ipts[0], ipts[1])
					fmt.Printf("isLess: %v, isCLess: %v\n", isLess, isCLess)
					idx := 1
					if isLess == isCLess {
						idx = 0
					}
					cpts = append(cpts, ipts[idx])
					fmt.Printf("添加交点 %d 到cpts: %v\n", idx, cpts)
				}
			} else {
				fmt.Printf("未找到交点\n")
			}
			cpts = append(cpts, line[i])
			fmt.Printf("添加当前点到cpts: %v\n", cpts)
		case !lptIsIn && !cptIsIn:
			fmt.Printf("情况2: 前一点和当前点都不在内\n")
			if ipts, ok := intersectPt(extent, [2][]float64{line[i-1], line[i]}); ok && len(ipts) > 1 {
				fmt.Printf("找到两个交点: %v\n", ipts)
				isLess := gen.PointLess(line[i-1], line[i])
				isCLess := gen.PointLess(ipts[0], ipts[1])
				fmt.Printf("isLess: %v, isCLess: %v\n", isLess, isCLess)
				f, s := 0, 1
				if isLess != isCLess {
					f, s = 1, 0
				}
				ls = append(ls, basic.NewLineFrom2Float64(ipts[f], ipts[s]))
				fmt.Printf("添加线段: %v\n", ls[len(ls)-1].Data())
			} else {
				fmt.Printf("未找到两个交点\n")
			}
			cpts = cpts[:0]
			fmt.Printf("清空cpts: %v\n", cpts)
		case lptIsIn && cptIsIn:
			fmt.Printf("情况3: 前一点和当前点都在内\n")
			cpts = append(cpts, line[i])
			fmt.Printf("添加当前点到cpts: %v\n", cpts)
		case lptIsIn && !cptIsIn:
			fmt.Printf("情况4: 前一点在内，当前点不在内\n")
			if ipts, ok := intersectPt(extent, [2][]float64{line[i-1], line[i]}); ok {
				fmt.Printf("找到交点: %v\n", ipts)
				lpt := cpts[len(cpts)-1]
				for _, ipt := range ipts {
					if !pointsEqual(ipt, lpt) {
						cpts = append(cpts, ipt)
						fmt.Printf("添加交点到cpts: %v\n", cpts)
					}
				}
			} else {
				fmt.Printf("未找到交点\n")
			}
			if len(cpts) > 1 {
				ls = append(ls, basic.NewLineFrom2Float64(cpts...))
				fmt.Printf("添加线段: %v\n", ls[len(ls)-1].Data())
			}
			cpts = cpts[:0]
			fmt.Printf("清空cpts: %v\n", cpts)
		}
		lptIsIn = cptIsIn
		fmt.Printf("更新lptIsIn: %v\n", lptIsIn)
	}
	if len(cpts) > 1 {
		ls = append(ls, basic.NewLineFrom2Float64(cpts...))
		fmt.Printf("添加最后线段: %v\n", ls[len(ls)-1].Data())
	}
	fmt.Printf("返回的线段数量: %d\n", len(ls))
	return ls, nil
}

// pointsEqual 比较点是否相等（带容差）
func pointsEqual(a, b []float64) bool {
	const tolerance = 1e-6
	if len(a) != 2 || len(b) != 2 {
		return false
	}
	return math.Abs(a[0]-b[0]) < tolerance && math.Abs(a[1]-b[1]) < tolerance
}

// Sutherland-Hodgman 算法实现
func clipLeft(ring [][]float64, xmin float64) [][]float64 {
	if len(ring) == 0 {
		return nil
	}
	var output [][]float64
	s := ring[len(ring)-1]
	for i := 0; i < len(ring); i++ {
		e := ring[i]
		if e[0] >= xmin {
			if s[0] < xmin {
				dx := e[0] - s[0]
				t := (xmin - s[0]) / dx
				y := s[1] + t*(e[1]-s[1])
				output = append(output, []float64{xmin, y})
			}
			output = append(output, e)
		} else if s[0] >= xmin {
			dx := e[0] - s[0]
			t := (xmin - s[0]) / dx
			y := s[1] + t*(e[1]-s[1])
			output = append(output, []float64{xmin, y})
		}
		s = e
	}
	return output
}

func clipRight(ring [][]float64, xmax float64) [][]float64 {
	if len(ring) == 0 {
		return nil
	}
	var output [][]float64
	s := ring[len(ring)-1]
	for i := 0; i < len(ring); i++ {
		e := ring[i]
		if e[0] <= xmax {
			if s[0] > xmax {
				dx := e[0] - s[0]
				t := (xmax - s[0]) / dx
				y := s[1] + t*(e[1]-s[1])
				output = append(output, []float64{xmax, y})
			}
			output = append(output, e)
		} else if s[0] <= xmax {
			dx := e[0] - s[0]
			t := (xmax - s[0]) / dx
			y := s[1] + t*(e[1]-s[1])
			output = append(output, []float64{xmax, y})
		}
		s = e
	}
	return output
}

func clipBottom(ring [][]float64, ymin float64) [][]float64 {
	if len(ring) == 0 {
		return nil
	}
	var output [][]float64
	s := ring[len(ring)-1]
	for i := 0; i < len(ring); i++ {
		e := ring[i]
		if e[1] >= ymin {
			if s[1] < ymin {
				dy := e[1] - s[1]
				t := (ymin - s[1]) / dy
				x := s[0] + t*(e[0]-s[0])
				output = append(output, []float64{x, ymin})
			}
			output = append(output, e)
		} else if s[1] >= ymin {
			dy := e[1] - s[1]
			t := (ymin - s[1]) / dy
			x := s[0] + t*(e[0]-s[0])
			output = append(output, []float64{x, ymin})
		}
		s = e
	}
	return output
}

func clipTop(ring [][]float64, ymax float64) [][]float64 {
	if len(ring) == 0 {
		return nil
	}
	var output [][]float64
	s := ring[len(ring)-1]
	for i := 0; i < len(ring); i++ {
		e := ring[i]
		if e[1] <= ymax {
			if s[1] > ymax {
				dy := e[1] - s[1]
				t := (ymax - s[1]) / dy
				x := s[0] + t*(e[0]-s[0])
				output = append(output, []float64{x, ymax})
			}
			output = append(output, e)
		} else if s[1] <= ymax {
			dy := e[1] - s[1]
			t := (ymax - s[1]) / dy
			x := s[0] + t*(e[0]-s[0])
			output = append(output, []float64{x, ymax})
		}
		s = e
	}
	return output
}

func ensureClosed(ring [][]float64) [][]float64 {
	if len(ring) == 0 {
		return ring
	}
	if !pointsEqual(ring[0], ring[len(ring)-1]) {
		return append(ring, ring[0])
	}
	return ring
}

// Polygon 裁剪多边形（使用Sutherland-Hodgman算法）
func Polygon(poly geom.Polygon, clipExtent *gen.Extent) ([]geom.Polygon, error) {
	if len(poly.Data()) == 0 {
		return nil, nil
	}

	polyExtent, err := polygonToExtent(poly)
	if err != nil {
		return nil, err
	}

	if clipExtent.Contains(polyExtent) {
		return []geom.Polygon{poly}, nil
	}

	if _, intersects := clipExtent.Intersect(polyExtent); !intersects {
		return nil, nil
	}

	var newPolygons []geom.Polygon
	for _, ring := range poly.Data() {
		if len(ring) < 3 {
			continue
		}

		currentRing := ensureClosed(ring)
		for _, clip := range []struct {
			fn  func([][]float64, float64) [][]float64
			val float64
		}{
			{clipLeft, clipExtent[0]},
			{clipRight, clipExtent[2]},
			{clipBottom, clipExtent[1]},
			{clipTop, clipExtent[3]},
		} {
			currentRing = clip.fn(currentRing, clip.val)
			if len(currentRing) == 0 {
				break
			}
			currentRing = ensureClosed(currentRing)
			if len(currentRing) < 4 {
				currentRing = nil
				break
			}
		}

		if len(currentRing) >= 4 {
			// 规范化环：以最小坐标点为起点，并确保顺时针顺序
			currentRing = normalizeRing(currentRing)
			newPoly := gen.NewPolygon([][][]float64{currentRing})
			newPolygons = append(newPolygons, newPoly)
		}
	}

	if len(newPolygons) == 0 {
		return []geom.Polygon{poly}, nil
	}

	return newPolygons, nil
}

// 规范化环：找到最小坐标点作为起点
func normalizeRing(ring [][]float64) [][]float64 {
	if len(ring) < 4 {
		return ring
	}

	// 找到最小坐标点
	minIndex := 0
	for i := 1; i < len(ring)-1; i++ { // 不考虑最后一个点（与第一个点相同）
		if ring[i][0] < ring[minIndex][0] || (ring[i][0] == ring[minIndex][0] && ring[i][1] < ring[minIndex][1]) {
			minIndex = i
		}
	}

	// 重新排列环，以最小点为起点
	newRing := make([][]float64, len(ring))
	for i := 0; i < len(ring)-1; i++ {
		newRing[i] = ring[(minIndex+i)%(len(ring)-1)]
	}
	newRing[len(ring)-1] = newRing[0] // 确保环闭合

	return newRing
}

// polygonToExtent 从多边形计算范围
func polygonToExtent(poly geom.Polygon) (*gen.Extent, error) {
	if len(poly.Data()) == 0 || len(poly.Data()[0]) == 0 {
		return nil, nil
	}

	// 打印多边形数据
	fmt.Printf("多边形数据: %v\n", poly.Data())

	// 初始化最小和最大坐标
	minX := poly.Data()[0][0][0]
	minY := poly.Data()[0][0][1]
	maxX := minX
	maxY := minY

	fmt.Printf("初始minX: %v, minY: %v, maxX: %v, maxY: %v\n", minX, minY, maxX, maxY)

	// 遍历所有点找到最小和最大坐标
	for _, ring := range poly.Data() {
		for _, pt := range ring {
			fmt.Printf("当前点: %v\n", pt)
			if pt[0] < minX {
				minX = pt[0]
				fmt.Printf("更新minX: %v\n", minX)
			}
			if pt[0] > maxX {
				maxX = pt[0]
				fmt.Printf("更新maxX: %v\n", maxX)
			}
			if pt[1] < minY {
				minY = pt[1]
				fmt.Printf("更新minY: %v\n", minY)
			}
			if pt[1] > maxY {
				maxY = pt[1]
				fmt.Printf("更新maxY: %v\n", maxY)
			}
		}
	}

	fmt.Printf("最终minX: %v, minY: %v, maxX: %v, maxY: %v\n", minX, minY, maxX, maxY)

	// 创建范围
	extent := &gen.Extent{minX, minY, maxX, maxY}
	fmt.Printf("创建的范围: %v\n", extent)
	return extent, nil
}
