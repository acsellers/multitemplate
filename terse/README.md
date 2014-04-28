terse: concise html templating
==========================

Terse is a syntax for writing templates using a syntax inspired by mustache and
slim.

_doctypes_

```
!!
// <!DOCTYPE html>
!! strict
// <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
```

_tags_

```
html
// <html> … </html>
h1  markup example 
// <h1>markup example</h1>
```

_tags with id/class_

```
#footer
// <div id=”footer”>...</div>
.clear
// <div class=”clear”></div>
table.striped
// <table class=”striped”> … </table>
```

_filters_

```
:js
  thing
// <script type=”text/javascript”> … </script>
```

_defines_

```
::mini.html
  html
    head > title= .Title
    body
      @content
[content]
  Content

/*
  {{ define "mini.html" }}
    <html>
      <head><title>{{ .Title }}</title></head>
      <body>{{ yield "content" }}</body>
    </html>
  {{ end }}

  {{ block "content" }}
    Content
  {{ end_block }}
*/
```

_yields_

```
@*
// {{ yield }}
@footer
// {{ yield “footer” }}
@footer=test.html 
// {{ content_for “footer” “test.html” }}
@footer|test.html
// {{ yield "footer" (fallback "test.html") }}
```

_extend_

```
@@layouts/app.html
// {{ extend "layouts/app.html" }}
```

_blocks_

```
// Regular block
[name]
   things
/*
  {{ block “name” }}
  things
  {{ end_block  }}
*/
``
```

```
// Define block
^name]
  things
/*
  {{ define_block “name” }}
  things
  {{ end_block }}
*/
```

```
// Exec Block
$name]
  things
/*
  {{ exec_block “name” }}
  things
  {{ end_block }}
*/
```

_if/else/range_
```
// If Else statement
?items
  things
!?
  no things
/*
  {{ if .items }}
    things
  {{ else }}
    no things
  {{ end }}
*/
```

```
// Range Else Statement
&.Items
  = name
!&
  no items

/*
  {{ range .items }}
    {{ .name }}
  {{ else }}
    no items
  {{ end }}
*/
```

_with_

```
// With Statement
>.User:$user
  = $user.Name
/*
  {{ with $user := .User }} 
    {{ $user.Name }}
  {{ end }}
*/

// With/Else Statement
>.User
  = .Name
!>
  Not logged in!
/*
  {{ with .User }}
    {{ .Name }}
  {{ else }}
    Not logged in!
  {{ end }}
*/
```

_template_

```
// Template call
>>layouts/header.html
{{ template "layouts/header.html" . }}

// Template call with specific data
>>layouts/footer.html $args
{{ template "layouts/footer.html" $args }}
```

_interpolation_

```
First Name: {{user.name}}
```
