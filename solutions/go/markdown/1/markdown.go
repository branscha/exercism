package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

// This exercise is waste of time. If I encountered such code, I would start from scratch.
// How far can you "refactor" the code without it being new, what is the limit here?
// The question should be: "rewrite this stuff completely"

type MdInput struct {
	markdown string
	lnPos    int
	mdPos    int
}

func (m *MdInput) Next() byte {
	if m.markdown[m.mdPos] == '\n' {
		m.lnPos = 0
	} else {
		m.lnPos++
	}
	m.mdPos++
	if m.mdPos < len(m.markdown) {
		return m.markdown[m.mdPos]
	} else {
		return ' '
	}
}

func (m *MdInput) Peek() byte {
	if m.mdPos < len(m.markdown) {
		return m.markdown[m.mdPos]
	} else {
		return '$'
	}
}

func (m *MdInput) HasNext() bool {
	return m.mdPos < len(m.markdown)
}

func (m *MdInput) IsStartOfLine() bool {
	return m.lnPos == 0
}

func (m *MdInput) Pos() int {
	return m.mdPos
}

var listItemNext *regexp.Regexp = regexp.MustCompile(`(?s)\*.*$`)

func (m *MdInput) ListItemFollows() bool {
	return listItemNext.MatchString(m.markdown[m.mdPos:])
}

type MdState struct {
	headerLevel int
	liNr        int
	liOpened    bool
	html        string
}

func (m *MdState) Emit(s string) {
	m.html += s
}

func (m *MdState) String() string {
	return m.html
}

func (m *MdState) ParseHeader(md *MdInput) {
	char := md.Peek()
	for char == '#' {
		m.headerLevel++
		char = md.Next()
	}

	if m.headerLevel >= 7 {
		// Unsupported header level, just copy the text
		m.Emit(fmt.Sprintf("<p>%s ", strings.Repeat("#", m.headerLevel)))
	} else {
		// Supported header level 1-6
		m.Emit(fmt.Sprintf("<h%d>", m.headerLevel))
	}

	char = md.Peek()
	if char == ' ' {
		char = md.Next()
	}
}

func (m *MdState) ParseListItem(md *MdInput) {
	if m.liNr == 0 {
		m.Emit("<ul>")
	}

	m.liNr++

	if !m.liOpened {
		m.Emit("<li>")
		m.liOpened = true
	}

	// Gobble space after *
	char := md.Next()
	if char == ' ' {
		char = md.Next()
	}
}

func (m *MdState) HandleLineEnd(md *MdInput) {
	// Eat the newline
	_ = md.Next()

	// Close list item
	if m.liOpened {
		m.Emit("</li>")
		m.liOpened = false
	}

	// Close list
	if !md.ListItemFollows() {
		m.Emit("</ul><p>")
		m.liOpened = false
		m.liNr = 0
	}

	// Close header
	if m.headerLevel > 0 {
		m.Emit(fmt.Sprintf("</h%d>", m.headerLevel))
		m.headerLevel = 0
	}
}

func (m *MdState) HandleDocumentEnd(md *MdInput) string {
	switch {
	case m.headerLevel >= 7:
		m.Emit("</p>")
		return m.String()
	case m.headerLevel > 0:
		m.Emit(fmt.Sprintf("</h%d>", m.headerLevel))
		return m.String()
	}

	if m.liNr > 0 {
		m.Emit("</li></ul>")
		return m.String()
	}

	if strings.Contains(m.String(), "<p>") {
		m.Emit("</p>")
		return m.String()
	}

	return "<p>" + m.String() + "</p>"
}

// Render translates markdown to HTML
func Render(markdown string) string {
	markdown = processStyles(markdown)
	md := MdInput{markdown, 0, 0}
	st := MdState{}
	for md.HasNext() {
		char := md.Peek()
		if char == '#' && md.IsStartOfLine() {
			st.ParseHeader(&md)
		} else if char == '*' && md.IsStartOfLine() {
			st.ParseListItem(&md)
		} else if char == '\n' || !md.HasNext() {
			st.HandleLineEnd(&md)
		} else {
			st.Emit(string(char))
			char = md.Next()
		}
	}
	return st.HandleDocumentEnd(&md)
}

func processStyles(markdown string) string {
	markdown = strings.Replace(markdown, "__", "<strong>", 1)
	markdown = strings.Replace(markdown, "__", "</strong>", 1)
	markdown = strings.Replace(markdown, "_", "<em>", 1)
	markdown = strings.Replace(markdown, "_", "</em>", 1)
	return markdown
}
