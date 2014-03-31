package multitemplate

import "bytes"

// A specialized Writer struct I'm using to make blocks work.
// I'm not sure whether I should keep this public, or make it private.
func newPouchWriter() *pouchWriter {
	p := &pouchWriter{}
	return p
}

type pouchWriter struct {
	names   []string
	root    bytes.Buffer
	buffers []bytes.Buffer
	discard bool
}

func (pw *pouchWriter) nesting() bool {
	return len(pw.buffers) > 0
}

func (pw *pouchWriter) Write(p []byte) (n int, err error) {
	if len(pw.buffers) > 0 {
		return pw.buffers[len(pw.buffers)-1].Write(p)
	}
	if !pw.discard {
		return pw.root.Write(p)
	}

	return len(p), nil
}

func (pw *pouchWriter) Nop() {
	pw.names = append(pw.names, "")
	pw.buffers = append(pw.buffers, bytes.Buffer{})
}

func (pw *pouchWriter) NoRoot() {
	pw.discard = true
}

func (pw *pouchWriter) Reset() {
	pw.discard = false
	pw.names = []string{}
	pw.buffers = []bytes.Buffer{}
}

func (pw *pouchWriter) Open(name string) {
	pw.names = append(pw.names, name)
	pw.buffers = append(pw.buffers, bytes.Buffer{})
}

func (pw *pouchWriter) Close() (name, content string) {
	if len(pw.names) > 0 {
		name = pw.names[len(pw.names)-1]
		content = pw.buffers[len(pw.buffers)-1].String()
		pw.names = pw.names[:len(pw.names)-1]
		pw.buffers = pw.buffers[:len(pw.buffers)-1]
	}
	return
}
