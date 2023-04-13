package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	logger *log.Logger
)

type Response struct {
	Headers     map[string][]string
	Method      string
	Body        any
	QueryString string
	Path        string
}

func init() {
	logger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)

}

func echo(w http.ResponseWriter, req *http.Request) {
	var respBody any
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		panic("error reading body")
	}
	respBody, err = base64.StdEncoding.DecodeString(string(reqBody))
	if err != nil {
		logger.Println("Request body was not base64 encoded")
		respBody = reqBody
	}

	if contains(strings.Split(req.Header.Get("Accept"), ","), "application/json") && len(reqBody) > 0 {
		jsonBody := make(map[string]interface{})
		err = json.Unmarshal(reqBody, &jsonBody)
		if err != nil {
			panic("request body was not json")
		}
		respBody = jsonBody
	}

	response := Response{
		Headers:     req.Header,
		Method:      req.Method,
		Body:        respBody,
		QueryString: req.URL.RawQuery,
		Path:        req.URL.Path,
	}

	value, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		panic("could not marshal to json")
	}

	logger.Printf("%s\n", value)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s\n", value)
}

func main() {
	http.HandleFunc("/", echo)
	fmt.Println("Listening on port 9001")
	http.ListenAndServe(":9001", nil)
}

func contains[T comparable](values []T, value T) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
