package helper

const (
	BASE_MAP_OSM = "o-osm"
)

func BaseLayers() map[string]string {
	base_layers := make(map[string]string)
	base_layers[BASE_MAP_OSM] = "OpenStreetMap"
	base_layers["o-osmbw"] = "OpenstreetMap Grayscale"
	base_layers["o-carto-light"] = "OpenStreetMap Carto Light"
	base_layers["o-carto-dark"] = "OpenStreetMap Carto Dark"
	base_layers["o-stamen-toner"] = "OpenStreetMap Stamen Toner"
	base_layers["o-stamen-toner-hybrid"] = "OpenStreetMap Stamen Toner Hybrid"
	base_layers["o-stamen-toner-labels"] = "OpenStreetMap Stamen Toner Labels"
	base_layers["o-stamen-toner-lines"] = "OpenStreetMap Stamen Toner Lines"
	base_layers["o-stamen-toner-background"] = "OpenStreetMap Stamen Toner Background"
	base_layers["o-stamen-toner-lite"] = "OpenStreetMap Stamen Toner Lite"
	base_layers["o-stamen-watercolor"] = "OpenStreetMap Stamen Water Color"
	base_layers["o-stamen-terrain"] = "OpenStreetMap Stamen Terrain"
	base_layers["o-stamen-terrain-background"] = "OpenStreetMap Stamen Terrain Background"

	base_layers["g-roadmap"] = "Google Maps - Road Map"
	base_layers["g-satellite"] = "Google Maps - Satellite"
	base_layers["g-terrain"] = "Google Maps - Terrain"
	base_layers["g-hybrid"] = "Google Maps - Hybrid"

	base_layers["b-Aerial"] = "Bing Maps - Aerial"
	// 2017-08-31: BirdsEye view is not working
	// At this line
	// var r = meta.resourceSets[0].resources[0];
	// in /assets/leaflet/Bing.js : initMetadata
	// r is empty when type/imageSet is BirdsEye
	// base_layers["b-BirdsEye"] = "Bing Maps - Birds Eye"

	base_layers["b-Road"] = "Bing Maps - Road"
	base_layers["b-CanvasDark"] = "Bing Maps - Canvas Dark"
	base_layers["b-CanvasLight"] = "Bing Maps - Canvas Light"
	base_layers["b-CanvasGray"] = "Bing Maps - Canvas Gray"

	base_layers["y-map"] = "Yandex - Map"
	base_layers["y-satellite"] = "Yandex - Satellite"
	base_layers["y-hybrid"] = "Yandex - Hybrid"
	base_layers["y-publicMap"] = "Yandex - Public Map"
	base_layers["y-publicMapHybrid"] = "Yandex - Public Map Hybrid"

	return base_layers
}
