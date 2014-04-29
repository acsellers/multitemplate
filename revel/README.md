Revel + Multitemplate
=====================

multitemplate is a package for executing templates from multiple template languages.
This package is the revel connector for multitemplate. Below is the details on how to
setup multitemplate so you can use it in your revel application.

app.conf
  // multitemplate must be added to your module list in your app.conf
  module.template=github.com/acsellers/multitemplate/revel
  module.jobs=github.com/revel/revel/modules/jobs

app/controllers/init.go

  func init() {
          // You can set DefaultLayouts for different content types, which are
          // found in your normal template paths
          multitemplate.DefaultLayout[multitemplate.HTML] = "layouts/app.html"

          // You must initialize the library in revel's OnAppStart function.
          revel.OnAppStart(multitemplate.Init)
  }

app/init.go

  // import the multitemplate/revel/app/controllers package and any sub-languages you want
  // to load
  import (
          _ "github.com/acsellers/multitemplate/bham"
          "github.com/acsellers/multitemplate/revel/app/controllers"
          "github.com/revel/revel"
  )



  // If you want to use the auto-refresh functionality in DevMode, you need to add the
  // multitemplate.ReloadFilter to the revel.Filters list.
  func init() {
          revel.Filters = []revel.Filter{
                  revel.PanicFilter,             // Recover from panics and display an error page instead.
                  revel.RouterFilter,            // Use the routing table to select the right Action
                  multitemplate.ReloadFilter,    // Filter for the multitemplate auto-refresh of templates
                  revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
                  ...

app/controllers/base_controller.go

  // Import the controllers area from multitemplate's revel connector.
  // Note that the package name that folder is multitemplate, so it is easier to know
  // where the Controller is from.
  import "github.com/acsellers/multitemplate/revel/app/controllers"

  // Simply replace any *revel.Controllers with multitemplate.Controllers
  type BaseController struct {
          multitemplate.Controller
          Txn *database.Transaction
  }


And that's it!

