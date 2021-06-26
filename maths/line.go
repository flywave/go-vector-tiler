package maths

import "math"

type Line [2]Pt

func NewLine(x1, y1, x2, y2 float64) Line {
	return Line{
		Pt{x1, y1},
		Pt{x2, y2},
	}
}

func NewLineFloat64(ln [2][]float64) Line {
	return Line{
		Pt{ln[0][0], ln[0][1]},
		Pt{ln[1][0], ln[1][1]},
	}
}
func NewLinesFloat64(ln ...[2][]float64) (lns []Line) {
	for i := range ln {
		lns = append(lns, Line{
			Pt{ln[i][0][0], ln[i][0][1]},
			Pt{ln[i][1][0], ln[i][1][1]},
		})
	}
	return lns
}

func NewLineWith2Float64(ln [2][]float64) Line {
	return Line{
		Pt{ln[0][0], ln[0][1]},
		Pt{ln[1][0], ln[1][1]},
	}
}

func (l Line) InBetween(pt Pt) bool {
	lx, gx := l[0].X, l[1].X
	if l[0].X > l[1].X {
		lx, gx = l[1].X, l[0].X
	}
	ly, gy := l[0].Y, l[1].Y
	if l[0].Y > l[1].Y {
		ly, gy = l[1].Y, l[0].Y
	}
	return lx <= pt.X && pt.X <= gx && ly <= pt.Y && pt.Y <= gy

}
func (l Line) ExInBetween(pt Pt) bool {
	lx, gx := l[0].X, l[1].X
	if l[0].X > l[1].X {
		lx, gx = l[1].X, l[0].X
	}
	ly, gy := l[0].Y, l[1].Y
	if l[0].Y > l[1].Y {
		ly, gy = l[1].Y, l[0].Y
	}

	goodx, goody := lx < pt.X && pt.X < gx, ly < pt.Y && pt.Y < gy
	if gx-lx == 0 {
		goodx = true
	}
	if gy-ly == 0 {
		goody = true
	}

	return goodx && goody

}

func (l Line) GetType() string {
	return "LineString"
}

func (l Line) IsVertical() bool {
	return l[0].X == l[1].X
}

func (l Line) IsHorizontal() bool {
	return l[0].Y == l[1].Y
}

func (l Line) Clamp(pt Pt) (p Pt) {
	p = pt
	lx, gx := l[0].X, l[1].X
	if l[0].X > l[1].X {
		lx, gx = l[1].X, l[0].X
	}
	ly, gy := l[0].Y, l[1].Y
	if l[0].Y > l[1].Y {
		ly, gy = l[1].Y, l[0].Y
	}

	if pt.X < lx {
		p.X = lx
	}
	if pt.X > gx {
		p.X = gx
	}
	if pt.Y < ly {
		p.Y = ly
	}
	if pt.Y > gy {
		p.Y = gy
	}
	return p
}

func (l Line) DistanceFromPoint(pt Pt) float64 {
	deltaX := l[1].X - l[0].X
	deltaY := l[1].Y - l[0].Y
	denom := math.Abs((deltaY * pt.X) - (deltaX * pt.Y) + (l[1].X * l[0].Y) - (l[1].Y * l[0].X))
	num := math.Sqrt(math.Pow(deltaY, 2) + math.Pow(deltaX, 2))
	if num == 0 {
		return 0
	}
	return denom / num
}

func (l Line) SlopeIntercept() (m, b float64, defined bool) {
	dx := l[1].X - l[0].X
	dy := l[1].Y - l[0].Y
	if dx == 0 || dy == 0 {
		return 0, l[0].Y, dx != 0
	}
	m = dy / dx
	b = l[0].Y - (m * l[0].X)
	return m, b, true
}

func (l Line) DeltaX() float64 { return l[1].X - l[0].X }

func (l Line) DeltaY() float64 { return l[1].Y - l[0].Y }

func (l Line) IsLeft(pt Pt) float64 {
	return (l.DeltaX() * (pt.Y - l[0].Y)) - ((pt.X - l[0].X) * l.DeltaY())
}

func (l Line) LeftRightMostPts() (Pt, Pt) {
	if XYOrder(l[0], l[1]) < 0 {
		return l[0], l[1]
	}
	return l[1], l[0]
}

func (l1 Line) LeftRightMostAsLine() Line {
	l, r := l1.LeftRightMostPts()
	return Line{l, r}
}

type Ring []Pt

type MultiLine []Line

func (l MultiLine) GetType() string {
	return "MultiLine"
}

type Polygon [][]Pt

func (l Polygon) GetType() string {
	return "Polygon"
}
