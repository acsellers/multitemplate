Multitemplate
=============

[![Build Status](https://travis-ci.org/acsellers/multitemplate.svg?branch=master)](https://travis-ci.org/acsellers/multitemplate)

While there are multiple template libraries for Go, I wanted one that would allow
me to mix different styles of templates, without having to create some kind of
lookup table to see which library a template was created in.

Multitemplate is a set of parsers, a html function library, template
interoperation functions, and an implementation of a Template similar to the
html/template Template struct in the standard library. The execution backend
of this is the html/template package from the standard library, multitemplate
is a set of functions, parsers, and glue that adds a larger set of functionality 
on top of the standard library.

Current documentation is available at [godoc.org](http://godoc.org/github.com/acsellers/multitemplate).

Bham documentation
------------------

I am continuing to add to the [godoc documentation](http://godoc.org/github.com/acsellers/multitemplate/bham),
this is a simple snippet from the booking examples written in bham.

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

Terse Documentation
-------------------

[Terse](http://godoc.org/github.com/acsellers/multitemplate/terse) is fairly well documented, but could always use some longer form examples.
This is the same snippet as bham, but in terse

```
  !!
  html
    head
      @head
    body
      [header]
        #header
          ?.user
            #options
              Connected as {{.user.Username}} |
              a href=(url "Hotels.Index") Search
              |
              a href=(url "Hotels.Settings") Settings
              |
              a href=(url "Application.Logout") Logout
      @content
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

The sample in the revel folder is the booking sample from
github.com/revel/revel, but with templates converted to use multitemplate languages.


Martini Integration
-------------------

You can find instruction on how to integrate multitemplate and revel at the godoc for
[github.com/acsellers/multitemplate/martini](http://godoc.org/github.com/acsellers/multitemplate/martini).
Information on the replacement controller struct to use with the revel integration is at the godoc for
[github.com/acsellers/multitemplate/martini](http://godoc.org/github.com/acsellers/multitemplate/martini).

The Sample in the martini folder is a port of a small Redis viewer I wrote a while ago,
but with the templates all in the terse language.


Versioning
----------

New versions are released when the features planned for that version 
are complete. If a feature looks to not be ready in the same 
timeframe as the other features earmarked for that version, then the
longer features may get bumped.


Future Version Plans
--------------------

_1.1_

* New things, maybe a language, helpers module, or something from left field.

_1.0_

* Stability, plan to keep same API at this point and add languages, subsystems

_0.6_

* Mustache partial CMTS
* More CMTS tests

_0.5_

* Terse parser gets ported to bham
* FormBuilder-like helper module
* Bham spec
* Common MultiTemplate Test Suite (set of tests that each language must pass)
* Stdlib CMTS
* Bham CMTS
* Terse CMTS
* Refactor mustache tests to use multitemplate.Template and Table Driven tests

Version History
---------------

_0.4_

* New language (terse)
* New sub-library for html helpers
* FormTag, Link, General, Simple Asset helper modules
* helpers modules that are enabled individually
* Better parser construction (working in terse)
* Martini integration (multirender)

_0.3_

* Blocks know how they've been escaped and will check for escaping ruleset matches when being rendered (security)
* Fix bugs in revel integration exposed by security fixes

_0.2_

* Fix issues with bham around function parsing
* Started documenting library
* Various fixes

_0.1_

* Move helpers, mussed, and bham into subdirectories of the same repo.
* Write a buildable version of multitemplate.Template
* Figure out how to set up Delims on standard library
* Write tests on multitemplate.Template
* Fix up yield function to take fallback template name
* Content_for function with template name
* Implement block and end_block
* Write a revel connector
