$(function() {

});

function loadBadge() {
    console.log(11)
    var form = $('#badgeForm');
    var interval = form.find('select[name="interval"]').val();
    var package = form.find('input[name="package"]').val();
    var url = siteURL + "/badges/downloads/"+interval+"/"+package;
    $('#badgeURL').text(url);
    $('#badgePreview').attr('src', url);
}