// This will leak memory
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {

	flag.Parse()

	ff, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer ff.Close()
	reader := bufio.NewReader(ff)

	start := time.Now()
	counter := 0

	// var container map[string]interface{}

	for {
		_, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// json.Unmarshal([]byte(line), &container)
		counter++
		if counter%1000000 == 0 {
			elapsed := time.Since(start)
			speed := float64(counter) / elapsed.Seconds()
			eta, _ := time.ParseDuration(fmt.Sprintf("%0.3fs", float64(72727729-counter)/speed))
			log.Printf("%d %0.3f %s\n", counter, speed, eta)
		}
	}

}
