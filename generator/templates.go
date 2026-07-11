package generator

import "html/template"

type memoTemplateData struct {
	Title           string
	Timestamp       string
	Excerpt         string
	ContentRendered template.HTML
}

var (
	memoTemplate = template.New("memo")
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
		<div>{{ .ContentRendered }}</div>
		<hr />
		<a href="../">back to index</a>
	</article>
</body>
</html>`))
}
