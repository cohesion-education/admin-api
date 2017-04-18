$(document).ready(function() {
  $.getJSON("/api/taxonomy", function(result) {
    $.each(result, function() {
      flattenTaxonomyOptions(this)
    });
  });
});

function flattenTaxonomyOptions(taxonomy){
  if(taxonomy == null || taxonomy.id == null){
    return
  }

  $.ajax({
    url: "/api/taxonomy/" + taxonomy.id + "/children"
  }).done(function(result){
    if(result.children){
      $.each(result.children, function(index, child) {
        child.name = taxonomy.name + " > " + child.name
        flattenTaxonomyOptions(child)
      })
    }else{
      $('#top-level-categories').append($("<option />").val(taxonomy.id).text(taxonomy.name))
    }
  });
}
