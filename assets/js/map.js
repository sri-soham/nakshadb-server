var Map = (function() {
  var thisClass = {};
  var _map, _url, _delete_feature_url, _marker_url, _extent, _tile_layer, _grid_url, _style_changed, _geometry_updated;
  var _geometry_type, _point_count, _points_group, _feature, _style, _update_hash, _adding_feature, _grid_cache, _base_layer;
  var _google_maps_key, _bing_maps_key, _base_layer_type;
  var _init = 0;
  var _feature_id = null;
  var POLYGON = 'polygon';
  var LINESTRING = 'linestring';
  var POINT = 'point';

  function set_map(extent) {
    if (extent.length > 0) {
      var tmp = nakshamap.Util.centerAndBoundsFromExtent(extent);
      _map = new nakshamap.Map('map', {center: tmp['center'], zoom: 15});
      return tmp['bounds'];
    }
    else {
      _map = new nakshamap.Map('map', {center: [0, 0], zoom: 4});
      return null;
    }
    
  }

  function add_base_layer() {
    _base_layer = nakshamap.Util.getBaseLayer(_base_layer_type, _bing_maps_key);
    _map.addLayer(_base_layer);
  }

  function add_layer(url) {
    var tile_url = url.replace('[ts]', _update_hash);
    if (_tile_layer === null) {
      _tile_layer = nakshamap.layer(tile_url);
      _tile_layer.onClick(layer_clicked);
      _map.addLayer(_tile_layer);
    }
    else {
      _tile_layer.setUrl(tile_url);
    }
  }

  function show_delete_layer(data, naksha_id) {
    if (!data.data) return;
    if (_feature_id !== null) return;

    var i, table, tr, td;

    _feature_id = naksha_id;
    table = $('#geom-infowindow').children('table').first();
    $(table).html('');
    if ($.isEmptyObject(data.data)) {
      tr = $('<tr>');
      td = $('<td>').css('text-align', 'left').html('Please select the fields to be shown in the "Info Window" tab');
      $(tr).append(td);
      $(table).append(tr);
    }
    else {
      for (i in data.data) {
        tr = $('<tr>');
        td = $('<td>').html(i + ':');
        $(tr).append(td);
        td = $('<td>').html(data.data[i]);
        $(tr).append(td);
        $(table).append(tr);
      }
    }
    $('#geom-infowindow').removeClass('disp-none');
  }

  function delete_feature() {
    if (_feature_id === null) return;

    var div = $('#geom-infowindow').children('div').first();
    var url = _delete_feature_url + _feature_id;
    var options = {
      type: 'POST',
      url: url,
      msg_div: div,
      success_callback: function(rt) {
        _update_hash = rt['update_hash'];
        refresh_tile_layer();
        Table.deleteRowOfGeometry(_feature_id);
        close_delete_feature();
      }
    };
    new naksha.Ajax(options);
  }

  function close_delete_feature() {
    _feature_id = null;
    $('#geom-infowindow').addClass('disp-none');
    $('#geom-infowindow').children('table').first().html('');
  }

  function base_layer_changed() {
    var type, parts
    
    type = $(this).val();
    parts = type.split('-');
    if (parts[0] == 'b' && _bing_maps_key.length === 0) {
      naksha.Alert.show('Please enter bing maps key in profile to use bing maps base layer', 'Bing Maps', 400);
      return;
    }

    if (parts[0] == 'g' && _google_maps_key.length === 0) {
      naksha.Alert.show('Please enter google maps key in profile to use google maps base layer', 'Google Maps', 400);
      return;
    }

    _map.removeLayer(_base_layer);
    _base_layer = nakshamap.Util.getBaseLayer(type, _bing_maps_key);
    _map.addLayer(_base_layer);
    _base_layer.bringToBack();
  }

  function parse_styles() {
    var xml_doc, xml, elm, frm, values, i, tr, id, spectrum_options;

    xml_doc = $.parseXML(_style);
    xml = $(xml_doc);
    $('#frm_styles').find('tr').addClass('disp-none');
    $('#frm_styles').find('input[type="text"]').prop('disabled', true);
    $('#tr-submit').removeClass('disp-none');
    elm = $(xml).find('MarkersSymbolizer');
    frm = document.getElementById('frm_styles');
    values = {};
    switch (_geometry_type) {
      case POLYGON:
        elm = $(xml).find('PolygonSymbolizer').first();
        values['fill'] = $(elm).attr('fill');
        values['fill_opacity'] = $(elm).attr('fill-opacity');
        elm = $(xml).find('LineSymbolizer').first();
        values['stroke'] = $(elm).attr('stroke');
        values['stroke_width'] = $(elm).attr('stroke-width');
        values['stroke_opacity'] = $(elm).attr('stroke-opacity');
      break;
      case LINESTRING:
        elm = $(xml).find('LineSymbolizer').first();
        values['stroke'] = $(elm).attr('stroke');
        values['stroke_width'] = $(elm).attr('stroke-width');
        values['stroke_opacity'] = $(elm).attr('stroke-opacity');
      break;
      case POINT:
        elm = $(xml).find('MarkersSymbolizer').first();
        values['fill'] = $(elm).attr('fill');
        values['fill_opacity'] = $(elm).attr('opacity');
        values['stroke'] = $(elm).attr('stroke');
        values['stroke_width'] = $(elm).attr('stroke-width');
        values['stroke_opacity'] = $(elm).attr('stroke-opacity');
        values['width'] = $(elm).attr('width');
        values['height'] = $(elm).attr('height');
      break;
    }
    spectrum_options = {
      preferredFormat: 'hex',
      className: 'full-spectrum',
      showInitial: true,
      showPalette: true,
      showSelectionPalette: true,
      maxSelectionSize: 10,
      hideAfterPaletteSelect: true,
          palette: [
        ["rgb(0, 0, 0)", "rgb(67, 67, 67)", "rgb(102, 102, 102)",
        "rgb(204, 204, 204)", "rgb(217, 217, 217)","rgb(255, 255, 255)"],
        ["rgb(152, 0, 0)", "rgb(255, 0, 0)", "rgb(255, 153, 0)", "rgb(255, 255, 0)", "rgb(0, 255, 0)",
        "rgb(0, 255, 255)", "rgb(74, 134, 232)", "rgb(0, 0, 255)", "rgb(153, 0, 255)", "rgb(255, 0, 255)"],
        ["rgb(230, 184, 175)", "rgb(244, 204, 204)", "rgb(252, 229, 205)", "rgb(255, 242, 204)", "rgb(217, 234, 211)",
        "rgb(208, 224, 227)", "rgb(201, 218, 248)", "rgb(207, 226, 243)", "rgb(217, 210, 233)", "rgb(234, 209, 220)",
        "rgb(221, 126, 107)", "rgb(234, 153, 153)", "rgb(249, 203, 156)", "rgb(255, 229, 153)", "rgb(182, 215, 168)",
        "rgb(162, 196, 201)", "rgb(164, 194, 244)", "rgb(159, 197, 232)", "rgb(180, 167, 214)", "rgb(213, 166, 189)",
        "rgb(204, 65, 37)", "rgb(224, 102, 102)", "rgb(246, 178, 107)", "rgb(255, 217, 102)", "rgb(147, 196, 125)",
        "rgb(118, 165, 175)", "rgb(109, 158, 235)", "rgb(111, 168, 220)", "rgb(142, 124, 195)", "rgb(194, 123, 160)",
        "rgb(166, 28, 0)", "rgb(204, 0, 0)", "rgb(230, 145, 56)", "rgb(241, 194, 50)", "rgb(106, 168, 79)",
        "rgb(69, 129, 142)", "rgb(60, 120, 216)", "rgb(61, 133, 198)", "rgb(103, 78, 167)", "rgb(166, 77, 121)",
        "rgb(91, 15, 0)", "rgb(102, 0, 0)", "rgb(120, 63, 4)", "rgb(127, 96, 0)", "rgb(39, 78, 19)",
        "rgb(12, 52, 61)", "rgb(28, 69, 135)", "rgb(7, 55, 99)", "rgb(32, 18, 77)", "rgb(76, 17, 48)"]
      ]
    };

    for (i in values) {
      frm.elements[i].value = values[i];
      frm.elements[i].disabled = false;
      if (i == 'fill' || i == 'stroke') {
        $('#frm_styles').find('input[name="' + i + '"]').spectrum(spectrum_options);
      }
      id = 'tr-' + i.split('_').join('-');
      $('#' + id).removeClass('disp-none');
    }
  }

  function set_form() {
    var options = {
      form: '#frm_styles',
      success_callback: function(rt) {
        naksha.Message.success($('#frm_styles'), 'Styles Updated');
        _style_changed = true;
        _update_hash = rt['update_hash'];
      }
    };
    new naksha.AjaxForm(options);
  }

  function enable_draw_feature() {
    if (_geometry_type === 'unknown') {
      $('#geometry_type').removeClass('disp-none');
    }
    else {
      $('#add-feature-buttons-div').removeClass('disp-none');
      $('#base-layer-div').addClass('disp-none');
      _adding_feature = true;
    }
  }

  function geometry_type_changed() {
    var geom_type = $(this).val();
    if (geom_type.length === 0) {
      return;
    }
    _geometry_type = geom_type;
    $('#add-feature-buttons-div').removeClass('disp-none');
    $('#base-layer-div').addClass('disp-none');
    parse_styles();
    _adding_feature = true;
  }

  function get_circle_marker(ll, draggable) {
    var icon = L.icon({
      iconUrl: _marker_url,
      iconSize: [20, 20],
      iconAnchor: [10, 10]
    });
    var marker =  L.marker(ll, {
      icon: icon,
      draggable: draggable
    });
    if (draggable) {
      marker.on('dragend', point_moved);
    }

    return marker;
  }

  function fetch_geometry_details(e) {
    if (e.naksha_id === null || e.naksha_id === undefined) return;

    var options = {
      type: 'GET',
      url: Table.getShowRowUrl(e.naksha_id),
      cache: false,
      success_callback: function(rt) {
        show_delete_layer(rt, e.naksha_id);
      }
    };
    new naksha.Ajax(options);
  }

  function point_moved() {
    var lls = [];
    // _points_group will have circle-markers
    _points_group.eachLayer(function(lyr) {
      lls.push(lyr.getLatLng());
    });
    // _feature will be a polyline
    _feature.setLatLngs(lls);
  }

  function map_clicked(e) {
    if (_adding_feature) {
      handle_add_feature(e);
    }
  }

  function layer_clicked(e) {
    if (!_adding_feature) {
      fetch_geometry_details(e);
    }
  }

  function handle_add_feature(e) {
    var marker;

    _point_count++;
    var point = '' + e.latlng.lng + ' ' + e.latlng.lat;
    if (_point_count == 1) {
      if (_geometry_type === POINT) {
        _adding_feature = false;
        _feature = get_circle_marker(e.latlng, false);
        $('#btn_done').prop('disabled', false);
        $('#geometry_type').addClass('disp-none');
      }
      else {
        _points_group = L.featureGroup();
        marker = get_circle_marker(e.latlng, true);
        _points_group.addLayer(marker);
        _points_group.addTo(_map);

        _feature = L.polyline([e.latlng], {
          stroke: true,
          color: '#33ff99',
          width: 5,
          opacity: 1.0,
          fill: false
        });
      }
      _map.addLayer(_feature);
    }
    else if (_point_count > 1) {
      if (_geometry_type != POINT) {
        $('#btn_done').prop('disabled', false);
        $('#geometry_type').addClass('disp-none');
        _feature.addLatLng(e.latlng);

        marker = get_circle_marker(e.latlng, true);
        _points_group.addLayer(marker);
      }
    }
  }

  function feature_added() {
    var ll, ewkt, parts, options;

    if (_geometry_type == POINT && _point_count === 0) {
      Table.errorMessage('Please click on map to add marker');
      return;
    }
    if (_geometry_type != POINT && _point_count < 2) {
      Table.errorMessage('Need two or more points');
      return;
    }

    _adding_feature = false;

    switch (_geometry_type) {
      case POINT:
        ewkt = 'SRID=4326;POINT(';
        ll = _feature.getLatLng();
        ewkt += ll.lng + ' ' + ll.lat + ')';
      break;
      case LINESTRING:
        parts = [];
        _points_group.eachLayer(function(lyr) {
          ll = lyr.getLatLng();
          parts.push(ll.lng + ' ' + ll.lat);
        });
        ewkt = 'SRID=4326;MULTILINESTRING((' + parts.join(',') + '))';
      break;
      case POLYGON:
        parts = [];
        _points_group.eachLayer(function(lyr) {
          ll = lyr.getLatLng();
          parts.push(ll.lng + ' ' + ll.lat);
        });
        parts.push(parts[0]);
        ewkt = 'SRID=4326;MULTIPOLYGON(((' + parts.join(',') + ')))';
      break;
    }

    options = {
      url: Table.getAddRowUrl(),
      type: 'POST',
      data: {with_geometry: '1', geometry: ewkt},
      success_callback: function(rt) {
        _update_hash = rt['row']['update_hash'];
        refresh_tile_layer();
        feature_discarded(true);
        Table.addRowFromMap(rt);
      }
    };
    new naksha.Ajax(options);
  }

  function feature_discarded(feature_added) {
    if (_feature != null) {
      _map.removeLayer(_feature);
    }
    if (_points_group != null) {
      _map.removeLayer(_points_group);
    }
    if (feature_added !== true) {
      if ($('#geometry_type').val().length > 0) {
        // discard button has been clicked and the
        // geometry_type select box value is not empty, then
        // feature type is not been determined yet.
        _geometry_type = 'unknown';
        $('#geometry_type').val('');
      }
    }
    $('#geometry_type').addClass('disp-none');
    $('#add-feature-buttons-div').addClass('disp-none');
    $('#base-layer-div').removeClass('disp-none');
    $('#btn_done').prop('disabled', true);
    _point_count = 0;
    _feature = null;
    _points_group = null;
    _adding_feature = false;
  }

  function refresh_tile_layer() {
    add_layer(_url);
    _style_changed = false;
    _geometry_updated = false;
  }

  function set_event_listeners() {
    $('#btn_draw_feature').on('click', enable_draw_feature);
    $('#geometry_type').on('change', geometry_type_changed);
    $('#btn_done').on('click', feature_added);
    $('#btn_discard').on('click', feature_discarded);
    $('#btn_delete_feature').on('click', delete_feature);
    $('#btn_close_delete_feature').on('click', close_delete_feature);
    $('#base_layer').bind('change', base_layer_changed);
  }

  thisClass.init = function() {
    if (_init == 0) {
      if (_map) {
        // feature/row added to an empty table.
        var tmp = nakshamap.Util.centerAndBoundsFromExtent(extent);
        _map.setView(tmp['center'], 15);
        refresh_tile_layer();
        setTimeout(function() { _map.fitBounds(tmp['bounds']); }, 1200);
      }
      else {
        // map being initiated for the first time.
        var bounds = set_map(_extent);
        add_base_layer();
        add_layer(_url);
        if (bounds != null) {
          // empty table.
          _map.fitBounds(bounds);
        }
        _map.on('click', map_clicked);
      }
      _init = 1;
    }
    else {
      if (_style_changed || _geometry_updated) {
        refresh_tile_layer();
      }
    }
  };

  thisClass.setValues = function(url, delete_feature_url, update_hash, marker_url, extent, geometry_type, style, google_maps_key, bing_maps_key, base_layer_type) {
    _url = url;
    _delete_feature_url = delete_feature_url;
    _update_hash = update_hash;
    _marker_url = marker_url;
    _extent = extent;
    _geometry_type = geometry_type;
    _style = style;
    _google_maps_key = google_maps_key;
    _bing_maps_key = bing_maps_key;
    _base_layer_type = base_layer_type;

    _tile_layer = null;
    var parts = url.split('/');
    parts = parts.slice(0, 5);
    var last_part = parts.pop().split('-');
    _grid_url = parts.join('/') + '/' + last_part[0];
    _adding_feature = false;
    _grid_cache = {};

    parse_styles();
    set_form();
    set_event_listeners();
    _style_changed = false;
    _geometry_updated = false;
    _point_count = 0;
    _points_group = null;
    _feature = null;
  };

  thisClass.geometryUpdated = function(update_hash) {
    if (update_hash !== null) {
      _update_hash = update_hash;
    }
    _geometry_updated = true;
  };

  thisClass.updateGeometryType = function(geom_type, data_geom) {
    var spos, epos, parts, lls, i, ll, geometry;

    if (data_geom.length > 0) {
      _geometry_type = geom_type;
      parse_styles();
      
      spos = data_geom.lastIndexOf('(');
      spos++;
      epos = data_geom.indexOf(')');
      data_geom = data_geom.substr(spos, (epos-spos));
      parts = data_geom.split(',');
      lls = [];
      for (i in parts) {
        ll = parts[i].split(' ');
        lls.push([ll[1], ll[0]]);
      }
      if (geom_type == POINT) {
        _extent = '(' + lls[0][1] + ' ' + lls[0][0] + ',' + lls[0][1] + ' ' + lls[0][0] + ')';
      }
      else {
        if (geom_type == LINESTRING) {
          geometry = L.polyline(lls);
        }
        else {
          geometry = L.polygon(lls);
        }
        parts = [];
        ll = geometry.getBounds().getSouthWest();
        parts.push(ll.lng + ' ' + ll.lat);
        ll = geometry.getBounds().getNorthEast();
        parts.push(ll.lng + ' ' + ll.lat);
        _extent = '(' + parts.join(',') + ')';
      }
      _init = 0;
    }
  };

  thisClass.getMap = function() {
    return _map;
  };

  return thisClass;
})();

