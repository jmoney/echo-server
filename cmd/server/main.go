package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	Info     *log.Logger
	Error    *log.Logger
	upgrader = websocket.Upgrader{}
)

type Response struct {
	Headers     map[string][]string
	Method      string
	Body        any
	QueryString string
	Path        string
}

type WebSocketResponse struct {
	Received string
	Sent     string
}

func init() {
	// infoLogWriter, eilw := syslog.New(syslog.LOG_NOTICE, "echo-server")

	// if eilw == nil {
	// 	Info = log.New(infoLogWriter,
	// 		"[INFO] ",
	// 		log.Ldate|log.Ltime|log.Lshortfile)
	// } else {
	Info = log.New(os.Stdout,
		"[INFO] ",
		log.Ldate|log.Ltime|log.Lshortfile)
	// }

	// errorLogWriter, eelw := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, "echo-server")
	// if eelw == nil {
	// 	Error = log.New(errorLogWriter,
	// 		"[ERROR] ",
	// 		log.Ldate|log.Ltime|log.Lshortfile)
	// } else {
	Error = log.New(os.Stderr,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile)
	// }
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

func echoWebsocket(w http.ResponseWriter, req *http.Request) {
	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer conn.Close()

	// Continuosly read and write message
	for {
		mt, received, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed:", err)
			break
		}
		sent := []byte(fmt.Sprintf("You sent: %s", string(received)))
		err = conn.WriteMessage(mt, sent)
		if err != nil {
			log.Println("write failed:", err)
			break
		}
		response := WebSocketResponse{
			Received: strings.TrimSuffix(string(received), "\n"),
			Sent:     strings.TrimSuffix(string(sent), "\n"),
		}
		value, err := json.Marshal(response)
		if err != nil {
			Error.Printf("could not marshal to json: %v\n", response)
			return
		}

		Info.Printf("%s\n", value)
	}
}

func main() {
	port := flag.Int("port", 9001, "The port to connect too.")
	server_type := flag.String("type", "all", "The type of server to run. Options are http, websocket, or all")
	flag.Parse()

	if *server_type == "http" {
		http.HandleFunc("/http", echo)
	} else if *server_type == "websocket" {
		http.HandleFunc("/websocket", echoWebsocket)
	} else if *server_type == "all" {
		http.HandleFunc("/http", echo)
		http.HandleFunc("/websocket", echoWebsocket)
	} else {
		Error.Printf("server_type %s is not supported\n", *server_type)
		os.Exit(1)
	}

	Info.Printf("The server type \"%s\" is listening on port %d\n", *server_type, *port)
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
