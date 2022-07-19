package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

const VERSION_KEY = "VERSION "

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
	return IPAddress
}

// ResponseLoggingHandler https://stackoverflow.com/questions/29319783/logging-responses-to-incoming-http-requests-inside-http-handlefunc
func ResponseLoggingHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// switch out response writer for a recorder
		// for all subsequent handlers
		c := httptest.NewRecorder()
		next(c, r)

		// copy everything from response recorder
		// to actual response writer
		for k, v := range c.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(c.Code)
		_, err := c.Body.WriteTo(w)
		ErrorHandler(err)
		log.Printf("%s %s %s %s %d", ReadUserIP(r), r.Method, r.URL, r.Proto, c.Code)
	}
}

func main() {
	setEnvErr := os.Setenv(VERSION_KEY, "1.0.0")
	ErrorHandler(setEnvErr)

	port := "8080"
	log.Println("Listening on :" + port)
	http.HandleFunc("/health", ResponseLoggingHandler(healthHandler))
	http.HandleFunc("/", ResponseLoggingHandler(mirrorHandler))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
