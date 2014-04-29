Martini + Multitemplate
=======================

Package multirender is a middleware for Martini that provides HTML
templates through multitemplate, it like to imitate and build on the
render package from martini-contrib.

```
package main

import (
  "github.com/codegangsta/martini"
  "github.com/cooldude/cache"

  // contians the multirender package.
  "github.com/acsellers/multitemplate/martini"
  // import any languages you want to use
  _ "github.com/acsellers/multitemplate/terse"
)

func main() {
  app := martini.Classic()
  app.Use(multirender.Renderer())

  app.Get("/", func (mr multirender.Render) {
    mr.HTML(200, "app/index.html", nil)
  })

  app.Get("/admins", func(mr multirender.Render) {
    ctx := mr.NewContext()
    ctx.RenderArgs["Users"] = AdminUsers
    mr.HTML(200, "app/user_list.html", ctx)
  })

  app.Run()
}
```

