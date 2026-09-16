# Ginko
Please note, this is an Alpha version

Go package for creating lightweight desktop apps simply and quickly in pure Go.
- Easy development (familiar HTML-like syntax for the UI).
- Easy cross-compiling for all OS and or architecture.
- Uses WebView for front-end rendering, but *does not require* JS or CSS.
- Compiles into a single binary executable, statically compiled.
- Distributed as a single executable file without dependencies.

![Screenshot 2023-11-16 at 10 06 37](https://github.com/serge-hulne/ginko/assets/303502/5cd2aeaf-3f0e-415e-854b-dc0f72b1feb3)


# Example

```go
package main

import (
	"fmt"
	"log"
	"sync"

	. "github.com/serge-hulne/ginko"

	. "maragu.dev/gomponents/html"
)

// State
var (
	counter   int = 0
	counterMu sync.Mutex
)

// example of Ajax call with HTMX syntax
func updateContent(w Response, req Request) {
	counterMu.Lock()
	counter++
	current := counter
	counterMu.Unlock()

	newContent := ButtonHTMX("/update-content",
		"#content",
		"content",
		fmt.Sprint(current),
	)
	Display(w, newContent)
}

// Home page (example of gomponents syntax for the UI)
func root(w Response, req Request) {
	page :=
		Doctype(
			HTML(
				HeadHTMX(),
				Body(
					ButtonHTMX("/update-content",
						"#content",
						"content",
						"0"),
				),
			),
		)
	Display(w, page)
}

// Registering actions
var action = ActionMap{
	"/update-content": updateContent,
	"/":               root,
}

// Running app
func main() {
	if err := Run_app("Basic App : Simple counter", "8090", action); err != nil {
		log.Fatal(err)
	}
}

```

# Dependencies
- Uses WebView.
- Requires Go and a C/C++ toolchain (for example Xcode on Mac), as Go connects to WebView via cgo.
- All the dependencies are installed automatically, via `go get` (see example, hereunder).

# Uses gomponents syntax for layout
https://github.com/maragudk/gomponents

# Uses HTMX syntax for Ajax calls
https://htmx.org/docs/

# Use
1. create a new directory : `mkdir MyApp`
2. `cd MyApp`
3. `go mod init App`
5. copy the example above in the current directory MyApp 
6. `go get -u github.com/serge-hulne/ginko`
7. `go build`
8. `./App`

# Cross compilation example
to compile from a Mac M1 to a "classical" Mac intel:

`CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o app-amd64-darwin app.go`

# licence 
MIT

# Additional examples

Todo-list app:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	. "github.com/serge-hulne/ginko"
)

// Endpoints
const (
	_addTodo = "/add-todo"
)

// State
var (
	todoList   []string
	todoListMu sync.Mutex
)

// root renders the home page with the to-do form and list
func root(w Response, req Request) {
	page :=
		Doctype(
			HTML(
				HeadHTMX(),
				Body(
					renderTodoList(),
				),
			),
		)
	Display(w, page)
}

// renderTodoList renders the current state of the to-do list
func renderTodoList() g.Node {
	todoListMu.Lock()
	items := append([]string(nil), todoList...)
	todoListMu.Unlock()

	var itemNodes g.Group
	for _, item := range items {
		itemNodes = append(itemNodes, Div(g.Text(item)))
	}

	return Div(
		ID("todo-list"),
		Form(
			Action(_addTodo), Method("post"),
			Input(Type("text"), Name("todoItem"), Placeholder("Add new item"), ID("todo-input")),
			Br(),
			ButtonHTMX(_addTodo, "#todo-list", "add", "Add a todo"),
		),
		itemNodes,
	)
}

// addTodo handles adding a new item to the to-do list
func addTodo(w Response, req Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	todoItem := req.FormValue("todoItem")
	if todoItem != "" {
		todoListMu.Lock()
		todoList = append(todoList, todoItem)
		todoListMu.Unlock()
	}

	// Print the todoList for debugging
	fmt.Println("Current Todo List:", todoList)

	// Display only the updated todo list, not the entire page
	Display(w, renderTodoList())
}

// Registering actions
var action = ActionMap{
	_addTodo: addTodo,
	"/":      root,
}

// Running the app
func main() {
	if err := Run_app("To-Do List App", "8090", action); err != nil {
		log.Fatal(err)
	}
}

```


