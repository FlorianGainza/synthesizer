package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/websocket"
)

type header struct {
	RiffMark      [4]byte
	FileSize      int32
	WaveMark      [4]byte
	FmtMark       [4]byte
	FormatSize    int32
	FormatType    int16
	NumChans      int16
	SampleRate    int32
	ByteRate      int32
	BytesPerFrame int16
	BitsPerSample int16
	DataMark      [4]byte
	DataSize      int32
}

var addr = flag.String("addr", ":8080", "http service address")
var upgrader = websocket.Upgrader{}

func sinWave(sampleRate int, freq int) []byte {
	soundData := make([]byte, sampleRate)
	cycleRate := float64(sampleRate) / float64(freq)

	var x float64
	for i := 0; i < sampleRate; i++ {
		x = float64(i) / cycleRate * (2 * math.Pi)
		soundData[i] = byte(((math.Sin(x) + 1) / 2) * 255)
	}
	return soundData
}

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

		numChannels := 1
		sampleRate := 44100
		// TODO what is precision ???
		precision := 1
		sound := sinWave(sampleRate, 262)

		h := header{
			RiffMark:      [4]byte{'R', 'I', 'F', 'F'},
			FileSize:      44100 + 44, // finalization
			WaveMark:      [4]byte{'W', 'A', 'V', 'E'},
			FmtMark:       [4]byte{'f', 'm', 't', ' '},
			FormatSize:    16, //16 for PCM. This is the size of the rest of the Subchunk which follows this number.
			FormatType:    1,  //PCM = 1 (i.e. Linear quantization) Values other than 1 indicate some form of compression.
			NumChans:      int16(numChannels),
			SampleRate:    int32(sampleRate),
			ByteRate:      int32(int(sampleRate) * numChannels * precision), // == SampleRate * NumChannels * BitsPerSample/8
			BytesPerFrame: int16(numChannels * precision),                   // == NumChannels * BitsPerSample/8 The number of bytes for one sample including all channels.
			BitsPerSample: int16(precision) * 8,
			DataMark:      [4]byte{'d', 'a', 't', 'a'},
			DataSize:      44100, // finalization
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
