package bham

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const (
	pse_text = iota
	pse_tag
	pse_if
	pse_else
	pse_end
	pse_range
	pse_with
	pse_decl
	pse_exe
)

var (
	tag      = regexp.MustCompile("^%([a-zA-Z0-9]+)")
	varDecl  = regexp.MustCompile("^\\$[a-zA-Z0-9]+ :=")
	idClass  = regexp.MustCompile("^([\\.#][a-zA-Z0-9-_]+)")
	idClass2 = regexp.MustCompile("([\\.#][a-zA-Z0-9-_]+)")
	varUse   = regexp.MustCompile("(\\$[a-zA-Z0-9]+)")
)

func (pt *protoTree) tokenize() error {
	posts := make([]token, 0, 64)
	scanner := bufio.NewScanner(bytes.NewBufferString(pt.source))
	var text, currentTag string
	var currentLevel, currentLine, lineLevel int
	for scanner.Scan() {
		currentLine++

		text = scanner.Text()
		if text == "" || strings.TrimSpace(text) == "" {
			continue
		}

		lineLevel, text = level(text)
		if currentLevel > lineLevel {
			for currentLevel >= lineLevel && currentLevel > 0 {
				pt.tokenList = append(
					pt.tokenList,
					posts[len(posts)-1],
				)
				posts = posts[:len(posts)-1]
				currentLevel--
			}
		}
		if lineLevel-1 > currentLevel {
			return fmt.Errorf("Line %d is indented more than necessary (%d) from the previous line %d", currentLine, lineLevel, currentLevel)
		}

		if tag.MatchString(text) {
			currentTag = tag.FindStringSubmatch(text)[1]
			pt.tokenList = append(
				pt.tokenList,
				token{
					content: "<" + currentTag + ">",
					purpose: pse_tag,
				},
			)

			posts = append(posts, token{
				content: "</" + currentTag + ">",
				purpose: pse_tag,
			})
			text = text[len(currentTag)+1:]
			if text != "" && idClass.MatchString(text) {
				pt.tokenList[len(pt.tokenList)-1].attrs = make(map[string][]string)
				for _, submatch := range idClass2.FindAllStringSubmatch(text, -1) {
					switch submatch[0][0] {
					case '.':
						pt.tokenList[len(pt.tokenList)-1].attrs["class"] = append(pt.tokenList[len(pt.tokenList)-1].attrs["class"], submatch[0][1:])
					case '#':
						pt.tokenList[len(pt.tokenList)-1].attrs["id"] = append(pt.tokenList[len(pt.tokenList)-1].attrs["id"], submatch[0][1:])
					}
				}
				for idClass.MatchString(text) {
					text = text[len(idClass.FindStringSubmatch(text)[0]):]
				}
			}
			if text != "" && text[0] == '(' && strings.Contains(text, ")") {
				var attrs string
				attrs, text = findAttrs(text)
				if attrs != "" {
					pt.tokenList[len(pt.tokenList)-1].extra = attrs
				}
			}
		} else {
			if text[0] == ':' {
				pt.tokenList = append(pt.tokenList, shortHandOpen(text[1:]))
				posts = append(posts, shortHandClose(text[1:]))
				currentLevel = lineLevel
				text = ""
			}

			if text != "" && idClass.MatchString(text) {
				pt.tokenList = append(pt.tokenList, token{content: "<div>", purpose: pse_tag})
				posts = append(posts, token{content: "</div>", purpose: pse_tag})
				pt.tokenList[len(pt.tokenList)-1].attrs = make(map[string][]string)
				for _, submatch := range idClass2.FindAllStringSubmatch(text, -1) {
					switch submatch[0][0] {
					case '.':
						pt.tokenList[len(pt.tokenList)-1].attrs["class"] = append(pt.tokenList[len(pt.tokenList)-1].attrs["class"], submatch[0][1:])
					case '#':
						pt.tokenList[len(pt.tokenList)-1].attrs["id"] = append(pt.tokenList[len(pt.tokenList)-1].attrs["id"], submatch[0][1:])
					}
				}
				text = text[len(idClass.FindStringSubmatch(text)[0]):]
			}
			if text != "" && text[0] == '(' && strings.Contains(text, ")") {
				var attrs string
				attrs, text = findAttrs(text)
				if attrs != "" {
					pt.tokenList[len(pt.tokenList)-1].extra = attrs
				}
			}

			if strings.HasPrefix(text, "=") {
				trimText := strings.TrimSpace(text[len("="):])
				switch {
				case strings.HasPrefix(trimText, "if "):
					pt.tokenList = append(
						pt.tokenList,
						token{
							content: strings.TrimPrefix(trimText, "if "),
							purpose: pse_if,
						},
					)
					posts = append(posts, token{purpose: pse_end, previous: len(pt.tokenList) - 1})
				case trimText == "else":
					pt.tokenList[len(pt.tokenList)-1].purpose = pse_else
					posts = append(posts, token{purpose: pse_end, previous: len(pt.tokenList) - 1})
				case strings.HasPrefix(trimText, "range "):
					pt.tokenList = append(
						pt.tokenList,
						token{
							content: strings.TrimPrefix(trimText, "range "),
							purpose: pse_range,
						},
					)
					posts = append(posts, token{purpose: pse_end, previous: len(pt.tokenList) - 1})
				case strings.HasPrefix(trimText, "with "):
					pt.tokenList = append(
						pt.tokenList,
						token{
							content: strings.TrimPrefix(trimText, "with "),
							purpose: pse_with,
						},
					)
					posts = append(posts, token{purpose: pse_end, previous: len(pt.tokenList) - 1})
				case varDecl.MatchString(trimText):
					pt.tokenList = append(
						pt.tokenList,
						token{
							content: trimText,
							purpose: pse_decl,
						},
					)
				default:
					pt.tokenList = append(pt.tokenList, token{content: trimText, purpose: pse_exe})
				}
			} else {
				pt.tokenList = append(pt.tokenList, token{content: text})
			}
			text = ""
		}
		if text == "" {
			currentLevel = lineLevel
			continue
		}
		if text[0] == ' ' {
			pt.tokenList = append(pt.tokenList, token{content: text[1:]})
			currentLevel = lineLevel
			continue
		}
	}
	for len(posts) > 0 {
		pt.tokenList = append(pt.tokenList, posts[len(posts)-1])
		posts = posts[:len(posts)-1]
	}

	return nil
}

