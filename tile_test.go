package tile

import (
	"math"
	"testing"

	"github.com/flywave/go-vector-tiler/maths/webmercator"
	"github.com/flywave/go-vector-tiler/util"
)

// 测试常量
const (
	testZ     uint32  = 10
	testX     uint32  = 512
	testY     uint32  = 512
	testLat   float64 = 40.7128
	testLong  float64 = -74.0060
	testDelta float64 = 0.0001 // 浮点数比较的容差
)

// 测试辅助函数：比较浮点数是否相等
func floatEquals(a, b, delta float64) bool {
	return math.Abs(a-b) < delta
}

// 测试辅助函数：比较两个点是否相等
func pointEquals(p1, p2 [2]float64, delta float64) bool {
	return floatEquals(p1[0], p2[0], delta) && floatEquals(p1[1], p2[1], delta)
}

// 测试 NewTile 函数
func TestNewTile(t *testing.T) {
	tile := NewTile(testZ, testX, testY)

	if tile == nil {
		t.Fatal("NewTile 返回了 nil")
	}

	if tile.Z != testZ || tile.X != testX || tile.Y != testY {
		t.Errorf("瓦片坐标错误，期望 (%d, %d, %d)，得到 (%d, %d, %d)",
			testZ, testX, testY, tile.Z, tile.X, tile.Y)
	}

	if tile.Buffer != DefaultTileBuffer {
		t.Errorf("瓦片缓冲区大小错误，期望 %f，得到 %f", DefaultTileBuffer, tile.Buffer)
	}

	if tile.Extent != DefaultExtent {
		t.Errorf("瓦片范围错误，期望 %f，得到 %f", DefaultExtent, tile.Extent)
	}

	if tile.Tolerance != DefaultEpislon {
		t.Errorf("瓦片容差值错误，期望 %f，得到 %f", DefaultEpislon, tile.Tolerance)
	}

	// 检查是否正确初始化了内部字段
	if tile.extent == nil {
		t.Error("瓦片地理范围未初始化")
	}

	if tile.bufpext == nil {
		t.Error("瓦片缓冲区范围未初始化")
	}
}

// 测试 NewTileLatLong 函数
func TestNewTileLatLong(t *testing.T) {
	tile := NewTileLatLong(testZ, testLat, testLong)

	if tile == nil {
		t.Fatal("NewTileLatLong 返回了 nil")
	}

	if tile.Z != testZ {
		t.Errorf("瓦片缩放级别错误，期望 %d，得到 %d", testZ, tile.Z)
	}

	if !floatEquals(tile.Lat, testLat, testDelta) {
		t.Errorf("瓦片纬度错误，期望 %f，得到 %f", testLat, tile.Lat)
	}

	if !floatEquals(tile.Long, testLong, testDelta) {
		t.Errorf("瓦片经度错误，期望 %f，得到 %f", testLong, tile.Long)
	}
}

// 测试 NewTileWithOptions 函数
func TestNewTileWithOptions(t *testing.T) {
	const (
		customBuffer    = 128.0
		customExtent    = 4096.0
		customTolerance = 5.0
	)

	tile := NewTileWithOptions(testZ, testX, testY, customBuffer, customExtent, customTolerance)

	if tile == nil {
		t.Fatal("NewTileWithOptions 返回了 nil")
	}

	if tile.Buffer != customBuffer {
		t.Errorf("瓦片缓冲区大小错误，期望 %f，得到 %f", customBuffer, tile.Buffer)
	}

	if tile.Extent != customExtent {
		t.Errorf("瓦片范围错误，期望 %f，得到 %f", customExtent, tile.Extent)
	}

	if tile.Tolerance != customTolerance {
		t.Errorf("瓦片容差值错误，期望 %f，得到 %f", customTolerance, tile.Tolerance)
	}
}

// 测试 XYZToStringId 和 XYZFromStringId 函数
func TestTileIdConversion(t *testing.T) {
	// 测试 XYZToStringId
	id := XYZToStringId(testX, testY, testZ)
	expectedId := "512.512.10"
	if id != expectedId {
		t.Errorf("XYZToStringId 错误，期望 %s，得到 %s", expectedId, id)
	}

	// 测试 XYZFromStringId
	x, y, z, err := XYZFromStringId(id)
	if err != nil {
		t.Errorf("XYZFromStringId 返回错误: %v", err)
	}
	if x != testX || y != testY || z != testZ {
		t.Errorf("XYZFromStringId 错误，期望 (%d, %d, %d)，得到 (%d, %d, %d)",
			testX, testY, testZ, x, y, z)
	}

	// 测试无效的ID
	_, _, _, err = XYZFromStringId("invalid")
	if err == nil {
		t.Error("XYZFromStringId 应该对无效ID返回错误")
	}
}

