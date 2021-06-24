package tile

import (
	"github.com/flywave/go-geom"
	geom "github.com/flywave/go-geom"
)

type Layer struct {
	// optional. if not set, the ProviderLayerName will be used
	Name     string
	Features []*geom.Feature
	SRID     int
}

// MVTName will return the value that will be encoded in the Name field when the layer is encoded as MVT
func (l *Layer) GetName() string {
	return l.Name
}
