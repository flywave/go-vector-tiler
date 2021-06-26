package plyg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	svg "github.com/ajstarks/svgo"
	"github.com/flywave/go-geom/general"
	"github.com/flywave/go-vector-tiler/convert"
	"github.com/flywave/go-vector-tiler/maths"
	"github.com/flywave/go-vector-tiler/maths/hitmap"
	"github.com/flywave/go-vector-tiler/maths/points"
)

var ColLenghtErr = errors.New("Col's need to have length of at least 2")

type Ring struct {
	Points []maths.Pt
	Label  maths.Label

	hasExtent bool
	extent    *general.Extent
}

func (r *Ring) initExtent() *general.Extent {
	if r.hasExtent {
		return r.extent
	}
	pts := convert.FromMathPoint(r.Points...)
	r.extent = general.NewExtent(pts...)
	r.hasExtent = true
	return r.extent
}

func (r *Ring) Extent() [4]float64 { return r.initExtent().Extent() }

func (r *Ring) MinX() float64       { return r.initExtent().MinX() }
func (r *Ring) MinY() float64       { return r.initExtent().MinY() }
func (r *Ring) MaxX() float64       { return r.initExtent().MaxX() }
func (r *Ring) ExtentArea() float64 { return r.initExtent().Area() }
func (r *Ring) MaxY() float64       { return r.initExtent().MaxY() }

func (r Ring) LineRing() (pts []maths.Pt) {
	pts = append(pts, r.Points...)
	wo := maths.WindingOrderOfPts(pts)
	if (r.Label == maths.Inside && wo == maths.CounterClockwise) ||
		(r.Label != maths.Inside && wo == maths.Clockwise) {
		points.Reverse(pts)
	}
	points.RotateToLowestsFirst(pts)
	return pts
}

type RingDesc struct {
	Idx   int
	PtIdx int
	Label maths.Label
}

type YEdge struct {
	Y     float64
	Descs []RingDesc
}

type EdgeByY []YEdge

func (s EdgeByY) Len() int           { return len(s) }
func (s EdgeByY) Less(i, j int) bool { return s[i].Y < s[j].Y }
func (s EdgeByY) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type RingCol struct {
	Rings       []Ring
	X1, X2      float64
	Y1s         []YEdge
	Y2s         []YEdge
	foundInside bool
}

func (rc *RingCol) String() string {
	if rc == nil {
		return "Ring: nil"
	}
	s := fmt.Sprintf("Ring: x1(%v), x2(%v)\n\tRings(%v){", rc.X1, rc.X2, len(rc.Rings))
	for i, r := range rc.Rings {
		s += fmt.Sprintf("\n\t\t(%v):%v", i, r)
	}
	if len(rc.Rings) > 0 {
		s += "\n\t"
	}
	s += "}\n"
	s += fmt.Sprintf("\n\tY1s(%v){", len(rc.Y1s))
	for i, r := range rc.Y1s {
		s += fmt.Sprintf("\n\t\t(%v):%v", i, r)
	}
	if len(rc.Y1s) > 0 {
		s += "\n\t"
	}
	s += "}\n"
	s += fmt.Sprintf("\tY2s(%v){", len(rc.Y2s))
	for i, r := range rc.Y2s {
		s += fmt.Sprintf("\n\t\t(%v):%v", i, r)
	}
	if len(rc.Y2s) > 0 {
		s += "\n\t"
	}
	s += "}\n"
	return s
}

