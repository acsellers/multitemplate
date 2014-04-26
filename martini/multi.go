/*
  Package multirender is a middleware for Martini that provides HTML
  templates through multitemplate, it like to imitate and build on the
  render package from martini-contrib.

  package main

  import (
    "github.com/codegangsta/martini"

    // contians the multirender package.
    "github.com/acsellers/multitemplate/martini"
    // import any languages you want to use
    _ "github.com/acsellers/multitemplate/terse"
  )

  func main() {
    app := martini.Classic()
    app.Use(multirender.Renderer())

    app.Get("/html", func (mr multirender.Render) {
      mr.HTMLContext(multirender.Context{
        Main: "app/index.html",
      })
    })

    app.Get("/expensive", func(mr multirender.Render) {
      ctx := multirender.NewContext("app/expensive.html")
      if content, ok := cache.Get("expensive_block"); ok {
        ctx.SetContent("content", ok)
      }
    })
  }

*/
package multirender

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/acsellers/multitemplate"
	"github.com/acsellers/multitemplate/helpers"
	"github.com/codegangsta/martini"
)

type Render interface {
	// JSON writes the status code and JSON version of the value to the
	// http.ResponseWriter
	JSON(status int, v interface{})
	// XML writes the status code and XML version of the value to the
	// http.ResponseWriter
	XML(status int, v interface{})
	// HTML render the template given by name, with the status and options
	// given.
	HTML(status int, name string, v interface{}, htmlOpt *Context)
	// Render the text to the output
	Text(status int, text string)
	// Write the given status to the Response
	Error(status int)
	// Redirect to the location with an optional status, default status is
	// 302.
	Redirect(location string, status ...int)
	// Return the Template by name
	Template() *multitemplate.Template
	// Sets the content type, will also append charset from Options
	SetContentType(string)
}

func (r *renderer) NewContext() *Context {
	return &Context{
		Layout:     r.opt.DefaultLayout,
		Yields:     make(map[string]string),
		blocks:     make(map[string]multitemplate.RenderedBlock),
		RenderArgs: make(map[string]interface{}),
	}
}

type Context struct {
	Status     int
	Layout     string
	NoLayout   bool
	Yields     map[string]string
	blocks     map[string]multitemplate.RenderedBlock
	RenderArgs map[string]interface{}
}

func (c *Context) SetContent(name, content string) {
	if c.blocks == nil {
		c.blocks = make(map[string]multitemplate.RenderedBlock)
	}
	c.blocks[name] = multitemplate.RenderedBlock{
		Content: template.HTML(content),
	}
}

func (c *Context) SetTemplate(blockName, templateName string) {
	if c.Yields == nil {
		c.Yields = make(map[string]string)
	}
	c.Yields[blockName] = templateName
}

func Renderer(opt Options) martini.Handler {
	if opt.Charset == "" {
		opt.Charset = "utf-8"
	}
	mt := multitemplate.New("martini").Funcs(opt.Funcs)
	var err error
	mt = mt.Funcs(helpers.GetHelpers(opt.Helpers...))
	mt, err = compile(opt, mt)
	if err != nil {
		fmt.Println(err)
	}
	return func(w http.ResponseWriter, r *http.Request, c martini.Context) {
		if martini.Env == martini.Dev {
			mt, err = compile(opt, mt)
		}
		c.MapTo(&renderer{w, r, mt, opt, err}, (*Render)(nil))
	}
}

type Options struct {
	// Directories to search for template files
	Directories []string
	// Layout to render by default
	DefaultLayout string
	// Helper modules to load from multitemplate helpers
	Helpers []string
	// Additional functions to add
	Funcs template.FuncMap
	// JSON & XML indentation, an empty string disables indentation
	IndentEncoding string
	// Default is set to utf-8
	Charset string
	// Note that you will need to set the delims for each multitemplate
	// language you are using, you cannot set it in this Options struct.
}

func compile(opt Options, mt *multitemplate.Template) (*multitemplate.Template, error) {
	var err error
	for _, dir := range opt.Directories {
		mt.Base = dir
		e := filepath.Walk(dir, func(path string, i os.FileInfo, e error) error {
			if !i.IsDir() {
				_, err := mt.ParseFiles(path)
				return err
			}
			return nil
		})
		if e != nil && err == nil {
			err = e
		}
	}

	return mt, err
}

type renderer struct {
	http.ResponseWriter
	r   *http.Request
	mt  *multitemplate.Template
	opt Options
	err error
}

func (r *renderer) SetContentType(ct string) {
	r.Header().Set("Content-Type", ct+";charset="+r.opt.Charset)
}

// this is pretty much the same as the render packages
func (r *renderer) JSON(status int, v interface{}) {
	var result []byte
	var err error
	if r.opt.IndentEncoding != "" {
		result, err = json.MarshalIndent(v, "", r.opt.IndentEncoding)
	} else {
		result, err = json.Marshal(v)
	}
	if err != nil {
		http.Error(r, err.Error(), 500)
	}
	r.SetContentType("application/json")
	r.WriteHeader(status)
	r.Write(result)
}
func (r *renderer) XML(status int, v interface{}) {
	var result []byte
	var err error
	if r.opt.IndentEncoding != "" {
		result, err = xml.MarshalIndent(v, "", r.opt.IndentEncoding)
	} else {
		result, err = xml.Marshal(v)
	}
	if err != nil {
		http.Error(r, err.Error(), 500)
	}
	r.SetContentType("application/xml")
	r.WriteHeader(status)
	r.Write(result)
}
func (r *renderer) HTML(status int, name string, htmlOpt *Context) {
	ctx := multitemplate.NewContext(htmlOpt.RenderArgs)
	ctx.Main = name
	if !htmlOpt.NoLayout {
		ctx.Layout = htmlOpt.Layout
	}
	if len(htmlOpt.Yields) > 0 {
		ctx.Yields = htmlOpt.Yields
	}
	if len(htmlOpt.blocks) > 0 {
		ctx.Blocks = htmlOpt.blocks
	}

	b := &bytes.Buffer{}
	e := r.mt.ExecuteContext(b, ctx)
	if e != nil {
		http.Error(r, e.Error(), 500)
	}

}

func (r *renderer) Error(status int) {
	r.WriteHeader(status)
}

func (r *renderer) Redirect(location string, status ...int) {
	code := http.StatusFound
	if len(status) > 0 {
		code = status[0]
	}
	http.Redirect(r, r.r, location, code)
}

func (r *renderer) Template() *multitemplate.Template {
	return r.mt
}

func (r *renderer) Text(status int, text string) {
	r.WriteHeader(status)
	io.WriteString(r, text)
}
