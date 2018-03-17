var nakshamap = {};

nakshamap.Util = L.Util.extend({
  getBaseLayer: function(base_layer_type, bing_maps_key) {
    var type, url, attribution, layer, provider;

    type = base_layer_type.split('-').slice(1).join('-');
    provider = base_layer_type[0];
    if (provider === 'g') {
      layer = L.gridLayer.googleMutant({
        type: type
      });
    }
    else if (provider === 'b') {
      layer = L.bingLayer(bing_maps_key, {type: type});
    }
    else if (provider === 'y') {
      layer = new L.Yandex(type);
    }
    else {
      attribution = '&copy; <a href="http://openstreetmap.org">OpenStreetMap</a> Contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>';
      switch (type) {
        case 'osmbw':
          url = 'http://{s}.tiles.wmflabs.org/bw-mapnik/{z}/{x}/{y}.png';
        break;
        case 'carto-light':
          url = 'http://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}.png';
          attribution = 'Map tiles by <a href="http://carto.com">Carto</a>, under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. Data by <a href="http://openstreetmap.org">OpenStreetMap</a> under ODbL';
        break;
        case 'carto-dark':
          url = 'http://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png';
          attribution = 'Map tiles by <a href="http://carto.com">Carto</a>, under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. Data by <a href="http://openstreetmap.org">OpenStreetMap</a> under ODbL';
        break;
        case 'stamen-toner':
          attribution = 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. - Map data &copy; <a href="http://www.openstreetmap.org">OpenStreetMap</a>';
          url = 'http://{s}.tile.stamen.com/toner/{z}/{x}/{y}.png';
        break;
        case 'stamen-toner-hybrid':
          attribution = 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. - Map data &copy; <a href="http://www.openstreetmap.org">OpenStreetMap</a>';
          url = 'http://{s}.tile.stamen.com/toner-hybrid/{z}/{x}/{y}.png';
        break;
        case 'stamen-toner-labels':
          attribution = 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. - Map data &copy; <a href="http://www.openstreetmap.org">OpenStreetMap</a>';
          url = 'http://{s}.tile.stamen.com/toner-labels/{z}/{x}/{y}.png';
        break;
        case 'stamen-toner-lines':
          attribution = 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. - Map data &copy; <a href="http://www.openstreetmap.org">OpenStreetMap</a>';
          url = 'http://{s}.tile.stamen.com/toner-lines/{z}/{x}/{y}.png';
        break;
        case 'stamen-toner-background':
          attribution = 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. - Map data &copy; <a href="http://www.openstreetmap.org">OpenStreetMap</a>';
          url = 'http://{s}.tile.stamen.com/toner-background/{z}/{x}/{y}.png';
        break;
        case 'stamen-toner-lite':
          attribution = 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. - Map data &copy; <a href="http://www.openstreetmap.org">OpenStreetMap</a>';
          url = 'http://{s}.tile.stamen.com/toner-lite/{z}/{x}/{y}.png';
        break;
        case 'stamen-watercolor':
          attribution = 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. - Map data &copy; <a href="http://www.openstreetmap.org">OpenStreetMap</a>';
          url = 'http://{s}.tile.stamen.com/watercolor/{z}/{x}/{y}.png';
        break;
        case 'stamen-terrain':
          attribution = 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. - Map data &copy; <a href="http://www.openstreetmap.org">OpenStreetMap</a>';
          url = 'http://{s}.tile.stamen.com/terrain/{z}/{x}/{y}.png';
        break;
        case 'stamen-terrain-background':
          attribution = 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. - Map data &copy; <a href="http://www.openstreetmap.org">OpenStreetMap</a>';
          url = 'http://{s}.tile.stamen.com/terrain-background/{z}/{x}/{y}.png';
        break;
        default:
          url = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
        break;
      }
      layer = L.tileLayer(url, {attribution: attribution, zIndex: 1});
    }

    return layer;
  },

  // extent is the string returned by postgis with ST_Extent
  centerAndBoundsFromExtent: function(extent) {
    var parts, part1, part, lat1, lng1, lat2, lng2, lat, lng;

    parts = extent.split(',')
    part1 = parts[0].split('(').pop();
    part2 = parts[1].split(')').shift();

    parts = part1.split(' ');
    lat1 = parseFloat(parts[1]);
    lng1 = parseFloat(parts[0]);
    parts = part2.split(' ');
    lat2 = parseFloat(parts[1]);
    lng2 = parseFloat(parts[0]);

    lat = lat1 + ((lat2 - lat1)/2);
    lng = lng1 + ((lng2 - lng1)/2);

    return {
      'center': [lat, lng],
      'bounds': [[lat1, lng1], [lat2, lng2]]
    };
  },

  showInfo: function(data) {
    var div, table, tr, td, span, i, f, dt;

    div = $('<div>').css({
      'position': 'absolute',
      'top': '100px',
      'left': '50%',
      'margin-left': '-150px',
      'width': '300px',
      'height': 'auto',
      'background-color': '#fff',
      'padding': '10px',
      'border': '1px #333 solid',
      'border-radius': '5px',
      'z-index': 1000
    });
    span = $('<span>').css({'cursor': 'pointer', 'float': 'right'}).html('x');
    $(span).bind('click', function() {
      $(this).parent('div').remove();
    });
    $(div).append(span);

    table = $('<table>');
    for (i in data) {
      dt = data[i];
      for (f in dt) {
        tr = $('<tr>');
        td = $('<td>').css({'text-align': 'right', 'vertical-align': 'top'}).html(f + ':' );
        $(tr).append(td);
        td = $('<td>').html(dt[f]);
        $(tr).append(td);
        $(table).append(tr);
      }
    }
    $(div).append(table);
    $(document.body).append(div);
  },

  SQL: function(url, query, success_callback, error_callback) {
    if (!url || url.length === 0) {
      url = window.location.protocol + '//' + window.location.host;
    }
    url += '/p/s';
    $.ajax({
      type: 'GET',
      url: url,
      data: {query: query},
      dataType: 'jsonp',
      cache: false
    }).done(function(rt) {
      if (rt['status']) {
        if (rt['status'] === 'success') {
          success_callback(rt['data']);
        }
        else {
          error_callback(rt['errors']);
        }
      }
      else {
        error_callback(['Invalid response from server']);
      }
    }).fail(function(rt, st) {
      error_callback(['Request failed']);
    });
  }
});

