package validator

import (
	"net/http"
	"net/url"
	"testing"
)

func getMapValidatorUserHash() map[string]string {
	user := make(map[string]string)
	user["google_maps_key"] = "asflkjdfasd"
	user["bing_maps_key"] = "sdfasfd"

	return user
}

func TestMapValidateBaseLayer(t *testing.T) {
	base_layers := make(map[string]string)
	base_layers["o-osm"] = "OpenStreetMap"
	base_layers["g-hybrid"] = "Google Maps - Hybrid"
	base_layers["b-Road"] = "Bing Maps - Road"
	base_layers["y-satellite"] = "Yandex - Satellite"

	user := getMapValidatorUserHash()
	testMapValidateBaseLayer(t, "", base_layers, user, 1)
	testMapValidateBaseLayer(t, "xyz", base_layers, user, 1)
	testMapValidateBaseLayer(t, "b-Road", base_layers, user, 0)
	testMapValidateBaseLayer(t, "g-hybrid", base_layers, user, 0)

	user = getMapValidatorUserHash()
	user["google_maps_key"] = ""
	testMapValidateBaseLayer(t, "g-hybrid", base_layers, user, 1)

	user = getMapValidatorUserHash()
	user["bing_maps_key"] = ""
	testMapValidateBaseLayer(t, "b-Road", base_layers, user, 1)

	user = getMapValidatorUserHash()
	user["google_maps_key"] = ""
	testMapValidateBaseLayer(t, "b-Road", base_layers, user, 0)

	user = getMapValidatorUserHash()
	user["bing_maps_key"] = ""
	testMapValidateBaseLayer(t, "g-hybrid", base_layers, user, 0)
}

func testMapValidateBaseLayer(t *testing.T, base_layer string, base_layers map[string]string, user map[string]string, expected_len int) {
	values := url.Values{}
	values.Set("base_layer", base_layer)
	request := makeRequest(http.MethodPost, values)
	validator := MakeMapValidator(request)
	errs := validator.ValidateBaseLayer(base_layers, user)
	err_len := len(errs)
	if err_len != expected_len {
		t.Errorf("Error count: expected(1), actual(%d)", err_len)
	}
}
