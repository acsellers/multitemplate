package mustache

import "strings"

type stash struct {
	tree       *protoTree
	content    string
	started    bool
	commenting bool
}

func (s *stash) needsMoreText() bool {
	normalOpen := strings.Index(s.content, s.tree.localLeft)
	normalUnescape := strings.Index(s.content, LeftEscapeDelim)

	if normalUnescape >= 0 && normalUnescape < normalOpen {
		closeIndex := strings.Index(
			s.content[normalUnescape+len(LeftEscapeDelim):],
			RightEscapeDelim,
		)
		return closeIndex == -1
	}
	if normalOpen >= 0 {
		closeIndex := strings.Index(
			s.content[normalOpen+len(s.tree.localLeft):],
			s.tree.localRight,
		)
		return closeIndex == -1
	}

	return false
}

func (s *stash) Append(t string) {
	ts := strings.TrimSpace(t)
	//standalone comments
	if s.commenting {
		if strings.HasSuffix(ts, s.tree.localRight) {
			s.commenting = false
		}
		return
	}
	if strings.HasPrefix(ts, s.tree.localLeft) && strings.Count(ts, s.tree.localLeft) == 1 {
		if strings.TrimSpace(ts[len(s.tree.localLeft):])[0] == '!' {
			if strings.HasSuffix(ts, s.tree.localRight) {
				return
			} else {
				s.commenting = true
				return
			}
		}
	}
	// sections opening and closing

	if strings.HasPrefix(ts, s.tree.localLeft) {
		if strings.Count(ts, s.tree.localLeft) == 1 ||
			(s.tree.localLeft == s.tree.localRight && strings.Count(ts, s.tree.localLeft) == 2) {
			if strings.HasSuffix(ts, s.tree.localRight) {
				switch strings.TrimSpace(ts[len(s.tree.localLeft):])[0] {
				case '#', '/', '^', '=':
					s.content = s.content + ts
					t = ""
				case '>':
					if strings.HasSuffix(t, ts+"\n") {
						s.content = s.content + t[:len(t)-1]
						t = ""
					}
				}
			}
		}
	}
	s.content = s.content + t
}
func (s *stash) hasAction() bool {
	return strings.Contains(s.content, s.tree.localLeft) ||
		strings.Contains(s.content, LeftEscapeDelim)
}

func (s *stash) pullToAction() (string, string) {
	var text, action string
	loc, abnormal := s.nextActionLocation()
	text = s.content[:loc]
	s.content = s.content[loc:]
	if abnormal {
		action = s.content[:len(LeftEscapeDelim)]
		s.content = s.content[len(LeftEscapeDelim):]
		closeLocation := strings.Index(s.content, RightEscapeDelim)
		action += s.content[:closeLocation+len(RightEscapeDelim)]
		s.content = s.content[closeLocation+len(RightEscapeDelim):]
	} else {
		action = s.content[:len(s.tree.localLeft)]
		s.content = s.content[len(s.tree.localLeft):]
		closeLocation := strings.Index(s.content, s.tree.localRight)
		action += s.content[:closeLocation+len(s.tree.localRight)]
		s.content = s.content[closeLocation+len(s.tree.localRight):]
	}
	return text, action
}

func (s *stash) nextActionLocation() (int, bool) {
	normalOpen := strings.Index(s.content, s.tree.localLeft)
	normalUnescape := strings.Index(s.content, LeftEscapeDelim)

	if normalUnescape >= 0 && normalUnescape <= normalOpen {
		return normalUnescape, true
	}
	if normalOpen >= 0 {
		return normalOpen, false
	}
	return -1, false
}
