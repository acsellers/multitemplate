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

  Django Functions

  - add_slashes: Wherever s has a single or double quote, add a backslash before the quote mark

  - cap_first: Capitalize the first letter in the string, ignoring further words and

  - center: Take the string value s and center it within a string of length num that is filled with spaces

  - cut: Remove all characters in cutset from the string s

  - default: When x is empty string, nil, etc., then output y otherwise output x

  - default_if_nil: When x is nil, return y, otherwise return x

  - escape: Return the string s after is was html escaped

  - escape_js: Return the string after it was escaped for javascript use

  - filesize_format: Given the size in bytes, return the smallest format that would be greater than 1 when output

  - first: Return the first value of a slice of items

  - float_format: Return a string where f has exactly n places after the decimal point

  - force_escape: Run html escape on a string

  - get_digit: Return the nth rightmost digit from i

  - join: Join the values in the slice l with j and return the resulting string

  - last: Returns the last item in a slice or array

  - length: Returns the length of an array, slice, map, or string

  - length_is: Tests that the length of l is equal, then returns a bool

  - link_to: Return a clickable link from an address of the first string with the text as the second string

  - ljust: Converts i to a string, then left aligns it within a area of length n

  - lower: Convert a string to lowercase

  - number_lines: Add the line number to each line in s in the form "1. "

  - pluralize: Pluralize a string

  - pprint: Print out a piece of data, for debugging purposes

  - pytime: Formats time according to Python time formatting syntax

  - quick_format: Escapes content, then converts newline characters into <br>'s

  - random: Returns a random item from a slice

  - rjust: Right align a value in a string of length n

  - safe: Mark a string to template.HTML so it will not be escaped

  - safe_seq: Convert all items in a slice to template.HTML marked strings

  - slugify: Returns a lowercased version of that removes non-alphabetic characters and spaces are now hyphens

  - title: Title case the string

  - truncate: Truncate a string to a certain number characters, the count includes the three-character ellipsis

  - upper: Convert string s to uppercase letters

  - url_encode: Encodes a value for use in a url

  - urlize: Return a clickable link from an url address string

  - urlize_truncate: Returns a clickable link from an address s where the clickable text is truncated to a number of characters
*/

package helpers
