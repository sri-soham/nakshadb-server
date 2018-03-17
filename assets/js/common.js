var naksha = {};

naksha.MsgBox = function(box_id) {
    if (box_id[0] !== '#') box_id = '#' + box_id;
    var _box = $(box_id);
    $(_box).css({
      'box-sizing': 'border-box',
      'position': 'relative'
    });

    function remove_message() {
        $(_box).html('');
        $(_box).css('display', 'none');
        return false;
    }

    function fill_box(msg, is_success) {
        if (msg instanceof Array) {
          msg = msg.join('<br />');
        }
        if (is_success === undefined) is_success = true;
        var div = $('<div>');
        $(div).append(msg);
        if (is_success) {
          $(div).css({'background-color': 'darkgreen', 'color': '#fff'});
        }
        else {
          $(div).css({'background-color': 'darkred', 'color': '#fff'});
        }
        $(div).css({
          'padding': '5px',
          'margin-bottom': '5px',
          'box-sizing': 'border-box',
          'position': 'absolute',
          'top': '0px',
          'left': '0px',
          'width': '100%'
        });
        $(div).html(msg);
        var span = $('<span>').css({
                    'box-sizing': 'border-box',
                    'padding': '5px',
                    'margin-left': '5px',
                    'position': 'absolute',
                    'right': '0px',
                    'top': '0px',
                    'z-index': '2',
                    'color': '#fff',
                    'cursor': 'pointer'
                   });
        $(span).html('X');
        $(_box).append(span);
        $(_box).append(div);
        $(_box).css('display', 'block');
        $(span).bind('click', remove_message);
    }

    this.clear = function() {
        remove_message();
    };

    this.show_message = function(msg) {
        fill_box(msg, true);
        setTimeout(remove_message, 5000);
    };

    this.show_errors = function(errors) {
        fill_box(errors, false);
    };
}

naksha.Message = (function() {
  var thisClass = {};

  function show_msg(par, msg, bg_color, auto_remove) {
    $(par).find('div.msg__child_').remove();

    var div = $('<div>');
    $(div).css({
      'padding': '5px',
      'margin-bottom': '5px',
      'box-sizing': 'border-box',
      'position': 'absolute',
      'top': '0px',
      'left': '0px',
      'width': '100%',
      'background-color': bg_color,
      'color': '#fff'
    });
    $(div).addClass('msg__child_');
    var msg_str;
    if (msg instanceof Array) {
      msg_str = msg.join('<br />');
    }
    else {
      msg_str = msg;
    }
    var span = $('<span>').css({
      'float': 'right',
      'padding': '0px 2px 2px 0px',
      'margin-left': '5px',
      'font-weight': 'bold',
      'cursor': 'pointer',
      'font-size': '110%',
      'font-style': 'normal'
    }).text('X')
    $(span).bind('click', function() {
      $(div).remove();
    });
    $(div).append(span);
    $(div).append(msg_str);
    /*
    $(div).bind('click', function() {
      $(div).remove();
    });
    */
    $(par).append(div);
    if (auto_remove) {
      setTimeout(function() { $(div).remove(); }, 10000);
    }
  }

  thisClass.success = function(par, msg) {
    show_msg(par, msg, 'darkgreen', true);
  };

  thisClass.error = function(par, msg) {
    show_msg(par, msg, 'darkred', false);
  };

  return thisClass;
})();

naksha.FlashMessage = (function() {
  var thisClass = {};

  thisClass.show = function(msg) {
    $('#flash-message-content').html(msg);
    $('#flash-message').css('display', 'block');
    setTimeout(naksha.FlashMessage.hide, 10000);
  };

  thisClass.hide = function() {
    $('#flash-message-content').html('');
    $('#flash-message').css('display', 'none');
  };

  $(document).ready(function() {
    $('#flash-message-close').bind('click', naksha.FlashMessage.hide);
  });

  return thisClass;
})();

