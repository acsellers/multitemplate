# bham - Blocky Hypertext Abstraction Markup

bham is a similar language to something like haml, but it implements
a subset of what haml does in order to keep the language neat and tidy.
I think of it as the haml syntax mixed with the command language of Go's
builtin template library, not that suprising since this library compiles 
into ParseTrees for the builtin template libraries. It can interoperate
with existing go templates, and benefits from the html/template's
escaping functionality.

----------------------

## Documentation

Full syntax documentation is on [Godoc.org](http://godoc.org/github.com/acsellers/unitemplate/bham)

## Working Markup Examples

For a web page with the title set to 'Example Page' and an h1 tag
with the content 'Whatever' you would do the following.

```
<!DOCTYPE html>
%html
  %head
    %title Example Page
  %body
    %h1 Whatever
```

If I wanted to display a default title if there wasn't a PageTitle 
attibute in the Render Arguments when Rendering the template, and
otherwise to use the PageTitle argument, then you would do the 
following. I'm not sure that the {{ }} bit works yet...

```
<!DOCTYPE html>
%html
  %head
    = if .PageTitle
      %title
        = .PageTitle
    = else
      %title No Title Set
```

And the big one, that shows off pretty much all the features...
```
<!DOCTYPE html>
%html(ng-app)
  %head
    = $current := .Current.Page.Name
    = javascript_include_tag "jquery" "angular"
    %title Web Introduction: {{ $current }}
    = if .ExtraJsFiles
      = javascript_include_tag .ExtraJsFiles
    = stylesheet_link_tag "ui-bootstrap"
    = with $current := .Current.Variables
      = template "Layouts/CurrentJs.html" $current
  %body
    %div(class="header {{ .HeaderType }}")
    .hello Welcome to the web {{ .User.Name }}
      You are in section {{ $current }}.

    .row-fluid
      .span3
        = template "Layouts/Navigation.html" .
      .span9
        = yield .
    .row-fluid
      .span10.offset1
        = range $index, $sponsor := .Sponsors
          .sponsor-mini(data-bg-image="sponsor-{{ $sponsor.Img }}")
            = link_to $sponsor.Name $sponsor.Url "class='name'"
```

## Implemented Features

* Plaintext passthrough
* %tag expansion
* If/Else Statements
* Tag Nesting
* Range statements for collection data structures
* = ... for Lines with pipelines on them
* Parentheses for HTML-like attributes
* Class and ID shorthand
* With statement for limited visibility variables
* Template Variables
* {{ }} For embedded pipeline output

## To Be Implemented Features

* Code Quality
* More Documentation

## Unlikely To Be Implemented Features

* Curly branch hashrocket syntax for attributes
* Multiple line prefixes for different visibility/escaping
