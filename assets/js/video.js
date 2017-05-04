$(document).ready(function() {
  $('#top-level-categories').prop('disabled', 'disabled')

  $.getJSON("/api/taxonomy/flatten", function(results) {
    $.each(results, function(i, taxonomy) {
      var option = $("<option />").val(taxonomy.id).text(taxonomy.name)
      if($('#selected_taxonomy_id').val() == taxonomy.id){
        option.attr("selected", true)
      }

      $('#top-level-categories').append(option)
    });

    $('#top-level-categories').prop('disabled', false)
  });

  $('form#video').submit(function(){
    $('#errors').hide()
    $('#progressbar').hide()

    var data = new FormData()
    $.each($(this).find("input[type='file']"), function(i, tag) {
      $.each($(tag)[0].files, function(i, file) {
          data.append(tag.name, file);
      });
    });
    var params = $(this).serializeArray();
    console.log("serialized form inputs: " + JSON.stringify(params))
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
      processData: false,
      xhr: function(){
        var xhr = new window.XMLHttpRequest()
        xhr.upload.addEventListener("progress", function(event) {
          if(event.lengthComputable){
            var percentComplete = event.loaded / event.total
            percentComplete = parseInt(percentComplete * 100)
            if(percentComplete === 100){
              $('#progressbar').html("Upload Complete, but video is still processing. You will be redirected when ready")
              console.log("upload complete")
            }else{
              $('#progressbar').html(percentComplete + "%")
              $('#progressbar').width((percentComplete - 10) + "%")
            }

            $('#progressbar').show()
          }
        }, false);
        return xhr
      }
    }).done(function(result){
      console.log("Success! " + JSON.stringify(result))
      window.location.replace(result.redirect_url)
    }).fail(function( jqXHR, textStatus, errorThrown ){
      $('#progressbar').hide()

      console.log("Failed :( " + textStatus + " " + errorThrown)
      console.log("jqXHR.responseText: " + jqXHR.responseText)

      try {
        var response = JSON.parse(jqXHR.responseText)
        if(response != null){
          var errors = response.error + '<ul>'
          for(i = 0; i < response.validation_errors.length; i++){
            errors += '<li>' + response.validation_errors[i].error + '</li>'
          }
          errors += '</ul>'

          $('#errors').html(errors)
          $('#errors').show()
        }
      } catch(e) {
        $('#errors').html(jqXHR.responseText)
        $('#errors').show()
      }
    });

    return false
  });
});
