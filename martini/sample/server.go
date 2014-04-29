// This example is converted from my redup package
// github.com/acsellers/redup
package main

import (
	"encoding/hex"
	"html/template"
	"os"
	"path/filepath"

	"github.com/codegangsta/martini"
	"github.com/simonz05/godis/redis"

	// contians the multirender package.
	"github.com/acsellers/multitemplate/martini"
	// import any languages you want to use
	_ "github.com/acsellers/multitemplate/terse"
)

var Conn *redis.Client

func main() {
	Conn = redis.New("tcp:127.0.0.1:6379", 0, "")

	os.Setenv("PORT", "3456")

	templateRoot := filepath.Join(os.Getenv("GOPATH"), "src/github.com/acsellers/multitemplate/martini/sample/templates")
	app := martini.Classic()
	app.Use(multirender.Renderer(multirender.Options{
		Directories:   []string{templateRoot},
		DefaultLayout: "layout.html",
		Helpers:       []string{"all"},
		Funcs: template.FuncMap{
			"to_hex": func(s string) string {
				return hex.EncodeToString([]byte(s))
			},
		},
	}))

	app.Use(
		martini.Static(
			filepath.Join(templateRoot, "assets"),
			martini.StaticOptions{Prefix: "assets"},
		),
	)

	app.Get("/", func(mr multirender.Render) {
		ctx := mr.NewContext()
		ctx.RenderArgs["Keys"] = AllKeys()
		ctx.RenderArgs["Values"] = KeyValues()
		mr.HTML(200, "index.html", ctx)
	})

	app.Get("/show/:name", func(params martini.Params, mr multirender.Render) {
		key := HtmlIdToKey(params["name"])
		ctx := mr.NewContext()
		ctx.RenderArgs["Content"] = ContentFor(key)
		mr.HTML(200, "show.html", ctx)
	})

	app.Run()
}

func AllKeys() []string {
	keys, err := Conn.Keys("*")
	if err == nil {
		return keys
	}
	return []string{}
}

func HtmlIdToKey(hash string) string {
	val, err := hex.DecodeString(hash)
	if err == nil {
		return string(val)
	}
	return ""
}
func KeyValues() []Value {
	res := []Value{}
	for _, key := range AllKeys() {
		res = append(res, ContentFor(key))
	}
	return res
}

type Value struct {
	Key      string
	IsList   bool
	Content  string
	Contents []string
}

func (v Value) Link() string {
	return hex.EncodeToString([]byte(v.Key))
}
func ContentFor(key string) Value {
	val := Value{Key: key}
	info, err := Conn.Type(key)
	if err == nil {
		switch info {
		case "string":
			v, e := Conn.Get(key)
			if e == nil {
				val.Content = v.String()
				return val
			}
			val.Content = e.Error()
			return val
		case "list":
			list, e := Conn.Lrange(key, 0, -1)
			if e == nil {
				val.IsList = true
				val.Contents = list.StringArray()
				return val
			}
			val.Content = e.Error()
			return val
		case "set":
			list, e := Conn.Smembers(key)
			if e == nil {
				val.IsList = true
				val.Contents = list.StringArray()
				return val
			}
			val.Content = e.Error()
			return val
		default:
			val.Content = info
			return val
		}
	}
	val.Content = err.Error()
	return val
}
