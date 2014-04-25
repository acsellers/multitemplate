/*
Terse is an html templating language inspired by slim, but
written for the features of multitemplate. The minimal logic
marks were inspired by mustache, but I'm not going to pretend
they're not template logic. Note that the logic syntax is the
same as the template logic syntax in the Go standard library.
Terse is not an acronym, it is a description of the design of
the language.

Features

- Significant Whitespace

- Auto-closing tags and functions

- Integrates with other multitemplate languages

- Uses same function calling syntax as stdlib templates

Using Terse

Active terse as a multitemplate languages by importing the
terse library like

    import _ "github.com/acsellers/multitemplate/terse"

At that point, Parse will now take a "terse" parser for templates
and ParseFiles/ParseGlob will detect files with a .terse extension
as terse formatted files.

An example terse file:

    #content
      .title
        h1= .Title
        = link_to "Home" (url_for "home")

The resulting HTML would look like:

  <div id="content">
    <div class="title">
      <h1>Welcome</h1>
      <a href="/home">Home</a>
    </div>
  </div>

Security and XSS Protection

Terse uses the built-in protection of html/template via multitemplate.

Simple HTML

Most of HTML documents is static text. Terse makes writing this static
text as simple as possible. There are three ways to write HTML tags, the
automatic way, the manual way, and the verbatim way. The automatic way
will detect tags by name, and will do a case sensitive search on the
elements listed in the ValidElements map. The manual way will detect
elements by a percentage sign before the element name, which is the same
behavior as Haml and Bham. Finally, you can just add verbatim HTML code,
and it wil be passed through to the final output. This is why the character
'<' has no special meaning in terse.

  // Source
  article
      <h1>My Title</h1>

  // Output
  <article>
    <h1>My Title</h1>
  </article>

Collapsing Tags

If you have lines that are simply declaring elements that are nesting,
you can collapse the lines using the '>' character. Note that you cannot
add attributes or content to the tags. At the moment, you can't use the
id/class shorthand here.

  // Source
  table > tr > td > %custom_element
    Gopher Freeman

  // Output
  <table><tr><td><custom_element>
    Gopher Freeman
  </custom_element></td></tr></table>

Class and Id Shorthand

Id's and Class attributes may be specified in shorthand. If you do
not specify an html element, terse will assume that you want to use
the div element. Multiple class attributes will be combined into the
class attribute, while multiple id attributes will take the final id
attribute. Classes are similar to css declarations, where '.'s precede
class attributes and '#'s precede id attributes. You can still provide
id and class attributes using the standard HTML attributes.

  // Source
  #sidebar
    ul#sidebar_list.vertical_list
      li.sidebar_element.active Thing

  // Output
  <div id="sidebar">
    <ul id="sidebar_list" class="vertical_list">
      <li class="sidebar_element active">Thing</li>
    </ul>
  </div>

HTML Attributes

There are two ways to specify attributes, one which allows for multiple
lines with attributes and another that requires all attributes to be
specied on the same line. By specifying the attribute name, then an
equals sign, then either a string, a $variable, a RenderArg which is
a . followed by words. Finally, if you just specify a helper function.
Note that if the helper function needs arguments, you should wrap the
function and arguments in parentheses. Boolean HTML attributes, that is
attributes that have no value, either require parentheses or can have
the value 'true' or 'false'.


  // Single Line
  div id="special_thing" class=$DivIds special thing content
  // Output
  <div id="special_thing" class="{{ $DivIds }}">
    special thing content
  </div>


  // Multi Line
  input(
    name="is_active"
    type="checkbox"
    value="yes"
    checked=.User.Active
  )
  // Output
  <input type="is_active" type="checkbox" value="yes" checkbox="{{.User.Active}}" />

  // All the ways to specify attributes
  input(
    name=$inputName
    type="checkbox"
    value=.User.Name
    checked=(if_print "checked" .User.ShowName)
    class=.Check.User.Class
  )

  // Boolean Attribute Option 1
  .content ng-view=true
  // Boolean Attribute Option 2
  .content(ng-view)
  // Output
  <div class="content" ng-view />

Doctype

Doctype lines are prefixed with two exclamation points (!!). There is
the default doctype, which is a line with just !!, terse will assume
you want to use the HTML5 Doctype (<!DOCTYPE html>). You can modify
the terse.Doctypes variable to change to another default by changing
the empty string key. Terse includes 8 different Doctypes by default,
which are Transitional (XHTML 1.0 Transitional), Strict (XHTML 1.0
Strict), Frameset (XHTML 1.0 Frameset), 5 (HTML5), 1.1 (XHTML 1.1),
Basic (XHTML 1.1 Basic), Mobile (XHTML Mobile 1.2), and RDFa
(XHTML+RDFa 1.0). These are the same as haml's Doctypes, except
that we don't worry about a Format option.

  // Source
  !!
  // Output
  <!DOCTYPE html>

Comments

By prefixing lines with a double slash (//), that line and any nested
lines will be marked as a comment. Comment lines do not show up in the
rendered template in any way.

Executing Code

There are preferred methods for if, range, with, and block functions,
but for template functions, there is a generic way to use code in
your templates. By adding an = at the start of the line, you can
mark that line as executable, for multiple lines there is a continuation
marker (/=) which will collapse multiple lines into a single line of code,
so long as the continuation lines are nested in the first code line. Note
that the syntax is the same as text/template in the standard library,
and the behavior about inserting code results into the output is same as
the text/template code as well.

  // Source
  = print .User
  // Output
  {{ print .User }}

  // Source
  = form_tag (url_for "user.new") (attrs
    /= "enctype" "multipart/form-data"
    /= "method" "PUT")
  // Output
  <form action="{{ url_for "user.new" }}" enctype="multipart/form-data" method="PUT">

Auto Closing Functions

When lines are nested, terse will look a function another function
that has the same name as the function that has the nested name, but
with the prefix 'end_'.

  // Source
  = fieldset_tag "Visibility"
    = checkbox_tag "admins"
    = checkbox_tag "users"
    = checkbox_tag "visitors"

  // Output
  {{ fieldset_tag "Visibility" }} {{ checkbox_tag "admins" }}
    {{ checkbox_tag "users" }}
    {{ checkbox_tag "visitors" }}
  {{ end_fieldset_tag }}

If and Else Statements

If lines begin with a question mark (?), then a statement to be evaluated. The
method of evaluation is the same as text/template, where if the statement
evaluates to false, 0, nil, or anything with a length of 0, the statement
is not executed. Else is an exclamation mark followed by a question mark,
or "not the question".

  // Source
  ?.User
    Welcome
    = .User.Name
  !?
    Please Login

  // Output
  {{ if .User }}
    Welcome
  {{ else }}
    Please Login
  {{ end }}


Range Statements

Range statements begin with an ampersand (&). Else statements for ranges
are an exclamation mark followed by an ampersand (!&). Variables for the
range statement are specified using colons after the statement you will
be ranging over.

  // Source
  &.Users:$user:$index
    {{ $index }}. {{ $user.Name }}
  !&
    No matching users.

  // Output
  {{ range $user, $index := .Users }}
    {{ $index }}. {{ $user.Name }}
  {{ else }}
    No matching users.
  {{ end }}


Block Statements

Blocks are central to multitemplate and terse features
them in a first-class manner. The simple block call is [name] where
name is the name you are using for that block. If you a special
form of the block statement, just replace the opening '['. For an
exec block you use a $, while define block uses a ^. The mnemonic
for those is that exec is the end of a block, so it is a regex endline
symbol ($), while define-block is more a beginning statement, and the
regex character for the beginning of line is ^. Blocks are automatically
closed in terse.

  // Source
  [content]
    h1
      Welcome
      = .User.Name

  // Output
  {{ block "content" }}
    <h1>Welcome {{.User.Name }}</h1>
  {{ end_block }}

  // Source
  ^content]
    I'm defining this block and not outputting here.
  // Output
  {{ define_block "content" }}
    I'm defining this block and not outputting here.
  {{ end_block }}

  // Source
  $content]
    We will definitely put some block content here.
  // Output
  {{ exec_block "content" }}
    We will definitely put some block content here.
  {{ end_block }}

Yield Statements

Yield statements (for blocks or templates) are written as @name,
where name is the name you are using for that block or template.

  // Source
  #content
    @content

  // Output
  <div id="content">
    {{ yield "content" }}
  </div>

Extend Statements

Extend statements (for inheriting from other templates) are written
using @@ followed by the template name (without quotation marks).

  // Source
  @@layouts/app.html

  // Output
  {{ extend "layouts/app.html" }}

With Statements

A with statement takes a pipeline and either sets the dot or a variable
to the value. With statments will only execute if the value is not a
falsey value (0, false, nil, empty string, array, slice or map). With
statements can have else statements, which are signified by a !>.

  // Source
  >.Current.User
    Welcome
    = .Name
    ?time_gt .LastLogin "3w"
      It's been a while
  !>
    Please login

  // Output
  {{ with .Current.User }}
    Welcome
    {{ .Name }}
    {{ if time_gt .LastLogin "3w" }}
      It's been a while
    {{ end }}
  {{ end }}

  // Source
  >.Current.Location:$loc
    = $loc.Name

  // Output
  {{ with $loc := .Current.Location }}
    {{ $loc.Name }}
  {{ end }}
Interpolation

In addition to the function lines, you can also embed functions into
regular lines, using the Delimeters set on terse (same as the standard
library delimeters by default).

  // Source
  Welcome {{ .User.Name }}.

  // Output (same)
  Welcome {{ .User.Name}}.

Filters

Filters are started with a :, then the name of a registered filter.
Filters can have interpolated code within them, if the filter may
have interpolated code, it needs to return true, else if it returns
false, the filtered text will not be checked for code. Note that the
processed text will be inserted into the template and any interpolated
code will be processed by html/template and escaped in a manner
consistent with the surrounding text.

  // Source
  :js
    $(document).ready(function{ setup_ajax(); });

  // Output
  <script type="text/javascript">
    $(document).ready(function{ setup_ajax(); });
  </script>

Currently there are 3 built-in filters, a plain filter, a javascript filter
(also aliased as js), and a css filter. The Plain filter will take any nested
lines and push the straight out to the the text/template parser. The javascript
filter will wrap the text in a script tag, then send it to be processed by the
text/template parser. Similarly, the css filter will wrap the text in a style
tag, then send it to the parser.

Notes

The format of the documentation and some examples were inspired by
Haml's REFERENCE file.

*/
package terse
