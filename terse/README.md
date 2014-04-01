trs: Terse html templating
==========================

Trs is a syntax for writing templates using a syntax inspired by mustache and
slim.

_doctypes_

```
!!!
// <!DOCTYPE html>
!!! strict
// <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
!!! xml
// <?xml version="1.0" encoding="utf-8" ?>
```

_tags_

```
html
// <html> … </html>
h1  markup example 
// <h1>markup example</h1>
```

_filters_

```
:js
  thing
// <script type=”text/javascript”> … </script>
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

blocks

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
&items
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

_tags with id/class_

```
#footer
// <div id=”footer”>...</div>
.clear
// <div class=”clear”></div>
table.striped
// <table class=”striped”> … </table>
```

_interpolation_

```
First Name: #{user.name}
```