nakshamap.Map = L.Map.extend({
  options: {
    naksha_tile_size: 256,
    naksha_resolution: 4
  }
});

nakshamap.Layer = L.LayerGroup.extend({
  _data: {},
  _grid_layer: null,
  _tile_layer: null,
  _ready: false,
  _naksha_listeners: {},

  my_options: {
    table: '',
    schema: ''
  },

  initialize: function(url, options) {
    L.Util.setOptions(this, options);

    if (options && options.hasOwnProperty('schema') && options.hasOwnProperty('table')) {
      // While adding layer with nakshamap library, current update hash might not
      // be available. So, fetch details with ajax call and add the layer.
      this._init_with_ajax(url, options);
    }
    else {
      // While modifying the map after signing in, update-hash is available
      // so, use the url directly.
      this._init_with_url(url, options);
    }
  },

  onClick: function(listener) {
    if (this._ready) {
      this._grid_layer.on('click', listener);
    }
    else {
      if (!this._naksha_listeners['click']) {
        this._naksha_listeners['click'] = [];
      }
      this._naksha_listeners['click'].push(listener);
    }
  },

  onMouseOver: function(listener) {
    if (this._ready) {
      this._grid_layer.on('mouseover', listener);
      this._grid_layer.hasMouseEvent(true);
    }
    else {
      if (!this._naksha_listeners['mouseover']) {
        this._naksha_listeners['mouseover'] = [];
      }
      this._naksha_listeners['mouseover'].push(listener);
    }
  },

  onMouseOut: function(listener) {
    if (this._ready) {
      this._grid_layer.on('mouseout', listener);
      this._grid_layer.hasMouseEvent(true);
    }
    else {
      if (!this._naksha_listeners['mouseout']) {
        this._naksha_listeners['mouseout'] = [];
      }
      this._naksha_listeners['mouseout'].push(listener);
    }
  },

  setUrl: function(url) {
    this._tile_layer.setUrl(url);
    var json_url = this._jsonUrl(url);
    this._grid_layer.setUrl(json_url, false);
  },

  _jsonUrl: function(url) {
    var json_url = url.replace(/\.png$/, '.json');
    json_url += '?callback={cb}'

    return json_url;
  },

  _init_with_ajax: function(orig_url, options) {
    var that = this;
    var url = orig_url + '/p/l/' + options.schema + '/' + options.table;
    $.ajax({
      type: 'GET',
      url: url,
      dataType: 'jsonp',
      cache: false
    }).done(function(rt) {
      if (rt['status']) {
        if (rt['status'] === 'success') {
          var json_url, i, j, layer_url, has_mouse_event, infowindow, grid_options;

          layer_url = rt['layer_url'];
          that._tile_layer = L.tileLayer(layer_url);
          json_url = that._jsonUrl(layer_url);
          infowindow = $.parseJSON(rt['infowindow']);
          grid_options = {
            'table': options.schema + '.' + options.table,
            'info_fields': infowindow['fields'].sort(),
            'sql_url': orig_url
          };
          that._grid_layer = new nakshamap.Grid(json_url, grid_options);
          that._layers = {};
          that.addLayer(that._tile_layer);
          that.addLayer(that._grid_layer);
          that._ready = true;
          has_mouse_event = false;
          for (i in that._naksha_listeners) {
            if (i == 'mouseover' || i == 'mouseout') {
              has_mouse_event = true;
            }
            for (j=0; j<that._naksha_listeners[i].length; ++j) {
              that._grid_layer.on(i, that._naksha_listeners[i][j]);
            }
          }
          that._grid_layer.hasMouseEvent(has_mouse_event);
        }
        else {
          console.log('errors: ' + rt['errors'].join('\n'));
        }
      }
      else {
        console.log('Invalid response from server');
      }
    }).fail(function(rt, st) {
      console.log('Request failed');
    });
  },

  _init_with_url: function(url, options) {
    this._tile_layer = L.tileLayer(url);
    var json_url = this._jsonUrl(url);
    if (!options) {
      options = {};
    }
    options['sql_url'] = '';
    this._grid_layer = new nakshamap.Grid(json_url, options);
    this._layers = {};
    this.addLayer(this._tile_layer);
    this.addLayer(this._grid_layer);
    this._ready = true;
  }
});