// 测试 ToString 方法
func TestToString(t *testing.T) {
	tile := NewTile(testZ, testX, testY)
	id := tile.ToString()
	expectedId := "512.512.10"
	if id != expectedId {
		t.Errorf("ToString 错误，期望 %s，得到 %s", expectedId, id)
	}
}

// 测试 Deg2Num 和 Num2Deg 方法
func TestCoordinateConversion(t *testing.T) {
	// 创建一个已知经纬度的瓦片
	tile := NewTileLatLong(testZ, testLat, testLong)

	// 测试 Deg2Num
	x, y := tile.Deg2Num()
	if uint32(x) != tile.X || uint32(y) != tile.Y {
		t.Errorf("Deg2Num 错误，期望 (%d, %d)，得到 (%d, %d)",
			tile.X, tile.Y, x, y)
	}

	// 测试 Num2Deg
	lat, lng := tile.Num2Deg()
	if !floatEquals(lat, tile.Lat, testDelta) || !floatEquals(lng, tile.Long, testDelta) {
		t.Errorf("Num2Deg 错误，期望 (%f, %f)，得到 (%f, %f)",
			tile.Lat, tile.Long, lat, lng)
	}
}

// 测试 Bounds 方法
func TestBounds(t *testing.T) {
	tile := NewTile(testZ, testX, testY)
	bounds := tile.Bounds()

	// 检查边界是否有4个值
	if len(bounds) != 4 {
		t.Errorf("Bounds 应该返回4个值，但返回了 %d 个", len(bounds))
	}

	// 检查边界的顺序：西经, 南纬, 东经, 北纬
	west := bounds[0]
	south := bounds[1]
	east := bounds[2]
	north := bounds[3]

	// 检查东经是否大于西经
	if east <= west {
		t.Errorf("东经应该大于西经，但得到 东经=%f, 西经=%f", east, west)
	}

	// 检查北纬是否大于南纬
	if north <= south {
		t.Errorf("北纬应该大于南纬，但得到 北纬=%f, 南纬=%f", north, south)
	}
}

// 测试 ToPixel 和 FromPixel 方法
func TestPixelConversion(t *testing.T) {
	tile := NewTile(testZ, testX, testY)

	// 测试 Web墨卡托坐标系
	testPoint := [2]float64{0, 0} // 原点

	// 转换为像素坐标
	pixelPoint, err := tile.ToPixel(util.WebMercator, testPoint)
	if err != nil {
		t.Errorf("ToPixel 返回错误: %v", err)
	}

	// 转换回地理坐标
	geoPoint, err := tile.FromPixel(util.WebMercator, pixelPoint)
	if err != nil {
		t.Errorf("FromPixel 返回错误: %v", err)
	}

	// 检查转换后的坐标是否接近原始坐标
	if !pointEquals(testPoint, geoPoint, testDelta) {
		t.Errorf("像素坐标转换不一致，原始点 (%f, %f)，转换后点 (%f, %f)",
			testPoint[0], testPoint[1], geoPoint[0], geoPoint[1])
	}

	// 测试无效的SRID
	_, err = tile.ToPixel(9999, testPoint)
	if err == nil {
		t.Error("ToPixel 应该对无效SRID返回错误")
	}
}

// 测试 ZRes 方法
func TestZRes(t *testing.T) {
	// 测试不同缩放级别的分辨率
	testCases := []struct {
		z        uint32
		expected float64
	}{
		{0, webmercator.MaxXExtent * 2 / DefaultExtent},
		{1, webmercator.MaxXExtent * 2 / (DefaultExtent * 2)},
		{10, webmercator.MaxXExtent * 2 / (DefaultExtent * math.Pow(2, 10))},
	}

	for _, tc := range testCases {
		tile := NewTile(tc.z, 0, 0)
		res := tile.ZRes()
		if !floatEquals(res, tc.expected, testDelta) {
			t.Errorf("ZRes 在缩放级别 %d 错误，期望 %f，得到 %f",
				tc.z, tc.expected, res)
		}
	}
}

