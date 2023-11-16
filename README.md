# Ginko
Go package for creating lightweight desktop apps simply and quickly.
- Easy development (familiar HTML-like syntax for the UI).
- Easy cross-compiling for all OS and or architechture.
- Uses WebView for front-end rendering, but *does not require* JS or CSS.
- Compiles into a single binary executable, statically compiled.
- Distributed as a single executable file without dependencies.

![Screenshot 2023-11-16 at 10 06 37](https://github.com/serge-hulne/ginko/assets/303502/5cd2aeaf-3f0e-415e-854b-dc0f72b1feb3)


# Example

```
package main

import (
	"fmt"

	. "github.com/julvo/htmlgo"
	. "github.com/serge-hulne/ginko"
)

// State
var (
	counter int = 0
)

// Actions
const (
	_update_content = "/update-content"
)

func updateContent(w Response, req Request) {
	counter++
	newContent := ButtonHTMX(_update_content,
		"#content",
		"content",
		fmt.Sprint(counter),
	)
	Display(w, newContent)
}

// Home page
func root(w Response, req Request) {
	page :=
		Html5_(
			HeadHTMX(),
			Body_(
				ButtonHTMX(_update_content,
					"#content",
					"content",
					"0"),
			),
		)
	Display(w, string(page))
}

// Registering actions
var action = ActionMap{
	_update_content: updateContent,
	"/":             root,
}

// Running app
func main() {
	Run_app("Basic App : Simple counter", "8090", action)
}

```

# Dependencies
- Uses WebView.
- Requires Go and a C/C++ toolchain (for example Xcode on Mac), as Go connects to WebView via cgo.
- All the dependencies are installed automatically, via `go get` (see example, hereunder).

# Use
1. create a new direcory : `mkdir MyApp`
2. `cd MyApp`
3. `go mod init App`
5. copy the example above in the current directory MyApp 
6. `go get -u github.com/serge-hulne/ginko`
7. `go build`
8. `./App`

# Cross compilation example
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o app-amd64-darwin app.go

# licence 
MIT

