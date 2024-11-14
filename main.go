//go:generate go run themes/include_themes.go
//go:generate go run scripts/include_scripts.go

package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	// Get user's option
	input_file_name := os.Args[1]
	output_file_name := os.Args[2]

	// Read the markdown file
	markdown, err := os.ReadFile(input_file_name)
	if err != nil {
		panic(err)
	}

	html, err := convertMarkdownToHtml(markdown)
	if err != nil {
		panic(err)
	}

	pdf, err := convertHTMLToPDF(html)

	// Write buffer contents to file on disk
	if err := os.WriteFile(
		output_file_name,
		pdf,
		fs.FileMode(0644)); err != nil {
		panic(err)
	}
}

func convertHTMLToPDF(html string) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	page := wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(html)))
	page.JavascriptDelay.Set(3000) // Wait for scripts to be loaded

	pdfg.AddPage(page)
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	if err := pdfg.Create(); err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}

func convertMarkdownToHtml(markdown []byte) (string, error) {
	markdown_converter := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	// Create buffer for HTML
	var html bytes.Buffer

	// Include the theme and set charset
	html.WriteString(fmt.Sprintf("<style type=text/css>%s</style>", css))
	html.WriteString("<meta charset=\"UTF-8\">\n")

	// Div for the content
	html.WriteString("<div class=\"markdown-body\">\n")

	// Convert Markdown to HTML and save it to the buffer
	err := markdown_converter.Convert(markdown, &html)
	if err != nil {
		return "", err
	}

	// Close the content div
	html.WriteString("</div>")

	htmlWithMathJax := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
        %s
	</head>
	<body>
		%s
	</body>
	</html>
	`, MATHJAX_SCRIPT, html.String())

	return htmlWithMathJax, nil
}
