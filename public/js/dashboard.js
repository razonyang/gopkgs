$(function() {
  // sidebar
  var currentURL = window.location.pathname;
  $('.sidebar a.nav-link').each(function() {
      if (currentURL.startsWith($(this).attr('href'))) {
          $(this).addClass('active');
      }
  });

  'use strict';
  window.addEventListener('load', function() {
    var forms = document.getElementsByClassName('needs-validation');
    var validation = Array.prototype.filter.call(forms, function(form) {
      form.addEventListener('submit', function(event) {
        if (form.checkValidity() === false) {
          event.preventDefault();
          event.stopPropagation();
        }
        form.classList.add('was-validated');
      }, false);
    });
  }, false);
});