func (rc *RingCol) appendToY1(ridx int, label maths.Label, ys []YPart) {
YLoop:
	for i := 0; i < len(ys); i++ {
		for j := range rc.Y1s {
			if rc.Y1s[j].Y == ys[i].Y {
				rc.Y1s[j].Descs = append(rc.Y1s[j].Descs,
					RingDesc{Idx: ridx, PtIdx: ys[i].Idx, Label: label},
				)
				continue YLoop
			}
		}
		rc.Y1s = append(rc.Y1s, YEdge{
			Y:     ys[i].Y,
			Descs: []RingDesc{{Idx: ridx, PtIdx: ys[i].Idx, Label: label}},
		})

	}
}
func (rc *RingCol) appendToY2(ridx int, label maths.Label, ys []YPart) {
YLoop:
	for i := 0; i < len(ys); i++ {
		for j := range rc.Y2s {
			if rc.Y2s[j].Y == ys[i].Y {
				rc.Y2s[j].Descs = append(rc.Y2s[j].Descs,
					RingDesc{Idx: ridx, PtIdx: ys[i].Idx, Label: label},
				)
				continue YLoop
			}
		}
		rc.Y2s = append(rc.Y2s, YEdge{
			Y:     ys[i].Y,
			Descs: []RingDesc{{Idx: ridx, PtIdx: ys[i].Idx, Label: label}},
		})

	}
}
func hitpoint(pt1, pt2, pt3 maths.Pt) maths.Pt {
	tri := maths.Triangle{pt1, pt2, pt3}
	sort.Sort(&tri)
	if tri[0].X == tri[1].X {
		return maths.Pt{X: tri[0].X + 1, Y: float64(int64((tri[0].Y - tri[1].Y) / 2))}
	}
	return maths.Pt{X: tri[1].X - 1, Y: float64(int64((tri[0].Y - tri[1].Y) / 2))}

}
func (rc *RingCol) addPts(hm hitmap.Interface, b *Builder, pts1, pts2 []maths.Pt) {
	pts := append(append([]maths.Pt{}, pts1...), pts2...)

	tri := maths.Triangle{pts[0], pts[1], pts[2]}
	label := hm.LabelFor(tri.Center())
	if ring, x1, y1s, x2, y2s, new := b.AddPts(label, pts1, pts2); new {
		ridx := len(rc.Rings)
		rc.X1, rc.X2 = x1, x2
		rc.Rings = append(rc.Rings, ring)
		rc.appendToY1(ridx, ring.Label, y1s)
		rc.appendToY2(ridx, ring.Label, y2s)
		if !rc.foundInside {
			rc.foundInside = ring.Label == maths.Inside
		}
	}
}

func (rc *RingCol) searchY1(y float64, fn func(idx int, ptIdx int, l maths.Label) bool) {
	if rc == nil {
		return
	}
	for _, yedge := range rc.Y1s {
		if y < yedge.Y {
			return
		}
		if y != yedge.Y {
			continue
		}
		for _, desc := range yedge.Descs {
			if !fn(desc.Idx, desc.PtIdx, desc.Label) {
				return
			}
		}
		return
	}
	return
}
func (rc *RingCol) searchY2(y float64, fn func(idx int, ptIdx int, l maths.Label) bool) {
	if rc == nil {
		return
	}
	for _, yedge := range rc.Y2s {
		if y < yedge.Y {
			return
		}
		if y != yedge.Y {
			continue
		}
		for _, desc := range yedge.Descs {
			if !fn(desc.Idx, desc.PtIdx, desc.Label) {
				return
			}
		}
		return
	}
	return
}

func (rc *RingCol) searchEdge(edge []YEdge, y1, y2 float64, fn func(idx int, ptIdx int, l maths.Label) bool) {

	if rc == nil {
		return
	}
	var wantn bool
	if y1 > y2 {
		y1, y2 = y2, y1
		wantn = true
	}
	switchfn := func(desc RingDesc, nptid int) bool {
		if wantn {
			return fn(desc.Idx, nptid, desc.Label)
		}
		return fn(desc.Idx, desc.PtIdx, desc.Label)
	}
	for i := range edge {
		if y1 < edge[i].Y {
			return
		}
		if y1 != edge[i].Y {
			continue
		}
		for _, desc := range edge[i].Descs {
			px := rc.Rings[desc.Idx].Points[desc.PtIdx].X
			pptid := desc.PtIdx - 1
			if pptid < 0 {
				pptid = len(rc.Rings[desc.Idx].Points) - 1
			}
			ppt := rc.Rings[desc.Idx].Points[pptid]

			if ppt.X == px && ppt.Y == y2 {
				if !switchfn(desc, pptid) {
					return
				}
			}
			nptid := desc.PtIdx + 1
			if nptid >= len(rc.Rings[desc.Idx].Points) {
				nptid = 0
			}
			npt := rc.Rings[desc.Idx].Points[nptid]
			if npt.X == px && npt.Y == y2 {
				if !switchfn(desc, nptid) {
					return
				}
			}

		}
	}

}

func (rc *RingCol) searchY1Edge(y1, y2 float64, fn func(idx int, ptIdx int, l maths.Label) bool) {
	rc.searchEdge(rc.Y1s, y1, y2, fn)
}
func (rc *RingCol) searchY2Edge(y1, y2 float64, fn func(idx int, ptIdx int, l maths.Label) bool) {
	rc.searchEdge(rc.Y2s, y1, y2, fn)
}

