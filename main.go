package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")
var upgrader = websocket.Upgrader{}

func synt(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		file, err := os.Open("./notes/C.wav")
		if err != nil {
			log.Println("open:", err)
			return
		}
		defer file.Close()

		stats, statsErr := file.Stat()
		if statsErr != nil {
			log.Println("stat:", statsErr)
			return
		}

		var size int64 = stats.Size()
		bytes := make([]byte, size)

		bufr := bufio.NewReader(file)
		_, err = bufr.Read(bytes)

		log.Printf("sending sound")
		c.WriteMessage(websocket.BinaryMessage, bytes)

		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	http.HandleFunc("/synt", synt)
	log.Println("starting server")
	log.Fatal(http.ListenAndServe(*addr, nil))
	log.Println("shutting down")
}
