$(document).ready(function() {
  $('#top-level-categories').prop('disabled', 'disabled')

  $.getJSON("/api/taxonomy", function(result) {
    $.each(result, function() {
      flattenTaxonomyOptions(this)
      $('#top-level-categories').prop('disabled', false)
    });
  });

  $('form#video').submit(function(){
    $('#errors').hide()

    var data = new FormData()
    $.each($(this).find("input[type='file']"), function(i, tag) {
      $.each($(tag)[0].files, function(i, file) {
          data.append(tag.name, file);
      });
    });
    var params = $(this).serializeArray();
    $.each(params, function (i, val) {
      data.append(val.name, val.value);
    })

    console.log($(this).attr('method') + ":" + $(this).attr('action'))

    $.ajax({
      type: $(this).attr('method'),
      url: $(this).attr('action'),
      data: data,
      cache: false,
      contentType: false,
      processData: false
    }).done(function(result){
      console.log("Success! " + JSON.stringify(result))
      window.location.replace(result.redirect_url)
    }).fail(function( jqXHR, textStatus, errorThrown ){
      console.log("Failed :( " + textStatus + " " + errorThrown)
      console.log("jqXHR.responseText: " + jqXHR.responseText)
      var response = $.parseJSON(jqXHR.responseText)
      if(response != null){
        var errors = response.error + '<ul>'
        for(i = 0; i < response.validation_errors.length; i++){
          errors += '<li>' + response.validation_errors[i].error + '</li>'
        }
        errors += '</ul>'

        $('#errors').html(errors)
        $('#errors').show()
      }
    });

    return false
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
      var option = $("<option />").val(taxonomy.id).text(taxonomy.name)
      if($('#selected_taxonomy_id').val() == taxonomy.id){
        option.attr("selected", true)
      }

      $('#top-level-categories').append(option)
    }
  });
}
