var ExportStatus = (function() {
  var thisClass = {};
  var _export_ids = [];
  var _checking = false;

  function export_error_message(id, msg) {
    add_export_message_box_iff(id);
    if (msg instanceof Array) {
      msg = msg.join('<br />');
    }
    var id_str = '#' + get_export_box_id(id);
    $(id_str).children('p').first().html(msg);
  }

  function export_success_message(id, msg) {
    add_export_message_box_iff(id);
    if (msg instanceof Array) {
      msg = msg.join('<br />');
    }
    var id_str = '#' + get_export_box_id(id);
    $(id_str).children('p').first().html(msg);
  }

  function add_export_message_box_iff(id) {
    var div, p, span, id_str;

    id_str = get_export_box_id(id);
    if (document.getElementById(id_str)) {
      return;
    }
    div = $('<div>').addClass('exim-message-box').attr('id', id_str);
    span = $('<span>').html('X').addClass('close');
    $(div).append(span);
    p = $('<p>');
    $(div).append(p);

    $('#export-message-container').append(div);
    $('#export-message-holder').removeClass('disp-none');
  }

  function get_export_box_id(id) {
    return 'export-message-box-' + id;
  }

  function check_export_status() {
    var id, url;

    id = get_next_export_id();
    if (id === 0) {
        _checking = false;
        return;
    }

    _checking = true;
    url = '/exports/' + id + '/status';

    $.ajax({
      type: 'GET',
      cache: false,
      url: url,
      dataType: 'json'
    }).done(function(rt) {
      if (rt['status']) {
        if (rt['status'] === 'success') {
          var export_status = parseInt(rt['export_status']);
          switch (export_status) {
            case -10:
              export_error_message(id, rt['export_name'] + ': Export failed');
              remove_id_from_export_ids(id);
              check_export_status();
            break;
            case 0:
              setTimeout(check_export_status, 5000);
            break;
            case 10:
              var msg = rt['export_name'] + ': Export completed. <a href="' + rt['download_url'] + '" target="_blank">Download Here</a>';
              export_success_message(id, msg);
              remove_id_from_export_ids(id);
              check_export_status();
            break;
          }
        }
        else {
          export_error_message(id, rt['errors']);
        }
      }
      else {
        export_error_message(id, 'Invalid response from server');
      }
    }).fail(function(rt, st) {
      export_error_message(id, 'Request failed');
    });
  }

  function show_loading_div() {
    $('#loading-div').removeClass('disp-none');
  }

  function hide_loading_div() {
    $('#loading-div').addClass('disp-none');
  }

  function show_message(msg) {
    if (msg instanceof Array) {
      msg = msg.join('<br />');
    }
    $('#msg-div').html(msg);
    $('#msg-div').removeClass('disp-none');
  }

  function hide_message() {
    $('#msg-div').html('');
    $('#msg-div').addClass('disp-none');
  }

  function get_next_export_id() {
    if (_export_ids.length > 0) {
      return parseInt(_export_ids[0]);
    }
    else {
      return 0;
    }
  }

  function remove_id_from_export_ids(id) {
    var tmp, i, tmp_id;

    id = parseInt(id);
    tmp = [];
    for (i=0; i<_export_ids.length; ++i) {
      tmp_id = parseInt(_export_ids[i]);
      if (id !== tmp_id) {
        tmp.push(tmp_id)
      }
    }

    _export_ids = tmp;
  }

  function close_export_message_div() {
    $(this).parent('div').remove();
    if ($('#export-message-container').children('div').length == 0) {
      $('#export-message-holder').addClass('disp-none');
    }
  }

  function toggle_export_message_container() {
    if ($('#export-message-container').hasClass('disp-none')) {
      $('#export-message-container').removeClass('disp-none');
    }
    else {
      $('#export-message-container').addClass('disp-none');
    }
  }

  function set_event_listeners() {
    $('#export-message-container').on('click', 'span.close', close_export_message_div);
    $('#export-message-holder').on('click', 'h5', toggle_export_message_container);
  }

  thisClass.init = function(ids) {
    $('#frm_export').ajaxForm({
      dataType: 'json',
      beforeSubmit: function() {
        show_loading_div();
        show_message('Starting...');
      },
      success: function(rt) {
        if (rt['status']) {
          if (rt['status'] === 'success') {
            show_message('Export started. You will be notified once it is completed.');
            hide_loading_div();	
            _export_ids.push(rt['id']);
            if (!_checking) {
              export_success_message(rt['id'], 'Exporting ...');
              check_export_status();
            }
          }
          else {
            show_message(rt['errors']);
            hide_loading_div();	
          }
        }
        else {
          show_message('Invalid response from server');
          hide_loading_div();	
        }
      },
      error: function(rt, st) {
        show_message('Request failed');
        hide_loading_div();	
      }
    });
  };

  thisClass.checkExports = function(ids) {
    set_event_listeners();
    _export_ids = ids;
    if (ids.length > 0) {
      $('#export-message-container').removeClass('disp-none');
      check_export_status();
    }
  };

  return thisClass;
})();