nakshamap.layer = function(url, options) {
  return new nakshamap.Layer(url, options);
};

nakshamap.Grid = L.UtfGrid.extend({
  _data: {},
  _data_ajax_queue: {},
  _has_mouse_event: false,
  options: {
    table: '',
    info_fields: []
  },

  hasMouseEvent: function(yes_or_no) {
    this._has_mouse_event = yes_or_no;
  },

  _click: function(e) {
    var naksha_id = this._getNakshaId(e);
    if (naksha_id !== null) {
      this._getDataAndFireEvent(e, naksha_id, 'click');
    }
  },

  _move: function (e) {
    var naksha_id = this._getNakshaId(e);

    if (naksha_id !== this._mouseOn) {
      if (this._mouseOn) {
        if (this._has_mouse_event) {
          this._fireNakshaEvent(e, naksha_id, 'mouseout');
        }
        if (this.options.pointerCursor) {
          this._container.style.cursor = '';
        }
      }
      if (naksha_id) {
        if (this._has_mouse_event) {
          this._getDataAndFireEvent(e, naksha_id, 'mouseover');
        }
        if (this.options.pointerCursor) {
          this._container.style.cursor = 'pointer';
        }
      }

      this._mouseOn = naksha_id;
    } else if (naksha_id) {
      if (this._has_mouse_event) {
        this._getDataAndFireEvent(e, naksha_id, 'mousemove');
      }
    }
  },

  _getNakshaId: function(e) {
    var map = this._map,
        point = map.project(e.latlng),
        tileSize = this.options.tileSize,
        resolution = this.options.resolution,
        x = Math.floor(point.x / tileSize),
        y = Math.floor(point.y / tileSize),
        gridX = Math.floor((point.x - (x * tileSize)) / resolution),
        gridY = Math.floor((point.y - (y * tileSize)) / resolution),
        max = map.options.crs.scale(map.getZoom()) / tileSize;

    x = (x + max) % max;
    y = (y + max) % max;

    var data = this._cache[map.getZoom() + '_' + x + '_' + y];
    var naksha_id = null;
    if (data && data.grid) {
        var idx = this._utfDecode(data.grid[gridY].charCodeAt(gridX)),
            key = data.keys[idx];

        if (key.length > 0) {
          naksha_id = key;
        }
    }

    return naksha_id;
  },

  _getDataAndFireEvent: function(e, naksha_id, to_fire_event) {
    if (this.options.table.length === 0 || this.options.info_fields.length === 0) {
      this._fireNakshaEvent(e, naksha_id, to_fire_event);
      return;
    }

    if (naksha_id in this._data) {
      this._fireNakshaEvent(e, naksha_id, to_fire_event);
      return;
    }
    if (naksha_id in this._data_ajax_queue) {
      return;
    }

    this._data_ajax_queue[naksha_id] = 'in-queue';
    var sql = 'SELECT ' + this.options.info_fields.join(', ') + ' FROM ' + this.options.table + ' WHERE naksha_id = ' + naksha_id;
    var that = this;
    nakshamap.Util.SQL(this.options.sql_url, sql, function(data) {
      that._data[naksha_id] = data;
      that._fireNakshaEvent(e, naksha_id, to_fire_event);
      delete that._data_ajax_queue[naksha_id];
    }, function(err) {
      console.log('Error: ' + err.join('<br />'));
      delete that._data_ajax_queue[naksha_id];
    });
  },

  _fireNakshaEvent: function(e, naksha_id, to_fire_event) {
    if (this.options.table.length === 0 || this.options.info_fields.length === 0) {
      e = L.Util.extend({naksha_id: naksha_id, data: {}}, e);
      this.fire(to_fire_event, e);
    }
    else {
      var data = this._data[naksha_id];
      nakshamap.Util.showInfo(data);
      e = L.Util.extend({naksha_id: naksha_id, data: data}, e);
      this.fire(to_fire_event, e);
    }
  }
});

