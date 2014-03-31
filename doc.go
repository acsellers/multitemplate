/*
Multitemplate is a library to allow you to write html templates in
multiple languages, then allowing those templates to work with each
other using either a Rails-like yield/content_for paradigm or a
Django style block/extends paradigm.

Multitemplate at the moment has 3 languages, the standard Go template
syntax, a simplified haml-like language called bham, and a simple
mustache implementation. Multitemplate has an open interface for
creating new syntaxes, so external languages can be easily used.

Terminology

Yield's are executing saved templates or blocks. You can add a fallback
template to yields, but not fallback content. Yielding with a name will return the
first template set for that name, or the content of the first block that
had that name.

ContentFor will set a template to be executed on a name. This is similar
to the template command built in to the go template library, but with
a layer of indirection.

Blocks are template content that is executed, then saved for a later time.
Blocks share names with ContentFor and yields, so a yield might outtput
the content from a block, or a template set with ContentFor.

Inherited templates are templates that use the inherits function to
define that after it is executed, then another template should be
executed as well. These templates should only be made up of non-writing
functions and blocks.

Context's are the way to use the more advanced features of the multitemplate
library. With Context's, you can set two templates to be executed, a Main
template (executed first), and a Layout template (executed after Main). You
can also set Templates for Yields and Block content. Since you can't pass
RenderArgs in the ExecuteContext, you should put your RenderArgs in the Dot
variable.

Layouts are templates executed after a Main template. Context's are the way
to define Layouts to be executed. The Main template should set up content
that can then be yielded using the Main template. Yielding without a name
will cause the main template's content to be output.

Integrations

While multitemplate is available to use as a library in all
Go applications, it also includes integrations libraries that
will integrate multitemplate into external frameworks without
requiring the user to learn how to integrate this library.

The Revel integration is (as far as I know) a drop-in
replacement for the Revel template library. Instructions on
how to integrate are available in the godoc for the
github.com/acsellers/multitemplate/revel subdirectory, while
the integration code is in
github.com/acsellers/multitemplate/revel/app/controllers due
to how revel deals with modules.

Yield functionality

The following code demonstrates the common types of yield statements.
The first yield will render the template assigned to the stylesheets
key, or render the template "include/javascripts.html" if there is
not a set argument. The second yield will render the template set on
"stylesheets" with the .Stylesheets as the data. The third will render
the template set for "sidebar" with the data originally put into
ExecuteTemplate. The fourth yield will render the main template. This will
be the template set for the Main attribute on the Context struct in this
case.


app_controller.go
  ctx := multitemplate.NewContext(renderArgs)
  ctx.Yields["sidebar"] = "sidebars/admin.html"
  ctx.Main = "app/main.html"
  ctx.Layout = "layouts/main.html"
  templates.ExecuteContext(writer, ctx)

layouts/main.html
  <html>
    <head>
      {{ yield "javascript" (fallback "include/javascripts.html") }}
      {{ yield "stylesheets" .Stylesheets}}
    </head>
    <body>
      <div class="row">
        <div class="span3">
          {{ yield "sidebar" }}
        </div>
        <div class="span9">
          {{ yield }}
        </div>
      </div>
    </body>
  </html>

Block functionality

The following code describes two templates that use the inherits
function to utilize template inheritence. It works similarly to
the yield example. Note that the inherits call should be the first
call of the view, any code before the inherits call may be sent to
the writer added to the Execute call.


app_controller.go
  templates.ExecuteTemplate(writer, "app/index.html", renderArgs)

layouts/main.html
  <html>
    <head>
      {{ block "javascript" }}
        <script src="/assets/app.js" type="text/javascript">
      {{ end_block }}
      <link rel="stylesheet" href="/assets/app.css">
    </head>
    <body>
      <div class="row">
        <div class="span3">
          {{ block "sidebar" }}
            <ul>
              <li class="home">Home</li>
            </ul>
          {{ end_block }}
        </div>
        <div class="span9">
          {{ block "content" }}
            Could not find your content.
          {{ end_block }}
        </div>
      </div>
    </body>
  </html>

app/index.html
  {{ inherits "layouts/main.html" }}
  {{ block "content" }}
    Amazing content!
  {{ end_block }}


Yield and Block functionality

As both yields and blocks are built on the same underlying mechanisms, they can be
combined in interesting ways. Implementation wise, blocks are like yields that have
embedded fallback content, while yields have to have separate template fallbacks. Both
Layouts and the Main template can extend other templates.

app_controller.go
  templates.ExecuteTemplate(writer, "app/index.html", renderArgs)

app/index.html
  {{ extends "layouts/main.html" }}

  {{ block "javascript" }}
    <script src="/assets/app.js">
  {{ end_block }}

  {{ content_for "more_stylesheets" "assets/beta-css.html" }}

layouts/main.html
  <html>
    <head>
      {{ yield "javascript" }}
      {{ block "stylesheets" }}
        {{ yield "more_stylesheets" }}
        <style>
          body { margin-left: 40px }
        </style>
      {{ end_block }}
    </head>
    <body>
      {{ block "content" }}
        Content goes here
      {{ end_block }}
    </body>
  </html>

Functions Reference

yield allows for rendering template aliases or simply rendering nothing. Rendering
the Main template without a Main template set is an error

  // yield the main template with the original RenderArgs
  {{ yield }}

  // Yield the main template with specific RenderArgs
  {{ yield .CurrentObject }}

  // Yield a pre-set template with original RenderArgs,
  // or render nothing
  {{ yield "hero_module" }}

  // Yield a pre-set template, or a fallback template
  // if the pre-set template is not present
  {{ yield "blurb" (fallback "demo/lorem_blurb.html" }}

  // Yield a pre-set template with specific RenderArgs
  {{ yield "carousel" .CarouselImages }}

  // Yield a pre-set template, or a fallback template, both
  // with specific RenderArgs
  {{ yield "results" .CurrentObject (fallback "errors/undefined.html") }}

content_for allows you to set a template to be rendered in a block or yield from
within a template.

  // Set the sidebar key to be a specific template
  {{ content_for "sidebar" .Current.Sidebar }}

block saves content inside the current template to a key. That key can be recalled
using yield or another block with the same key in the final template with a block
of that key. Keys are claimed by the first block to render to them.

  {{ block "sidebar" }}
    <ul>
      <li>One</li>
      <li>Two</li>
      <li>Three</li>
    </ul>
  {{ end_block }}

end_block ends the content are started by block

extends marks that the current template is made up of blocks that will be executed
in the context of another template. Template inheritence can be carried to arbitrary
levels, you are not limited to using extends only once in template execution.

  {{ extends "include/main.html" }}

root_dot is the orignal RenderArgs passed in to the ExecuteTemplate call

  {{ $title = root_dot.Title }}

exec execute an arbitrary template with the passed name and data

  {{ exec .Header.Path . }}

fallback sets a specific template to be rendered in the case that a yield call finds
that there is no content set for the key of the yield.

  {{ yield "footer" . (fallback "include/old_footer.html") }}

Things to know about

yield . is ambiguous when dot is currently a string. It could be either a request
to output a pre-set template or block, or to render the main template with the dot
as the data.

Assigning the same key in a Context for both Yields and Content means that the
Content will be ignored. Calling content_for and block (in templates) with the same key has lock-out
protection within the template functions. In this case, we will use the template
named in the Yields map. Within the templates, the rule is the first to claim the key,
wins. Any integrations that hide the Context, will operate under the assumption that the
last claim before template execution should win.

Getting an error about a stack overflow during template execution is most likely a
template that is yielding itself.

The block function has two related functions, define_block and exec_block. If you need
to define a block, even when it would normally execute the block (for instance, if you
are in the main layout, and which to ensure a block exists before yielding or executing
a template), define_block will save the content of the block and not output the content.
This will not override any content already saved for that block name. exec_block is the
reverse function, it will force the block to be executed. If you need to start the main
template with a block, and you are using a template, exec_block will cause the block
to be executed correctly. You should use the standard block call in nearly all situations.

There are tests in integration_test.go that spell out how yields and block interact in
all the situations that I could think of. The test cases are spelled out with minimal
templates, names for each test case, and a description of what the situation that the case
is testing.

Bham is a beta-quality library. I've tried to fix the bugs that I'm aware of, but I'm
sure that there's more lurking out there.

The Mustache implementation here is alpha quality. It's low on the totem pole for
improvements.

Version plans

First release is 0.1, which has bham and html/templates available as first-class languages.
Blocks and yields are supported, along with layouts.

Second release is 0.5, which adds the helpers library, and the super_block, main_block and
define_block functions. Also a whole bunch of new tests for yields, blocks and their
interactions.

Third release will be either a 0.6 or a 1.0, adding things I forgot, fixing bugs
discovered and things that need to be fixed. Mustache will get integration tests, function
calling, blocks. If there were relatively few bugs to fix, then this will be 1.0.

Releases after 1.0, will be adding functions or languages. Template languages I'm interested in
investigating adding are: jade/slim, full haml, some sort of lispy thing, handlebars,
jinja2, and Razor.
*/
package multitemplate