type mplysByArea struct {
	pmap map[int]int
	ply  [][][]maths.Pt
}

func (mp mplysByArea) Len() int { return len(mp.ply) }
func (mp mplysByArea) Swap(i, j int) {
	li := mp.pmap[i]
	mp.pmap[i] = mp.pmap[j]
	mp.pmap[j] = li
	mp.ply[i], mp.ply[j] = mp.ply[j], mp.ply[i]
}
func (mp mplysByArea) Less(i, j int) bool {
	return points.SinArea(mp.ply[i][0]) < points.SinArea(mp.ply[j][0])
}

func (rc *RingCol) MultiPolygon() [][][]maths.Pt {
	if rc == nil || rc.Rings == nil {
		return nil
	}

	var discardPlys = make([]bool, len(rc.Rings))
	var outsidePlys []int
	var rings [][][]maths.Pt
	var miny, maxy float64

	if len(rc.Y1s) > 0 {
		miny, maxy = rc.Y1s[0].Y, rc.Y1s[0].Y
	} else if len(rc.Y2s) > 0 {
		miny, maxy = rc.Y2s[0].Y, rc.Y2s[0].Y
	}

	for _, yedge := range rc.Y1s {
		if miny > yedge.Y {
			miny = yedge.Y
		}
		if maxy < yedge.Y {
			maxy = yedge.Y
		}
		for _, desc := range yedge.Descs {
			if desc.Label == maths.Outside {
				discardPlys[desc.Idx] = true
				continue
			}
		}

	}

	for _, yedge := range rc.Y2s {
		if miny > yedge.Y {
			miny = yedge.Y
		}
		if maxy < yedge.Y {
			maxy = yedge.Y
		}
		for _, desc := range yedge.Descs {
			if desc.Label == maths.Outside {
				discardPlys[desc.Idx] = true
				continue
			}
		}
	}

	idxmap := make(map[int]int)
	segmap := make(map[int]hitmap.Segment)

	for i, ring := range rc.Rings {

		if discardPlys[i] {
			continue
		}

		if ring.Label == maths.Outside {
			e := ring.Extent()
			if e[1] == miny || e[3] == maxy {
				continue
			}
			outsidePlys = append(outsidePlys, i)
			continue
		}
		idxmap[len(rings)] = i
		lnring := ring.LineRing()
		segmap[len(rings)] = hitmap.NewSegmentFromRing(maths.Inside, ring.Points)
		rings = append(rings, [][]maths.Pt{lnring})
	}

	for _, i := range outsidePlys {

		for j := len(rings) - 1; j >= 0; j-- {
			pts := convert.FromMathPoint(rings[j][0]...)
			ibb := general.NewExtent(pts...)

			if ibb.Area() <= rc.Rings[i].ExtentArea() {
				continue
			}
			if !ibb.Contains(&(rc.Rings[i])) {
				continue
			}

			lnring := rc.Rings[i].LineRing()
			if !segmap[j].Contains(lnring[0]) {
				continue
			}
			rings[j] = append(rings[j], lnring)
			break
		}
	}
	return rings

}

type tri [4]int

func getTriangles(pt2maxy map[maths.Pt]int64, col1, col2 []maths.Pt) (tris []tri, col1idx int, col2idx int, err error) {
	clen1, clen2 := len(col1), len(col2)
	switch {
	case clen1 == 0 || clen2 == 0:
		return nil, 0, 0, ColLenghtErr
	case clen1 < 2 && clen2 < 2:
		return nil, 0, 0, ColLenghtErr
	case clen1 == 1:
		// col1      col2
		//          + 0
		//         /|
		//        / |
		//       /  |
		//      /   |
		//     /    |
		//  0 +-----+ 1
		return []tri{{0, 1, 0, 2}}, 0, 1, nil
	case clen2 == 1:
		// col1      col2
		//  0 +
		//    |\
		//    | \
		//    |  \
		//    |   \
		//    |    \
		//  1 +-----+ 0
		return []tri{{0, 2, 0, 1}}, 1, 0, nil

	}

	// try to draw a line from col2[0] to col1[1]:
	// col1      col2
	//  0 +-----+ 0
	//    |    /|
	//    |   / |
	//    |  /  |
	//    | /   |
	//    |/    |
	//  1 +-----+ 1
	maxy, ok := pt2maxy[col1[0]]
	if !ok || maxy <= int64(col2[0].Y*100) {
		// We can draw the line, so let's return the simple triangles.
		tris = append(tris, tri{0, 2, 0, 1})
		idx := 0
		// check that col2[1].Y is >= col1[1].Y
		if int64(col2[1].Y*100) <= int64(col1[1].Y*100) {
			idx = 1
			tris = append(tris, tri{1, 1, 0, 2})
		}
		return tris, 1, idx, nil
	}
	// we can not if there is a line from col1[0] headed below col2[0].Y
	// 0 +-----+ 0
	//   |\   /|
	//   | \/  |
	// 1 +  x  + 1
	//   |   \ |
	//   |    \|
	// 2 +     + 2

	idx := 1
	for ; idx <= len(col2) && int64(col2[idx].Y*100) < maxy; idx++ {
		tris = append(tris, tri{0, 1, idx - 1, 2})
	}
	tris = append(tris, tri{0, 1, idx - 1, 2}, tri{0, 2, idx, 1})
	return tris, 1, idx, nil
}

