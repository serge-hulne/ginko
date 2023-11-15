# Ginko
Go package for creating lightweight desktop apps simply and quickly.

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

# Use
1. go mod init App
2. go get -u github.com/serge-hulne/ginko
3. go build
