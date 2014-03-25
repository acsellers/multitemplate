# mussed


Mussed is a nearly mustache template library that maps to text/template/parse Trees instead of being a complete template library.

## Divergences from Mustache

* Quote characters are escaped with the Code instead of the Entity Name
* Templates that aren't found are treated as fatal errors instead of empty strings
* On the third partials test, Go is more proactive than mustache and escaped '<'s where an average mustache would not
* Partials do not inherit the indentation of their caller, this was found on partial specs 7-9.
