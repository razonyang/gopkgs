$(function() {
    new ClipboardJS('.btn-copy');
});

function loadBadge() {
    var form = $('#badgeForm');
    var package = form.find('input[name="package"]').val();
    if (package == '') {
        $('#badgeURL').val('');
        return;
    }
    var interval = form.find('select[name="interval"]').val();
    var style = form.find('select[name="style"]').val();
    var url = siteURL + "/badges/downloads/"+interval+"/"+package;
    if (style != "") {
        url += "?style="+style;
    }
    $('#badgeURL').val(url);
    $('#badgePreview').attr('src', url);
}