package audio

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/bwmarrin/discordgo" // transfer the dependencies to private repo
	"layeh.com/gopus"               // transfer the dependencies to private repo
)

type Audio struct {
	// file info
	filePath string
	fileName string
	fileSize int64
	// ffmpeg output
	output io.ReadCloser
	buffer *bufio.Reader
	// encode options
	encodeOptions *EncodeOptions
	// audio stream
	AudioStream chan []int16
}

func New(filePath string) (*Audio, error) {
	info, err := os.Lstat(filePath)

	if info != nil && info.IsDir() {
		return nil, ErrFileNotFound
	}
	switch {
	case os.IsNotExist(err):
		return nil, ErrFileNotFound
	case os.IsPermission(err):
		return nil, ErrInPermission
	case info != nil && info.IsDir():
		return nil, ErrNeedsToBeFile
	case err == nil:
		break
	default:
		return nil, err // unknown error
	}
	return &Audio{
			filePath:    filePath,
			fileName:    info.Name(),
			fileSize:    info.Size(),
			AudioStream: make(chan []int16, 2),
		},
		nil
}

func (audio *Audio) SetEncodeOptions(options *EncodeOptions) {
	audio.encodeOptions = options
}

func (audio *Audio) Load() error {
	if audio.encodeOptions == nil {
		return ErrNoEncodeOptions
	}
	cmd, err := audio.FFmpegCommand()
	if err != nil {
		return err
	}

	audio.output, err = cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	defer cmd.Process.Kill()
	defer audio.output.Close()

	// create a buffer to read the ffmpeg output
	audio.buffer = bufio.NewReaderSize(audio.output, 16384)

	// read the ffmpeg output and send it to the audio stream
	for {

		// read the audio buffer
		audioBuffer := make([]int16, audio.encodeOptions.FrameSize*audio.encodeOptions.Channels)

		// validate the audio buffer
		err = binary.Read(audio.buffer, binary.LittleEndian, &audioBuffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}

		// handle errors
		if err != nil {
			log.Println("Error reading ffmpeg stdout:", err)
			break
		}

		// skip empty audio buffers
		if len(audioBuffer) == 0 {
			continue
		}

		// send the audio buffer to the audio stream
		audio.AudioStream <- audioBuffer
	}

	return nil
}

func (audio *Audio) SendPCM(
	voice *discordgo.VoiceConnection,
) {

	opusEncoder, err := gopus.NewEncoder(
		audio.encodeOptions.FrameRate,
		audio.encodeOptions.Channels,
		gopus.Audio,
	)

	if err != nil {
		log.Println("01 - Error creating opus encoder", err)
		return
	}

	for {
		recv, ok := <-audio.AudioStream
		if !ok {
			log.Println("02 - PCM Channel closed")
			return
		}

		opus, err := opusEncoder.Encode(
			recv,
			audio.encodeOptions.FrameSize,
			audio.encodeOptions.GetMaxBytes(),
		)

		if err != nil {
			log.Println("03 - Error encoding pcm to opus", err)
			return
		}

		voice.OpusSend <- opus
	}
}

func (audio *Audio) FFmpegCommand() (*exec.Cmd, error) {
	args := []string{
		"-i", audio.filePath,
		"-f", "s16le",
		"-map", "0:a",
		"-reconnect", "1",
		"-reconnect_at_eof", "1",
		"-b:a", strconv.Itoa(audio.encodeOptions.Bitrate),
		"-ac", strconv.Itoa(audio.encodeOptions.Channels),
		"-ar", strconv.Itoa(audio.encodeOptions.FrameRate),
		"-threads", strconv.Itoa(audio.encodeOptions.Threads),
		"-packet_loss", strconv.Itoa(audio.encodeOptions.PacketLoss),
		"-frame_duration", strconv.Itoa(audio.encodeOptions.FrameDuration),
		"-compression_level", strconv.Itoa(audio.encodeOptions.CompressionLevel),
		"-application", "audio",
		"pipe:1",
	}
	cmd := exec.Command("ffmpeg", args...)
	return cmd, nil
}

func (audio *Audio) Close() {
	close(audio.AudioStream)
}

func (audio *Audio) GeneratePCM(
	voice *discordgo.VoiceConnection,
) {
	opusEncoder, err := gopus.NewEncoder(
		audio.encodeOptions.FrameRate,
		audio.encodeOptions.Channels,
		gopus.Audio,
	)
	if err != nil {
		log.Println("Error creating opus encoder", err)
		return
	}

	for audioBuffer := range audio.AudioStream {
		opus, err := opusEncoder.Encode(audioBuffer, audio.encodeOptions.FrameSize, 960)
		if err != nil {
			log.Println("Error encoding audio", err)
			return
		}
		voice.OpusSend <- opus
	}
}
