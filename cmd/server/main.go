package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Headers     map[string][]string
	Method      string
	Body        any
	QueryString string
	Path        string
}

func echo(w http.ResponseWriter, req *http.Request) {
	var respBody any
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		panic("error reading body")
	}
	respBody = reqBody

	if req.Header.Get("Accept") == "application/json" && len(reqBody) > 0 {
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

	value, err := json.Marshal(response)
	if err != nil {
		panic("could not marshal to json")
	}

	fmt.Printf("%s\n", value)
	fmt.Fprintf(w, "%s\n", value)
}

func main() {
	http.HandleFunc("/", echo)
	http.ListenAndServe(":9001", nil)
}
