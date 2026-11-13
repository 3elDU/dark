package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/3elDU/dark/dom/parser"
	"github.com/3elDU/dark/http"
	"github.com/grafana/sobek"
)

var ua = flag.String("u", "firefox", "User agent to use for all HTTP requests.\nOne of: firefox, chrome, custom")

func createConsole(vm *sobek.Runtime) {
	c := vm.NewObject()

	c.Set("log", func(call sobek.FunctionCall) sobek.Value {
		var b strings.Builder
		for i, arg := range call.Arguments {
			b.WriteString(arg.String())
			if i < len(call.Arguments)-1 {
				b.WriteString(", ")
			}
		}

		fmt.Println(b.String())

		return vm.ToValue("meow!")
	})

	vm.Set("console", c)
}

func main() {
	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		panic("expected a url as the last argument")
	}

	page, err := http.FetchWebpage(url, *ua)
	if err != nil {
		panic(err)
	}

	dom, err := parser.ParseDocument(page)
	if err != nil {
		panic(err)
	}

	dom.Print()

	// no event loop yet
	vm := sobek.New()
	createConsole(vm)
}
