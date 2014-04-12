package filters

import (
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"strings"
	"text/template"
	"unicode"

	"github.com/acsellers/inflections"
	"github.com/davecgh/go-spew/spew"
	"github.com/tebeka/strftime"
)

// These functions come from Django

var filterFuncs = template.FuncMap{
	"add_slashes": func(s string) string {
		output := make([]byte, len(s)*2)
		index := 0
		for _, val := range []byte(s) {
			if val == '"' || val == '\'' {
				output[index] = '\\'
				index++
			}
			output[index] = val
			index++
		}
		return string(output[:index])
	},
	"cap_first": func(s string) string {
		var seenLetter bool
		output := []rune{}
		for _, r := range s {
			if !seenLetter && unicode.IsLetter(r) {
				seenLetter = true
				r = unicode.ToUpper(r)
			}
			output = append(output, r)
		}

		return string(output)
	},
	"center": func(s string, num int) string {
		sRunes := []rune(s)
		if len(sRunes) >= num {
			return s
		}
		nRunes := make([]rune, num)
		for i := 0; i < num; i++ {
			nRunes[i] = ' '
		}
		offset := (len(nRunes) - len(sRunes)) / 2
		for i := 0; i < len(sRunes); i++ {
			nRunes[i+offset] = sRunes[i]
		}

		return string(nRunes)
	},
	"cut": func(s, cutset string) string {
		output := []rune{}
		var put bool
		for _, r := range s {
			put = true
			for _, c := range cutset {
				if c == r {
					put = false
				}
			}
			if put {
				output = append(output, r)
			}
		}
		return string(output)
	},
	"default": func(x, y interface{}) interface{} {
		switch v := x.(type) {
		case string:
			if v == "" {
				return y
			}
		case bool:
			if v == false {
				return y
			}
		default:
			rv := reflect.ValueOf(x)
			if !rv.IsValid() {
				return y
			}
			if rv.Kind() == reflect.Ptr && rv.IsNil() {
				return y
			}
			if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Map {
				if rv.Len() == 0 {
					return y
				}
			}
		}
		return x
	},
	"default_if_nil": func(x, y interface{}) interface{} {
		rv := reflect.ValueOf(x)
		if !rv.IsValid() {
			return y
		}
		if rv.Kind() == reflect.Ptr && rv.IsNil() {
			return y
		}
		return x
	},
	"escape": func(i interface{}) template.HTML {
		switch v := i.(type) {
		case template.HTML:
			return v
		case string:
			return template.HTML(template.HTMLEscapeString(v))
		}
		return ""
	},
	"escapejs": template.JSEscapeString,
	"filesize_format": func(i int64) string {
		var r int64
		if i == 1 {
			return fmt.Sprintf("1 Byte")
		}
		if i < 1024 {
			return fmt.Sprintf("%d Bytes", i)
		}
		r = i % 1024
		i = i / 1024
		if i < 1024 {
			return fmt.Sprintf("%d.%d KB", i, r/102)
		}
		r = i % 1024
		i = i / 1024
		if i < 1024 {
			return fmt.Sprintf("%d.%d MB", i, r/102)
		}
		r = i % 1024
		i = i / 1024
		if i < 1024 {
			return fmt.Sprintf("%d.%d GB", i, r/102)
		}
		r = i % 1024
		i = i / 1024
		if i < 1024 {
			return fmt.Sprintf("%d.%d TB", i, r/102)
		}
		r = i % 1024
		i = i / 1024
		if i < 1024 {
			return fmt.Sprintf("%d.%d EB", i, r/102)
		}
		return "0 Bytes"
	},

	"first": func(l interface{}) interface{} {
		lv := reflect.ValueOf(l)
		if lv.Kind() == reflect.Slice {
			if lv.Len() > 0 {
				return lv.Index(0).Interface()
			} else {
				return nil
			}
		}
		return l
	},
	"float_format": func(f float64, n int) string {
		return fmt.Sprintf("%."+fmt.Sprint(n)+"f", f)
	},

	"force_escape": template.HTMLEscapeString,
	"get_digit": func(i interface{}, num int) string {
		s := fmt.Sprint(i)
		if num <= 0 {
			return ""
		}
		if len(s) >= num {
			return string(s[len(s)-num])
		}
		return ""
	},
	"join": func(l interface{}, j string) string {
		lv := reflect.ValueOf(l)
		if lv.Kind() == reflect.Slice && lv.Len() > 0 {
			output := make([]string, lv.Len())
			for i := 0; i < lv.Len(); i++ {
				output[i] = fmt.Sprint(lv.Index(i).Interface())
			}
			return strings.Join(output, j)
		}
		return fmt.Sprint(l)
	},
	"last": func(l interface{}) interface{} {
		lv := reflect.ValueOf(l)
		if lv.Kind() == reflect.Slice {
			if lv.Len() > 0 {
				return lv.Index(lv.Len() - 1).Interface()
			} else {
				return nil
			}
		}
		return l
	},
	"length": func(l interface{}) int {
		lv := reflect.ValueOf(l)
		k := lv.Kind()
		if k == reflect.Slice || k == reflect.Array || k == reflect.Map || k == reflect.String {
			return lv.Len()
		}
		return 1
	},
	"length_is": func(l interface{}, n int) bool {
		lv := reflect.ValueOf(l)
		k := lv.Kind()
		if k == reflect.Slice || k == reflect.Array || k == reflect.Map || k == reflect.String {
			return lv.Len() == n
		}
		return 1 == n
	},
	"link_to": func(link, name string) template.HTML {
		u, e := url.Parse(link)
		if e != nil {
			return ""
		}

		return template.HTML(
			fmt.Sprintf(`<a href="%s">%s</a>`, u.String(), name),
		)
	},
	"ljust": func(s string, num int) string {
		sRunes := []rune(s)
		if len(sRunes) >= num {
			return s
		}
		nRunes := make([]rune, num)
		for i := 0; i < num; i++ {
			nRunes[i] = ' '
		}
		for i := 0; i < len(sRunes); i++ {
			nRunes[i] = sRunes[i]
		}

		return string(nRunes)
	},
	"lower": strings.ToLower,
	"number_lines": func(s string) string {
		output := strings.Split(s, "\n")
		for i, line := range output {
			output[i] = fmt.Sprintf("%d. %s", i+1, line)
		}
		return strings.Join(output, "\n")
	},
	"pluralize": inflections.Pluralize,
	"pprint":    spew.Sprint,
	"pytime":    strftime.Format,
	"quick_format": func(s string) template.HTML {
		return template.HTML(strings.Replace(string(template.HTMLEscapeString(s)), "\n", "<br>", -1))
	},
	"random": func(l interface{}) interface{} {
		lv := reflect.ValueOf(l)
		if lv.Kind() == reflect.Slice && lv.Len() > 0 {
			return lv.Index(rand.Int() % lv.Len()).Interface()
		}
		return nil
	},
	"rjust": func(i interface{}, num int) string {
		s := fmt.Sprint(i)
		sRunes := []rune(s)
		if len(sRunes) >= num {
			return s
		}
		nRunes := make([]rune, num)
		for i := 0; i < num; i++ {
			nRunes[i] = ' '
		}
		offset := len(nRunes) - len(sRunes)
		for i := 0; i < len(sRunes); i++ {
			nRunes[i+offset] = sRunes[i]
		}

		return string(nRunes)
	},
	"safe": func(s string) template.HTML {
		return template.HTML(s)
	},
	"safeseq": func(l []string) []template.HTML {
		output := make([]template.HTML, len(l))
		for i, v := range l {
			output[i] = template.HTML(v)
		}
		return output
	},

	"title": strings.ToTitle,

	"truncate": func(s string, n int) string {
		if len(s) <= n {
			return s
		}
		output := make([]rune, n)
		var current int
		for _, r := range s {
			output[current] = r
			current++
			if current == n {
				for i := current - 3; i < current; i++ {
					output[i] = '.'
				}
				return string(output)
			}
		}
		return s[:n]
	},
	"upper":      strings.ToUpper,
	"url_encode": url.QueryEscape,
	"urlize": func(link string) template.HTML {
		u, e := url.Parse(link)
		if e != nil {
			return ""
		}

		return template.HTML(
			fmt.Sprintf(`<a href="%s">%s</a>`, u.String(), link),
		)
	},
	"urlize_truncate": func(link string, num int) template.HTML {
		chars := []rune(link)
		u, e := url.Parse(link)
		if e != nil {
			return ""
		}
		if num >= len(chars) {

			return template.HTML(
				fmt.Sprintf(`<a href="%s">%s</a>`, u.String(), link),
			)
		}
		return template.HTML(
			fmt.Sprintf(`<a href="%s">%s</a>`, u.String(), chars[:num]),
		)
	},
}
