var TableAdmin = (function() {
  var thisClass = {};

  function set_infowindow(infowindow) {
    var i, col, div, input, span, columns, options;

    infowindow = $.parseJSON(infowindow);
    columns = Table.getColumns();
    for (i in columns) {
      col = columns[i];
      if (col === 'the_geom' || col === 'the_geom_webmercator') {
        continue;
      }

      input = $('<input>').attr({
        'type': 'checkbox',
        'name': 'columns',
        'value': col
      });
      if ($.inArray(col, infowindow['fields']) > -1) {
        $(input).prop('checked', true);
      }
      div = $('<div>').addClass('infowindow-column');
      $(div).append(input);
      span = $('<span>').addClass('ml10').html(col);
      $(div).append(span);
      $('#infowindow-columns').append(div);
    }

    options = {
      form: '#frm_infowindow',
      success_message: 'Updated infowindow fields'
    };
    new naksha.AjaxForm(options);
  }

  function set_api_access_form() {
    var options = {
      form: '#frm_api_access',
      success_message: 'Api access modified'
    };
    new naksha.AjaxForm(options);
  }

  thisClass.init = function(infowindow) {
    set_infowindow(infowindow);
    set_api_access_form();
  };

  thisClass.removeColumnFromInfowindow = function(column) {
    var input, div;

    input = $('input[type="checkbox"][name="columns"][value="' + column + '"]');
    if (input) {
      div = $(input).parent('div.infowindow-column');
      $(div).remove();
    }
  };

  return thisClass;
})();

