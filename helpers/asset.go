package helpers

import (
	"html/template"
	"net/url"
	"path/filepath"
)

var assetFuncs = template.FuncMap{
	"atom_link": func(link string) template.HTML {
		al := AttrList{
			"rel":   "alternate",
			"type":  "application/atom+xml",
			"title": "ATOM",
			"href":  AppInfo.FullLink(link),
		}
		return buildTag("link", "", al)
	},
	"favicon_link": func(filename string) template.HTML {
		al := AttrList{"rel": "shortcut icon"}
		switch filepath.Ext(filename) {
		case ".png":
			al["type"] = "image/png"
			al["href"] = AppInfo.ImageLink(filename)
		case ".ico":
			al["type"] = "image/vnd.microsoft.icon"
			al["href"] = AppInfo.ImageLink(filename)
		default:
			al["href"] = AppInfo.ImageLink(filename)
		}
		return buildTag("link", "", al)
	},
	"image_tag": func(filename string, options ...AttrList) template.HTML {
		al := combine("", "", options)
		al["src"] = AppInfo.ImageLink(filename)
		return buildTag("img", "", al)
	},
	"javascript_link": func(filenames ...string) template.HTML {
		content := ""
		al := AttrList{}
		for _, filename := range filenames {
			al["src"] = AppInfo.JavascriptLink(filename)
			content += string(buildTag("script", " ", al))
		}
		return template.HTML(content)
	},
	"root_url": func(link string) template.HTML {
		return template.HTML(AppInfo.RootURL.String())
	},
	"rss_link": func(link string) template.HTML {
		al := AttrList{
			"rel":   "alternate",
			"type":  "application/rss+xml",
			"title": "RSS",
			"href":  AppInfo.FullLink(link),
		}
		return buildTag("link", "", al)
	},
	"stylesheet_link": func(filenames ...string) template.HTML {
		content := ""
		al := AttrList{
			"media": "screen",
			"rel":   "stylesheet",
		}
		for _, filename := range filenames {
			al["href"] = AppInfo.StylesheetLink(filename)
			content += string(buildTag("link", "", al))
		}

		return template.HTML(content)
	},
}

type AssetInfo struct {
	RootURL            *url.URL
	DocRoot            string
	ImageRelative      string
	JavascriptRelative string
	StylesheetRelative string
}

/*
Example AppInfo
    helpers.AppInfo = AssetInfo{
      RootURL:        url.Parse("http://localhost"),
      DocRoot:            "assets",
      ImageRelative:      "images",
      JavascriptRelative: "javascripts",
      StylesheetRelative: "stylesheets",
    }
*/
var AppInfo = AssetInfo{}

// FullLink will determine whether the url given is a full path
// and if not will add the root url for the app to the path.
func (ai AssetInfo) FullLink(path string) string {
	tu, _ := url.Parse(path)
	if tu.Host == "" {
		tu, _ = url.Parse(ai.RootURL.String())
		tu.Path = path
	}
	return tu.String()
}

// AssetLink takes a path and a prefix, then puts them
// together without doubling up on /'s
func (ai AssetInfo) AssetLink(path, prefix string) string {
	tu, _ := url.Parse(path)
	// Use relative urls for local urls
	if tu.Host == "" {
		if len(prefix) == 0 || prefix[0] != '/' {
			prefix = "/" + prefix
		}
		if path[0] != '/' {
			path = "/" + path
		}
		return prefix + path
	}

	// return the url in path otherwise
	return tu.String()
}

// ImageLink returns a link for the filename, you must include
// the extension for the file (png, jpg, gif).
func (ai AssetInfo) ImageLink(filename string) string {
	return ai.AssetLink(filename, ai.ImageRelative)
}

// JavascriptLink returns a link for the filename, if you do
// not include the .js extension it will be appended for you.
func (ai AssetInfo) JavascriptLink(filename string) string {
	if filepath.Ext(filename) != ".js" {
		filename += ".js"
	}
	return ai.AssetLink(filename, ai.JavascriptRelative)
}

// StylesheetLink returns a link for the filename, if you do
// not include the .css extension, it will be appended.
func (ai AssetInfo) StylesheetLink(filename string) string {
	if filepath.Ext(filename) != ".css" {
		filename += ".css"
	}
	return ai.AssetLink(filename, ai.StylesheetRelative)
}
