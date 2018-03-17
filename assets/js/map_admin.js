var MapAdmin = (function() {
  var thisClass = {};
  var _maps_url, _layer_id;

  function set_edit_form() {
    new naksha.AjaxForm({
      form: '#frm_map_edit',
      success_message: 'Details updated'
    });
  }

  function set_delete_form() {
    new naksha.AjaxForm({
      form: '#frm_map_delete',
      success_callback: function(rt) {
        window.location.href = rt['redir_url'];
      }
    });
  }

  function set_base_layer_form() {
    new naksha.AjaxForm({
      form: '#frm_map_base_layer',
      success_message: 'Base map updated'
    });
  }

  function set_autocomplete() {
    $('#table_name').autocomplete({
      minLength: 2,
      delay: 500,
      source: function(request, response) {
        var url = _maps_url + '/search_tables?table_name=' + $('#table_name').val();
        var td = $('#table_name').parent('td');
        var options = {
          url: url,
          type: 'GET',
          msg_div: td,
          success_callback: function(rt) {
            return response(rt['tables']);
          }
        };
        new naksha.Ajax(options);
      },
      select: function(evt, ui) {
        _layer_id = ui.item.layer_id;
      }
    });
  }

  function add_row() {
    var tr, td, input;

    tr = $('<tr>');
    td = $('<td>').addClass('pos-rel').html($('#table_name').val());
    $(tr).append(td);
    td = $('<td>').addClass('col100');
    input = $('<input>').attr({
      'type': 'button',
      'name': 'btn_delete_' + _layer_id,
      'value': 'Delete Layer'
    }).addClass('btn-red del-layer')
    $(td).append(input);
    $(tr).append(td);
    $('#existing-tables').append(tr);
  }

  function set_event_listeners() {
    $('#btn_add_layer').bind('click', add_layer_to_map);
    $('#existing-tables').on('click', '.del-layer', delete_layer_from_map);
  }

  function add_layer_to_map() {
    var td, options;

    td = $('#table_name').parent('td');
    options = {
      url: _maps_url + '/add_layer',
      type: 'POST',
      data: {'layer_id': _layer_id},
      msg_div: $(td),
      success_callback: function(rt) {
        add_row();
        naksha.Message.success($(td), 'Layer added');
        _layer_id = null;
        $('#table_name').val('');
      }
    };
    new naksha.Ajax(options);
  }

  function delete_layer_from_map() {
    var layer_id, td, tr, options;

    layer_id = $(this).attr('name').split('_').pop();
    td = $(this).parent('td').prev('td');
    tr = $(td).parent('tr');
    options = {
      url: _maps_url + '/delete_layer',
      type: 'POST',
      data: {'layer_id': layer_id},
      msg_div: td,
      success_callback: function(rt) {
        $(tr).fadeOut(2000, function() {
          $(tr).remove();
        });
      }
    };
    new naksha.Ajax(options);
  }

  function set_url_form() {
    var options = {
      form: '#frm_map_hash',
      success_callback: function(rt) {
        naksha.Message.success($('#frm_map_hash'), 'URL updated');
        var parts = $('#view-map-link').attr('href').split('/');
        var parts1 = parts.pop();
        parts1 = parts1.split('-');
        parts.push(parts1[0] + '-' + document.getElementById('frm_map_hash').elements['hash'].value)
        $('#view-map-link').attr('href', parts.join('/'));
      }
    };
    new naksha.AjaxForm(options);
  }

  thisClass.init = function(map_url) {
    _maps_url = map_url;
    set_edit_form();
    set_delete_form();
    set_base_layer_form();
    set_autocomplete();
    set_event_listeners();
    set_url_form();
  };

  return thisClass;
})();

