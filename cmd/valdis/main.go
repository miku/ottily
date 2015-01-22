// This will leak memory
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/idada/v8.go"
)

func main() {

	flag.Parse()

	ff, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer ff.Close()
	reader := bufio.NewReader(ff)

	vm := v8.NewEngine()
	script := vm.Compile([]byte("input.length"), nil)
	context := vm.NewContext(nil)

	counter := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		context.Scope(func(cs v8.ContextScope) {
			s := vm.Compile([]byte(fmt.Sprintf("var input = %s", strconv.Quote(line))), nil)
			cs.Run(s)
			result := cs.Run(script)
			fmt.Println(result.ToString())
		})
		counter++
		if counter%100000 == 0 {
			log.Println("Running GC.")
			runtime.GC()
			log.Println("Done.")
		}
	}

}
