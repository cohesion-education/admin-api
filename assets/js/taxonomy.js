$(document).ready(function() {
  $('a.taxonomy').click(getTaxonomyChildren);
  $('.add-taxonomy').click(addTaxonomyHandler);
});

var addTaxonomyFormTemplate  = '<li><form id="[form-id]">'
  + '<input type="text" name="name" placeholder="Name">'
  + '<input type="hidden" name="parent_id" value="[parent-id]" />'
  + '<input type="submit" value="Save" />'
  + '</form></li>'

var addTaxonomyLITemplate = '<li class="add"><a class="add-taxonomy" href="/api/taxonomy/add/[id]">Add</a></li>'
var taxonomyLITemplate = '<li><a class="taxonomy" href="/api/taxonomy/[id]/children" class="taxonomy">[name]</a></li>'
var newChildTaxonomyLITemplate = '<li><a class="taxonomy" href="/api/taxonomy/[id]/children">[name]</a>'
  + '<ul>' + addTaxonomyLITemplate + '</ul></li>'

function getAddTaxonomyFormID(parentID){
  var formID = 'add-' + parentID
  return formID
}

function reinitOnclickBindings(){
  $('a.taxonomy').unbind('click', getTaxonomyChildren).bind('click', getTaxonomyChildren)
  $('.add-taxonomy').unbind('click', addTaxonomyHandler).bind('click', addTaxonomyHandler)
}

function getTaxonomyChildren(){
  var url = $(this).attr('href')
  var current = $(this)

  if($(this).siblings("ul").length != 0){
    if($(this).siblings("ul").css("visibility") != 'hidden'){
      $(this).siblings("ul").slideToggle()
      return false
    }else if($(this).siblings("ul").css("visibility") == 'hidden'){
      $(this).parent().remove(current.next())
    }
  }

  $.ajax({
     type: "GET",
     url: url,
     success: function(result){
       var childrenHtml = '<ul style="display:none">'
       if (result.children != null) {
         for (var i in result.children) {
           var child = result.children[i]
           var li = taxonomyLITemplate
                    .replace("[id]", child.id)
                    .replace("[name]", child.name)

           childrenHtml += li
         }
       }

       childrenHtml += addTaxonomyLITemplate.replace("[id]", result.parent_id)
       childrenHtml += '</ul>'
       current.parent().append(childrenHtml)
       current.siblings("ul").slideToggle()
       reinitOnclickBindings()
     },
     error: function(jqXHR, message){
       alert("Failed to get children " + message)
     }
   });

   return false
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
     data: data
   }).done(function(result){
     taxonomyID = result

     var li = newChildTaxonomyLITemplate
                .replace("[id]", taxonomyID)
                .replace("[name]", taxonomy.name)
                .replace("[id]", taxonomyID)

     $('form#' + getAddTaxonomyFormID(taxonomy.parent_id)).parent().replaceWith(li)
     reinitOnclickBindings()
   }).fail(function( jqXHR, textStatus, errorThrown ){
     alert("Failed to Add Taxonomy. " + errorThrown + ": " + jqXHR.responseText)
   });

  return false
}
