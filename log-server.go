package main

import (
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	port     int
	filename string
	f        *os.File
)

type LogMsg struct {
	Message string
}

type OkResp struct {
	Ok        bool      `json:"ok"`
	CreatedAt time.Time `json:"createdAt"`
}

func PostLogHandler(res rest.ResponseWriter, req *rest.Request) {

	content, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		rest.Error(res, "Error reading request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		rest.Error(res, "Error writing log message to file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//ioutil.WriteFile(filename, "\n", 0644)

	res.WriteHeader(http.StatusOK)
	res.WriteJson(OkResp{Ok: true, CreatedAt: time.Now()})
}

func main() {
	// configure command line flags
	flag.IntVar(&port, "port", 8080, "HTTP Server Port")
	flag.StringVar(&filename, "filename", "output", "Output file path")
	flag.Parse()

	// set up logging
	l, err := os.Create("log-server.log")
	defer l.Close()
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	logger := log.New(io.Writer(l), "", 0)

	// set up output file to save log messages to
	logger.Printf("Saving to output file %s", filename)
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}

	// set up HTTP server
	httpAddr := fmt.Sprintf(":%v", port)
	logger.Printf("Listening on %v", httpAddr)

	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		EnableGzip:               true,
		EnableLogAsJson:          true,
		Logger:                   logger,
		PreRoutingMiddlewares: []rest.Middleware{
			&rest.CorsMiddleware{
				RejectNonCorsRequests: false,
				OriginValidator: func(origin string, request *rest.Request) bool {
					return origin == "*"
				},
				AllowedHeaders: []string{"Accept", "Content-Type", "X-Requested-With"},
			},
		},
	}
	handler.SetRoutes(
		&rest.Route{"POST", "/log", PostLogHandler},
	)

	log.Fatal(http.ListenAndServe(httpAddr, &handler))
}
