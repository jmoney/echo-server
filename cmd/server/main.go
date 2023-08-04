package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strings"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

type Response struct {
	Headers     map[string][]string
	Method      string
	Body        any
	QueryString string
	Path        string
}

func init() {
	infoLogWriter, eilw := syslog.New(syslog.LOG_NOTICE, "echo-server")
	if eilw == nil {
		Info = log.New(infoLogWriter,
			"[INFO] ",
			log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Info = log.New(os.Stdout,
			"[INFO] ",
			log.Ldate|log.Ltime|log.Lshortfile)
	}

	errorLogWriter, eelw := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, "echo-server")
	if eelw == nil {
		Error = log.New(errorLogWriter,
			"[ERROR] ",
			log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Error = log.New(os.Stderr,
			"[ERROR] ",
			log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func echo(w http.ResponseWriter, req *http.Request) {
	var respBody any
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		panic("error reading body")
	}
	respBody, err = base64.StdEncoding.DecodeString(string(reqBody))
	if err != nil {
		Info.Println("Request body was not base64 encoded")
		respBody = reqBody
	}

	if contains(strings.Split(req.Header.Get("Accept"), ","), "application/json") && len(reqBody) > 0 {
		jsonBody := make(map[string]interface{})
		err = json.Unmarshal(reqBody, &jsonBody)
		if err != nil {
			Error.Println("Request body was not json")
		} else {
			respBody = jsonBody
		}
	}

	response := Response{
		Headers:     req.Header,
		Method:      req.Method,
		Body:        respBody,
		QueryString: req.URL.RawQuery,
		Path:        req.URL.Path,
	}

	value, err := json.Marshal(response)
	if err != nil {
		Error.Printf("could not marshal to json: %v\n", response)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not marshal to json")
		return
	}

	Info.Printf("%s\n", value)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s\n", value)
}

func main() {
	port := flag.Int("port", 9002, "The port to connect too.")
	flag.Parse()

	http.HandleFunc("/", echo)
	log.Printf("Listening on port %d\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

// function that checks if a value is in a slice using comparable generics
func contains[T comparable](values []T, value T) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
