package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

const VERSION_KEY = "VERSION"

func mirrorHandler(writer http.ResponseWriter, request *http.Request) {
	respHeader := writer.Header()
	for k, v := range request.Header {
		for _, vv := range v {
			respHeader.Add(k, vv)
		}
	}
	respHeader.Set("version", os.Getenv(VERSION_KEY))
}

func healthHandler(writer http.ResponseWriter, _ *http.Request) {
	_, err := writer.Write([]byte("200"))
	ErrorHandler(err)
}

func ErrorHandler(write2 error) {
	if write2 != nil {
		log.Printf("Error: %s", write2)
	}
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	if ip, _, err := net.SplitHostPort(IPAddress); err == nil {
		IPAddress = ip
	}
	return IPAddress
}

func Logger(out io.Writer, h http.Handler) http.Handler {
	logger := log.New(out, "", 0)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := &responseObserver{ResponseWriter: w, status: 200}
		h.ServeHTTP(o, r)
		addr := ReadUserIP(r)
		logger.Printf("%s %s %s %s %d %d", addr, r.Method, r.URL.Path, r.Proto, o.status, o.written)
	})
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func main() {
	setEnvErr := os.Setenv(VERSION_KEY, "1.0.0")
	ErrorHandler(setEnvErr)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	log.Println("Listening on :" + httpPort)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", mirrorHandler)
	log.Fatal(http.ListenAndServe(":"+httpPort, Logger(os.Stderr, http.DefaultServeMux)))
}
