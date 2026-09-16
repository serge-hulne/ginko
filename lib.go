package ginko

import (
	_ "embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	webview "github.com/webview/webview_go"
)

type Response = http.ResponseWriter
type Request = *http.Request
type ActionMap = map[string]func(Response, Request)

// Display renders a gomponents Node to the given ResponseWriter.
func Display(w Response, n g.Node) {
	if err := n.Render(w); err != nil {
		log.Println("ginko: error rendering:", err)
	}
}

//go:embed assets/w3.css
var w3CSS []byte

//go:embed assets/htmx.min.js
var htmxJS []byte

const (
	w3CSSPath  = "/ginko-assets/w3.css"
	htmxJSPath = "/ginko-assets/htmx.min.js"
)

func registerAssets(mux *http.ServeMux) {
	mux.HandleFunc(w3CSSPath, func(w Response, req Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(w3CSS)
	})
	mux.HandleFunc(htmxJSPath, func(w Response, req Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(htmxJS)
	})
}

// csrfGuard builds a ServeMux that only serves the loopback origin the app
// itself navigates to, rejecting state-changing requests whose Origin or
// Referer header points elsewhere (defense against a remote/malicious page
// driving the app's local HTTP endpoints).
func csrfGuard(port string, action ActionMap) *http.ServeMux {
	mux := http.NewServeMux()
	registerAssets(mux)
	selfOrigin := "http://127.0.0.1:" + port

	for k, v := range action {
		handler := v
		mux.HandleFunc(k, func(w Response, req Request) {
			if req.Method != http.MethodGet && req.Method != http.MethodHead {
				if origin := req.Header.Get("Origin"); origin != "" && origin != selfOrigin {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				if referer := req.Header.Get("Referer"); referer != "" &&
					referer != selfOrigin && !strings.HasPrefix(referer, selfOrigin+"/") {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}
			handler(w, req)
		})
	}
	return mux
}

// Run_app starts the local HTTP server and opens a WebView window pointed
// at it. It returns an error immediately if the server cannot bind to the
// given port; once the server is up, the WebView runs on the calling
// goroutine until the window is closed.
func Run_app(title, port string, Action ActionMap) error {
	mux := csrfGuard(port, Action)
	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		return fmt.Errorf("ginko: error starting server: %w", err)
	}

	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Println("ginko: server error:", err)
		}
	}()

	w := webview.New(false)
	url := fmt.Sprintf("http://127.0.0.1:%s", port)
	defer w.Destroy()
	w.SetTitle(title)
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate(url)
	w.Run()

	return nil
}

func ButtonHTMX(action, target, id, text string) g.Node {
	return Button(
		Class("w3-button w3-blue w3-round"),
		g.Attr("hx-put", action),
		g.Attr("hx-target", target),
		g.Attr("hx-swap", "outerHTML"),
		ID(id),
		g.Text(text),
	)
}

func HeadHTMX() g.Node {
	return Head(
		Link(Rel("stylesheet"), Href(w3CSSPath)),
		StyleEl(g.Text("body { margin: 20px; }")),
		Script(Src(htmxJSPath)),
	)
}
