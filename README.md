Multitemplate
=============

While there are multiple template libraries for Go, I wanted one that would allow
me to mix different styles of templates, without having to create some kind of
lookup table to see which library a template was created in.

Multitemplate is a set of parsers, a html function library, template
interoperation functions, and an implementation of a Template similar to the
html/template Template struct in the standard library.

_Done:_

* Move helpers, mussed, and bham into subdirectories of the same repo.
* Write a buildable version of multitemplate.Template
* Figure out how to set up Delims on standard library
* Write tests on multitemplate.Template
* Implement the template interoperation library
* Fix up yield function to take fallback template name
* Content_for function with template name
* Implement block and end_block

_To Do:_

* Write a revel connector
* Implement the helpers library
* Figure out what other libraries (martini?) could use a connector