naksha.Alert = (function() {
  var thisClass = {};

  thisClass.show = function(message, title, width) {
    var div, div1, h4, input, margin_diff;

    if (width === undefined) {
      width = 300;
    }
    margin_diff = parseInt(width / 2);
    margin_diff = 0 - margin_diff;
    div = $('<div>').css({
      'position': 'fixed',
      'width': width + 'px',
      'left': '50%',
      'top': '200px',
      'margin-left': margin_diff + 'px',
      'border': '2px #fff solid',
      'text-align': 'center',
      'padding-bottom': '20px',
      'border-radius': '10px',
      'background-color': '#000',
      'color': '#fff',
      'z-index': '2000'
    });
    if (title && title.length > 0) {
      h4 = $('<p>').css({
        'font-size': '110%',
        'margin-bottom': '10px',
        'padding': '20px',
        'border-bottom': '1px #333 solid'
      }).html(title);
      $(div).append(h4);

    }
    div1 = $('<div>').css({
      'text-align': 'left',
      'padding': '20px',
      'margin-bottom': '10px'
    }).html(message);
    input = $('<input>').attr({
      'type': 'button',
      'name': 'btn1',
      'value': 'OK'
    }).css({
      'width': '120px',
      'background-color': '#fff',
      'border': '0px',
      'padding': '5px 10px',
      'font-weight': 'bold',
      'color': '#000'
    });
    $(input).bind('click', function() {
      $(this).parent('div').remove();
    });
    $(div).append(div1);
    $(div).append(input);
    $('#container').append(div);
    $(input).focus();
  };

  return thisClass;
})();

/*
  _options.form:
    id of the form for which ajaxForm is to be enabled

  _options.msg_div:
    id of the html element in which messages have to be displayed

  _options.success_callback:
    Function to be called when request is successful. Response from ajax call is
    passed on as first argument to this function. If null, a message 'success'
    will be shown.

  _options.error_callback:
    Function to be called when request fails. JqXHR object, status-text and error-
    thrown are passed as first, second and third arguments respectively to this
    function. If null, a message 'Request failed' will be shown.

  _options.data_type:
    Default is 'json'

  _options.url:
    Default is the 'action' attribute of the form
*/
naksha.AjaxForm = function(options) {
  var _submit_btn, _submit_text, _options;

  _options = {
    form: null,
    msg_div: null,
    success_callback: null,
    error_callback: null,
    success_message: null,
    before_submit: null,
    data_type: 'json',
    url: null
  };
  $.extend(_options, options);

  if (_options.form === null) {
    throw 'Form is required';
    return;
  }

  if (_options.msg_div === null) {
    _options.msg_div = _options.form;
  }

  if (_options.form[0] !== '#') {
    _options.form = '#' + _options.form;
  }

  if (_options.msg_div !== null && $.type(_options.msg_div) === 'string' && _options.msg_div[0] !== '#') {
    _options.msg_div = '#' + _options.msg_div;
  }

  _submit_btn = $(_options.form).find('input[type="submit"]').first();
  _submit_text = $(_submit_btn).val();

  function saving_button() {
    $(_submit_btn).prop('disabled', true);
    $(_submit_btn).val('Processing ...');
  }

  function save_button() {
    $(_submit_btn).prop('disabled', false);
    $(_submit_btn).val(_submit_text);
  }

  function error(err) {
    naksha.Message.error($(_options.msg_div), err);
  }

  function success(msg) {
    naksha.Message.success($(_options.msg_div), msg);
  }

  function init() {
    var form_options = {};
    form_options['dataType'] = _options.data_type;
    if (_options.url !== null) {
      form_options['url'] = _options.url;
    }
    form_options['beforeSubmit'] = function() {
      if (_options.before_submit) {
        var ret = _options.before_submit();
        if (ret === false) {
          return false;
        }
      }
      saving_button();
    };
    form_options['success'] = function(rt) {
      save_button();
      if (rt['status']) {
        if (rt['status'] === 'success') {
          if (_options.success_callback === null) {
            var message = (_options.success_message === null) ? 'Success' : _options.success_message;
            success(message);
          }
          else {
            _options.success_callback(rt);
          }
        }
        else {
          error(rt['errors']);
        }
      }
      else {
        error('Invalid response from server');
      }
    };
    form_options['error'] = function(rt, st, err) {
      save_button();
      if (_options.error_callback === null) {
        error('Request failed');
      }
      else {
        _options.error_callback(rt, st, err);
      }
    };

    $(_options.form).ajaxForm(form_options);
  };

  init();
};

