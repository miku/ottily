// memory leak
// time go run cmd/gaffer/main.go datasets/1k.ldj  > /dev/null

// real	0m1.036s
// user	0m0.979s
// sys	0m0.115s
//
// time go run cmd/gaffer/main.go datasets/100k.ldj  > /dev/null

// real	1m33.336s
// user	1m32.135s
// sys	0m5.268s
//
// 60% of 16G RAM used
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olebedev/go-duktape"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("LDJ file required")
	}

	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	ctx := duktape.NewContext()

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)

		ctx.PushString(fmt.Sprintf(`var input = %s;`, strconv.Quote(line)))
		// ctx.PushString("obj = JSON.parse(input); JSON.stringify(obj);")
		// ctx.PushString(`println("Hello")`)
		ctx.Eval()
		ctx.PushString(`obj = JSON.parse(input); obj["001"] = "sample_" + obj["001"]; JSON.stringify(obj);`)
		ctx.Eval()
		result := ctx.GetString(-1)
		ctx.Pop()
		ctx.Pop()
		fmt.Println(result)

		// ctx.EvalString(fmt.Sprintf(`var input = %s`, strconv.Quote(line)))
		// ctx.EvalString(`obj = JSON.parse(input)`)
		// ctx.EvalString(`obj["001"] = "SAMPLE_" + obj["001"]`)
		// ctx.EvalString(`JSON.stringify(obj)`)
		// result := ctx.GetString(-1)
		// ctx.Pop()
		// fmt.Println(result)
	}

}
