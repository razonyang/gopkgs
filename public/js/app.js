$(function() {
    // ajax
    $.ajaxSetup({
        beforeSend: function (xhr)
        {
            // sets CSRF token header.
            xhr.setRequestHeader("X-CSRF-Token", $('meta[name="csrf-token"]').attr('content'));
        },
        error: function(xhr, status, err) {
            console.error(err)
        }
    });

    $('.toast').toast('show');
    $('.toast').on('hidden.bs.toast', function () {
        $(this).remove();
    });
});

function notify(message, type) {
    var title = '<i class="fas fa-check-circle text-success"></i> SUCCESS'
    switch (type) {
        case 'warning':
            title = '<i class="fas fa-exclamation-circle text-warning"></i> WARNING'
            break;
        case 'error':
            title = '<i class="fas fa-exclamation-circle text-danger"></i> ERROR'
            break;
    }
    var eleID = 'toast-' + $('.toast').length; 
    var toast = '<div id="' + eleID + '" class="toast" role="alert" aria-live="assertive" aria-atomic="true" data-delay="5000">'+
        '<div class="toast-header">'+
        '<strong class="mr-auto">' + title + '</strong>'+
        '<button type="button" class="ml-2 mb-1 close" data-dismiss="toast" aria-label="Close">'+
        '<span aria-hidden="true">&times;</span>'+
        '</button>'+
        '</div>'+
        '<div class="toast-body">' + message + '</div>'+
        '</div>';
    $('#toasts').append(toast);
    $('#'+eleID).toast('show');
}