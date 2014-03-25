//
package assets

import (
	"html/template"
	"path"
	"strings"

	"github.com/acsellers/helpers/utils"
)

var (
	Functions = map[string]interface{}{
		"audio_path":             AudioPath,
		"audio_tag":              AudioTag,
		"favicon_link_tag":       FaviconLinkTag,
		"font_path":              FontPath,
		"image_path":             ImagePath,
		"image_tag":              ImageTag,
		"javascript_path":        JavascriptPath,
		"javascript_src_tag":     JavascriptSrcTag,
		"javascript_include_tag": JavascriptIncludeTag,
		"stylesheet_link_tag":    StylesheetLinkTag,
		"stylesheet_path":        StylesheetPath,
		//		"video_path":             VideoPath,
		//		"video_tag":              VideoTag,
	}

	Paths = map[string]string{
		"javascript": "/public/js/",
		"css":        "/public/css",
		"image":      "/public/img",
		"audio":      "/public/audio/",
		"font":       "/public/fonts/",
		"video":      "/public/video/",
	}
	PrefixedPaths = map[string]map[string]string{
		"javascript": map[string]string{
			"asset": "/assets/js",
		},
		"css": map[string]string{
			"asset": "/assets/css",
		},
		"image": map[string]string{
			"kittens": "http://placekitten.com/",
			"dogs":    "http://placedog.com/",
		},
	}

	defaultMaps = map[string]map[string]string{
		"javascript": map[string]string{
			"type": "text/javascript",
		},
		"css": map[string]string{
			"type":  "text/css",
			"rel":   "stylesheet",
			"media": "screen",
		},
	}
)

func AudioPath(filename string) string {
	return pathFor(filename, "audio")
}

func AudioTag(filename string, options ...string) template.HTML {
	optionsMap := utils.Optionize(options)
	optionsMap["src"] = AudioPath(filename)
	return utils.Tag("audio", optionsMap)
}

func FaviconLinkTag(name string, options ...string) template.HTML {
	optionsMap := utils.Optionize(options)
	optionsMap["href"] = name
	switch path.Ext(name) {
	case "ico":
		optionsMap["rel"] = "shortcut icon"
	case "png", "apng":
		optionsMap["rel"] = "icon"
		optionsMap["type"] = "image/png"
	case "gif":
		optionsMap["rel"] = "icon"
		optionsMap["type"] = "image/gif"
	case "":
		optionsMap["rel"] = "shortcut icon"
		optionsMap["href"] = optionsMap["href"] + ".ico"
	case "svg":
		optionsMap["rel"] = "icon"
		optionsMap["type"] = "image/svg"
	case "jpg", "jpeg":
		optionsMap["rel"] = "icon"
		optionsMap["type"] = "image/jpeg"
	default:
		optionsMap["rel"] = "shortcut icon"
	}

	return utils.Tag("link", optionsMap)
}

func FontPath(name string) string {
	return pathFor(name, "font")
}

func ImagePath(name string) string {
	return pathFor(name, "image")
}

func ImageTag(name string, options ...string) template.HTML {
	optionsMap := utils.Optionize(options)
	optionsMap["src"] = ImagePath(name)
	return utils.Tag("img", optionsMap)
}

func JavascriptPath(name string) string {
	return pathFor(name, "javascript")
}

func JavascriptSrcTag(name string, options ...string) template.HTML {
	optionsMap := utils.OptionizeWithDefaults(options, defaultMaps["javascript"])
	optionsMap["src"] = JavascriptPath(name)

	return utils.Tag("script", optionsMap)
}

func JavascriptIncludeTag(names ...string) template.HTML {
	tags := make([]string, len(names))
	options := make(map[string]string)
	for key, val := range defaultMaps["javascript"] {
		options[key] = val
	}

	for i, name := range names {
		options["src"] = JavascriptPath(names)
		tags[i] = string(utils.Tag("script", options))
	}

	return template.HTML(strings.Join(tags))
}

func StylesheetPath(name string) string {
	return pathFor(name, "css")
}

func StylesheetLinkTag(items ...string) template.HTML {
	sheets, optionsMap := utils.StrictOptionizeWithDefaults(more, defaultMaps["css"])
	tags := make([]string, len(sheets))

	for i, sheet := range sheets {
		optionsMap["href"] = StylesheetPath(sheet)
		tags[i] = string(utils.Tag("link", optionsMap))
	}

	return template.HTML(strings.Join(tags))
}

/*
Video elements need some more work because of sources/tracks,
I need to decide whether to support that type of functionality,
ignore that type, or to leave that to users.

func VideoPath(name string) string {
	return pathFor(name, "video")
}

func VideoTag(items ...string) template.HTML {

}
*/

// interal utils
func pathFor(name, pather string) string {
	if strings.Contains(name, "://") {
		return name
	}
	prefix, unprefixed := isPrefixed(pather, name)
	if prefix != "" {
		return path.Join(PrefixedPaths[pather][prefix], unprefixed)
	}
	if len(name) > 0 && name[0] == '/' {
		return name
	}
	if _, ok := Paths[pather]; ok {
		return path.Join(Paths[pather], name)
	}

	return path.Join("/", name)
}
