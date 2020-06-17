package oscillator_test

import (
	"os"
	"testing"

	"github.com/FlorianGainza/synthesizer/pkg/oscillator"
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
	if !sliceAreEqual(oscillator.Sin(44100, 262, 169, 0), expected) {
		t.Errorf("Wrong sin wave generated")
	}
}

func TestSinGenerateCwithOffset(t *testing.T) {
	expected := loadBinaryContent("sin-c-offset.dat")
	if !sliceAreEqual(oscillator.Sin(44100, 262, 169, 85), expected) {
		t.Errorf("Wrong square wave with offset generated")
	}
}

func TestSquareGenerateC(t *testing.T) {
	expected := loadBinaryContent("square-c.dat")
	if !sliceAreEqual(oscillator.Square(44100, 262, 169, 0), expected) {
		t.Errorf("Wrong square wave generated")
	}
}

func TestSquareGenerateCwithOffset(t *testing.T) {
	expected := loadBinaryContent("square-c-offset.dat")
	if !sliceAreEqual(oscillator.Square(44100, 262, 169, 85), expected) {
		t.Errorf("Wrong square wave with offset generated")
	}
}
