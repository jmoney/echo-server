package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var (
	Info     *log.Logger
	Error    *log.Logger
	upgrader = websocket.Upgrader{}
)

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

func main() {
	port := flag.Int("port", 9002, "The port to connect too.")
	flag.Parse()

	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		defer conn.Close()

		// Continuosly read and write message
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read failed:", err)
				break
			}
			message = []byte(fmt.Sprintf("You sent: %s", string(message)))
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write failed:", err)
				break
			}
		}
	})

	Info.Printf("Listening on port %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		Error.Fatal(err)
	}
}
