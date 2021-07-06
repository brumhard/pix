package viewer

import (
	"fmt"
	"html/template"
	"net/http"
	"ogframe/assets"
)

var _ http.Handler = (*Viewer)(nil)

type Viewer struct {
	template *template.Template
}

func NewViewer(templateName string) *Viewer {
	return &Viewer{
		template: template.Must(template.ParseFS(assets.HTMLTemplates, fmt.Sprintf("%s.gohtml", templateName))),
	}
}

func (h *Viewer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.template.Execute(w, nil)
}
