var ShowMap = (function() {
  var thisClass = {};
  var _map, _bing_maps_key, _schema_name;

  function extents_map_from_array(extents_arr) {
    var i, extents;

    extents = {};
    for (i=0; i<extents_arr.length; i+=2) {
      extents[extents_arr[i]] = extents_arr[i+1];
    }

    return extents;
  }

  function add_tile_layers(extents, layer_data) {
    var ll_bounds, i, ld, layer, tmp;

    ll_bounds = L.latLngBounds();
    for (i in layer_data) {
      ld = layer_data[i];

      layer = nakshamap.layer('', {schema: _schema_name, table: ld['table_name']});
      layer.addTo(_map);

      tmp = nakshamap.Util.centerAndBoundsFromExtent(extents[ld['table_name']]);
      ll_bounds.extend(tmp['bounds']);
    }
    _map.fitBounds(ll_bounds);
  }

  thisClass.init = function(map_div, schema_name, base_layer, layer_data, extents_arr) {
    _schema_name = schema_name;
    var extents = extents_map_from_array(extents_arr);
    var tmp = nakshamap.Util.centerAndBoundsFromExtent(extents[layer_data[0]['table_name']]);
    _map = new nakshamap.Map(map_div, {center: tmp['center'], zoom: 10});
    nakshamap.Util.getBaseLayer(base_layer, _bing_maps_key).addTo(_map);
    add_tile_layers(extents, layer_data);
  };

  thisClass.setBingMapsKey = function(key) {
    _bing_maps_key = key;
  };

  return thisClass;
})();

