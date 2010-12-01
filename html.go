package hubbard

import "io"
import "html"

type htmlWriter struct {
	w         io.Writer
	closeTags []string
	first     bool
}

func newHtmlWriter(w io.Writer) *htmlWriter {
	return &htmlWriter{w, make([]string, 0), false}
}

func (self *htmlWriter) raw(html string) *htmlWriter {
	self.checkFirst()
	io.WriteString(self.w, html)
	return self
}

func (self *htmlWriter) start(name string) *htmlWriter {
	self.checkFirst()
	io.WriteString(self.w, `<` + name)
	self.push(name)
	self.first = true
	return self
}

func (self *htmlWriter) with(name string, value string) *htmlWriter {
	if !self.first {
		panic("attributes aren't valid here")
	}
	io.WriteString(self.w, ` ` + name + `="` + html.EscapeString(value) + `"` )
	return self
}

func (self *htmlWriter) div(class string) *htmlWriter {
	return self.start("div")
}

func (self *htmlWriter) table() *htmlWriter {
	return self.start("table")
}

func (self *htmlWriter) tr() *htmlWriter {
	return self.start("tr")
}

func (self *htmlWriter) td() *htmlWriter {
	return self.start("td")
}

func (self *htmlWriter) text(text string) *htmlWriter {
	self.checkFirst()
	io.WriteString(self.w, html.EscapeString(text))
	return self
}

func (self *htmlWriter) end() *htmlWriter {
	self.checkFirst()
	io.WriteString(self.w, `</` + self.pop() + `>`)
	return self
}

func (self *htmlWriter) checkFirst() {
	if self.first {
		io.WriteString(self.w, ">")
		self.first = false
	}
}

func (self *htmlWriter) push(closeTag string) {
	self.closeTags = append(self.closeTags, closeTag)
}

func (self *htmlWriter) pop() string {
	pos := len(self.closeTags) - 1
	closeTag := self.closeTags[pos]
	self.closeTags[pos] = ""
	self.closeTags = self.closeTags[0:pos]
	return closeTag
}

