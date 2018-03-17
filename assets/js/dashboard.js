var Dashboard = (function() {
  var thisClass = {};

  function handle_delete_table(e) {
    e.preventDefault();
    e.stopPropagation();
    var options = {
      title: 'Delete Table',
      message: 'Are you sure that you want to delete this table?',
      yes_callback: delete_table,
      yes_params: {
        div: $(this).parent('p').parent('div'),
        href: $(this).attr('href')
      }
    };
    new naksha.ConfirmBox(options);
  }

  function delete_table(data) {
    var div = data['div'];
    var p = $(div).children('p').first();
    var options = {
      url: data['href'],
      type: 'POST',
      msg_div: p,
      success_callback: function(rt) {
        $(div).fadeOut(1000, function() {
          $(div).remove();
        });
      }
    };
    new naksha.Ajax(options);
  }

  function set_event_listeners() {
    $('.del-table').on('click', handle_delete_table);
  }

  thisClass.init = function() {
    set_event_listeners();
  };

  return thisClass;
})();

