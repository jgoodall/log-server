package main

import (
	"bufio"
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

// output file to save logs to
var f *os.File

// response type for logs successfully saved
type OkResp struct {
	Ok        bool      `json:"ok"`
	CreatedAt time.Time `json:"createdAt"`
}

// response type for logs that errored
type ErrResp struct {
	Ok        bool      `json:"ok"`
	ErroredAt time.Time `json:"erroredAt"`
	Error     string    `json:"error"`
}

// handler for POST requests to write logs
func PostLogHandler(res rest.ResponseWriter, req *rest.Request) {

	content, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		e := "Error reading request body: " + err.Error()
		res.WriteHeader(http.StatusInternalServerError)
		res.WriteJson(ErrResp{Ok: false, ErroredAt: time.Now(), Error: e})
		return
	}

	_, err = f.Write(content)
	if err != nil {
		e := "Error writing log message to file: " + err.Error()
		res.WriteHeader(http.StatusInternalServerError)
		res.WriteJson(ErrResp{Ok: false, ErroredAt: time.Now(), Error: e})
		return
	}

	res.WriteHeader(http.StatusOK)
	res.WriteJson(OkResp{Ok: true, CreatedAt: time.Now()})
}

// handler for POST requests to write logs
func GetLogsHandler(res rest.ResponseWriter, req *rest.Request) {
	res.WriteHeader(http.StatusOK)
	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		res.WriteJson(scanner.Text())
		res.(http.ResponseWriter).Write([]byte("\n"))
	}
	res.(http.Flusher).Flush()
}

func main() {
	var port int
	var filepath, logpath string

	// configure command line flags
	flag.IntVar(&port, "port", 8080, "HTTP Server Port")
	flag.StringVar(&filepath, "filepath", "output.json", "Output JSON file path")
	flag.StringVar(&logpath, "logpath", "log-server.log", "Log file path")
	flag.Parse()

	// set up logging
	l, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE, 0644)
	defer l.Close()
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	logger := log.New(io.Writer(l), "", 0)

	// set up output file to save log messages to
	logger.Printf("Saving to output file %s", filepath)
	f, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		log.Println(err)
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
		&rest.Route{"GET", "/logs", GetLogsHandler},
	)

	log.Fatal(http.ListenAndServe(httpAddr, &handler))
}
