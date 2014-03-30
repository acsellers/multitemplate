/*
multitemplate is a package that allows multiple templates written in multiple
template languages to call between each other. This is the revel connector for
multitemplate. See godoc.org/github.com/acsellers/multitemplate/revel for the
integration instructions.
*/
package multitemplate

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	mt "github.com/acsellers/multitemplate"
	"github.com/revel/revel"
)

var (
	// DefaultLayout allows you to set a layout that will automatically be added
	// per content type.
	DefaultLayout = make(map[RequestFormat]string)
	// Template is the template loader used by multitemplate. It will be replaced each
	// time the templates a refreshed if you are using auto-refresh.
	Template *mt.Template
	// In DevMode, multitemplate will automatically
	// refresh templates using revel's Watcher struct,
	// in ProductionMode, you can change ProdRefresh to
	// true to get that behavior. Maybe you're doing
	// something special with FUSE and template folders
	// so you would need this. RefreshPaths will not start
	// watching folders in DevMode without this.
	ProdRefresh bool
	// CurrentError will record any errors encountered when loading templates, so it
	// can be displayed when rendering the page.
	CurrentError error
	extraPaths   []string
	refresh      *templateRefresher
	watch        *revel.Watcher
)

type RequestFormat string

// Init must be called from within revel's OnAppStart method.
func Init() {
	revel.INFO.Println("Multitemplate starting")

	if revel.DevMode || ProdRefresh {
		watch = revel.NewWatcher()
		refresh = &templateRefresher{}
		revel.INFO.Println("multitemplate watches start")
		watch.Listen(refresh, revel.TemplatePaths...)
		revel.INFO.Println("multitemplate watches stop")
	}

	CurrentError = RefreshTemplates()
	if CurrentError != nil {
		revel.ERROR.Println(CurrentError)
	}
}

// If you have added paths for templates beyond the paths accessible
// in the array revel.TemplatePaths, call this function to pick up
// those directories autmatically. You can also pass extra paths to
// this function if you prefer it, and multitemplate will also
// save and watch those paths.
func RefreshPaths(morePaths ...string) {
MoreLoop:
	for _, p := range morePaths {
		for _, wp := range extraPaths {
			if p == wp {
				continue MoreLoop
			}
		}
		extraPaths = append(extraPaths, p)
	}
	if revel.DevMode || ProdRefresh {
		watch.Listen(refresh, revel.TemplatePaths...)
		watch.Listen(refresh, extraPaths...)
	}
	CurrentError = RefreshTemplates()
}

// RefreshTemplates is automatically called by the auto-refresh system
// or when you call RefreshPaths.
func RefreshTemplates() error {
	revel.INFO.Println("Start multitemplate refresh")
	Template = mt.New("revel_root")
	Template.Funcs(revel.TemplateFuncs)

	var err error
	for _, rp := range revel.TemplatePaths {
		Template.Base = rp
		filepath.Walk(rp, func(path string, i os.FileInfo, e error) error {
			if !i.IsDir() {
				_, err = Template.ParseFiles(path)
				if err != nil {
					revel.ERROR.Printf("Parse Error in %s: %v", path, err)
					CurrentError = err
					return err
				}
			}
			return nil
		})
	}
	for _, ep := range extraPaths {
		Template.Base = ep
		filepath.Walk(ep, func(path string, i os.FileInfo, e error) error {
			if !i.IsDir() {
				_, err = Template.ParseFiles(path)
				if err != nil {
					revel.ERROR.Printf("Parse Error in %s: %v", path, err)
					CurrentError = err
					return err
				}
			}
			return nil
		})
	}
	if CurrentError != nil {
		return CurrentError
	}
	revel.INFO.Println("Multitemplate refresh completed successfully")
	return nil
}

type templateRefresher struct{}

func (tr *templateRefresher) Refresh() *revel.Error {
	revel.INFO.Println("multitemplate: refreshing templates")
	CurrentError = RefreshTemplates()
	return nil
}

const (
	HTML RequestFormat = "html"
	XML  RequestFormat = "xml"
	JSON RequestFormat = "json"
	TXT  RequestFormat = "txt"
)

// Controller can be added to your controller structs in your revel app
// to use the multitemplate system.
type Controller struct {
	*revel.Controller
	layout   string
	nolayout bool
	yields   map[string]string
	content  map[string]template.HTML
}

// SetLayout sets the layout to be executed for this action. Set it to
// the empty string to disable the layout for this action.
func (c *Controller) SetLayout(name string) {
	if name == "" {
		c.nolayout = true
	} else {
		c.layout = name
	}
}

// Block sets a pre-rendered HTML string that can show up in a yield
// or block.
func (c *Controller) Block(name string, content template.HTML) {
	if c.content == nil {
		c.content = make(map[string]template.HTML)
	}
	c.content[name] = content
}

