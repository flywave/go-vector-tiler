package clip

import (
	"sort"

	geom "github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"

	"github.com/flywave/go-vector-tiler/basic"
	"github.com/flywave/go-vector-tiler/maths"
	"github.com/flywave/go-vector-tiler/maths/lines"
)

type byxy [][]float64

func (b byxy) Len() int      { return len(b) }
func (b byxy) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byxy) Less(i, j int) bool {
	if b[i][0] != b[j][0] {
		return b[i][0] < b[j][0]
	}
	return b[i][1] < b[j][1]
}

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

func LineString(linestr geom.LineString, extent *gen.Extent) (ls []basic.Line, err error) {
	line := lines.FromTLineString(linestr)
	if len(line) == 0 {
		return ls, nil
	}

	var cpts [][]float64
	lptIsIn := extent.ContainsPoint(line[0])
	if lptIsIn {
		cpts = append(cpts, line[0])
	}

	for i := 1; i < len(line); i++ {
		cptIsIn := extent.ContainsPoint(line[i])
		switch {
		case !lptIsIn && cptIsIn:
			if ipts, ok := intersectPt(extent, [2][]float64{line[i-1], line[i]}); ok && len(ipts) > 0 {
				if len(ipts) == 1 {

					cpts = append(cpts, ipts[0])

				} else {
					isLess := gen.PointLess(line[i-1], line[i])
					isCLess := gen.PointLess(ipts[0], ipts[1])
					idx := 1
					if isLess == isCLess {
						idx = 0
					}
					cpts = append(cpts, ipts[idx])
				}

			}
			cpts = append(cpts, line[i])
		case !lptIsIn && !cptIsIn:
			if ipts, ok := intersectPt(extent, [2][]float64{line[i-1], line[i]}); ok && len(ipts) > 1 {
				isLess := gen.PointLess(line[i-1], line[i])
				isCLess := gen.PointLess(ipts[0], ipts[1])
				f, s := 0, 1
				if isLess != isCLess {
					f, s = 1, 0
				}
				ls = append(ls, basic.NewLineFrom2Float64(ipts[f], ipts[s]))

			}
			cpts = cpts[:0]
		case lptIsIn && cptIsIn:
			cpts = append(cpts, line[i])
		case lptIsIn && !cptIsIn:
			if ipts, ok := intersectPt(extent, [2][]float64{line[i-1], line[i]}); ok {
				_ = ipts
				lpt := cpts[len(cpts)-1]
				for _, ipt := range ipts {
					if ipt[0] != lpt[0] || ipt[1] != lpt[1] {
						cpts = append(cpts, ipt)
					}
				}
			}
			ls = append(ls, basic.NewLineFrom2Float64(cpts...))
			cpts = cpts[:0]
		}
		lptIsIn = cptIsIn
	}
	if len(cpts) > 0 {
		ls = append(ls, basic.NewLineFrom2Float64(cpts...))
	}
	return ls, nil
}
