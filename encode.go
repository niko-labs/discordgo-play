package audio

type EncodeOptions struct {
	Volume           int // change audio volume (256=normal)
	FrameRate        int // audio sampling rate (ex 48000)
	FrameSize        int // audio frame size (ex 960)
	FrameDuration    int // audio frame duration (default 20ms)
	Bitrate          int // audio encoding bitrate in kb/s
	Threads          int // Number of threads to use (0 for auto)
	Channels         int // Number of audio channels
	PacketLoss       int // expected packet loss percentage
	CompressionLevel int // Compression level, higher is better quality but slower encoding (0 - 10)
	maxBytes         int // maximum size of the audio buffer
}

func NewDefaultEncodeOptions() *EncodeOptions {
	eo := &EncodeOptions{
		Volume:           256,
		FrameRate:        48000,
		FrameSize:        960,
		FrameDuration:    20,
		Bitrate:          96000,
		Threads:          0,
		Channels:         2,
		PacketLoss:       0,
		CompressionLevel: 10,
	}
	eo.maxBytes = (eo.FrameSize * 2) * eo.Channels
	return eo
}
func NewEncodeOptions(volume, frameRate, frameSize, frameDuration, bitrate, threads, channels, packetLoss, compressionLevel int) *EncodeOptions {
	// You really know what you're doing, right?
	eo := &EncodeOptions{
		Volume:           volume,
		FrameRate:        frameRate,
		FrameSize:        frameSize,
		FrameDuration:    frameDuration,
		Bitrate:          bitrate,
		Threads:          threads,
		Channels:         channels,
		PacketLoss:       packetLoss,
		CompressionLevel: compressionLevel,
	}
	eo.maxBytes = (eo.FrameSize * 2) * eo.Channels
	return eo
}

func (e *EncodeOptions) GetMaxBytes() int {
	return e.maxBytes
}
