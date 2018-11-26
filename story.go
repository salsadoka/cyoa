package cyoa

import (
	"net/http"
	"html/template"
	"log"
	"strings"
)

func init(){
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template


var defaultHandlerTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Choose Your Own Adventure</title>
</head>
<body>
<h1> {{.Title}}</h1>
{{range .Paragraphs}}
    <p> {{.}}</p>
{{end}}
<ul>
    {{range .Options}}
        <li><a href="/{{.Arc}}">{{.Text}}</a></li>
    {{end}}
</ul>
</body>
</html>`

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption{
	return func(h *handler){
		h.t = t
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl}
	for _, opt := range opts {
		opt(&h)
	}

	return h
}

type handler struct{
	s Story
	t *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	path := strings.TrimSpace(r.URL.Path)

	if "" == path || "/" == path {
		path = "/intro"
	}

	path = path[1:] //getting rid of the slash

	chapter, ok := h.s[path]
	if !ok {
		http.Error(w, "Page doesn't exist", http.StatusNotFound)
	}

	err := h.t.Execute(w, chapter)
	if nil != err{
		log.Fatalf("Failed to load story with err: %s", err)
		http.Error(w, "Something went wrong...", http.StatusBadRequest)
	}
}


type Story map[string]StoryArc

type StoryArc struct{
	Title string `json:"title"`
	Paragraphs []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct{
	Text string `json:"text"`
	Arc string `json:"arc"`
}