var ImportStatus = (function() {
  var thisClass = {};
  var _import_ids = [];
  var _checking = false;

  function import_error_message(id, msg) {
    add_import_message_box_iff(id);
    if (msg instanceof Array) {
      msg = msg.join('<br />');
    }
    var id_str = '#' + get_import_box_id(id);
    $(id_str).children('p').first().html(msg);
  }

  function import_success_message(id, msg) {
    add_import_message_box_iff(id);
    if (msg instanceof Array) {
      msg = msg.join('<br />');
    }
    var id_str = '#' + get_import_box_id(id);
    $(id_str).children('p').first().html(msg);
  }

  function add_import_message_box_iff(id) {
    var div, p, span, id_str;

    id_str = get_import_box_id(id);
    if (document.getElementById(id_str)) {
      return;
    }
    div = $('<div>').addClass('exim-message-box').attr('id', id_str);
    span = $('<span>').html('X').addClass('close');
    $(div).append(span);
    p = $('<p>');
    $(div).append(p);

    $('#import-message-container').append(div);
    $('#import-message-holder').removeClass('disp-none');
  }

  function get_import_box_id(id) {
    return 'import-message-box-' + id;
  }

  function check_import_status() {
    var id, url;

    id = get_next_import_id();
    if (id === 0) {
        _checking = false;
        return;
    }

    _checking = true;
    url = '/tables/' + id + '/status';

    $.ajax({
      type: 'GET',
      cache: false,
      url: url,
      dataType: 'json'
    }).done(function(rt) {
      if (rt['status']) {
        if (rt['status'] === 'success') {
          switch (rt['import_status']) {
            case 'importing':
              import_success_message(id, 'Importing ....');
              setTimeout(check_import_status, 5000);
            break;
            case 'success':
              var msg = rt['import_name'] + ': Import completed. <a href="' + rt['table_url'] + '">View table</a>';
              import_success_message(id, msg);
              remove_id_from_import_ids(id);
              check_import_status();
            break;
            case 'error':
              import_error_message(id, rt['import_name'] + ': Import failed');
              remove_id_from_import_ids(id);
              check_import_status();
            break;
          }
        }
        else {
          import_error_message(id, rt['errors']);
        }
      }
      else {
        import_error_message(id, 'Invalid response from server');
      }
    }).fail(function(rt, st) {
      import_error_message(id, 'Request failed');
    });
  }

  function get_next_import_id() {
    if (_import_ids.length > 0) {
      return parseInt(_import_ids[0]);
    }
    else {
      return 0;
    }
  }

  function remove_id_from_import_ids(id) {
    var tmp, i, tmp_id;

    id = parseInt(id);
    tmp = [];
    for (i=0; i<_import_ids.length; ++i) {
      tmp_id = parseInt(_import_ids[i]);
      if (id !== tmp_id) {
        tmp.push(tmp_id)
      }
    }

    _import_ids = tmp;
  }

  function close_import_message_div() {
    $(this).parent('div').remove();
    if ($('#import-message-container').children('div').length == 0) {
      $('#import-message-holder').addClass('disp-none');
    }
  }

  function toggle_import_message_container() {
    if ($('#import-message-container').hasClass('disp-none')) {
      $('#import-message-container').removeClass('disp-none');
    }
    else {
      $('#import-message-container').addClass('disp-none');
    }
  }

  function set_event_listeners() {
    $('#import-message-container').on('click', 'span.close', close_import_message_div);
    $('#import-message-holder').on('click', 'h5', toggle_import_message_container);
  }

  thisClass.trackImport = function(id) {
    _import_ids.push(parseInt(id));
    if (!_checking) {
      import_success_message(id, 'Importing ...');
      check_import_status();
    }
  };

  thisClass.checkImports = function(ids) {
    set_event_listeners();
    _import_ids = ids;
    if (ids.length > 0) {
      $('#import-message-container').removeClass('disp-none');
      check_import_status();
    }
  };

  return thisClass;
})();

