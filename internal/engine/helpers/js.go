package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseJSLiteral(s string) (any, error) {
	p := &jsParser{s: s}
	v, err := p.value()
	if err != nil {
		return nil, err
	}
	p.skipWS()
	if p.i < len(s) {
		return nil, fmt.Errorf("unexpected trailing data at offset %d", p.i)
	}
	return v, nil
}

func MatchObject(s string, i int) (int, bool) {
	if i >= len(s) || (s[i] != '{' && s[i] != '[') {
		return 0, false
	}
	open := s[i]
	close := byte('}')
	if open == '[' {
		close = ']'
	}
	depth := 0
	inStr := false
	var quote byte
	for ; i < len(s); i++ {
		c := s[i]
		if inStr {
			if c == '\\' {
				i++
			} else if c == quote {
				inStr = false
			}
			continue
		}
		switch c {
		case '"', '\'':
			inStr = true
			quote = c
		case open:
			depth++
		case close:
			depth--
			if depth == 0 {
				return i + 1, true
			}
		}
	}
	return 0, false
}

type jsParser struct {
	s string
	i int
}

func (p *jsParser) value() (any, error) {
	p.skipWS()
	if p.i >= len(p.s) {
		return nil, errors.New("unexpected end of input")
	}
	switch c := p.s[p.i]; {
	case c == '{':
		return p.object()
	case c == '[':
		return p.array()
	case c == '"' || c == '\'':
		return p.readString()
	case c == '-' || c == '+' || c == '.' || (c >= '0' && c <= '9'):
		return p.readNumber()
	default:
		return p.keyword()
	}
}

func (p *jsParser) object() (any, error) {
	p.i++
	p.skipWS()
	obj := make(map[string]any)
	if p.i < len(p.s) && p.s[p.i] == '}' {
		p.i++
		return obj, nil
	}
	for {
		p.skipWS()
		var key string
		if p.i < len(p.s) && (p.s[p.i] == '"' || p.s[p.i] == '\'') {
			k, err := p.readString()
			if err != nil {
				return nil, err
			}
			key = k
		} else {
			key = p.readIdentifier()
		}
		p.skipWS()
		if p.i >= len(p.s) || p.s[p.i] != ':' {
			return nil, fmt.Errorf("expected ':' after key %q at offset %d", key, p.i)
		}
		p.i++
		v, err := p.value()
		if err != nil {
			return nil, err
		}
		obj[key] = v
		p.skipWS()
		if p.i >= len(p.s) {
			return nil, errors.New("unterminated object")
		}
		switch p.s[p.i] {
		case ',':
			p.i++
		case '}':
			p.i++
			return obj, nil
		default:
			return nil, fmt.Errorf("unexpected %q in object at offset %d", p.s[p.i], p.i)
		}
	}
}

func (p *jsParser) array() (any, error) {
	p.i++
	p.skipWS()
	arr := make([]any, 0)
	if p.i < len(p.s) && p.s[p.i] == ']' {
		p.i++
		return arr, nil
	}
	for {
		v, err := p.value()
		if err != nil {
			return nil, err
		}
		arr = append(arr, v)
		p.skipWS()
		if p.i >= len(p.s) {
			return nil, errors.New("unterminated array")
		}
		if p.s[p.i] == ',' {
			p.i++
			p.skipWS()
			if p.i < len(p.s) && p.s[p.i] == ']' {
				p.i++
				return arr, nil
			}
			continue
		}
		if p.s[p.i] == ']' {
			p.i++
			return arr, nil
		}
		return nil, fmt.Errorf("unexpected %q in array at offset %d", p.s[p.i], p.i)
	}
}

func (p *jsParser) skipWS() {
	for p.i < len(p.s) {
		switch p.s[p.i] {
		case ' ', '\t', '\n', '\r':
			p.i++
		default:
			return
		}
	}
}

func (p *jsParser) readIdentifier() string {
	start := p.i
	for p.i < len(p.s) {
		c := p.s[p.i]
		if c == '_' || c == '$' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(p.i > start && (c >= '0' && c <= '9')) {
			p.i++
			continue
		}
		break
	}
	return p.s[start:p.i]
}

func (p *jsParser) keyword() (any, error) {
	for _, kw := range []struct {
		text string
		val  any
	}{
		{"void 0", nil},
		{"undefined", nil},
		{"null", nil},
		{"true", true},
		{"false", false},
		{"NaN", 0.0},
	} {
		if strings.HasPrefix(p.s[p.i:], kw.text) {
			p.i += len(kw.text)
			return kw.val, nil
		}
	}
	return nil, fmt.Errorf("unexpected token at offset %d", p.i)
}

func (p *jsParser) readNumber() (any, error) {
	start := p.i
	for p.i < len(p.s) {
		c := p.s[p.i]
		ok := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') ||
			c == '.' || c == 'x' || c == 'X' || c == 'e' || c == 'E' || c == '+' || c == '-'
		if !ok {
			break
		}
		p.i++
	}
	tok := p.s[start:p.i]
	if len(tok) > 2 && tok[0] == '0' && (tok[1] == 'x' || tok[1] == 'X') {
		if v, err := strconv.ParseInt(tok[2:], 16, 64); err == nil {
			return v, nil
		}
	}
	if strings.ContainsAny(tok, ".eE") {
		if v, err := strconv.ParseFloat(tok, 64); err == nil {
			return v, nil
		}
		return nil, fmt.Errorf("invalid number %q", tok)
	}
	if v, err := strconv.ParseInt(tok, 10, 64); err == nil {
		return v, nil
	}
	if v, err := strconv.ParseFloat(tok, 64); err == nil {
		return v, nil
	}
	return nil, fmt.Errorf("invalid number %q", tok)
}

func (p *jsParser) readString() (string, error) {
	quote := p.s[p.i]
	p.i++
	var b strings.Builder
	for p.i < len(p.s) {
		c := p.s[p.i]
		if c == quote {
			p.i++
			return b.String(), nil
		}
		if c == '\\' {
			p.i++
			if p.i >= len(p.s) {
				return "", errors.New("unterminated escape")
			}
			e := p.s[p.i]
			p.i++
			switch e {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case 'b':
				b.WriteByte('\b')
			case 'f':
				b.WriteByte('\f')
			case '\\', '"', '\'', '/':
				b.WriteByte(e)
			case 'u':
				r, err := p.readUnicodeEscape()
				if err != nil {
					return "", err
				}
				b.WriteRune(r)
			default:
				b.WriteByte(e)
			}
			continue
		}
		b.WriteByte(c)
		p.i++
	}
	return "", errors.New("unterminated string")
}

func (p *jsParser) readUnicodeEscape() (rune, error) {
	if p.i+4 > len(p.s) {
		return 0, errors.New("short unicode escape")
	}
	v, err := strconv.ParseUint(p.s[p.i:p.i+4], 16, 32)
	if err != nil {
		return 0, err
	}
	p.i += 4
	r := rune(v)
	// combine UTF-16 surrogate pairs
	if r >= 0xD800 && r <= 0xDBFF && strings.HasPrefix(p.s[p.i:], `\u`) {
		save := p.i
		p.i += 2
		if p.i+4 <= len(p.s) {
			lo, err := strconv.ParseUint(p.s[p.i:p.i+4], 16, 32)
			if err == nil {
				l := rune(lo)
				if l >= 0xDC00 && l <= 0xDFFF {
					p.i += 4
					return 0x10000 + (r-0xD800)<<10 + (l - 0xDC00), nil
				}
			}
		}
		p.i = save
	}
	return r, nil
}
