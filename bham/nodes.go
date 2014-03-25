package bham

const (
	identRaw = iota
	identFilter
	identText
	identExecutable
	identTag
	identTagOpen
	identTagClose
	identIf
	identRange
	identWith
)

func (pt *protoTree) insertRaw(content string, level int) {
	pt.currNodes = append(pt.currNodes, protoNode{
		level:      level,
		identifier: identRaw,
		content:    content,
	})
}

func (pt *protoTree) insertFilter(content string, level int, handler FilterHandler) {
	pt.currNodes = append(pt.currNodes, protoNode{
		level:      level,
		identifier: identFilter,
		content:    content,
		filter:     handler,
	})
}

func (pt *protoTree) insertText(line templateLine) {
	pt.currNodes = append(pt.currNodes, protoNode{
		level:      line.indentation,
		identifier: identText,
		content:    line.content,
	})
}

func (pt *protoTree) insertIf(statement string, level int, ifNodes, elseNodes []protoNode) {
	pt.currNodes = append(pt.currNodes, protoNode{
		level:      level,
		identifier: identIf,
		content:    statement,
		list:       ifNodes,
		elseList:   elseNodes,
	})
}
func (pt *protoTree) insertRange(statement string, level int, rangeNodes, elseNodes []protoNode) {
	pt.currNodes = append(pt.currNodes, protoNode{
		level:      level,
		identifier: identRange,
		content:    statement,
		list:       rangeNodes,
		elseList:   elseNodes,
	})
}
func (pt *protoTree) insertWith(statement string, level int, innerNodes []protoNode) {
	pt.currNodes = append(pt.currNodes, protoNode{
		level:      level,
		identifier: identWith,
		content:    statement,
		list:       innerNodes,
	})
}
func (pt *protoTree) insertExecutable(statement string, level int) {
	pt.currNodes = append(pt.currNodes, protoNode{
		level:      level,
		identifier: identExecutable,
		content:    statement,
	})
}
