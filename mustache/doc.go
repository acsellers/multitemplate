/*
mustache is a package to parse the mustache template format into the
Tree format needed by multitemplate. In general, this implementation
follows the base mustache specification, with no significant changes.
However this implementation does diverge in respect to templates
functions, (lambdas in mustache speak). Fill in this section when lambda's
are completed.

Lambdas that are called as if statements will have their returns taken
and used to set the dot value, or quite similar to the with statement
from Go's stdlib template libraries.
*/
package mustache
