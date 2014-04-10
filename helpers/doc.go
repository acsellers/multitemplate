/*

  Helpers is a package that adds a set of functions to multitemplate
  to simplify various common html operations. It is organized into a
  modules that are loaded using the LoadHelpers function. All modules
  depend on a "core" module that will always be loaded.

  Current Modules are: "forms"

  Upcoming Modules are: "js", "assets", "form_builder",

  Meta-Modules are: "all"

  Core functions

  - attr: link together a string key and a value, can be given as an argument to the attrs function

      {{ $disable := attr "disabled" (call .User.AllowedTo "edit_important_thing") }}
      <div id="my_form">
        {{ check_box_tag "Published", true, (attrs $disable) }}
        {{ text_field_tag "Title", true, (attrs $disable) }}
      </div>

  - attrs: turn an array of items into a map that can be passed to functions. You can mix attr's and
  string keys followed values. There is not a built-in attrs merge yet, but functions take multiple
  attrs results and will merge them together at runtime.

      {{ text_field_tag "name" .User.Name (attrs "class" "span50" "autocomplete" "off") }}

  - data: A special case of the attrs function, this will append "data-" before the keys to simplify
  the creation of HTML5 data sttributes.

      {{ select_tag "reviewer" .AuthorListing (data "placeholder" "Select a reviewer for the article") }}


  FormTag functions

  - button_tag

  - check_box_tag

  - email_field_tag

  - fieldset_tag

  - end_fieldset_tag

  - file_field_tag

  - form_tag

  - end_form_tag

  - hidden_field_tag

  - label_tag

  - number_field_tag

  - password_field_tag

  - phone_field_tag

  - radio_button_tag

  - range_field_tag

  - search_field_tag

  - submit_tag

  - text_area_tag

  - text_field_tag

  - url_field_tag

  - utf8_tag

*/
package helpers
