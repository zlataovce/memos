package generator

import (
	"bytes"
	"fmt"
	"html/template"
)

type memoTemplateData struct {
	Title     string
	Timestamp string
	Excerpt   string
	Content   template.HTML
}

type indexTemplateData struct {
	Memos []*Memo
}

var (
	memoTemplate  = template.New("memo")
	indexTemplate = template.New("index")
)

func init() {
	template.Must(memoTemplate.Parse(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{ .Title }}</title>
	<meta name="og:description" content="{{ .Excerpt }}">
</head>
<body>
	<article>
		<h1>{{ .Title }}</h1>
		<p>{{ .Timestamp }}</p>
		<hr />
		<div>{{ .Content }}</div>
		<hr />
		<a href="../">back to index</a>
	</article>
</body>
</html>`))

	template.Must(indexTemplate.Parse(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>memos</title>
</head>
<body>
	<ul>
	{{ range .Memos }}
		<li><a href="./{{ .ID }}.html">{{ .Meta.Title }}</a></li>
	{{ end }}
	</ul>
</body>
</html>`))
}

func GenerateIndex(memos []*Memo) (string, error) {
	var buf bytes.Buffer
	if err := indexTemplate.Execute(&buf, indexTemplateData{Memos: memos}); err != nil {
		return "", &ErrRender{cause: fmt.Errorf("index template - %w", err)}
	}

	return buf.String(), nil
}
