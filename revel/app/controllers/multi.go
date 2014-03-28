/*
 */
package multitemplate

import (
	"bytes"
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
	DefaultLayout map[RequestFormat]string
	Template      *mt.Template
	// In DevMode, multitemplate will automatically
	// refresh templates using revel's Watcher struct,
	// in ProductionMode, you can change ProdRefresh to
	// true to get that behavior. Maybe you're doing
	// something special with FUSE and template folders
	// so you would need this. RefreshPaths will not start
	// watching folders in DevMode without this.
	ProdRefresh bool
	extraPaths  []string
	refresh     *templateRefresher
	watch       *revel.Watcher
)

type RequestFormat string

func init() {
	if revel.DevMode || ProdRefresh {
		watch := revel.NewWatcher()
		refresh = &templateRefresher{}
		watch.Listen(refresh, revel.TemplatePaths...)
	}
	e := RefreshTemplates()
	if e != nil {
		revel.ERROR.Println(e)
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
	e := RefreshTemplates()
	if e != nil {
		revel.ERROR.Println(e)
	}
}

// This should happen automatically in
func RefreshTemplates() error {
	revel.WARN.Println("Start refresh of templates")
	Template = mt.New("revel_root")
	Template.Funcs(revel.TemplateFuncs)

	var err error
	for _, rp := range revel.TemplatePaths {
		Template.Base = rp
		filepath.Walk(rp, func(path string, i os.FileInfo, e error) error {
			if !i.IsDir() {
				_, err = Template.ParseFiles(path)
				if err != nil {
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
					return err
				}
			}
			return nil
		})
	}

	revel.WARN.Println("End refresh of templates")
	return nil
}

type templateRefresher struct{}

func (tr *templateRefresher) Refresh() *revel.Error {
	RefreshTemplates()
	return nil
}

const (
	HTML RequestFormat = "html"
	XML  RequestFormat = "xml"
	JSON RequestFormat = "json"
	TXT  RequestFormat = "txt"
)

func init() {
	DefaultLayout = make(map[RequestFormat]string)
}

type Controller struct {
	*revel.Controller
	layout          string
	nolayout        bool
	yields, content map[string]string
}

func (c *Controller) SetLayout(name string) {
	if name == "" {
		c.nolayout = true
	} else {
		c.layout = name
	}
}

func (c *Controller) ContentFor(name, templateName string) {
	if c.yields == nil {
		c.yields = make(map[string]string)
	}
	c.yields[name] = templateName
}

func (c *Controller) Render(extraRenderArgs ...interface{}) revel.Result {
	RefreshTemplates()
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
		ctx.Content[key] = template.HTML(content)
	}

	if ctx.Layout == "" && DefaultLayout[RequestFormat(c.Request.Format)] != "" && !c.nolayout {
		ctx.Layout = DefaultLayout[RequestFormat(c.Request.Format)]
	}

	ctx.Main = c.Name + "/" + c.MethodType.Name + "." + c.Request.Format

	return &MultiTemplateResult{ctx}
}

type MultiTemplateResult struct {
	ctx *mt.Context
}

func (mtr *MultiTemplateResult) Apply(req *revel.Request, resp *revel.Response) {
	// Handle panics when rendering templates.
	/*
		defer func() {
			if err := recover(); err != nil {
				revel.ERROR.Println(err)
				revel.PlaintextErrorResult{fmt.Errorf("Template Execution Panic in %s:\n%s",
					mtr.ctx.Main, err)}.Apply(req, resp)
			}
		}()
	*/

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

	revel.WARN.Println("start execute")
	revel.WARN.Println(mtr.ctx)
	if Template.Lookup(mtr.ctx.Layout) == nil {
		mtr.ctx.Layout = strings.ToLower(mtr.ctx.Layout)
	}
	if Template.Lookup(mtr.ctx.Main) == nil {
		mtr.ctx.Main = strings.ToLower(mtr.ctx.Main)
	}
	e := Template.ExecuteContext(&b, mtr.ctx)
	revel.WARN.Println(e)
	revel.WARN.Println("end execute")

	if !chunked {
		resp.Out.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	}
	resp.WriteHeader(http.StatusOK, "text/html")
	b.WriteTo(out)

}
