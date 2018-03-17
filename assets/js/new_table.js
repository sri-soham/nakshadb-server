var NewTable = (function() {
  var thisClass = {};

  function set_upload_form() {
    var options = {
      form: '#frm_new_table',
      success_callback: function(rt) {
        naksha.Alert.show('Import started. You will be notified after completion.');
        ImportStatus.trackImport(rt['id']);
      }
    };
    new naksha.AjaxForm(options);
  }

  function set_empty_form() {
    var options = {
      form: '#frm_empty_table',
      success_callback: function(rt) {
        window.location.href = rt['url'];
      }
    };
    new naksha.AjaxForm(options);
  }

  thisClass.init = function() {
    $('#accordion').accordion();
    set_upload_form();
    set_empty_form();
  };

  return thisClass;
})();

