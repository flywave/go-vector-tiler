package tile

type Provider interface {
	GetDataByTile(*Tile) []*Layer
	GetSrid() uint64
}
