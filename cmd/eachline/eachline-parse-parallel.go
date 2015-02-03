package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

func Worker(lines chan string, ticker chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var container map[string]interface{}
	for line := range lines {
		json.Unmarshal([]byte(line), &container)
		ticker <- true
	}
}

func Collector(ticker, done chan bool) {
	start := time.Now()
	counter := 0
	for _ = range ticker {
		counter++
		if counter%1000000 == 0 {
			elapsed := time.Since(start)
			speed := float64(counter) / elapsed.Seconds()
			eta, _ := time.ParseDuration(fmt.Sprintf("%0.3fs", float64(72727729-counter)/speed))
			log.Printf("%d %0.3f %s\n", counter, speed, eta)
		}
	}
	elapsed := time.Since(start)
	speed := float64(counter) / elapsed.Seconds()
	eta, _ := time.ParseDuration(fmt.Sprintf("%0.3fs", float64(72727729-counter)/speed))
	log.Printf("%d %0.3f %s\n", counter, speed, eta)
	done <- true
}

func main() {

	numWorkers := flag.Int("w", runtime.NumCPU(), "workers")

	flag.Parse()

	runtime.GOMAXPROCS(*numWorkers)

	ff, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer ff.Close()
	reader := bufio.NewReader(ff)

	lines := make(chan string)
	ticker := make(chan bool)
	done := make(chan bool)

	go Collector(ticker, done)

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go Worker(lines, ticker, &wg)
	}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		lines <- line
	}
	close(lines)
	wg.Wait()
	close(ticker)
	<-done
}
