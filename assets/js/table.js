var Table = (function() {
  var thisClass = {};
  var _url, _page, _columns, _del_col_url, _order_column, _order_type;
  var _row_id = 0, _col_name = '';
  var _displayed_rows = 0;

  function get_row_id(id) {
    return 'naksha-' + id;
  }

  function error_message(msg) {
    var div = $('<div>').css('color', '#f00');
    if (msg instanceof Array) {
      msg = msg.join('<br />');
    }
    
    $(div).html(msg);
    $('#geom-message').append(div);
    $('#geom-message').removeClass('disp-none');
    setTimeout(clear_message, 10000);
  }

  function clear_message() {
    $('#geom-message').html('');
    $('#geom-message').addClass('disp-none');
  }

  function fetch_data() {
    _page++;
    var options = {
      type: 'GET',
      data: {'order_column': _order_column, 'order_type': _order_type},
      url: _url + 'data/' + _page,
      success_callback: function(rt) {
        if (_displayed_rows === 0) {
          $('#gt-body').html('');
        }
        add_data(rt);
      },
      error_callback: error_message
    };
    new naksha.Ajax(options);
  }

  function add_data(data) {
    var rows, i, j, row, col, tr, td, geom_type, count, id;

    for (i=0; i<data['rows'].length; ++i) {
      row = data['rows'][i];
      // When a feature is added on map or a new record is added in table view
      // a new row is added to the table view which may be out of order.
      // Out of order in the sense: Rows are fetched and displayed in batches of 40.
      // If a table has 75 rows and first 40 are being displayed now. A row is added
      // (either through drawing-feature-on-map or adding-row-on-table), it will show
      // up in display with id 76. Now, when the next set of 40 is fetched, this 76th
      // will be fetched and added again; effectively two rows displayed for one row
      // in the database. To avoid this, delete the row with id, if it already exists.
      id = get_row_id(row['naksha_id']);
      if (document.getElementById(id)) {
        $('#' + id).remove();
      }
      tr = $('<tr>').attr('id', id);
      for (j=0; j<_columns.length; ++j) {
        col = _columns[j];
        td = $('<td>').attr('data-col', col);
        switch (col) {
          case 'naksha_id':
            $(td).addClass('del-row').append(row[col]);
          break;
          case 'the_geom':
            geom_type = geom_type_from_val(row[col]);
            $(td).addClass('upd-cell').attr('data-geom', row[col]).html(geom_type);
          break;
          case 'created_at':
          case 'updated_at':
            $(td).html(row[col]);
          break;
          default:
            $(td).addClass('upd-cell').html(row[col]);
          break;
        }
        $(tr).append(td);
      }
      $('#gt-body').append(tr);
    }
    _displayed_rows += data['rows'].length;
    count = parseInt(data['count']);
    if (count > _displayed_rows) {
      $('#more-rows-td').removeClass('disp-none');
    }
    else {
      $('#more-rows-td').addClass('disp-none');
    }
  }

  function geom_type_from_val(val) {
    var part, parts;

    if (val.length === 0) {
      part = '';
    }
    else {
      part = val.split('(')[0];
      parts = part.split(';');
      if (parts.length === 2) {
        part = parts[1];
      }
      else {
        part = '';
      }
    }

    return part;
  }

  function set_headers(columns) {
    var i, tr, th, rem, span;

    columns = columns.split(',');
    rem = ['naksha_id', 'the_geom', 'created_at', 'updated_at'];
    _columns = $(columns).not(rem).get();
    _columns.unshift('the_geom');
    _columns.unshift('naksha_id');
    _columns.push('created_at');
    _columns.push('updated_at');

    tr = $('<tr>');
    for (i=0; i<_columns.length; ++i) {
      th = $('<th>').html(desc_from_name(_columns[i])).attr('data-col', _columns[i]);
      span = $('<span>').html('&dtrif;').addClass('pl5 curs-pointer sorter');
      switch (_columns[i]) {
        case 'naksha_id':
        case 'the_geom':
        case 'created_at':
        case 'updated_at':
        break;
        default:
          $(th).addClass('can-del-col');
        break;
      }
      if (_columns[i] != 'the_geom') {
        $(th).append(span);
      }
      $(tr).append(th);
    }
    $('#gt-header').append(tr);
  }

  function desc_from_name(col) {
    var desc, parts, i, part;

    desc = col.replace(/[_ -]+/g, ' ');
    parts = desc.split(' ');
    desc = '';
    for (i in parts) {
      part = $.trim(parts[i]);
      if (part.length > 0) {
        desc += parts[i][0].toUpperCase() + parts[i].substr(1) + ' ';
      }
    }
    desc = $.trim(desc);

    return desc;
  }

  function set_add_column_form() {
    var options = {
      form: '#frm_add_column',
      before_submit: function() {
        var errors = [];
        var frm = document.getElementById('frm_add_column');
        if (frm.elements['name'].value.length === 0) {
          errors.push('Please enter the name');
        }
        if (frm.elements['data_type'].value.length === 0) {
          errors.push('Please select the data type');
        }
        if (errors.length > 0) {
          naksha.Message.error($('#frm_add_column'), errors);
          return false;
        }
      },
      success_callback: function(rt) {
        var name = document.getElementById('frm_add_column').elements['name'].value;
        add_column(name);
        close_add_column_form();
      }
    };
    new naksha.AjaxForm(options);
  }

  function add_column(name) {
    var desc, td, new_td;

    desc = desc_from_name(name);
    td = $('#gt-header').children('tr').first().children('th').last().prev('th');
    new_td = $('<th>').html(desc).attr('data-col', name).addClass('can-del-col');
    $(new_td).insertBefore($(td));
    $('#gt-body').children('tr').each(function() {
      td = $(this).children('td').last().prev('td');
      new_td = $('<td>').attr('data-col', name).html('').addClass('upd-cell');
      $(new_td).insertBefore($(td));
    });
  }

  function close_add_column_form() {
    var frm = document.getElementById('frm_add_column');
    $('#geom-add-column').addClass('disp-none');
    frm.elements['name'].value = '';
    frm.elements['data_type'].value = '';
  }

  function show_add_column_form() {
    $('#geom-add-column').removeClass('disp-none');
    $('#geom-add-column').find('input[name="name"]').focus();
  }

  function handle_delete_row() {
    if (_row_id > 0) return;

    _row_id = $(this).parent('tr').attr('id').split('-').pop();
    var options = {
      title: 'Delete Row',
      message: 'Are you sure that you want to delete this row?',
      yes_callback: delete_row,
      no_callback: function(data) {
        _row_id = 0;
      }
    };
    new naksha.ConfirmBox(options);
  }

  function delete_row() {
    var url = _url + 'delete/' + _row_id;
    var options = {
      url: _url + 'delete/' + _row_id,
      type: 'POST',
      success_callback: function(rt) {
        var tr = '#' + get_row_id(_row_id);
        $(tr).remove();
        Map.geometryUpdated(rt['update_hash']);
        _row_id = 0;
      },
      error_callback: function(rt) {
        _row_id = 0;
        error_message(rt);
      }
    };
    new naksha.Ajax(options);
  }

  function handle_delete_column() {
    var col_name = $(this).attr('data-col');
    var name = $(this).html();
    var pos = name.indexOf('<span');
    name = name.substr(0, pos);
    name = $.trim(name);
    var options = {
      title: 'Delete column',
      message: 'Are you sure that you want to delete <br />"' + name + '"?',
      yes_callback: delete_column,
      yes_params: {column_name: col_name}
    };
    new naksha.ConfirmBox(options);
  }

  function delete_column(data) {
    var div = $('#geom-delete-column').children('div').first().children('div').last();
    var column_name = data['column_name'];
    var options = {
      url: _del_col_url,
      data: {'column_name': column_name},
      type: 'POST',
      success_callback: function(rt) {
        remove_column(column_name);
        var is_part = TableAdmin.removeColumnFromInfowindow(column_name);
        if (is_part) {
          var msg = 'This column is part of infowindow columns. In the "Info Window" tab ' +
              'uncheck the box for the deleted column and click "Save"';
          naksha.Alert.show(msg);
        }
      }
    };
    new naksha.Ajax(options);
  }

  function remove_column(col_name) {
    var pos, i, td;

    // remove the th and td's that contain the column values
    i = 0;
    $('#gt-header').children('tr').first().children('th').each(function() {
      ++i;
      if ($(this).attr('data-col') == col_name) {
        pos = i;
      }
    });
    td = $('#gt-header').children('tr').first().children('th:nth-child(' + pos + ')');
    $(td).remove();
    $('#gt-body').children('tr').each(function() {
      td = $(this).children('td:nth-child(' + pos + ')');
      $(td).remove();
    });

    // remove the column name from the _columns array
    pos = -1;
    for (i=0; i<_columns.length; ++i) {
      if (col_name == _columns[i]) {
        pos = i;
        break;
      }
    }
    if (pos > -1) {
      _columns.splice(pos, 1);
    }
  }

  function edit_cell() {
    // Only one edit at a time.
    if (_row_id > 0) return;

    var val, inp, desc;

    $(this).addClass('active-cell');
    _row_id = $(this).parent('tr').attr('id').split('-').pop();
    _col_name = $(this).data('col');
    desc = desc_from_name(_col_name);
    $('#geom-edit').children('div').first().html('');
    if (_col_name == 'the_geom') {
      val = $(this).attr('data-geom');
      inp = $('<textarea>').attr({'name': 'updated_value', 'id': 'updated_value'}).addClass('col100p-np row200 padding10').val(val);
    }
    else {
      val = $.trim($(this).html());
      inp = $('<input>').attr({
              'type': 'text',
              'name': 'updated_value',
              'id': 'updated_value'
            }).addClass('col100p-np padding10').val(val);
    }
    $('#geom-edit').children('p').html(desc + ': ');
    $('#geom-edit').children('div').first().append(inp);
    $('#geom-edit').removeClass('disp-none');
    $('#updated_value').focus();
  }

  function update_value() {
    var geom_type;
    var value = $('#updated_value').val();
    var options = {
      type: 'POST',
      url: _url + 'update/' + _row_id,
      data: {'column': _col_name, 'value': value},
      success_callback: function(rt) {
        if (_col_name != 'the_geom') {
          $('#gt-body').find('td.active-cell').html(value);
        }
        else {
          geom_type = geom_type_from_val(value);
          $('#gt-body').find('td.active-cell').attr('data-geom', value).html(geom_type);
          Map.geometryUpdated(rt['update_hash']);
          if ($('#gt-body').children('tr').length == 1) {
            update_geometry_type();
          }
        }
        close_update_form();
      },
      error_callback: error_message
    };
    new naksha.Ajax(options);
  }

  function close_update_form() {
    $('#geom-edit').addClass('disp-none');
    $('#geom-edit').children('div').first().html('');
    $('#gt-body').find('td.active-cell').removeClass('active-cell');
    _row_id = 0;
    _col_name = '';
  }

  function add_row() {
    var url = thisClass.getAddRowUrl();
    var options = {
     url: thisClass.getAddRowUrl(),
     type: 'POST',
     success_callback: add_empty_row,
     error_callback: error_message
    };
    new naksha.Ajax(options);
  }

  function add_empty_row(rt) {
    var row = rt['row'];
    var data = {}, j;

    if ('update_hash' in row) {
        delete row['update_hash'];
    }
    for (j=0; j<_columns.length; ++j) {
      switch (_columns[j]) {
        case 'naksha_id':
        case 'created_at':
        case 'updated_at':
        break;
        default:
          row[_columns[j]] = '';
        break;
      }
    }
    data['rows'] = [];
    data['rows'].push(row)
    data['count'] = 1;
    add_data(data);
  }

  function update_geometry_type() {
    var geom_type, td, type, data_geom;

    td = $('#gt-body').children('tr').first().children('td:nth-child(2)');
    type = $(td).html();
    type = $.trim(type);
    switch (type) {
      case 'POINT':           geom_type = 'point';      break;
      case 'MULTILINESTRING': geom_type = 'linestring'; break;
      case 'MULTIPOLYGON':    geom_type = 'polygon';    break;
      default:                geom_type = 'unknown';    break;
    }
    data_geom = $(td).attr('data-geom');
    Map.updateGeometryType(geom_type, data_geom);
  }

  function show_order_by_div() {
    var th, div, offset, left, ttop, span, col;

    th = $(this).parent('th');
    col = $(th).attr('data-col');
    if (document.getElementById('sort_by_div')) {
        if ($('#sort_by_div').attr('data-col') == col) {
          $('#sort_by_div').remove();
          return;
        }
        else {
          $('#sort_by_div').remove();
        }
    }

    offset = $(th).offset();
    div = $('<div>');
    $(div).append('Order: ');
    span = $('<span>').attr('data-order', 'asc')
                      .html('asc')
                      .css({
                        'padding': '5px',
                        'cursor': 'pointer',
                        'margin-left': '5px',
                        'border': '1px #333 solid'
                      });
    if (_order_column == col && _order_type == 'asc') {
      $(span).css('background-color', '#69f');
    }
    $(span).bind('click', fetch_sorted_data);
    $(div).append(span);
    span = $('<span>').attr('data-order', 'desc')
                      .html('desc')
                      .css({
                        'padding': '5px',
                        'cursor': 'pointer',
                        'margin-left': '5px',
                        'border': '1px #333 solid'
                      });
    if (_order_column == col && _order_type == 'desc') {
      $(span).css('background-color', '#69f');
    }
    $(span).bind('click', fetch_sorted_data);
    $(div).append(span);

    left = parseInt(offset.left) + 'px';
    ttop = parseInt(offset.top + (2 * $(th).height())) + 'px';
    $(div).css({
      'left': left,
      'top': ttop,
      'z-index': '10',
      'position': 'absolute',
      'background-color': '#ff9',
      'border': '1px #000 solid',
      'padding': '20px'
    }).attr({
      'id': 'sort_by_div',
      'data-col': col
    });
    $(document.body).append(div);
  }

  function fetch_sorted_data() {
    _order_type = $(this).attr('data-order');
    _order_column = $(this).parent('div').attr('data-col');
    _page = 0;
    _displayed_rows = 0;
    $('#sort_by_div').remove();
    fetch_data();
  }

  function set_event_listeners() {
    $('#geom-message').bind('click', clear_message);
    $('#gt-body').on('click', '.del-row', handle_delete_row);
    $('#gt-body').on('dblclick', '.upd-cell', edit_cell);
    $('#btn_update').on('click', update_value);
    $('#btn_cancel').on('click', close_update_form);
    $('#btn_fetch_more').on('click', fetch_data);
    $('#btn_add_row').on('click', add_row);
    $('#btn_delete').on('click', delete_row);
    $('#btn_cancel_add_column').on('click', close_add_column_form);
    $('#btn_add_column').on('click', show_add_column_form);
    $('#gt-header').on('dblclick', 'th.can-del-col', handle_delete_column);
    $('#gt-header').on('click', 'span.sorter', show_order_by_div);
  }

  thisClass.init = function(url, del_col_url, columns, infowindow) {
    _url = url;
    _del_col_url = del_col_url;
    _page = 0;
    set_event_listeners();
    $('#tabs').tabs({
      activate: function(evt, ui) {
        if ($(ui.newPanel).attr('id') == 'maps-tab') {
          Map.init();
        }
      }
    });
    set_headers(columns);
    set_add_column_form();
    _order_column = 'naksha_id';
    _order_type = 'asc';
    fetch_data();
  };

  thisClass.errorMessage = function(msg) {
    error_message(msg);
  };

  thisClass.deleteRowOfGeometry = function(naksha_id) {
    var id = get_row_id(naksha_id);
    tr = $('#' + id);
    $(tr).remove();
  };

  thisClass.getAddRowUrl = function() {
    return _url + 'add';
  };

  thisClass.getShowRowUrl = function(id) {
    return _url + 'show/' + id;
  };

  thisClass.getColumns = function() {
    return _columns;
  };

  thisClass.addRowFromMap = function(rt) {
    var data = {};
    data['rows'] = [];
    delete rt['update_hash'];
    data['rows'].push(rt['row']);
    add_data(data);
  };

  return thisClass;
})();

