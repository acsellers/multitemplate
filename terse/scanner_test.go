package terse

import (
	"fmt"
	"testing"
)

func TestScanner(t *testing.T) {
	for _, test := range scannerTests {
		tree := scan(test.Content)
		if test.RootLen != len(tree.Children) {
			fmt.Println("In Test: " + test.Name)
			fmt.Println("Tree Root Nodes did not match expected: %v vs %v", test.RootLen, len(tree.Children))
			t.Fail()
			continue
		}
		if !test.Tabs && tree.String() != test.Content {
			fmt.Println("In Test: " + test.Name)
			fmt.Println("Tree was not output the same as the input")
			fmt.Println("Expected:")
			fmt.Println(test.Content)
			fmt.Println("Returned:")
			fmt.Println(tree.String())
			t.Fail()
			continue
		}
	}
}

type scannerTest struct {
	Name    string
	Content string
	Tabs    bool
	RootLen int
}

var scannerTests = []scannerTest{
	scannerTest{
		Name:    "Blank Tree",
		Content: "",
		RootLen: 0,
	},
	scannerTest{
		Name:    "Single Line",
		Content: "!!",
		RootLen: 1,
	},
	scannerTest{
		Name:    "Two Lines Unindented",
		Content: "!!\nhtml",
		RootLen: 2,
	},
	scannerTest{
		Name:    "Two Lines Indented",
		Content: "html\n  head",
		RootLen: 1,
	},
	scannerTest{
		Name:    "Two Lines with Tabs",
		Tabs:    true,
		Content: "html\n\thead",
		RootLen: 1,
	},
	scannerTest{
		Name:    "Three Level Indentation",
		Content: "html\n  head\n    title= .Title",
		RootLen: 1,
	},
	scannerTest{
		Name:    "Mixed spaces and tabs",
		Tabs:    true,
		Content: "html\n  head\n\tbody",
		RootLen: 1,
	},
	scannerTest{
		Name:    "Mixed Indentations with multiple roots",
		Content: "[one]\n  blah\n[two]\n  first\n    second",
		RootLen: 2,
	},
}
