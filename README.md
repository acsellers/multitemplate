Multitemplate
=============

[![Build Status](https://travis-ci.org/acsellers/multitemplate.svg?branch=master)](https://travis-ci.org/acsellers/multitemplate)

While there are multiple template libraries for Go, I wanted one that would allow
me to mix different styles of templates, without having to create some kind of
lookup table to see which library a template was created in.

Multitemplate is a set of parsers, a html function library, template
interoperation functions, and an implementation of a Template similar to the
html/template Template struct in the standard library.

Current documentation is available at [godoc.org](http://godoc.org/github.com/acsellers/multitemplate).

Bham documentation
------------------

I am continuing to add to the [godoc documentation](http://godoc.org/github.com/acsellers/multitemplate/bham),
but a quick taste of bham (as of 0.2) is as follows:

```
  !!!
  %html
    %head
      = yield "head"
    %body
      = block "header"
        #header
          = if .user
            #options
              Connected as {{.user.Username}} |
              %a(href="{{ url "Hotels.Index" }}") Search
              |
              %a(href="{{ url "Hotels.Settings" }}") Settings
              |
              %a(href="{{ url "Application.Logout" }}") Logout
      = end_block
      = yield "content"
```


Using the library
-----------------

Like database/sql, you can load in dialects (or parsers in this case) like the following

```
import (
  _ "github.com/acsellers/multitemplate/bham"
  "github.com/acsellers/multitemplate"
)
```

multitemplate uses extensions to detect which parser to use, a file named layout.bham.html
would be shortened to layout.html as a template and parsed using the bham parser.

The simplest way to parse files, it to simply pass an array of file names to ParseFiles or 
Template.ParseFiles (which will determine parsers, remove extenstions for you), but you can 
also pass the name, source and parser name to the Parse function, but Parse will not remove
extensions or detect parsers for you.

In addition to Execute and ExecuteTemplate, there is also an ExecuteContext, which it the
way to configure layouts, and pre-set blocks and template to yield or output during execution.



Revel integration
-----------------

You can find instruction on how to integrate multitemplate and revel at the godoc for
[github.com/acsellers/multitemplate/revel](http://godoc.org/github.com/acsellers/multitemplate/revel).
Information on the replacement controller struct to use with the revel integration is at the godoc for
[github.com/acsellers/multitemplate/revel/app/controllers](http://godoc.org/github.com/acsellers/multitemplate/revel/app/controllers).

Samples in the revel folder are a selection of samples available from
github.com/revel/revel, but with templates converted to use multitemplate. (Yes I
know there's only one there now, I'll add more soon).


Development Status
------------------

_Done:_

* Move helpers, mussed, and bham into subdirectories of the same repo.
* Write a buildable version of multitemplate.Template
* Figure out how to set up Delims on standard library
* Write tests on multitemplate.Template
* Implement the template interoperation library
* Fix up yield function to take fallback template name
* Content_for function with template name
* Implement block and end_block
* Write a revel connector

_To Do:_

* Implement the helpers library
* Figure out what other libraries (martini?) could use a connector
