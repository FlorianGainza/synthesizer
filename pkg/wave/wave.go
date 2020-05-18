package wave

import (
	"math"
)

// Header describes the wave format header
type Header struct {
	RiffMark      [4]byte // 0  4  Contains the letters "RIFF"
	FileSize      int32   // 4  4  Size of the rest of the chunk
	WaveMark      [4]byte // 8  4  Contains the letters "WAVE"
	FmtMark       [4]byte // 12 4  Contains the letters "fmt "
	FormatSize    int32   // 16 4  Size of the rest of the subchunk
	FormatType    int16   // 20 2  Linear quantization (PCM = 1)
	NumChans      int16   // 22 2  Number of channels
	SampleRate    int32   // 24 4  Number of samples per second
	ByteRate      int32   // 28 4  Number of bytes per second
	BytesPerFrame int16   // 32 2  Number of bytes per sample * channels
	BitsPerSample int16   // 34 2  Number of bits per sample
	DataMark      [4]byte // 36 4  Contains the letters "data"
	DataSize      uint32  // 40 4  Size of the data
}

// Sin genererates a sinusoidal signal in a binary form
// Wave sample value will be between 0 and 255
// Will encode according to :
//     - The sample rate per seconde
//     - The signal frequency per seconde
//     - The number of samples
func Sin(sampleRate uint16, freq uint16, nbSamples uint32) []byte {
	// Number of samples in one cycle of signal
	sigCycle := float64(sampleRate) / float64(freq)
	samples := make([]byte, nbSamples)

	var x float64
	for i := uint32(0); i < nbSamples; i++ {
		x = float64(i) / sigCycle * (2 * math.Pi)
		samples[i] = byte(((math.Sin(x) + 1) / 2) * 255)
	}

	return samples
}

// Square genererates a square signal in a binary form
// Wave sample value will be between 0 and 255
// Will encode according to :
//     - The sample rate per seconde
//     - The signal frequency per seconde
//     - The number of samples
func Square(sampleRate uint16, freq uint16, nbSamples uint32) []byte {
	// Number of samples in one cycle of signal
	sigCycle := float64(sampleRate) / float64(freq)
	samples := make([]byte, nbSamples)

	var isPositive bool
	for i := uint32(0); i < nbSamples; i++ {
		isPositive = math.Mod((float64(i)/sigCycle), 1) < 0.5
		if isPositive {
			samples[i] = byte(255)
		} else {
			samples[i] = byte(0)
		}
	}

	return samples
}