// 测试 ZEpislon 方法
func TestZEpislon(t *testing.T) {
	// 测试最大缩放级别
	maxZTile := NewTile(MaxZ, 0, 0)
	if maxZTile.ZEpislon() != 0 {
		t.Errorf("在最大缩放级别 %d，ZEpislon 应该返回 0，但得到 %f",
			MaxZ, maxZTile.ZEpislon())
	}

	// 测试容差为0
	zeroTolTile := NewTileWithOptions(10, 0, 0, DefaultTileBuffer, DefaultExtent, 0)
	if zeroTolTile.ZEpislon() != 0 {
		t.Errorf("当容差为 0 时，ZEpislon 应该返回 0，但得到 %f",
			zeroTolTile.ZEpislon())
	}

	// 测试正常情况
	normalTile := NewTile(10, 0, 0)
	expected := normalTile.Tolerance / (math.Exp2(float64(normalTile.Z)) * normalTile.Extent)
	if !floatEquals(normalTile.ZEpislon(), expected, testDelta) {
		t.Errorf("ZEpislon 错误，期望 %f，得到 %f",
			expected, normalTile.ZEpislon())
	}
}

// 测试 IsNeighbor 方法
func TestIsNeighbor(t *testing.T) {
	center := NewTile(10, 10, 10)

	// 测试相邻的瓦片
	neighbors := [][2]uint32{
		{9, 10},  // 左
		{11, 10}, // 右
		{10, 9},  // 上
		{10, 11}, // 下
		{9, 9},   // 左上
		{11, 9},  // 右上
		{9, 11},  // 左下
		{11, 11}, // 右下
	}

	for _, n := range neighbors {
		neighbor := NewTile(10, n[0], n[1])
		if !center.IsNeighbor(neighbor) {
			t.Errorf("瓦片 (%d, %d, %d) 应该是 (%d, %d, %d) 的相邻瓦片",
				neighbor.Z, neighbor.X, neighbor.Y, center.Z, center.X, center.Y)
		}
	}

	// 测试非相邻的瓦片
	nonNeighbors := [][3]uint32{
		{10, 10, 10}, // 自身
		{10, 8, 10},  // 距离太远
		{10, 10, 8},  // 距离太远
		{9, 10, 10},  // 不同缩放级别
	}

	for _, n := range nonNeighbors {
		nonNeighbor := NewTile(n[0], n[1], n[2])
		if center.IsNeighbor(nonNeighbor) {
			t.Errorf("瓦片 (%d, %d, %d) 不应该是 (%d, %d, %d) 的相邻瓦片",
				nonNeighbor.Z, nonNeighbor.X, nonNeighbor.Y, center.Z, center.X, center.Y)
		}
	}
}

// 测试 GetParent 方法
func TestGetParent(t *testing.T) {
	// 测试正常情况
	child := NewTile(10, 512, 512)
	parent := child.GetParent()

	if parent == nil {
		t.Fatal("GetParent 返回了 nil")
	}

	if parent.Z != 9 || parent.X != 256 || parent.Y != 256 {
		t.Errorf("父瓦片坐标错误，期望 (9, 256, 256)，得到 (%d, %d, %d)",
			parent.Z, parent.X, parent.Y)
	}

	// 测试根瓦片
	root := NewTile(0, 0, 0)
	rootParent := root.GetParent()
	if rootParent != nil {
		t.Error("根瓦片的父瓦片应该是 nil")
	}
}

// 测试 GetChildren 方法
func TestGetChildren(t *testing.T) {
	parent := NewTile(9, 256, 256)
	children := parent.GetChildren()

	if len(children) != 4 {
		t.Fatalf("GetChildren 应该返回4个子瓦片，但返回了 %d 个", len(children))
	}

	// 检查子瓦片的坐标
	expectedChildren := [][2]uint32{
		{512, 512}, // 左上
		{513, 512}, // 右上
		{512, 513}, // 左下
		{513, 513}, // 右下
	}

	for i, child := range children {
		if child.Z != 10 || child.X != expectedChildren[i][0] || child.Y != expectedChildren[i][1] {
			t.Errorf("子瓦片 %d 坐标错误，期望 (10, %d, %d)，得到 (%d, %d, %d)",
				i, expectedChildren[i][0], expectedChildren[i][1], child.Z, child.X, child.Y)
		}
	}

	// 测试最大缩放级别
	maxZTile := NewTile(MaxZ, 0, 0)
	maxZChildren := maxZTile.GetChildren()
	if maxZChildren != nil {
		t.Error("最大缩放级别的瓦片不应该有子瓦片")
	}
}

