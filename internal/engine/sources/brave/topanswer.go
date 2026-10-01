package brave

import (
	"html"
	"strings"

	"metasearch/internal/engine"
	"metasearch/internal/engine/helpers"
)

const (
	maxTopAnswerBlocks = 40
	maxTopAnswerText   = 8000
)

func parseTopAnswer(data map[string]any) *engine.RichAnswer {
	body, _ := data["body"].(map[string]any)
	response, _ := body["response"].(map[string]any)
	web, _ := response["web"].(map[string]any)
	raw, _ := web["results"].([]any)

	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok || richStr(m["subtype"]) != "qa" {
			continue
		}
		qa, ok := m["qa"].(map[string]any)
		if !ok {
			continue
		}
		ans, ok := qa["answer"].(map[string]any)
		if !ok {
			continue
		}
		question := helpers.CleanText(richStr(qa["question"]))
		blocks := answerBlocks(richStr(ans["text"]))
		if question == "" || len(blocks) == 0 {
			continue
		}
		ta := &engine.RichTopAnswer{
			Question: question,
			Blocks:   blocks,
			Author:   helpers.CleanText(richStr(ans["author"])),
			Title:    helpers.CleanText(richStr(m["title"])),
			URL:      richStr(m["url"]),
		}
		if n, ok := intVal(ans["upvoteCount"]); ok && n > 0 {
			ta.Upvotes = int(n)
		}
		if mu, ok := m["meta_url"].(map[string]any); ok {
			ta.MetaURL = &engine.RichMetaURL{Hostname: richStr(mu["hostname"]), Favicon: richStr(mu["favicon"])}
		}
		return &engine.RichAnswer{Type: "topAnswer", TopAnswer: ta}
	}
	return nil
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
