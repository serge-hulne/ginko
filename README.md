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

// example of Ajax call with HTMX syntax
func updateContent(w Response, req Request) {
	counter++
	newContent := ButtonHTMX(_update_content,
		"#content",
		"content",
		fmt.Sprint(counter),
	)
	Display(w, newContent)
}

// Home page (example of htmlgo syntax for the UI)
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

# Uses Htmlgo syntax for layout
https://github.com/julvo/htmlgo

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
	"net/http"
	"strings"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
	g "github.com/serge-hulne/ginko"
)

// Endpoints
const (
	_addTodo = "/add-todo"
)

// State
var todoList []string

// root renders the home page with the to-do form and list
func root(w g.Response, req g.Request) {
	page :=
		Html5_(
			g.HeadHTMX(),
			Body_(
				Div(
					Attr(a.Id("todo-list")),
					renderTodoList(),
				),
			),
		)
	g.Display(w, string(page))
}

// renderTodoList renders the current state of the to-do list
func renderTodoList() HTML {
	var sb strings.Builder
	for _, item := range todoList {
		sb.WriteString(fmt.Sprintf("<div>%s</div>", item))
	}
	return Div(
		Attr(a.Id("todo-list")),
		Form(
			Attr(a.Action(_addTodo), a.Method("post")),
			Input(Attr(a.Type("text"), a.Name("todoItem"), a.Placeholder("Add new item"), a.Id("todo-input"))),
			Br_(),
			g.ButtonHTMX(_addTodo, "#todo-list", "add", "Add a todo"),
		),
		HTML(sb.String()),
	)

}

// addTodo handles adding a new item to the to-do list
func addTodo(w g.Response, req g.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	todoItem := req.FormValue("todoItem")
	if todoItem != "" {
		todoList = append(todoList, todoItem)
	}

	// Print the todoList for debugging
	fmt.Println("Current Todo List:", todoList)

	// Display only the updated todo list, not the entire page
	g.Display(w, renderTodoList())
}

// Registering actions
var action = g.ActionMap{
	_addTodo: addTodo,
	"/":      root,
}

// Running the app
func main() {
	g.Run_app("To-Do List App", "8090", action)
}

```


