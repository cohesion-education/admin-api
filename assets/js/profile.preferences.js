$(document).ready(function() {
  $('form#preferences').submit(savePreferencesHandler)
});

function savePreferencesHandler(){
  var action = $(this).attr('action')
  var data = $(this).serialize()
  console.log('action: ' + action)
  console.log('json data: ' + data)

  $.ajax({
     type: "POST",
     url: action,
     data: data
   }).done(function(result){
     alert("Your preferences have been updated")
   }).fail(function( jqXHR, textStatus, errorThrown ){
     alert("Failed to save your preferences. " + errorThrown + ": " + jqXHR.responseText)
   });

  return false
}
