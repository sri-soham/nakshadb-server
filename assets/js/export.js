var Export = (function() {
  var thisClass = {};
  var _export_id;
  
  function delete_export(data) {
    var url, td, span, msg_td, export_id;

    id = '#' + data['id'];
    export_id = id.split('_').pop();
    td = $(id).parent('td');
    msg_td = $(id).parent('td').parent('tr').children('td').first();
    span = $('<span>').html('.....');
    $(td).append(span);
    $(id).addClass('disp-none');

    var url = '/exports/' + export_id + '/delete';
    var options = {
      url: url,
      type: 'POST',
      msg_div: msg_td,
      success_callback: function(rt) {
        var tr = $(td).parent('tr');
        $(tr).fadeOut(2000, function() {
          $(tr).remove();
        });
      },
      error_callback: function(msg) {
        restore_link(td);
        naksha.Message.error($(msg_td), msg);
      }
    };
    naksha.Ajax(options);
  }

  function restore_link(td) {
    $(td).children('span').remove();
    $(td).children('a').removeClass('disp-none');
  }

  function handle_delete_export(e) {
    e.preventDefault();
    e.stopPropagation();

    var options = {
      title: 'Delete  Export',
      message: 'Are you sure that you want to delete this export?',
      yes_callback: delete_export,
      yes_params: {id: $(this).attr('id')}
    };
    new naksha.ConfirmBox(options);
  }

  thisClass.init = function() {
    $('#exports-list').on('click', '.del-export', handle_delete_export);
  };

  return thisClass;
})();

