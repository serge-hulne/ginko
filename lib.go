package ginko

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
	webview "github.com/webview/webview_go"
)

type Response = http.ResponseWriter
type Request = *http.Request
type ActionMap = map[string]func(Response, Request)

var display = fmt.Fprint

func run_server(action ActionMap) {
	for k, v := range action {
		http.HandleFunc(k, v)
	}
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func run_app(title, port string, Action ActionMap) {
	go func() {
		run_server(Action)
	}()
	w := webview.New(false)
	url := fmt.Sprintf("http://127.0.0.1:%s", port)
	defer w.Destroy()
	w.SetTitle(title)
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate(url)
	w.Run()
}

func ButtonHTMX(action, target, id, text string) HTML {
	b := fmt.Sprintf(`
		<button 
			class="w3-button w3-blue w3-round" 
			hx-put="%s" 
			hx-target=%s 
			hx-swap="outerHTML"
			id=%s>
			%s
		</button>
	`, action, target, id, text)
	return HTML(b)
}

func HeadHTMX() HTML {
	return Head_(
		Link(Attr(a.Rel("stylesheet"), a.Href("https://www.w3schools.com/w3css/4/w3.css"))),
		Style_(Text("body { margin: 20px; }")),
		Script(
			Attr(
				a.Src("https://unpkg.com/htmx.org@1.9.8"),
			),
			JavaScript(""),
		))
}
