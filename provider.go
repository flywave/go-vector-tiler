package tile

type Provider interface {
	GetDdataByTile(*Tile) []*Layer
}
