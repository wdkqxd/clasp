package brave

import (
	"html"
	"regexp"
	"strings"

	"clasp/internal/engine"
	"clasp/internal/engine/helpers"
)

const (
	maxTopAnswerBlocks = 40
	maxTopAnswerText   = 8000
)

var (
	reHref        = regexp.MustCompile(`href="(https?://[^"]+)"`)
	reTitleAnchor = regexp.MustCompile(`(?s)<a[^>]*class="title[^"]*"[^>]*href="([^"]+)"[^>]*>(.*?)</a>`)
	reCite        = regexp.MustCompile(`(?s)<cite[^>]*>.*?</cite>`)
	reComments    = regexp.MustCompile(`(?s)<!--.*?-->`)
	reTags        = regexp.MustCompile(`(?s)<[^>]*>`)
)

func parseTopAnswer(doc, query string) *engine.RichAnswer {
	box := elementByID(doc, `id="codebox"`)
	if box == "" {
		return nil
	}
	answer := stripCite(elementByID(box, `id="codebox-answer"`))
	blocks := answerBlocks(answer)
	if len(blocks) == 0 {
		return nil
	}

	ta := &engine.RichTopAnswer{
		Question: helpers.CleanText(query),
		Blocks:   blocks,
	}
	if m := reTitleAnchor.FindStringSubmatch(box); m != nil {
		ta.Title = stripTags(m[2])
		if ta.Title != "" {
			ta.Question = ta.Title
		}
		ta.URL = html.UnescapeString(m[1])
	} else if m := reHref.FindStringSubmatch(box); m != nil {
		ta.URL = html.UnescapeString(m[1])
	}
	if host := hostOf(ta.URL); host != "" {
		ta.MetaURL = &engine.RichMetaURL{Hostname: host}
	}
	return &engine.RichAnswer{Type: "topAnswer", TopAnswer: ta}
}

func stripTags(s string) string {
	s = reComments.ReplaceAllString(s, "")
	s = reTags.ReplaceAllString(s, " ")
	return collapseText(html.UnescapeString(s))
}

func stripCite(s string) string {
	if loc := reCite.FindAllStringIndex(s, -1); len(loc) > 0 {
		last := loc[len(loc)-1]
		if strings.TrimSpace(s[last[1]:]) == "" {
			return s[:last[0]]
		}
	}
	return s
}

func elementByID(doc, id string) string {
	i := strings.Index(doc, id)
	if i < 0 {
		return ""
	}
	open := strings.LastIndex(doc[:i], "<div")
	if open < 0 {
		return ""
	}
	return elementInner(doc, open, "div")
}

func hostOf(raw string) string {
	i := strings.Index(raw, "://")
	if i < 0 {
		return ""
	}
	rest := raw[i+3:]
	if j := strings.IndexAny(rest, "/?#"); j >= 0 {
		rest = rest[:j]
	}
	return strings.ToLower(rest)
}

func elementInner(doc string, start int, name string) string {
	i := strings.Index(doc[start:], "<"+name)
	if i < 0 {
		return ""
	}
	i += start
	j := strings.IndexByte(doc[i:], '>')
	if j < 0 {
		return ""
	}
	from := i + j + 1
	body := from
	depth := 1
	for from < len(doc) {
		k := strings.IndexByte(doc[from:], '<')
		if k < 0 {
			break
		}
		k += from
		e := strings.IndexByte(doc[k:], '>')
		if e < 0 {
			break
		}
		e += k
		if tagName(doc[k+1:e]) == name {
			if strings.HasPrefix(doc[k+1:e], "/") {
				depth--
				if depth == 0 {
					return doc[body:k]
				}
			} else if !strings.HasSuffix(strings.TrimSpace(doc[k+1:e]), "/") {
				depth++
			}
		}
		from = e + 1
	}
	return doc[body:]
}

func answerBlocks(s string) []engine.RichAnswerBlock {
	if s == "" {
		return nil
	}
	var (
		out    []engine.RichAnswerBlock
		inline strings.Builder
		code   strings.Builder
		items  []string
		item   strings.Builder
		inPre  bool
		inItem bool
		quote  bool
		list   string
		skip   int
	)

	flushInline := func() {
		t := collapseText(inline.String())
		inline.Reset()
		if t == "" {
			return
		}
		kind := "text"
		if quote {
			kind = "quote"
		}
		out = append(out, engine.RichAnswerBlock{Kind: kind, Text: t})
	}
	flushCode := func() {
		inPre = false
		t := strings.Trim(code.String(), "\n \t")
		code.Reset()
		if t == "" {
			return
		}
		out = append(out, engine.RichAnswerBlock{Kind: "code", Text: t})
	}
	flushItem := func() {
		inItem = false
		if t := collapseText(item.String()); t != "" {
			items = append(items, t)
		}
		item.Reset()
	}
	flushList := func() {
		if list == "" {
			return
		}
		flushItem()
		list = ""
		if len(items) > 0 {
			out = append(out, engine.RichAnswerBlock{Kind: "list", Items: items})
		}
		items = nil
	}
	addText := func(t string) {
		t = html.UnescapeString(t)
		switch {
		case inPre:
			code.WriteString(t)
		case inItem:
			item.WriteString(t)
		default:
			inline.WriteString(t)
		}
	}

	for i := 0; i < len(s); {
		if s[i] != '<' {
			j := strings.IndexByte(s[i:], '<')
			if j < 0 {
				j = len(s) - i
			}
			addText(s[i : i+j])
			i += j
			continue
		}
		j := strings.IndexByte(s[i:], '>')
		if j < 0 {
			break
		}
		tag := s[i+1 : i+j]
		i += j + 1

		name := tagName(tag)
		closing := strings.HasPrefix(tag, "/")
		empty := strings.HasSuffix(tag, "/")

		if skip > 0 {
			switch name {
			case "script", "style", "svg", "iframe":
				if closing {
					skip--
				} else if !empty {
					skip++
				}
			}
			continue
		}

		switch name {
		case "script", "style", "svg", "iframe":
			if !closing && !empty {
				skip++
			}
		case "br":
			if !closing {
				if inPre {
					code.WriteString("\n")
				} else {
					flushInline()
				}
			}
		case "pre":
			if closing {
				flushCode()
			} else {
				flushInline()
				flushList()
				inPre = true
			}
		case "ul", "ol":
			if closing {
				flushList()
			} else {
				flushInline()
				flushList()
				list = name
			}
		case "li":
			if closing {
				flushItem()
			} else {
				flushInline()
				flushItem()
				inItem = true
			}
		case "blockquote":
			quote = !closing
			flushInline()
		case "hr":
			flushInline()
		case "p", "div", "h1", "h2", "h3", "h4", "tr":
			flushInline()
		}
		if len(out) >= maxTopAnswerBlocks {
			break
		}
	}
	flushInline()
	flushCode()
	flushList()

	var total int
	trimmed := out[:0]
	for _, b := range out {
		total += len(b.Text)
		for _, it := range b.Items {
			total += len(it)
		}
		if total > maxTopAnswerText {
			break
		}
		trimmed = append(trimmed, b)
	}
	if len(trimmed) == 0 {
		return nil
	}
	return trimmed
}

func tagName(tag string) string {
	name := strings.TrimLeft(tag, "/")
	for i := 0; i < len(name); i++ {
		switch name[i] {
		case ' ', '\t', '\n', '\r', '/':
			return strings.ToLower(name[:i])
		}
	}
	return strings.ToLower(name)
}

func collapseText(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