// ContentFor sets a template to be rendered for a key, this can be used
// by either a block call or a yield call.
func (c *Controller) ContentFor(name, templateName string) {
	if c.yields == nil {
		c.yields = make(map[string]string)
	}
	c.yields[name] = templateName
}

// Standard Render call, this renders the default template for this action.
func (c *Controller) Render(extraRenderArgs ...interface{}) revel.Result {
	// Get the calling function name.
	_, _, line, ok := runtime.Caller(1)
	if !ok {
		revel.ERROR.Println("Failed to get Caller information")
	}

	// Get the extra RenderArgs passed in.
	if renderArgNames, ok := c.MethodType.RenderArgNames[line]; ok {
		if len(renderArgNames) == len(extraRenderArgs) {
			for i, extraRenderArg := range extraRenderArgs {
				c.RenderArgs[renderArgNames[i]] = extraRenderArg
			}
		} else {
			revel.ERROR.Println(len(renderArgNames), "RenderArg names found for",
				len(extraRenderArgs), "extra RenderArgs")
		}
	} else {
		revel.ERROR.Println("No RenderArg names found for Render call on line", line,
			"(Method", c.MethodType.Name, ")")
	}

	ctx := mt.NewContext(c.RenderArgs)
	ctx.Layout = c.layout
	if len(c.yields) > 0 {
		ctx.Yields = c.yields
	}
	for key, content := range c.content {
		ctx.Content[key] = content
	}

	if ctx.Layout == "" && DefaultLayout[RequestFormat(c.Request.Format)] != "" && !c.nolayout {
		ctx.Layout = DefaultLayout[RequestFormat(c.Request.Format)]
	}

	ctx.Main = c.Name + "/" + c.MethodType.Name + "." + c.Request.Format

	if CurrentError != nil {
		return c.RenderError(CurrentError)
	}

	return &templateResult{ctx}
}

// RenderTemplate renders a specific template by path. If a DefaultLayout value
// is available for this content type, it will be filled in automatically.
func (c *Controller) RenderTemplate(templateName string) revel.Result {
	ctx := mt.NewContext(c.RenderArgs)
	ctx.Layout = c.layout
	if len(c.yields) > 0 {
		ctx.Yields = c.yields
	}
	for key, content := range c.content {
		ctx.Content[key] = content
	}

	if ctx.Layout == "" && DefaultLayout[RequestFormat(c.Request.Format)] != "" && !c.nolayout {
		ctx.Layout = DefaultLayout[RequestFormat(c.Request.Format)]
	}

	ctx.Main = templateName

	if CurrentError != nil {
		return c.RenderError(CurrentError)
	}

	return &templateResult{ctx}
}

type templateResult struct {
	ctx *mt.Context
}

func (mtr *templateResult) Apply(req *revel.Request, resp *revel.Response) {
	// Handle panics when rendering templates.
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Println(err)
			revel.PlaintextErrorResult{fmt.Errorf("Template Execution Panic in %s:\n%s",
				mtr.ctx.Main, err)}.Apply(req, resp)
		}
	}()

	chunked := revel.Config.BoolDefault("results.chunked", false)

	// If it's a HEAD request, throw away the bytes.
	out := io.Writer(resp.Out)
	if req.Method == "HEAD" {
		out = ioutil.Discard
	}

	// In a prod mode, write the status, render, and hope for the best.
	// (In a dev mode, always render to a temporary buffer first to avoid having
	// error pages distorted by HTML already written)
	if chunked && !revel.DevMode {
		resp.WriteHeader(http.StatusOK, "text/html")

		Template.ExecuteContext(resp.Out, mtr.ctx)
		return
	}

	// Render the template into a temporary buffer, to see if there was an error
	// rendering the template.  If not, then copy it into the response buffer.
	// Otherwise, template render errors may result in unpredictable HTML (and
	// would carry a 200 status code)
	var b bytes.Buffer

	if Template.Lookup(mtr.ctx.Layout) == nil {
		mtr.ctx.Layout = strings.ToLower(mtr.ctx.Layout)
	}
	if Template.Lookup(mtr.ctx.Main) == nil {
		mtr.ctx.Main = strings.ToLower(mtr.ctx.Main)
	}
	e := Template.ExecuteContext(&b, mtr.ctx)
	if e != nil {
		er := &revel.ErrorResult{mtr.ctx.Dot.(map[string]interface{}), e}
		er.Apply(req, resp)
		return
	}

	if !chunked {
		resp.Out.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	}
	resp.WriteHeader(http.StatusOK, "text/html")
	b.WriteTo(out)

}

var ReloadFilter = func(c *revel.Controller, fc []revel.Filter) {
	if watch != nil {
		watch.Notify()
	}

	fc[0](c, fc[1:])
}
