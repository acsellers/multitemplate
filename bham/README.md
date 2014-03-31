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

## Features being worked on

* More Documentation

## Unlikely To Be Implemented Features

* Curly branch hashrocket syntax for attributes
* Multiple line prefixes for different visibility/escaping
