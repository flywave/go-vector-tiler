module github.com/flywave/go-vector-tiler

go 1.23.0

toolchain go1.24.4

require (
	github.com/ajstarks/svgo v0.0.0-20210406150507-75cfd577ce75
	github.com/flywave/go-geo v0.0.0-20250314091853-e818cb9de299
	github.com/flywave/go-geom v0.0.0-20250607125323-f685bf20f12c
	github.com/flywave/go-mapbox v0.0.0-20250802054600-bc350f35abb4
	github.com/flywave/go3d v0.0.0-20250314015505-bf0fda02e242
	github.com/gdey/tbltest v0.0.0-20180914212833-1865222d591f
	github.com/go-test/deep v1.0.7
	github.com/pborman/uuid v1.2.1
)

require (
	github.com/flywave/go-geoid v0.0.0-20210705014121-cd8f70cb88bb // indirect
	github.com/flywave/go-geos v0.0.0-20210924031454-d16b758e2026 // indirect
	github.com/flywave/go-pbf v0.0.0-20230306063816-5e5b0da27bbd // indirect
	github.com/flywave/go-proj v0.0.0-20211220121303-46dc797a5cd0 // indirect
	github.com/google/uuid v1.2.0 // indirect
)

replace github.com/flywave/go-geo => ../go-geo

replace github.com/flywave/go-proj => ../go-proj

replace github.com/flywave/go-geoid => ../go-geoid

replace github.com/flywave/go-geos => ../go-geos