func _getTrianglesForCol(ctx context.Context, pt2maxy map[maths.Pt]int64, col1, col2 []maths.Pt) (tris []tri, err error) {
	i := 0
	for j := 0; j < len(col2); {
		if ctx.Err() != nil {
			return nil, context.Canceled
		}
		ttris, col1idx, col2idx, err := getTriangles(pt2maxy, col1[i:], col2[j:])
		if err != nil {
			return nil, err
		}
		for t := range ttris {
			tris = append(tris, tri{ttris[t][0] + i, ttris[t][1], ttris[t][2] + j, ttris[t][3]})
		}
		i, j = i+col1idx, j+col2idx
		if i == len(col1)-1 && j == len(col2)-1 {
			break
		}
	}
	return tris, nil
}

func BuildRingCol(ctx context.Context, hm hitmap.Interface, col1, col2 []maths.Pt, pt2my map[maths.Pt]int64) (col RingCol, err error) {
	var len1, len2 = len(col1), len(col2)
	_, _ = len1, len2

	var b Builder

	tris, err := _getTrianglesForCol(ctx, pt2my, col1, col2)
	if err != nil {
		return col, err
	}
	for _, t := range tris {
		col.addPts(hm, &b, col1[t[0]:t[0]+t[1]], col2[t[2]:t[2]+t[3]])
	}

	ring, x1, y1s, x2, y2s := b.CurrentRing()
	if len(ring.Points) == 0 {
		if !col.foundInside {
			col.Rings = nil
		}
		sort.Sort(EdgeByY(col.Y1s))
		sort.Sort(EdgeByY(col.Y2s))
		return col, nil
	}
	col.X1 = x1
	col.X2 = x2
	ridx := len(col.Rings)
	col.Rings = append(col.Rings, ring)
	col.appendToY1(ridx, ring.Label, y1s)
	col.appendToY2(ridx, ring.Label, y2s)
	if !col.foundInside {
		col.foundInside = ring.Label == maths.Inside
	}
	if !col.foundInside {
		col.Rings = nil
	}
	if ctx.Err() != nil {
		return col, ctx.Err()
	}
	sort.Sort(EdgeByY(col.Y1s))
	sort.Sort(EdgeByY(col.Y2s))
	return col, nil
}

func slopeCheck(pt1, pt2, pt3 maths.Pt, x1, x2 float64) bool {
	if pt1.X == x1 && pt2.X == x2 && pt3.X == x2 {
		return false
	}
	if pt1.Y == pt2.Y && pt1.Y == pt3.Y {
		return true
	}

	m1, _, d1 := maths.Line{pt1, pt2}.SlopeIntercept()
	m2, _, d2 := maths.Line{pt1, pt3}.SlopeIntercept()
	return d1 && d2 && m1 == m2
}

