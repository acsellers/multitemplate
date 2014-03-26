/*
Multitemplate is a library to allow you to write html templates in
multiple languages, then allowing those templates to work with each
other using either a Rails-like yield/content_for paradigm or a
Django style block/extends paradigm.

Multitemplate at the moment has 3 languages, the standard Go template
syntax, a simplified haml-like language called bham, and a simple
mustache implementation. Multitemplate has an open interface for
creating new syntaxes, so external languages can be easily used.

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
the yield example.

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
embedded fallback content, while yields have to have separate template fallbacks.

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

Known Issues

yield . is ambiguous when dot is currently a string. It could be either a request
to output a pre-set template or block, or to render the main template with the dot
as the data.

Using the same name for a block and yield will lead to ambiguous results when rendering.
Depending on how you are outputting the code (yield "name" vs block "name") will determine
which version outputs. This should only appear when you are setting the same key using
both a block and a yield in non-final templates, that are then rendered in the final
template.

Using a fallback argument for the main template will not work. This is due to the fact
that the fallback function is an experimental feature that may show up quite a bit in
the helpers library.

Fallback may not be parsed correctly in the bham library.
*/
package multitemplate
