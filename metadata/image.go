package metadata

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed nft.tmpl
var svgTemplate string
var tmpl *template.Template

type svgVariables struct {
	Domain   string
	FontSize int
}

func init() {
	tmpl = template.Must(template.New("svg").Parse(svgTemplate))
}

// SVGImage generates an svg image for the given name.
func SVGImage(name string) ([]byte, error) {
	var b bytes.Buffer
	err := tmpl.Execute(&b, svgVariables{
		Domain:   name,
		FontSize: fontSize(name),
	})
	if err != nil {
		return nil, fmt.Errorf("error rendering image: %v", err)
	}
	return b.Bytes(), nil
}

func fontSize(name string) int {
	if len(name) <= 13 {
		return 21
	}
	if len(name) > 26 {
		return 11
	}
	return 21 + 13 - len(name)
}
