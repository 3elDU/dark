package main

import (
	"fmt"

	"github.com/grafana/sobek"
)

func createConsole(vm *sobek.Runtime) {
	c := vm.NewObject()

	c.Set("log", func(call sobek.FunctionCall) sobek.Value {
		fmt.Println(call.Arguments)

		return vm.ToValue("meow!")
	})

	vm.Set("console", c)
}

func main() {
	vm := sobek.New()
	createConsole(vm)

	v, err := vm.RunString("console.log('Hello from JS!')")
	if err != nil {
		panic(err)
	}

	if str, ok := v.Export().(string); ok {
		fmt.Println("Got string from JS:", str)
	}
}
