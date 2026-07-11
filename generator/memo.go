package generator

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

var ErrMissingMeta = errors.New("missing frontmatter")

type ErrParse struct {
	cause error
}

func (ep *ErrParse) Error() string {
	return fmt.Sprintf("failed to parse: %s", ep.cause.Error())
}

func (ep *ErrParse) Unwrap() error {
	return ep.cause
}

type ErrRender struct {
	cause error
}

func (er *ErrRender) Error() string {
	return fmt.Sprintf("failed to render: %s", er.cause.Error())
}

func (er *ErrRender) Unwrap() error {
	return er.cause
}

type Memo struct {
	Meta    *Frontmatter
	Content string
}

func ParseMemo(v string) (*Memo, error) {
	parts := strings.SplitN(v, "---", 3)
	if len(parts) < 2 {
		return nil, &ErrParse{cause: ErrMissingMeta}
	}

	fm, err := ParseFrontmatter(parts[1])
	if err != nil {
		return nil, &ErrParse{cause: fmt.Errorf("frontmatter - %w", err)}
	}

	return &Memo{Meta: fm, Content: parts[2]}, nil
}

func formatExcerpt(content string) string {
	exc := strings.Trim(strings.ReplaceAll(content, "\n", " "), " ")
	if len(exc) > 100 {
		return exc[:100] + "..."
	}

	return exc
}

func (m *Memo) Generate() (string, error) {
	var contentBuf bytes.Buffer
	err := goldmark.Convert([]byte(m.Content), &contentBuf)
	if err != nil {
		return "", &ErrRender{cause: fmt.Errorf("markdown - %w", err)}
	}

	exc := formatExcerpt(m.Content)
	if len(m.Meta.Excerpt) > 0 {
		exc = m.Meta.Excerpt
	}

	var buf bytes.Buffer
	err = memoTemplate.Execute(&buf, memoTemplateData{
		Title:           m.Meta.Title,
		Timestamp:       m.Meta.Timestamp.Format(time.DateOnly),
		Excerpt:         exc,
		ContentRendered: template.HTML(contentBuf.String()),
	})
	if err != nil {
		return "", &ErrRender{cause: fmt.Errorf("template - %w", err)}
	}

	return buf.String(), nil
}
