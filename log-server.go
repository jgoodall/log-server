package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// output file to save logs to
var filepath string
var logger log.Logger

// log message format
// message is required, all other fields are optional
// createdAt is automatically assigned
type LogMessage struct {
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	Tags      []string  `json:"tags"`
	Type      string    `json:"type"`
}

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

// error handler - write to log and return error to client
func Error(res rest.ResponseWriter, code int, err string) {
	res.WriteHeader(code)
	res.WriteJson(ErrResp{Ok: false, ErroredAt: time.Now(), Error: err})
	logger.Println(err)
}

// handler for POST requests to write a log entry
// returns a 200 if successfully saved, else an error code
func PostLogHandler(res rest.ResponseWriter, req *rest.Request) {
	outfile, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer outfile.Close()
	if err != nil {
		Error(res, http.StatusInternalServerError, "Error opening output file: "+err.Error())
		return
	}

	logEntry := LogMessage{}
	err = req.DecodeJsonPayload(&logEntry)
	if err != nil {
		Error(res, http.StatusInternalServerError, "Error decoding request body: "+err.Error())
		return
	}

	if logEntry.Message == "" {
		Error(res, http.StatusInternalServerError, "Logs must include a message field")
		return
	}
	logEntry.CreatedAt = time.Now()

	js, err := json.Marshal(logEntry)
	if err != nil {
		Error(res, http.StatusInternalServerError, "Error encoding request as JSON: "+err.Error())
		return
	}

	_, err = outfile.Write(js)
	_, err = outfile.Write([]byte("\n"))
	if err != nil {
		Error(res, http.StatusInternalServerError, "Error writing log message to file: "+err.Error())
		return
	}

	res.WriteHeader(http.StatusOK)
	res.WriteJson(OkResp{Ok: true, CreatedAt: time.Now()})
}

// handler for GET requests to retrieve logs
// returns the newline delimited logs
func GetLogsHandler(res rest.ResponseWriter, req *rest.Request) {
	outfile, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDONLY, 0644)
	defer outfile.Close()
	if err != nil {
		Error(res, http.StatusInternalServerError, "Error opening output file: "+err.Error())
		return
	}

	res.WriteHeader(http.StatusOK)
	r := bufio.NewReader(outfile)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		res.(http.ResponseWriter).Write(line)
		res.(http.ResponseWriter).Write([]byte("\n"))
	}
	if scanner.Err() != nil {
		Error(res, http.StatusInternalServerError, "Error reading output file: "+scanner.Err().Error())
		return
	}
	res.(http.Flusher).Flush()
}

func main() {
	var port int
	var logpath string

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
					return origin == "http://localhost:9000"
				},
				AllowedMethods:                []string{"GET", "POST"},
				AllowedHeaders:                []string{"Accept", "Content-Type", "X-Requested-With"},
				AccessControlAllowCredentials: true,
				AccessControlMaxAge:           3600,
			},
		},
	}
	handler.SetRoutes(
		&rest.Route{"POST", "/log", PostLogHandler},
		&rest.Route{"GET", "/logs", GetLogsHandler},
	)

	log.Fatal(http.ListenAndServe(httpAddr, &handler))
}
