$(function() {
  var btnVerify = $('#btnVerify')
  btnVerify.click(function() {
    btnVerify.attr("disabled", true).html('Verifying <i class="fas fa-fw fa-spinner fa-spin"></i>')
    $.get('/domain/verify/' + $(this).attr('data-id'), function(resp) {
      if (resp.data && resp.data.verified) {
        btnVerify.removeClass('btn-warning').addClass('btn-success').text("Verified");
        btnVerify.off('click');
      } else {
        btnVerify.text("Verify");
        notify(resp.message, "error");
      }
      btnVerify.removeAttr("disabled");
    }, 'json');
  });

  btnVerify.trigger('click');
});
