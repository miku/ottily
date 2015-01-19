package main

import "github.com/idada/v8.go"

func main() {
	engine := v8.NewEngine()
	script := engine.Compile([]byte("'Hello ' + 'World!'"), nil)
	context := engine.NewContext(nil)

	context.Scope(func(cs v8.ContextScope) {
		result := cs.Run(script)
		println(result.ToString())
	})
}
