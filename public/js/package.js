function deletePackage(id) {
    if (!confirm("Are you sure to delete this package")) {
        return
    }

    $.ajax({
        url: '/package/' +id,
        method: 'DELETE',
        success: function(resp) {
            if (resp.status != 'success') {
                alert(resp.message);
                return;
            }
            window.location.reload();
        },
    })
}