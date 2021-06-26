package region

import (
	"errors"
	"fmt"

	"github.com/flywave/go-vector-tiler/container/singlelist/point/list"
	"github.com/flywave/go-vector-tiler/maths"
)

type IntersectionCode uint8

type IntersectionPt struct {
	Code IntersectionCode
}

type Axis struct {
	region      *Region
	idx         int
	downOrRight bool
	pt0, pt1    *list.Pt
	winding     maths.WindingOrder
}

func (a *Axis) GoString() string {
	return fmt.Sprintf("[%v,%v]-(%v)-(%v){%v}", a.pt0, a.pt1, a.downOrRight, a.idx, a.region.GoString())
}

func (a *Axis) Next() *Axis {
	if a.idx == 3 {
		return nil
	}
	return a.region.Axis(a.idx + 1)
}

func (a *Axis) PushInBetween(pt list.ElementerPointer) bool {
	return a.region.PushInBetween(a.pt0, a.pt1, pt)
}

func (a *Axis) AsLine() maths.Line { return maths.Line{a.pt0.Point(), a.pt1.Point()} }

func (a *Axis) Intersect(line maths.Line) (pt maths.Pt, doesIntersect bool) {
	axisLine := a.AsLine()
	var ok bool
	pt, ok = maths.Intersect(axisLine, line)
	if !ok {
		return pt, ok
	}

	if !line.InBetween(pt) {
		return pt, false
	}
	if !axisLine.ExInBetween(pt) {
		if (axisLine.IsHorizontal() && line.IsVertical()) ||
			(axisLine.IsVertical() && line.IsHorizontal()) ||
			(!axisLine.InBetween(pt)) {
			return pt, false
		}

	}
	return pt, true
}

func (a *Axis) inside(pt maths.Pt) bool {
	switch a.idx % 4 {
	case 0:
		return pt.X > a.pt0.X
	case 1:
		if a.winding.IsClockwise() {
			return pt.Y > a.pt0.Y
		}
		return pt.Y < a.pt0.Y
	case 2:
		return pt.X < a.pt0.X
	case 3:
		if a.winding.IsClockwise() {
			return pt.Y < a.pt0.Y
		}
		return pt.Y > a.pt0.Y
	}
	return false
}

var ErrNoDirection = errors.New("Line does not have direction on that coordinate.")

type PlacementCode uint8

const (
	PCInside      PlacementCode = 0x00               // 0000
	PCBottom                    = 0x01               // 0001
	PCTop                       = 0x02               // 0010
	PCRight                     = 0x04               // 0100
	PCLeft                      = 0x08               // 1000
	PCTopRight                  = PCTop | PCRight    // 0110
	PCTopLeft                   = PCTop | PCLeft     // 1010
	PCBottomRight               = PCBottom | PCRight // 0101
	PCBottomLeft                = PCBottom | PCLeft  // 1001

	PCAllAround = PCTop | PCLeft | PCRight | PCBottom // 1111
)

func (a *Axis) Placement(pt maths.Pt) PlacementCode {

	idx := a.idx % 4
	switch {
	case idx == 0 && pt.X <= a.pt0.X:
		return PCLeft
	case idx == 2 && pt.X >= a.pt0.X:
		return PCRight

	case ((a.winding.IsClockwise() && a.idx == 3) || a.idx == 1) && pt.Y <= a.pt0.Y:
		return PCTop
	case ((a.winding.IsClockwise() && a.idx == 1) || a.idx == 3) && pt.Y >= a.pt0.Y:
		return PCBottom
	default:
		return PCInside
	}

}

func (a *Axis) IsInward(line maths.Line) (bool, error) {
	p1, p2 := line[0], line[1]

	dx := p2.X - p1.X
	dy := p2.Y - p1.Y

	idx := a.idx % 4
	switch idx {
	case 0, 2:
		if dx == 0 {
			return false, ErrNoDirection
		}
		if idx == 0 {
			return dx > 0, nil
		}
		return dx < 0, nil

	case 1, 3:
		if dy == 0 {
			return false, ErrNoDirection
		}

		if a.winding.IsCounterClockwise() {
			if idx == 1 {
				idx = 3
			} else {
				idx = 1
			}
		}
		if idx == 1 {
			return dy > 0, nil
		}
		return dy < 0, nil
	}
	return false, ErrNoDirection
}
