/*
Package webmercator does the translation to and from WebMercator and WGS84
Gotten from: http://wiki.openstreetmap.org/wiki/Mercator#C.23
*/
package webmercator

import (
	"errors"
	"fmt"
	"math"
)

const (
	RMajor = 6378137.0
	RMinor = 6356752.3142
	Ratio  = RMinor / RMajor
)

const (
	SRID        = 3857
	EarthRadius = RMajor
	Deg2Rad     = math.Pi / 180
	Rad2Deg     = 180 / math.Pi
	PiDiv2      = math.Pi / 2.0
	PiDiv4      = math.Pi / 4.0

	M_PIby360           = math.Pi / 360
	EARTH_CIRCUMFERENCE = EarthRadius * 2 * math.Pi
	MAXEXTENT           = EARTH_CIRCUMFERENCE / 2.0
	MAXEXTENTby180      = MAXEXTENT / 180
)

var MAX_LATITUDE = Rad2Deg * (2*math.Atan(math.Exp(180*Deg2Rad)) - PiDiv2)

var Extent = [4]float64{MinXExtent, MinYExtent, MaxXExtent, MaxYExtent}

var ErrCoordsRequire2Values = errors.New("Coords should have at least 2 coords")

func RadToDeg(rad float64) float64 {
	return rad * Rad2Deg
}

func DegToRad(deg float64) float64 {
	return deg * Deg2Rad
}

var Eccent float64
var Com float64

func init() {
	Eccent = math.Sqrt(1.0 - (Ratio * Ratio))
	Com = 0.5 * Eccent
}

func con(phi float64) float64 {
	v := Eccent * math.Sin(phi)
	return math.Pow(((1.0 - v) / (1.0 + v)), Com)
}

func LonToX(lon float64) float64 {
	return RMajor * DegToRad(lon)
}

func LatToY(lat float64) float64 {
	lat = math.Min(MAX_LATITUDE, math.Max(lat, -MAX_LATITUDE))
	y := math.Log(math.Tan((90+lat)*M_PIby360)) * Rad2Deg
	y = y * MAXEXTENTby180
	return y
}

var (
	MaxXExtent = LonToX(180)
	MaxYExtent = LatToY(MAX_LATITUDE)
	MinXExtent = -MaxXExtent
	MinYExtent = -MaxYExtent
)

func XToLon(x float64) float64 {
	return RadToDeg(x) / RMajor
}

func YToLat(y float64) float64 {
	ts := math.Exp(-y / RMajor)
	phi := PiDiv2 - 2*math.Atan(ts)
	dphi := 1.0
	i := 0
	for (math.Abs(dphi) > 0.000000001) && (i < 15) {
		dphi = PiDiv2 - 2*math.Atan(ts*con(phi)) - phi
		phi += dphi
		i++
	}
	return RadToDeg(phi)
}

func ToLonLat(c ...float64) ([]float64, error) {
	if len(c) < 2 {
		return c, fmt.Errorf("Coords should have at least 2 coords")
	}
	crds := []float64{XToLon(c[0]), YToLat(c[1])}
	crds = append(crds, c[2:]...)
	return crds, nil
}

func ToXY(c ...float64) ([]float64, error) {
	if len(c) < 2 {
		return c, fmt.Errorf("Coords should have at least 2 coords")
	}
	crds := []float64{LonToX(c[0]), LatToY(c[1])}
	crds = append(crds, c[2:]...)
	return crds, nil
}