naksha.Ajax = function(options) {
  var _options;

  _options = {
    url: null,
    type: 'POST',
    data: null,
    data_type: 'json',
    success_callback: null,
    error_callback: null,
    success_message: null,
    msg_div: null,
    cache: false
  };
  $.extend(_options, options);

  if (_options.msg_div !== null && $.type(_options.msg_div) === 'string' && _options.msg_div[0] !== '#') {
    _options.msg_div = '#' + _options.msg_div;
  }

  function show_loading_image() {
    $('#loading-image-holder').removeClass('disp-none');
  }

  function hide_loading_image() {
    $('#loading-image-holder').addClass('disp-none');
  }

  function error_message(msg) {
    if (_options.msg_div) {
      naksha.Message.error($(_options.msg_div), msg);
    }
    else {
      if (msg instanceof Array) {
        msg = msg.join('<br />');
      }
      naksha.Alert.show(msg);
    }
  }

  function success_message(msg) {
    if (_options.msg_div) {
      naksha.Message.success($(_options.msg_div), msg);
    }
    else {
      if (msg instanceof Array) {
        msg = msg.join('<br />');
      }
      naksha.Alert.show(msg);
    }
  }

  function init() {
    var ajax_options = {};
    ajax_options['type'] = _options.type;
    ajax_options['url'] = _options.url;
    if (_options.data !== null) {
      ajax_options['data'] = _options.data;
    }
    ajax_options['dataType'] = _options.data_type;
    ajax_options['cache'] = false;

    show_loading_image();
    $.ajax(ajax_options).done(function(rt) {
      hide_loading_image();
      if (rt['status']) {
        if (rt['status'] === 'success') {
          if (_options.success_callback === null) {
            var message = (_options.success_message === null) ? 'Success' : _options.success_message;
            success_message(message);
          }
          else {
            _options.success_callback(rt);
          }
        }
        else {
          if (_options.error_callback) {
            _options.error_callback(rt['errors']);
          }
          else {
            error_message(rt['errors']);
          }
        }
      }
      else {
        if (_options.error_callback) {
          _options.error_callback('Invalid message from server');
        }
        else {
          error_message('Invalid message from server');
        }
      }
    }).fail(function(rt, st) {
      hide_loading_image();
      if (_options.error_callback) {
        _options.error_callback('Request failed');
      }
      else {
        error_message('Request failed');
      }
    });
  };

  init();
};

//
//  yes_params: object containing parameters that have to be passed to yes_callback
//
naksha.ConfirmBox = function(options) {
  var _options, _div;

  _options = {
    title: 'Delete',
    message: 'Are you sure that you want to delete this?',
    yes_callback: null,
    yes_params: {},
    no_callback: null,
    no_params: {}
  };

  $.extend(_options, options);

  if (_options.yes_callback === null) {
    throw 'Callback for "Yes" is required';
    return;
  }

  function show_html() {
    var h4, div1, yes_btn, no_btn, p;

    _div = $('<div>').addClass('modal-div small confirm');
    h4 = $('<h4>').html(_options.title);
    $(_div).append(h4);

    div1 = $('<div>').addClass('pos-rel');
    p = $('<p class="mb10 tCenter">').html(_options.message);
    $(div1).append(p);

    p = $('<p class="mt10 tCenter">');
    yes_btn = $('<input>').attr({
      type: 'button',
      name: '_confirm_btn_yes',
      value: 'Yes'
    }).addClass('mr10');
    $(yes_btn).bind('click', function() {
      $(_div).remove();
      _options.yes_callback(_options.yes_params);
    }).addClass('col50 bd1').css({
      'padding': '5px 10px',
      'border': '2px #fff solid'
    });
    $(p).append(yes_btn);

    no_btn = $('<input>').attr({
      type: 'button',
      name: '_confirm_btn_no',
      value: 'No'
    }).addClass('col50 bd1').css({
      'padding': '5px 10px',
      'border': '2px #fff solid'
    });
    $(no_btn).bind('click', function() {
      $(_div).remove();
      if (_options.no_callback) {
        _options.no_callback(_options.no_params);
      }
    });
    $(p).append(no_btn);

    $(div1).append(p);
    $(_div).append(div1);

    $(document.body).append(_div);
    $(no_btn).focus();
  }

  show_html();
};

