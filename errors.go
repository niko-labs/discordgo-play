package audio

import (
	"errors"
)

var (
	// ErrFileNotFound is returned when the file is not found
	ErrFileNotFound       = errors.New("file not found")
	ErrNeedsToBeFile      = errors.New("needs to be a file")
	ErrInPermission       = errors.New("permission denied")
	ErrFFmpegNotInstalled = errors.New("ffmpeg is not found")
	ErrFFmpegFailed       = errors.New("ffmpeg failed to start")
	ErrFFmpegKilled       = errors.New("ffmpeg was killed")
	ErrFFmpegEOF          = errors.New("ffmpeg reached EOF")
	ErrFFmpegRead         = errors.New("ffmpeg read error")
	ErrNoEncodeOptions    = errors.New("no encode options, please set encode options")
)
