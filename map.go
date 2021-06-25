package tile

import "github.com/flywave/go-geom/general"

// NewMap creates a new map with the necessary default values
func NewWebMercatorMap(name string) Map {
	return Map{
		Name:       name,
		Bounds:     nil,
		SRID:       3857,
		TileExtent: 4096,
	}
}

type Map struct {
	Name string
	// Contains an attribution to be displayed when the map is shown to a user.
	// 	This string is sanitized so it can't be abused as a vector for XSS or beacon tracking.
	Attribution string
	// The maximum extent of available map tiles in WGS:84
	// latitude and longitude values, in the order left, bottom, right, top.
	// Default: [-180, -85, 180, 85]
	Bounds *general.Extent
	// The first value is the longitude, the second is latitude (both in
	// WGS:84 values), the third value is the zoom level.
	Center [3]float64
	Layers []Layer

	SRID uint64
	// MVT output values
	TileExtent uint64
	TileBuffer uint64
}

// func (m Map) Tiler(ctx context.Context) ([]byte, error) {

// 	// layer stack
// 	mvtLayers := make([]*Layer, len(m.Layers))

// 	// set our waitgroup count
// 	wg.Add(len(m.Layers))

// 	// iterate our layers
// 	for i, layer := range m.Layers {

// 		// go routine for fetching the layer concurrently
// 		go func(i int, l Layer) {
// 			mvtLayer := mvt.Layer{
// 				Name: l.MVTName(),
// 			}

// 			// on completion let the wait group know
// 			defer wg.Done()

// 			ptile := provider.NewTile(tile.Z, tile.X, tile.Y,
// 				uint(m.TileBuffer), uint(m.SRID))

// 			// fetch layer from data provider
// 			err := l.Provider.TileFeatures(ptile, func(f *provider.Feature) error {
// 				// skip row if geometry collection empty.
// 				g, ok := f.Geometry.(geom.Collection)
// 				if ok && len(g.Geometries()) == 0 {
// 					return nil
// 				}

// 				geo := f.Geometry

// 				// check if the feature SRID and map SRID are different. If they are then reporject
// 				if f.SRID != m.SRID {
// 					// TODO(arolek): support for additional projections
// 					g, err := basic.ToWebMercator(f.SRID, geo)
// 					if err != nil {
// 						return fmt.Errorf("unable to transform geometry to webmercator from SRID (%v) for feature %v due to error: %w", f.SRID, f.ID, err)
// 					}
// 					geo = g
// 				}

// 				// add default tags, but don't overwrite a tag that already exists
// 				for k, v := range l.DefaultTags {
// 					if _, ok := f.Tags[k]; !ok {
// 						f.Tags[k] = v
// 					}
// 				}

// 				// TODO (arolek): change out the tile type for VTile. tegola.Tile will be deprecated
// 				tegolaTile := tegola.NewTile(tile.ZXY())

// 				sg := tegolaGeo
// 				// multiple ways to turn off simplification. check the atlas init() function
// 				// for how the second two conditions are set
// 				if !l.DontSimplify && simplifyGeometries && tile.Z < simplificationMaxZoom {
// 					sg = simplify.SimplifyGeometry(tegolaGeo, tegolaTile.ZEpislon())
// 				}

// 				// check if we need to clip and if we do build the clip region (tile extent)
// 				var clipRegion *geom.Extent
// 				if !l.DontClip {
// 					// CleanGeometry is expecting to operate in pixel coordinates so the clipRegion
// 					// will need to be in this same coordinate system. this will change when the new
// 					// make valid routing is implemented
// 					pbb, err := tegolaTile.PixelBufferedBounds()
// 					if err != nil {
// 						return fmt.Errorf("err calculating tile pixel buffer bounds: %w", err)
// 					}

// 					clipRegion = geom.NewExtent([2]float64{pbb[0], pbb[1]}, [2]float64{pbb[2], pbb[3]})
// 				}

// 				geo = mvt.PrepareGeo(geo, tile.Extent3857(), float64(mvt.DefaultExtent))

// 				tegolaGeo, err = validate.CleanGeometry(ctx, sg, clipRegion)
// 				if err != nil {
// 					return fmt.Errorf("err making geometry valid: %w", err)
// 				}

// 				geo, err = convert.ToGeom(tegolaGeo)
// 				if err != nil {
// 					return nil
// 				}

// 				mvtLayer.AddFeatures(mvt.Feature{
// 					ID:       &f.ID,
// 					Tags:     f.Tags,
// 					Geometry: geo,
// 				})

// 				return nil
// 			})
// 			if err != nil {
// 				switch {
// 				case errors.Is(err, context.Canceled):
// 					// Do nothing if we were cancelled.

// 				default:
// 					z, x, y := tile.ZXY()
// 					// TODO (arolek): should we return an error to the response or just log the error?
// 					// we can't just write to the response as the waitgroup is going to write to the response as well
// 					log.Printf("err fetching tile (z: %v, x: %v, y: %v) features: %v", z, x, y, err)
// 				}
// 				return
// 			}

// 			// add the layer to the slice position
// 			mvtLayers[i] = &mvtLayer
// 		}(i, layer)
// 	}

// 	// wait for the waitgroup to finish
// 	wg.Wait()

// 	// stop processing if the context has an error. this check is necessary
// 	// otherwise the server continues processing even if the request was canceled
// 	// as the waitgroup was not notified of the cancel
// 	if ctx.Err() != nil {
// 		return nil, ctx.Err()
// 	}

// 	// add layers to our tile
// 	mvtTile.AddLayers(mvtLayers...)

// 	// generate the MVT tile
// 	vtile, err := mvtTile.VTile(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// encode our mvt tile
// 	return proto.Marshal(vtile)
// }
