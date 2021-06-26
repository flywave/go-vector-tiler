package intersect

import (
	"fmt"
	"log"

	list "github.com/flywave/go-vector-tiler/container/singlelist"
	ptList "github.com/flywave/go-vector-tiler/container/singlelist/point/list"
	"github.com/flywave/go-vector-tiler/maths"
)

type Point struct {
	ptList.Pt

	Inward bool

	subject list.Sentinel
	region  list.Sentinel
}

func NewPt(pt maths.Pt, inward bool) *Point {
	return &Point{Pt: *ptList.NewPt(pt), Inward: inward}
}
func NewPoint(x, y float64, inward bool) *Point {
	return &Point{Pt: *ptList.NewPoint(x, y), Inward: inward}
}

func (p *Point) String() string {
	return fmt.Sprintf("Intersec{ X: %v, Y: %v, Inward: %v}", p.Pt.X, p.Pt.Y, p.Inward)
}
func (p *Point) AsSubjectPoint() *SubjectPoint { return (*SubjectPoint)(p) }
func (p *Point) AsRegionPoint() *RegionPoint   { return (*RegionPoint)(p) }

func (p *Point) NextWalk() list.Elementer {
	if p.Inward {
		return p.subject.Next()
	}
	return p.region.Next()
}

func (p *Point) PrintNeighbors() {
	log.Println("Me:", p.String())
	log.Println("\tRegion Neighbor. ->", p.region.Next())
	log.Println("\tSubject Neighbor. ->", p.subject.Next())
}

type RegionPoint Point

func (i *RegionPoint) Next() list.Elementer { return i.region.Next() }
func (i *RegionPoint) SetNext(e list.Elementer) list.Elementer {
	return i.region.SetNext(e)
}

func (i *RegionPoint) List() *list.List                { return i.region.List() }
func (i *RegionPoint) SetList(l *list.List) *list.List { return i.region.SetList(l) }
func (i *RegionPoint) AsSubjectPoint() *SubjectPoint {
	return (*SubjectPoint)(i)
}
func (i *RegionPoint) AsIntersectPoint() *Point { return (*Point)(i) }
func (i *RegionPoint) Point() maths.Pt          { return i.Pt.Point() }
func (i *RegionPoint) GoString() string {
	return fmt.Sprintf("%T(%[1]p)[%v;%v]", i, i.Point(), i.Inward)
}

type SubjectPoint Point

func (i *SubjectPoint) Next() list.Elementer { return i.subject.Next() }
func (i *SubjectPoint) SetNext(e list.Elementer) list.Elementer {
	return i.subject.SetNext(e)
}

func (i *SubjectPoint) List() *list.List                { return i.subject.List() }
func (i *SubjectPoint) SetList(l *list.List) *list.List { return i.subject.SetList(l) }
func (i *SubjectPoint) AsRegionPoint() *RegionPoint {
	return (*RegionPoint)(i)
}
func (i *SubjectPoint) AsIntersectPoint() *Point { return (*Point)(i) }
func (i *SubjectPoint) Point() maths.Pt          { return i.Pt.Point() }
func (i *SubjectPoint) GoString() string {
	return fmt.Sprintf("%T(%[1]p)[%v;%v]", i, i.Point(), i.Inward)
}
