package wave_test

import (
	"os"
	"testing"

	"github.com/FlorianGainza/synthesizer/pkg/wave"
)

func sliceAreEqual(a []byte, b []byte) bool {
	areEqual := true
	if len(a) == len(b) {
		for i, v := range a {
			if v != b[i] {
				areEqual = false
				break
			}
		}
	} else {
		areEqual = false
	}
	return areEqual
}

func loadBinaryContent(fileName string) []byte {
	f, err := os.Open("../../fixtures/test/" + fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	stats, statsErr := f.Stat()
	if statsErr != nil {
		panic(statsErr)
	}

	var size int64 = stats.Size()
	content := make([]byte, size)

	if _, err := f.Read(content); err != nil {
		panic("File " + fileName + " not readable")
	}

	return content
}

func TestSinGenerateC(t *testing.T) {
	expected := loadBinaryContent("sin-c.dat")
	if !sliceAreEqual(wave.Sin(44100, 262, 44100), expected) {
		t.Errorf("Wrong sin wave generated")
	}
}

func TestSquareGenerateC(t *testing.T) {
	expected := loadBinaryContent("square-c.dat")
	if !sliceAreEqual(wave.Square(44100, 262, 44100), expected) {
		t.Errorf("Wrong square wave generated")
	}
}
