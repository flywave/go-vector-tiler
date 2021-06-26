package points

import (
	"sort"

	"github.com/flywave/go-vector-tiler/maths"
)

func SortAndUnique(pts []maths.Pt) []maths.Pt {
	if len(pts) == 0 {
		return pts
	}

	sort.Sort(ByXY(pts))

	count := 0
	for i := range pts {
		if pts[count].IsEqual(pts[i]) {
			continue
		}

		count++
		pts[count] = pts[i]
	}

	return pts[:count+1]
}
