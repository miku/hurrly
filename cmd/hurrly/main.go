package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
)

const Version = "0.1.1"

type Result struct {
	Status    string
	URL       string
	Took      float64
	Locations []string
	Epoch     int64
}

func (r Result) String() string {
	return fmt.Sprintf("%s\t%0.4f\t%d\t%s\t%s\t", r.Status, r.Took, r.Epoch, r.URL, strings.Join(r.Locations, "|"))
}

type URLValue struct {
	Format string `json:"format"`
	Value  string `json:"value"`
}

type Value struct {
	Index     int             `json:"index"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	TTL       int             `json:"ttl"`
	Timestamp string          `json:"timestamp"`
}

type APIResponse struct {
	Code   int     `json:"responseCode"`
	Handle string  `json:"handle"`
	Values []Value `json:"values"`
}

func worker(queue chan *url.URL, out chan Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for u := range queue {
		r := retrieve(u)
		out <- r
	}
}

func sink(out chan Result, done chan bool) {
	for r := range out {
		fmt.Println(r)
	}
	done <- true
}

// retrieve will try to GET and parse a DOI API response and will always
// return a Result, which will contain status (either HTTP or internal error designations)
func retrieve(target *url.URL) Result {

	rt := http.DefaultTransport
	var req *http.Request

	err := backoff.Retry(func() (e error) {
		req, e = http.NewRequest("GET", target.String(), nil)
		return
	}, backoff.NewExponentialBackOff())

	if err != nil {
		return Result{Status: "E_REQ", URL: target.String(), Took: 0, Epoch: time.Now().Unix()}
	}

	var resp *http.Response

	start := time.Now()
	err = backoff.Retry(func() (e error) {
		resp, e = rt.RoundTrip(req)
		if e != nil {
			log.Printf("retrying %s", req.URL.String())
		}
		return e
	}, backoff.NewExponentialBackOff())
	elapsed := time.Since(start)

	if err != nil {
		return Result{Status: "E_REQ", URL: target.String(), Took: elapsed.Seconds(), Epoch: time.Now().Unix()}
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Result{Status: "E_READ", URL: target.String(), Took: elapsed.Seconds(), Epoch: time.Now().Unix()}
	}

	resp.Body.Close()

	var ar APIResponse
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return Result{Status: "E_JSON", URL: target.String(), Took: elapsed.Seconds(), Epoch: time.Now().Unix()}
	}

	result := Result{Status: resp.Status, URL: target.String(), Took: elapsed.Seconds(), Epoch: time.Now().Unix()}

	for _, value := range ar.Values {
		if value.Type == "URL" {
			var v URLValue
			err := json.Unmarshal(value.Data, &v)
			if err != nil {
				return Result{Status: "E_JSON", URL: target.String(), Took: elapsed.Seconds(), Epoch: time.Now().Unix()}
			}
			result.Locations = append(result.Locations, v.Value)
		}
	}
	return result
}

func main() {

	numWorkers := flag.Int("w", runtime.NumCPU(), "number of workers")
	version := flag.Bool("v", false, "prints current program version")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	runtime.GOMAXPROCS(*numWorkers)

	reader := bufio.NewReader(os.Stdin)

	queue := make(chan *url.URL)
	out := make(chan Result)
	done := make(chan bool)

	go sink(out, done)

	var wg sync.WaitGroup

	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go worker(queue, out, &wg)
	}

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		target := strings.TrimSpace(line)

		if target == "" {
			continue
		}

		if !strings.HasPrefix(target, "http://doi.org/api/handles/") {
			log.Fatal("only http://doi.org/api/handles/... URLs supported for now")
		}

		parsed, err := url.Parse(target)
		if err != nil {
			continue
		}

		queue <- parsed
	}

	close(queue)
	wg.Wait()
	close(out)
	<-done
}
