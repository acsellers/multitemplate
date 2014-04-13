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
      AssetRoot:      filepath.Join(os.Getwd(), "public"),
      ImageRelative:      "images",
      JavascriptRelative: "javascripts",
      StylesheetRelative: "stylesheets",
    }
*/
var AppInfo = AssetInfo{}

func (ai AssetInfo) FullLink(path string) string {
	tu, _ := url.Parse(path)
	if tu.Host == "" {
		tu, _ = url.Parse(ai.RootURL.String())
		tu.Path = path
	}
	return tu.String()
}

func (ai AssetInfo) AssetLink(path, prefix string) string {
	tu, _ := url.Parse(path)
	// Use relative urls for local urls
	if tu.Host == "" {
		if prefix[0] != '/' {
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

func (ai AssetInfo) ImageLink(filename string) string {
	return ai.AssetLink(filename, ai.ImageRelative)
}

func (ai AssetInfo) JavascriptLink(filename string) string {
	if filepath.Ext(filename) != ".js" {
		filename += ".js"
	}
	return ai.AssetLink(filename, ai.JavascriptRelative)
}

func (ai AssetInfo) StylesheetLink(filename string) string {
	if filepath.Ext(filename) != ".css" {
		filename += ".css"
	}
	return ai.AssetLink(filename, ai.StylesheetRelative)
}