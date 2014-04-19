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

      {{ text_field_tag "reviewer" .Author (data "placeholder" "Select a reviewer for the article") }}


  FormTag functions

  - button_tag: Create a button tag

  - check_box_tag: Create an input tag with a type of checkbox

  - email_field_tag: Create an input tag with a type of email

  - fieldset_tag: Open a fieldset tag with a legend with the name you passed

  - end_fieldset_tag: Close a fieldset tag, for auto-closing template languages

  - file_field_tag: Input field with a type of file, make sure that the form tag has an enctype of "multipart/form-data"

  - form_tag: Open a form tag, this does not yet add any csrf protection, though it is planned.

  - end_form_tag: Closes a form tag, for auto-closing template languages

  - hidden_field_tag: Creates a hidden input field, this can be modified client-side

  - label_tag: Creates a label tag

  - number_field_tag: Creates an input field with a type of number

  - password_field_tag: Creates a password field input, this will hide the characters users type in.

  - phone_field_tag: Creates an input with a type of tel, this is great for mobile

  - radio_button_tag: Creates a radio button input element.

  - range_field_tag: Creates a input with a type of range. That means it's a slider.

  - search_field_tag: Createa an input tag with a type of search

  - submit_tag: Create an input tag with a type of submit

  - text_area_tag: Create a textarea tag

  - text_field_tag: Create a normal text field input

  - url_field_tag: Create an input for a url.

  - utf8_tag: Add a hidden input with a value that is a valid UTF8 character ouside of ASCII.

  General Functions

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

  Select Tag Functions

  The functions for select tags are loaded as part of the Form Tags module.

  - option: Turn two values into an Option pair, the first value is the name for the option while the second value becomes the value for the option.

  - options: Turns a set of values into a list of options. Each value is set to both then name and value of the option. Any Option's from the option function are passed through.

  - group_options: Take an OptionList from options or options_with_values and enclose them with an optgroup with the label of the first argument.

  - options_with_values: Turns a variable number of values into a list of options. This is equivalent to calling Option on each pair then passing the results to options.

  - select_tag: First value is a string for the name of the select tag, second is as Option List of a mix Options and/or OptionGroups, followed by AttrList's to set other options about the select tag.

  Link Functions

  - link_to: Return a clickable link from an address of the first string with the text as the second string

  - link_to_function: Return a clickable link that runs a javascript function instead of a web address.

  - url_encode: Encodes a value for use in a url

  - urlize: Return a clickable link from an url address string

  - urlize_truncate: Returns a clickable link from an address s where the clickable text is truncated to a number of characters

  Asset Functions

  Note that Asset Functions don't use integrations to discover assets, they are configured with the
  package variable AppInfo which contains the information to create the necessary links.

  - atom_link: Returns a link tag for a rss feed based on the RootURL + the path you send.

  - favicon_link: Returns the link to a favicon. Note that the favicon should be located in the ImageRoot.

  - image_tag: Returns the img tag for an image in the ImageRoot.

  - javascript_link: Returns the script tags for one or more javascript files in JavascriptRoot

  - rss_link: Returns a link tag for a rss feed based on the RootURL + the path you send.

  - stylesheet_link: Return the link tags for one or more stylesheets base on the StylesheetRoot + paths given.
*/
package helpers
