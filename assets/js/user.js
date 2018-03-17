var UserProfile = (function() {
  var thisClass = {};

  function set_form() {
    var options = {
      form: '#frm_change_password',
      success_callback: function(rt) {
        naksha.Message.success($('#frm_change_password'), rt['message']);
      }
    };
    new naksha.AjaxForm(options);
  }

  function set_event_listeners() {
    $('#btn_google_maps').bind('click', update_maps_key);
    $('#btn_bing_maps').bind('click', update_maps_key);
  }

  function update_maps_key() {
    var inp, key, value, par_div, options;

    inp = $(this).parent('div').children('input').first();
    key = $(inp).attr('name');
    value = $(inp).val();
    par_div = $(this).parent('div').parent('div');

    options = {
      url: window.location.href,
      type: 'POST',
      data: {key: key, value: value},
      msg_div: par_div,
      success_message: 'Value updated'
    };
    new naksha.Ajax(options);
  }

  thisClass.init = function() {
    set_form();
    set_event_listeners();
  };

  return thisClass;
})();

