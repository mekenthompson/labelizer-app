package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func main() {
	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("_templates/*.html")),
	}
	e.Renderer = renderer

	// Named route "foobar"
	e.GET("/template/test", func(c echo.Context) error {
		return c.Render(http.StatusOK, "test.html", map[string]interface{}{})
	})

	e.Static("/assets/js", "app/dist")

	// Serve incoming routes from spa
	for _, route := range []string{"setup", "/"} {
		e.File(route, "public/index.html")
	}

	log.Fatal(e.Start(":8000"))
}
