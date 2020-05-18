package signal

import "math"

func sinWave(sampleRate uint16, samples uint32, freq uint16) []byte {
	soundData := make([]byte, samples)
	cycleRate := float64(sampleRate) / float64(freq)

	var x float64
	for i := uint32(0); i < samples; i++ {
		x = float64(i) / cycleRate * (2 * math.Pi)
		soundData[i] = byte(((math.Sin(x) + 1) / 2) * 255)
	}
	return soundData
}
