package multitemplate

import (
	"bytes"
	"fmt"
	"html/template"
)

// A specialized Writer struct I'm using to make blocks work.
// I'm not sure whether I should keep this public, or make it private.
func newPouchWriter() *pouchWriter {
	p := &pouchWriter{}
	return p
}

type pouchWriter struct {
	names     []string
	root      bytes.Buffer
	buffers   []bytes.Buffer
	rulesets  []Ruleset
	discard   bool
	next      RenderedBlock
	check     bool
	immediate bool
	err       error
}

func (pw *pouchWriter) nesting() bool {
	return len(pw.buffers) > 0
}

func (pw *pouchWriter) Write(p []byte) (n int, err error) {
	if pw.check {
		pw.check = false
		var remainder []byte
		var rl Ruleset

		switch {
		case len(p) < 8:
			return 0, fmt.Errorf("Sentinel not received")
		case string(p[0:8]) == string(CSS):
			remainder = p[8:]
			rl = CSS
		case len(p) < 12:
			return 0, fmt.Errorf("Sentinel not received")
		case string(p[0:12]) == string(JS):
			remainder = p[12:]
			rl = JS
		case len(p) < 15:
			return 0, fmt.Errorf("Sentinel not received")
		case string(p[0:15]) == string(HTML):
			remainder = p[15:]
			rl = HTML
		default:
			return 0, fmt.Errorf("Sentinel not received")
		}
		pw.rulesets = append(pw.rulesets, rl)
		if rl == pw.next.Type || pw.next.Type == User {
			if pw.immediate {
				if len(pw.buffers) > 0 {
					return pw.buffers[len(pw.buffers)-1].Write([]byte(pw.next.Content))
				} else {
					pw.root.Write([]byte(pw.next.Content))
				}
			} else {
				if len(pw.buffers) > 1 {
					return pw.buffers[len(pw.buffers)-2].Write([]byte(pw.next.Content))
				} else {
					pw.root.Write([]byte(pw.next.Content))
				}
			}
		} else {
			pw.err = fmt.Errorf("Mismatched block contexts for block content: %s", pw.next.Content)
			return 0, pw.err
		}

		if len(remainder) != 0 {
			pw.Write(remainder)
		}

		return len(p), nil
	}

	if len(pw.buffers) > 0 {
		return pw.buffers[len(pw.buffers)-1].Write(p)
	}
	if !pw.discard {
		return pw.root.Write(p)
	}

	return len(p), nil
}

func (pw *pouchWriter) Nop(rb RenderedBlock) {
	pw.names = append(pw.names, "")
	pw.buffers = append(pw.buffers, bytes.Buffer{})
	pw.check = true
	pw.next = rb
}

func (pw *pouchWriter) Immediate(rb RenderedBlock) {
	pw.check = true
	pw.immediate = true
	pw.next = rb
}

func (pw *pouchWriter) NoRoot() {
	pw.discard = true
}

func (pw *pouchWriter) Reset() {
	pw.discard = false
	pw.names = []string{}
	pw.buffers = []bytes.Buffer{}
	pw.rulesets = []Ruleset{}
}

func (pw *pouchWriter) Open(name string) {
	pw.names = append(pw.names, name)
	pw.buffers = append(pw.buffers, bytes.Buffer{})
	pw.check = true
}

func (pw *pouchWriter) Close() (name string, rb RenderedBlock) {
	if len(pw.names) > 0 {
		name = pw.names[len(pw.names)-1]
		content := pw.buffers[len(pw.buffers)-1].String()
		rb = RenderedBlock{Content: template.HTML(content), Type: pw.rulesets[len(pw.rulesets)-1]}
		pw.names = pw.names[:len(pw.names)-1]
		pw.buffers = pw.buffers[:len(pw.buffers)-1]
	}
	return
}
