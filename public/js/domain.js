$(function() {
  $('.domain-item[data-verified="false"]').each(function() {
    var item = $(this)
    $.get('/domain/verify/' + $(this).attr('data-id'), function(resp) {
      var classes = '';
      var title = '';
      if (resp.data && resp.data.verified) {
        classes = 'fa-check-circle text-success';
        item.find('.btn-verify').remove();
      } else {
        classes = 'fa-exclamation-circle text-danger';
        title = 'unable to verify domain ownership';
      }
      item.find('.fa-spinner').removeClass('fa-spinner fa-spin text-warning').addClass(classes).attr('title', title)
    }, 'json');
  });
});

function deleteDomain(id) {
    if (!confirm("Are you sure to delete this domain")) {
        return
    }

    $.ajax({
      url: '/domain/' +id,
      method: 'DELETE',
      success: function(resp) {
        if (resp.status != 'success') {
          alert(resp.message);
          return;
        }

        window.location.reload();
      },
    });
}
