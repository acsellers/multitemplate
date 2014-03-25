package bham

import (
	"testing"

	"github.com/acsellers/assert"
)

func TestLex1(t *testing.T) {
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
		test.AreEqual(0, pt.lineList[0].indentation)
		test.AreEqual("!!!", pt.lineList[0].content)
		test.AreEqual(0, pt.lineList[1].indentation)
		test.AreEqual("%html", pt.lineList[1].content)
		test.AreEqual(1, pt.lineList[2].indentation)
		test.AreEqual("%head", pt.lineList[2].content)
	})
}

func TestLex2(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		pt := &protoTree{
			name: "example.html",
			source: `!!!
%html
  %head(one="two" \
    three="four")`,
		}
		pt.lex()
		test.IsNil(pt.err)
		test.AreEqual(3, len(pt.lineList))
		test.AreEqual(0, pt.lineList[0].indentation)
		test.AreEqual("!!!", pt.lineList[0].content)
		test.AreEqual(0, pt.lineList[1].indentation)
		test.AreEqual("%html", pt.lineList[1].content)
		test.AreEqual(1, pt.lineList[2].indentation)
		test.AreEqual("%head(one=\"two\" three=\"four\")", pt.lineList[2].content)
	})
}

func TestLex3(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		pt := &protoTree{
			name: "example.html",
			source: `!!!
%html
  = insert_head "ng-app" \
    "ng-controller:superCtrl"`,
		}
		pt.lex()
		test.IsNil(pt.err)
		test.AreEqual(3, len(pt.lineList))
		test.AreEqual(0, pt.lineList[0].indentation)
		test.AreEqual("!!!", pt.lineList[0].content)
		test.AreEqual(0, pt.lineList[1].indentation)
		test.AreEqual("%html", pt.lineList[1].content)
		test.AreEqual(1, pt.lineList[2].indentation)
		test.AreEqual("= insert_head \"ng-app\" \"ng-controller:superCtrl\"", pt.lineList[2].content)
	})
}

func TestLex4(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		pt := &protoTree{
			name: "example.html",
			source: `!!!
%html
  insert_head "ng-app" \
  "ng-controller:superCtrl"`,
		}
		pt.lex()
		test.IsNil(pt.err)
		test.AreEqual(4, len(pt.lineList))
		test.AreEqual(0, pt.lineList[0].indentation)
		test.AreEqual("!!!", pt.lineList[0].content)
		test.AreEqual(0, pt.lineList[1].indentation)
		test.AreEqual("%html", pt.lineList[1].content)
		test.AreEqual(1, pt.lineList[2].indentation)
		test.AreEqual("insert_head \"ng-app\" \\", pt.lineList[2].content)
	})
}