func merge2AdjectRC(c1, c2 RingCol) (col RingCol) {
	seenRings := make(map[[2]int]bool)
	xc := c1.X2
	cols := [2]RingCol{c1, c2}

	col.X1 = c1.X1
	col.X2 = c2.X2
	var ocoli, ccoli, ptid, nptid int

	var searchCol = func(coli int, y1, y2 float64, fn func(idx int, pidx int, l maths.Label) bool) {
		if coli == 0 {
			cols[0].searchY2Edge(y1, y2, fn)
			return
		}
		cols[1].searchY1Edge(y1, y2, fn)
	}

	var ringsToProcess [][2]int

	for i := range c1.Y2s {
		for _, d := range c1.Y2s[i].Descs {
			if _, ok := seenRings[[2]int{0, d.Idx}]; ok {
				continue
			}
			seenRings[[2]int{0, d.Idx}] = false
			ringsToProcess = append(ringsToProcess, [2]int{0, d.Idx})
		}
	}
	for i := range c1.Y1s {
		for _, d := range c1.Y1s[i].Descs {
			if _, ok := seenRings[[2]int{0, d.Idx}]; ok {
				continue
			}
			seenRings[[2]int{0, d.Idx}] = false
			col.Rings = append(col.Rings, c1.Rings[d.Idx])
		}
	}
	for i := range c1.Rings {
		if _, ok := seenRings[[2]int{0, i}]; ok {
			continue
		}
		col.Rings = append(col.Rings, c1.Rings[i])
	}

	for i := range c2.Y1s {
		for _, d := range c2.Y1s[i].Descs {
			if _, ok := seenRings[[2]int{1, d.Idx}]; ok {
				continue
			}
			seenRings[[2]int{1, d.Idx}] = false
			ringsToProcess = append(ringsToProcess, [2]int{1, d.Idx})
		}
	}
	for i := range c2.Y2s {
		for _, d := range c2.Y2s[i].Descs {
			if _, ok := seenRings[[2]int{1, d.Idx}]; ok {
				continue
			}
			seenRings[[2]int{1, d.Idx}] = false
			col.Rings = append(col.Rings, c2.Rings[d.Idx])
		}
	}
	for i := range c2.Rings {
		if _, ok := seenRings[[2]int{1, i}]; ok {
			continue
		}
		col.Rings = append(col.Rings, c2.Rings[i])
	}

	stime := time.Now()

	for p := range ringsToProcess {

		c, r := ringsToProcess[p][0], ringsToProcess[p][1]
		if seenRings[[2]int{c, r}] {
			continue
		}
		seenRings[[2]int{c, r}] = true

		var nring Ring
		nring.Label = cols[c].Rings[r].Label
		ptid = 0
		nptid = 1
		ccoli = c
		if ccoli == 1 {
			ocoli = 0
		} else {
			ocoli = 1
		}
		cri := r
		pt := cols[ccoli].Rings[cri].Points[ptid]
		npt := cols[ccoli].Rings[cri].Points[nptid]
		ptmap := make(map[maths.Pt]int)
		ptcounter := make(map[maths.Pt]int)
		walkedRings := [][2]int{{c, r}}
		for {
			etime := time.Now()
			elapsed := etime.Sub(stime)
			if elapsed.Minutes() > 10 {
				fn := genWriteoutCols(c1, c2)
				log.Println("Taking too long, writing file to ", fn)

				panic("Took too long")
			}
			if ptcounter[pt] > 5 {
				log.Println("Col1:", c1.String())
				log.Println("Col2:", c2.String())
				log.Println("On ring:", ccoli, cri)
				log.Println(cols[ccoli].Rings[cri].Points)
				pi := walkedRings[len(walkedRings)-2]
				log.Println("Previous ring:", pi[0], pi[1])
				log.Println(cols[pi[0]].Rings[pi[1]].Points)
				log.Println("Processing ", p, "(", ringsToProcess[p], ") of the following rings that needed to be processed.:", ringsToProcess)
				log.Println(cols[ringsToProcess[p][0]].Rings[ringsToProcess[p][1]].Points)
				log.Println("Walked rings:", walkedRings)
				fn := genWriteoutCols(c1, c2)
				log.Println("Wrote out columns info to:", fn)
				writeOutSVG(fn, cols[:], walkedRings)

				panic("Inif loop?")
			}

			if idx, ok := ptmap[pt]; ok {
				for _, pt1 := range nring.Points[idx:] {
					delete(ptmap, pt1)
				}
				nring.Points = nring.Points[:idx]
			}
			if len(nring.Points) > 1 && slopeCheck(nring.Points[len(nring.Points)-2], nring.Points[len(nring.Points)-1], pt, xc, xc) {
				delete(ptmap, nring.Points[len(nring.Points)-1])
				nring.Points[len(nring.Points)-1] = pt
			} else {
				nring.Points = append(nring.Points, pt)
				ptcounter[pt]++
			}
			ptmap[pt] = len(nring.Points) - 1
			if pt.X != xc || npt.X != xc {
				goto NextPoint
			}
			searchCol(ocoli, pt.Y, npt.Y, func(idx int, pidx int, l maths.Label) bool {
				if l != nring.Label {
					return true
				}

				ocri := cri
				ptid = pidx
				nptid = ptid + 1

				ccoli, ocoli = ocoli, ccoli
				cri = idx
				if nptid >= len(cols[ccoli].Rings[cri].Points) {
					nptid = 0
				}
				seenRings[[2]int{ccoli, idx}] = true
				walkedRings = append(walkedRings, [2]int{ccoli, idx})
				cols[ccoli].Rings[cri].Extent()

				pt := cols[ccoli].Rings[cri].Points[ptid]
				npt := cols[ccoli].Rings[cri].Points[nptid]
				ptcounter[pt]++

				if npt.X != pt.X {
					return false
				}

				searchCol(ocoli, pt.Y, npt.Y, func(idx int, pidx int, l maths.Label) bool {
					if l != nring.Label {
						return true
					}
					if idx == ocri {
						return true
					}
					ptid = pidx
					nptid = ptid + 1
					ccoli, ocoli = ocoli, ccoli
					cri = idx
					seenRings[[2]int{ccoli, idx}] = true
					walkedRings = append(walkedRings, [2]int{ccoli, idx})
					cols[ccoli].Rings[cri].Extent()
					if nptid >= len(cols[ccoli].Rings[cri].Points) {
						nptid = 0
					}
					return false

				})
				return false
			})

		NextPoint:
			ptid, nptid = nptid, nptid+1
			if nptid >= len(cols[ccoli].Rings[cri].Points) {
				nptid = 0
			}
			pt = cols[ccoli].Rings[cri].Points[ptid]
			npt = cols[ccoli].Rings[cri].Points[nptid]
			if pt.IsEqual(nring.Points[0]) {
				break
			}
		}
		plen := len(nring.Points)
		if plen > 3 {
			switch {
			case slopeCheck(nring.Points[plen-2], nring.Points[plen-1], nring.Points[0], col.X1, col.X2):
				nring.Points = nring.Points[:plen-1]
			case slopeCheck(nring.Points[plen-1], nring.Points[0], nring.Points[1], col.X1, col.X2):
				nring.Points = nring.Points[1:]
			}
		}
		if plen < 3 {
			fn := genWriteoutCols(c1, c2)
			log.Println("Generated a ring with fewer then 3 points: ", fn, nring)

			panic("Generated a ring with fewer then 3 points. ")
		}
		points.RotateToLowestsFirst(nring.Points)

		col.Rings = append(col.Rings, nring)

	}
	for i, r := range col.Rings {
		for j, pt := range r.Points {
			switch pt.X {
			case col.X1:
				col.appendToY1(i, r.Label, []YPart{{Y: pt.Y, Idx: j}})
			case col.X2:
				col.appendToY2(i, r.Label, []YPart{{Y: pt.Y, Idx: j}})
			}
		}
	}
	sort.Sort(EdgeByY(col.Y1s))
	sort.Sort(EdgeByY(col.Y2s))

	for i := range col.Y1s {
		cpt := maths.Pt{X: col.X1, Y: col.Y1s[i].Y}
		for j, d := range col.Y1s[i].Descs {
			ring := col.Rings[d.Idx]
			if d.Label != ring.Label {
				col.Y1s[i].Descs[j].Label = ring.Label
			}

			pt := ring.Points[d.PtIdx]
			if !cpt.IsEqual(pt) {
				var found bool
				for r := range ring.Points {
					if cpt.IsEqual(ring.Points[r]) {
						col.Y1s[i].Descs[j].PtIdx = r
						found = true
						break
					}
				}
				if !found {
					log.Println("col", col.String())
					log.Println("Did not find r when trying to fix up Y1.", i, j)
					panic("Did not find r when trying to fix up Y1.")
				}
			}

		}

	}
	for i := range col.Y2s {
		cpt := maths.Pt{X: col.X2, Y: col.Y2s[i].Y}
		for j, d := range col.Y2s[i].Descs {
			ring := col.Rings[d.Idx]
			if d.Label != ring.Label {
				col.Y2s[i].Descs[j].Label = ring.Label
			}
			pt := ring.Points[d.PtIdx]
			if !cpt.IsEqual(pt) {
				var found bool
				for r := range ring.Points {
					if cpt.IsEqual(ring.Points[r]) {
						col.Y2s[i].Descs[j].PtIdx = r
						found = true
						break
					}
				}
				if !found {
					log.Println("col", col.String())
					log.Println("Did not find r when trying to fix up Y2.", i, j)
					panic("Did not find r when trying to fix up Y2.")
				}
			}
		}
	}
	return col
}