// 测试 Clone 方法
func TestClone(t *testing.T) {
	original := NewTile(testZ, testX, testY)
	clone := original.Clone()

	// 检查基本属性是否相同
	if clone.Z != original.Z || clone.X != original.X || clone.Y != original.Y {
		t.Errorf("克隆瓦片坐标错误，期望 (%d, %d, %d)，得到 (%d, %d, %d)",
			original.Z, original.X, original.Y, clone.Z, clone.X, clone.Y)
	}

	if clone.Extent != original.Extent || clone.Buffer != original.Buffer || clone.Tolerance != original.Tolerance {
		t.Error("克隆瓦片的属性与原始瓦片不同")
	}

	// 检查内部字段是否正确克隆
	if clone.extent == nil || clone.bufpext == nil {
		t.Error("克隆瓦片的内部字段未正确初始化")
	}

	// 检查是否是深拷贝
	originalExt := original.extent.Extent()
	cloneExt := clone.extent.Extent()

	for i := 0; i < 4; i++ {
		if originalExt[i] != cloneExt[i] {
			t.Errorf("克隆瓦片的地理范围与原始瓦片不同，索引 %d: 期望 %f，得到 %f",
				i, originalExt[i], cloneExt[i])
		}
	}

	// 修改克隆瓦片，检查是否影响原始瓦片
	clone.SetBuffer(clone.Buffer * 2)
	if original.Buffer == clone.Buffer {
		t.Error("修改克隆瓦片影响了原始瓦片，不是深拷贝")
	}
}

// 测试 SetBuffer、SetExtent 和 SetTolerance 方法
func TestSetters(t *testing.T) {
	tile := NewTile(testZ, testX, testY)

	// 测试 SetBuffer
	newBuffer := 128.0
	tile.SetBuffer(newBuffer)
	if tile.Buffer != newBuffer {
		t.Errorf("SetBuffer 错误，期望 %f，得到 %f", newBuffer, tile.Buffer)
	}

	// 测试 SetExtent
	newExtent := 4096.0
	tile.SetExtent(newExtent)
	if tile.Extent != newExtent {
		t.Errorf("SetExtent 错误，期望 %f，得到 %f", newExtent, tile.Extent)
	}

	// 测试 SetTolerance
	newTolerance := 5.0
	tile.SetTolerance(newTolerance)
	if tile.Tolerance != newTolerance {
		t.Errorf("SetTolerance 错误，期望 %f，得到 %f", newTolerance, tile.Tolerance)
	}
}

// 测试 GetExtent 和 GetBufferedExtent 方法
func TestGetExtents(t *testing.T) {
	tile := NewTile(testZ, testX, testY)

	// 测试 GetExtent
	extent := tile.GetExtent()
	if extent == nil {
		t.Fatal("GetExtent 返回了 nil")
	}

	// 测试 GetBufferedExtent
	bufExt := tile.GetBufferedExtent()
	if bufExt == nil {
		t.Fatal("GetBufferedExtent 返回了 nil")
	}

	// 检查缓冲区范围是否正确
	bufExtArr := bufExt.Extent()
	expectedBufExt := [4]float64{
		0 - tile.Buffer,
		0 - tile.Buffer,
		tile.Extent + tile.Buffer,
		tile.Extent + tile.Buffer,
	}

	for i := 0; i < 4; i++ {
		if bufExtArr[i] != expectedBufExt[i] {
			t.Errorf("缓冲区范围错误，索引 %d: 期望 %f，得到 %f",
				i, expectedBufExt[i], bufExtArr[i])
		}
	}
}

// 测试 PixelBufferedBounds 方法
func TestPixelBufferedBounds(t *testing.T) {
	tile := NewTile(testZ, testX, testY)

	bounds, err := tile.PixelBufferedBounds()
	if err != nil {
		t.Errorf("PixelBufferedBounds 返回错误: %v", err)
	}

	// 检查边界是否有4个值
	if len(bounds) != 4 {
		t.Errorf("PixelBufferedBounds 应该返回4个值，但返回了 %d 个", len(bounds))
	}

	// 检查边界值是否正确
	expectedBounds := [4]float64{
		0 - tile.Buffer,
		0 - tile.Buffer,
		tile.Extent + tile.Buffer,
		tile.Extent + tile.Buffer,
	}

	for i := 0; i < 4; i++ {
		if bounds[i] != expectedBounds[i] {
			t.Errorf("像素缓冲区边界错误，索引 %d: 期望 %f，得到 %f",
				i, expectedBounds[i], bounds[i])
		}
	}
}
