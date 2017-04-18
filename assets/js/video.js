$(document).ready(function() {
  $.getJSON("/api/taxonomy", function(result) {
    var videoCategories = $('#top-level-categories')
    $.each(result, function() {
      videoCategories.append($("<option />").val(this.id).text(this.name));
    });
  });

  $('#top-level-categories').change(taxonomySelectedChangeHandler)
});

function taxonomySelectedChangeHandler(){
  var taxonomyID = $(this).val()
  alert(taxonomyID)
  $.getJSON("/api/taxonomy/" + taxonomyID + "/children", function(result) {
    $.each(result.children, function() {
      alert(this.id + ' ' + this.name)
    });
  });
}
