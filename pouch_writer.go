package multitemplate

import "bytes"

func NewPouchWriter() *PouchWriter {
	p := &PouchWriter{}
	return p
}

type PouchWriter struct {
	names   []string
	root    bytes.Buffer
	buffers []bytes.Buffer
	discard bool
}

func (pw *PouchWriter) Write(p []byte) (n int, err error) {
	if len(pw.buffers) > 0 {
		return pw.buffers[len(pw.buffers)-1].Write(p)
	}
	if !pw.discard {
		return pw.root.Write(p)
	}

	return len(p), nil
}

func (pw *PouchWriter) Nop() {
	pw.names = append(pw.names, "")
	pw.buffers = append(pw.buffers, bytes.Buffer{})
}

func (pw *PouchWriter) NoRoot() {
	pw.discard = true
}

func (pw *PouchWriter) Reset() {
	pw.discard = false
	pw.names = []string{}
	pw.buffers = []bytes.Buffer{}
}

func (pw *PouchWriter) Open(name string) {
	pw.names = append(pw.names, name)
	pw.buffers = append(pw.buffers, bytes.Buffer{})
}

func (pw *PouchWriter) Close() (name, content string) {
	if len(pw.names) > 0 {
		name = pw.names[len(pw.names)-1]
		content = pw.buffers[len(pw.buffers)-1].String()
		pw.names = pw.names[:len(pw.names)-1]
		pw.buffers = pw.buffers[:len(pw.buffers)-1]
	}
	return
}
