package bham

import (
	"testing"

	"github.com/acsellers/assert"
)

func TestAnalysis1(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		pt := &protoTree{
			name: "example.html",
			source: `!!!
%html
  %head`,
		}
		pt.lex()
		test.IsNil(pt.err)
		test.AreEqual(3, len(pt.lineList))
		pt.analyze()
		test.IsNil(pt.err)
		test.AreEqual(4, len(pt.nodes))
		test.AreEqual("<!DOCTYPE html>", pt.nodes[0].content)
		test.AreEqual("%html", pt.nodes[1].content)
		test.AreEqual(identTagOpen, pt.nodes[1].identifier)
		test.AreEqual("%head", pt.nodes[2].content)
		test.AreEqual(identTag, pt.nodes[2].identifier)
		test.AreEqual("%html", pt.nodes[3].content)
	})
}

func TestAnalysis2(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		pt := &protoTree{
			name: "example.html",
			source: `!!!
= if .HTML
  %html
    %head
      %title Test
= else
  %xhtml`,
		}
		pt.lex()
		test.IsNil(pt.err)
		if pt.err == nil {
			test.AreEqual(7, len(pt.lineList))
			pt.analyze()
			test.IsNil(pt.err)
			test.AreEqual(2, len(pt.nodes))
			test.AreEqual("<!DOCTYPE html>", pt.nodes[0].content)
			test.AreEqual(".HTML", pt.nodes[1].content)
			test.AreEqual(identIf, pt.nodes[1].identifier)
			test.AreEqual(5, len(pt.nodes[1].list))
			test.AreEqual("%html", pt.nodes[1].list[0].content)
			test.AreEqual("%head", pt.nodes[1].list[1].content)
			test.AreEqual("%title Test", pt.nodes[1].list[2].content)
			test.AreEqual("%head", pt.nodes[1].list[3].content)
			test.AreEqual("%html", pt.nodes[1].list[4].content)
			test.AreEqual(1, len(pt.nodes[1].elseList))
		}
	})
}

func TestAnalysis3(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		pt := &protoTree{
			name: "example.html",
			source: `!!!
%html
  = .Script`,
		}
		pt.lex()
		test.IsNil(pt.err)
		test.AreEqual(3, len(pt.lineList))
		pt.analyze()
		test.IsNil(pt.err)
		test.AreEqual(4, len(pt.nodes))
		test.AreEqual("<!DOCTYPE html>", pt.nodes[0].content)
		test.AreEqual("%html", pt.nodes[1].content)
		test.AreEqual(identTagOpen, pt.nodes[1].identifier)
		test.AreEqual(".Script", pt.nodes[2].content)
		test.AreEqual("%html", pt.nodes[3].content)
	})
}

func TestAnalysis4(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		pt := &protoTree{
			name: "example.html",
			source: `!!!
= if .HTML
  %html
    %head
      %title Test`,
		}
		pt.lex()
		test.IsNil(pt.err)
		if pt.err == nil {
			test.AreEqual(5, len(pt.lineList))
			pt.analyze()
			test.IsNil(pt.err)
			test.AreEqual(2, len(pt.nodes))
			test.AreEqual("<!DOCTYPE html>", pt.nodes[0].content)
			test.AreEqual(".HTML", pt.nodes[1].content)
			test.AreEqual(identIf, pt.nodes[1].identifier)
			test.AreEqual(5, len(pt.nodes[1].list))
			test.AreEqual("%html", pt.nodes[1].list[0].content)
			test.AreEqual("%head", pt.nodes[1].list[1].content)
			test.AreEqual("%title Test", pt.nodes[1].list[2].content)
			test.AreEqual("%head", pt.nodes[1].list[3].content)
			test.AreEqual("%html", pt.nodes[1].list[4].content)
		}
	})
}

func TestAnalysis5(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		pt := &protoTree{
			name: "example.html",
			source: `= range .List
  .name= .
%head
  %title Test`,
		}
		pt.lex()
		test.IsNil(pt.err)
		if pt.err == nil {
			test.AreEqual(4, len(pt.lineList))
			pt.analyze()
			test.IsNil(pt.err)
			test.AreEqual(4, len(pt.nodes))
			test.AreEqual(".List", pt.nodes[0].content)
			test.AreEqual(identRange, pt.nodes[0].identifier)
			test.AreEqual("%head", pt.nodes[1].content)
			test.AreEqual("%title Test", pt.nodes[2].content)
			test.AreEqual("%head", pt.nodes[3].content)
		}
	})
}

func TestAnalysis6(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		pt := &protoTree{
			name: "example.html",
			source: `= $var := hello "andrew"
%title= $var`,
		}
		pt.lex()
		test.IsNil(pt.err)
		if pt.err == nil {
			test.AreEqual(2, len(pt.lineList))
			test.AreEqual(0, pt.lineList[1].indentation)
			pt.analyze()
			test.IsNil(pt.err)
			test.AreEqual(2, len(pt.nodes))
			test.AreEqual("$var := hello \"andrew\"", pt.nodes[0].content)
			test.AreEqual("%title= $var", pt.nodes[1].content)
		}
	})
}

func TestAnalysis7(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		pt := &protoTree{
			name: "example.html",
			source: `%first
  = if true
    %now Hello
%end`,
		}
		pt.lex()
		test.IsNil(pt.err)
		if pt.err == nil {
			test.AreEqual(4, len(pt.lineList))
			test.AreEqual(1, pt.lineList[1].indentation)
			test.AreEqual(2, pt.lineList[2].indentation)
			test.AreEqual(0, pt.lineList[3].indentation)
			pt.analyze()
			test.IsNil(pt.err)
			test.AreEqual(4, len(pt.nodes))
			test.AreEqual("%first", pt.nodes[0].content)
			test.AreEqual("true", pt.nodes[1].content)
			test.AreEqual(1, len(pt.nodes[1].list))
			test.AreEqual("%first", pt.nodes[2].content)
			test.AreEqual("%end", pt.nodes[3].content)
		}
	})
}
