// Ottily executes a javascript snippet on each line of an input file in parallel.
//
// Noop:
//
//     $ ottily datasets/simple.ldj
//     {"name": "ottily", "language": "Golang"}
//
// Inline script with -e:
//
//     $ ottily -e 'output=input.length' datasets/simple.ldj
//	   40
//
//     $ ottily -e 'o=JSON.parse(input); o["language"] = "Go"; output=JSON.stringify(o);' datasets/simple.ldj
//     {"language":"Go","name":"ottily"}
//
// Pass a script file:
//
//     $ ottily -s scripts/classified.js datasets/simple.ldj
//     CLASSIFIED
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"

	"github.com/robertkrimen/otto"
)

const NOOP_SCRIPT = "output = input"
const VERSION = "0.1.0"

func Worker(lines, out chan *string, script, preload string, wg *sync.WaitGroup) {
	defer wg.Done()
	vm := otto.New()

	if preload != "" {
		_, err := vm.Run(preload)
		if err != nil {
			log.Fatal(err)
		}
	}

	compiled, err := vm.Compile("", script)
	if err != nil {
		log.Fatal(err)
	}

	for line := range lines {
		vm.Set("input", *line)
		_, err := vm.Run(compiled)
		if err != nil {
			log.Fatal(err)
		}
		result, err := vm.Get("output")
		if result == otto.NullValue() {
			continue
		}
		if err != nil {
			log.Fatal(err)
		}
		r := result.String()
		out <- &r
	}
}

// FanInWriter writes the channel content to the writer
func FanInWriter(writer io.Writer, in chan *string, done chan bool) {
	for s := range in {
		writer.Write([]byte(*s))
		writer.Write([]byte("\n"))
	}
	done <- true
}

func main() {
	preload := flag.String("p", "", "include this js file before processing")
	script := flag.String("s", "", "script to execute on each line of input")
	execute := flag.String("e", "", "execute argument on each line of input")
	numWorkers := flag.Int("w", runtime.NumCPU(), "number of workers")
	version := flag.Bool("v", false, "prints current program version")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		log.Fatal("input file required")
	}

	preloadContent := ""

	if *preload != "" {
		ff, err := os.Open(*preload)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(ff)
		if err != nil {
			log.Fatal(err)
		}
		preloadContent = string(b)
	}

	content := NOOP_SCRIPT

	if *script != "" {
		ff, err := os.Open(*script)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(ff)
		if err != nil {
			log.Fatal(err)
		}
		content = string(b)
	}

	if *execute != "" {
		content = *execute
	}

	ff, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer ff.Close()
	reader := bufio.NewReader(ff)

	if *numWorkers > 0 {
		runtime.GOMAXPROCS(*numWorkers)
	}

	queue := make(chan *string)
	out := make(chan *string)
	done := make(chan bool)
	var wg sync.WaitGroup

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	go FanInWriter(writer, out, done)

	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go Worker(queue, out, content, preloadContent, &wg)
	}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		queue <- &line
	}
	close(queue)
	wg.Wait()
	close(out)
	<-done
}
