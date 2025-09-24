package tile

import "errors"

// 集中定义所有错误变量
var (
	// ErrInvalidTile 表示无效的瓦片
	ErrInvalidTile = errors.New("invalid tile")
	// ErrInvalidPath 表示无效的路径
	ErrInvalidPath = errors.New("invalid path")
	// ErrEmptyLayers 表示空图层
	ErrEmptyLayers = errors.New("empty layers")
)
