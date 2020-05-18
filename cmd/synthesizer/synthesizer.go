package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"net/http"

	"github.com/FlorianGainza/synthesizer/pkg/wave"
	"github.com/gorilla/websocket"
)

var frequencies = map[string]uint16{
	"c":  262,
	"c#": 278,
	"d":  294,
	"d#": 311,
	"e":  330,
	"f":  349,
	"f#": 370,
	"g":  392,
	"g#": 415,
	"a":  440,
	"a#": 466,
	"b":  494,
}

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
		_, pitch, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", pitch)

		numChannels := 1
		sampleRate := uint16(44100)
		// TODO what is precision ???
		precision := 1
		samples := uint32(44100) // 0.002 secondes - 0.2 sample

		freq, _ := frequencies[string(pitch)]
		sound := wave.Square(sampleRate, freq, samples)

		h := wave.Header{
			RiffMark:      [4]byte{'R', 'I', 'F', 'F'},
			FileSize:      int32(samples) + 44,
			WaveMark:      [4]byte{'W', 'A', 'V', 'E'},
			FmtMark:       [4]byte{'f', 'm', 't', ' '},
			FormatSize:    16,
			FormatType:    1,
			NumChans:      int16(numChannels),
			SampleRate:    int32(sampleRate),
			ByteRate:      int32(int(sampleRate) * numChannels * precision),
			BytesPerFrame: int16(numChannels * precision),
			BitsPerSample: int16(precision) * 8,
			DataMark:      [4]byte{'d', 'a', 't', 'a'},
			DataSize:      samples,
		}
		var buf bytes.Buffer
		if err := binary.Write(&buf, binary.LittleEndian, &h); err != nil {
			log.Println("error:", err)
			return
		}
		if err := binary.Write(&buf, binary.LittleEndian, &sound); err != nil {
			log.Println("error:", err)
			return
		}
		log.Printf("sending sound")
		c.WriteMessage(websocket.BinaryMessage, buf.Bytes())

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