func (pt *protoTree) compact() {
	newTokenList := make([]token, 0, len(pt.tokenList))
	var appendToken token
	appendToken.purpose = pse_text
	for _, currentToken := range pt.tokenList {
		if currentToken.purpose == pse_text {
			appendToken.content += currentToken.content
		} else {
			newTokenList = append(newTokenList,
				appendToken,
				currentToken,
			)
			appendToken.content = ""
		}
	}
	pt.tokenList = newTokenList
}

type token struct {
	content  string
	purpose  int
	previous int
	extra    string
	attrs    map[string][]string
}

func (t token) parent() int {
	switch t.purpose {
	case pse_end:
		return t.previous
	case pse_else:
		return t.previous
	default:
		return -1
	}
}

func (t token) textual() bool {
	return t.purpose == pse_text || t.purpose == pse_tag
}

func (t token) strcontent() string {
	if t.purpose == pse_tag {
		if t.attrs != nil {
			if vals, ok := t.attrs["id"]; ok {
				if strings.Contains(t.extra, `id="`) {
					idIndex := strings.Index(t.extra, `id="`)
					t.extra = t.extra[:idIndex+4] + strings.Join(vals, IdJoin) + IdJoin + t.extra[idIndex+4:]
				} else {
					t.extra = t.extra + " id=\"" + strings.Join(vals, IdJoin) + "\""
				}
			}
			if vals, ok := t.attrs["class"]; ok {
				if strings.Contains(t.extra, "class=\"") {
					classIndex := strings.Index(t.extra, `class="`)
					t.extra = t.extra[:classIndex+7] + strings.Join(vals, " ") + " " + t.extra[classIndex+7:]
				} else {
					t.extra = t.extra + " class=\"" + strings.Join(vals, " ") + "\""
				}
			}
		}
		if t.extra != "" {
			return t.content[0:len(t.content)-1] + " " + strings.TrimSpace(t.extra) + ">"
		}
	}
	return t.content
}
