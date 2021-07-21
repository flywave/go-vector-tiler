package tile

import (
	"fmt"
	"testing"

	"github.com/flywave/go-geom/general"
)

func TestGrid(t *testing.T) {
	ext := NewMercGrid(&general.Extent{1.3546579391466225e+07, 4.4001531394751705e+06, 1.3621401362447761e+07, 4.47109306088063e+06})
	fmt.Println(ext.TileBound(10))
}
