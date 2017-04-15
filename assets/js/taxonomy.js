$(document).ready(function() {
  $('.add-taxonomy').click(addTaxonomyHandler);
});

var addTaxonomyFormTemplate  = '<li><form id="[form-id]">'
  + '<input type="text" name="name" placeholder="Name">'
  + '<input type="hidden" name="parent_id" value="[parent-id]" />'
  + '</form></li>'

var taxonomyLITemplate = '<li><a href="/taxonomy/[id]">[name]</a>'
+ '<ul><li class="add"><a class="add-taxonomy" href="/taxonomy/add/[id]">Add</a></li></ul>'
+ '</li>'

function getAddTaxonomyFormID(parentID){
  var formID = 'add-' + parentID
  return formID
}

function addTaxonomyHandler(){
  var url = $(this).attr('href')
  var parentID = url.substr(url.lastIndexOf('/') + 1)
  if(!$.isNumeric(parentID)){
    parentID = 0
  }

  var formID = getAddTaxonomyFormID(parentID)
  if($('form#' + formID).length != 0){
    return false
  }

  var form = addTaxonomyFormTemplate
              .replace("[form-id]", formID)
              .replace("[parent-id]", parentID)

  $(form).insertBefore($(this).parent())
  $('#' + formID).submit(addTaxonomyFormSubmitHandler);

  return false
}

function addTaxonomyFormSubmitHandler(){
  var taxonomy = new Object();
  taxonomy.name = $(this).find('input[name="name"]').val()
  taxonomy.parent_id = $(this).find('input[name="parent_id"]').val()
  taxonomy.parent_id = Number(taxonomy.parent_id)

  var data = JSON.stringify(taxonomy)
  console.log('json data: ' + data)

  $.ajax({
     type: "POST",
     url: "/api/taxonomy",
     data: data,
     success: function(result){
       taxonomyID = result

       var li = taxonomyLITemplate
                  .replace("[id]", taxonomyID)
                  .replace("[name]", taxonomy.name)
                  .replace("[id]", taxonomyID)

       $('form#' + getAddTaxonomyFormID(taxonomy.parent_id)).parent().replaceWith(li)

       $('.add-taxonomy').click(addTaxonomyHandler);
     },
     error: function(jqXHR, message){
       alert("Failed to add Taxonomy " + message)
     }
   });

  return false
}