func MergeCols(cols []RingCol) RingCol {
	lcol := cols[0]
	for i := 1; i < len(cols); i++ {
		lcol = merge2AdjectRC(lcol, cols[i])
	}
	return lcol
}

func GenerateMultiPolygon(cols []RingCol) (plys []maths.Polygon) {
	var lock sync.Mutex
	var wg sync.WaitGroup
	var wChan = make(chan [2]int)
	var numWorkers = runtime.NumCPU()

	li := -1
	var worker = func(id int) {
		for i := range wChan {
			wcol := MergeCols(cols[i[0]:i[1]])
			wply := wcol.MultiPolygon()
			lock.Lock()
			for i := range wply {
				plys = append(plys, wply[i])
			}
			lock.Unlock()
		}
		wg.Done()
	}
	for i := 0; i < numWorkers; i++ {
		go worker(i)
	}
	wg.Add(numWorkers)

	for i := range cols {
		if len(cols[i].Rings) == 0 {
			if li != -1 {
				wChan <- [2]int{li, i}
				li = -1
			}
			continue
		}
		if li == -1 {
			li = i
		}
	}
	if li != -1 {
		wChan <- [2]int{li, len(cols)}
	}
	close(wChan)
	wg.Wait()
	return plys
}

func writeOutSVG(fn string, cols []RingCol, onlyRings [][2]int) {
	var filter bool
	ringFilter := make(map[[2]int]bool)
	if len(onlyRings) > 0 {
		filter = true
		for i := range onlyRings {
			ringFilter[onlyRings[i]] = true
		}
	}
	f, err := os.Create(fn + ".svg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	canvas := svg.New(f)
	canvas.Startview(786, 1024, int(cols[0].X1)-10, 2000, int(cols[1].X2)+10, 2200)
	defer canvas.End()

	canvas.Def()
	canvas.Marker("markerCircle", 1, 1, 0, 0)
	canvas.Circle(5, 5, 1, "stroke:none;fill:#8a8a8a;fill-opacity:0.3")
	canvas.MarkerEnd()
	canvas.DefEnd()

	style := func(l maths.Label, i int) string {
		if i == 0 {

			if l == maths.Inside {
				return "fill:#0000ff;fill-opacity:0.3;stroke:none;marker-mid: url(#markerCircle)"
			}
			return "fill:#ff0000;fill-opacity:0.3;stroke:none; marker-mid: url(#markerCircle)"
		}
		if l == maths.Inside {
			return "fill:#0088ff;fill-opacity:0.3;stroke:none; marker-mid: url(#markerCircle)"
		}
		return "fill:#ff8800;fill-opacity:0.3;stroke:none; marker-mid: url(#markerCircle)"

	}

	pointsToIntArray := func(pts []maths.Pt) (xs []int, ys []int) {
		for _, pt := range pts {
			xs = append(xs, int(pt.X))
			ys = append(ys, int(pt.Y))
		}
		return xs, ys
	}
	pointmap := make(map[maths.Pt]struct{})

	for i, col := range cols {
		for j, r := range col.Rings {
			if filter && ringFilter[[2]int{i, j}] {
				continue
			}
			for _, pt := range r.Points {
				pointmap[pt] = struct{}{}
			}
		}
	}

	canvas.Scale(1.5)

	canvas.Line(int(cols[0].X1), -20, int(cols[0].X1), 4126, "stroke:#8a8a8a")
	canvas.Line(int(cols[0].X2), -20, int(cols[0].X2), 4126, "stroke:#8a8a8a")
	canvas.Line(int(cols[1].X2), -20, int(cols[1].X2), 4126, "stroke:#8a8a8a")

	for i, c := range cols {
		for j, r := range c.Rings {
			if filter && ringFilter[[2]int{i, j}] {
				continue
			}
			xs, ys := pointsToIntArray(r.Points)
			canvas.Polygon(xs, ys, style(r.Label, i))
		}
	}
	canvas.Gend()
}
