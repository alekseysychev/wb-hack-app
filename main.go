// Example static file server.
//
// Serves static files from the given directory.
// Exports various stats at /stats .
package main

import (
	"expvar"
	"flag"
	"log"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
)

var addr = flag.String("addr", ":80", "TCP address to listen to")

func main() {
	// Parse command-line flags.
	flag.Parse()

	dir := "./web"

	// Setup FS handler
	fs := &fasthttp.FS{
		Root:       dir,
		IndexNames: []string{"index.html"},
	}
	fsHandler := fs.NewRequestHandler()

	// Create RequestHandler serving server stats on /stats and files
	// on other requested paths.
	// /stats output may be filtered using regexps. For example:
	//
	//   * /stats?r=fs will show only stats (expvars) containing 'fs'
	//     in their names.
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/api":
			// type ResponceType []string
			// responce := make(ResponceType, 0)
			responce := []string{
				"Поисковая подсказка #0",
				"Поисковая подсказка #1",
				"Поисковая подсказка #2",
				"Поисковая подсказка #3",
				"Поисковая подсказка #4",
				"Поисковая подсказка #5",
				"Поисковая подсказка #6",
				"Поисковая подсказка #7",
				"Поисковая подсказка #8",
				"Поисковая подсказка #9",
			}
			data, err := jsoniter.Marshal(responce)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
			if _, err := ctx.Write(data); err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
		case "/stats":
			expvarhandler.ExpvarHandler(ctx)
		default:
			fsHandler(ctx)
			updateFSCounters(ctx)
		}
	}

	// Start HTTP server.
	if len(*addr) > 0 {
		log.Printf("Starting HTTP server on %q", *addr)
		go func() {
			if err := fasthttp.ListenAndServe(*addr, requestHandler); err != nil {
				log.Fatalf("error in ListenAndServe: %s", err)
			}
		}()
	}

	log.Printf("Serving files from directory %q", dir)
	log.Printf("See stats at http://%s/stats", *addr)

	// Wait forever.
	select {}
}

func updateFSCounters(ctx *fasthttp.RequestCtx) {
	// Increment the number of fsHandler calls.
	fsCalls.Add(1)

	// Update other stats counters
	resp := &ctx.Response
	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		fsOKResponses.Add(1)
		fsResponseBodyBytes.Add(int64(resp.Header.ContentLength()))
	case fasthttp.StatusNotModified:
		fsNotModifiedResponses.Add(1)
	case fasthttp.StatusNotFound:
		fsNotFoundResponses.Add(1)
	default:
		fsOtherResponses.Add(1)
	}
}

// Various counters - see https://golang.org/pkg/expvar/ for details.
var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)
