# DiscordGo Play Audio Library

## Overview

This library provides a robust solution for loading, encoding, and streaming audio files in Discord applications using `discordgo` and `gopus`. It supports various audio formats, configuration options for encoding, and efficient audio streaming to a Discord voice connection.

## Key Features

1. **Audio Loading**:
    - Load audio files using `ffmpeg`.
    - Convert audio files to PCM format (`s16le`) for processing.

2. **Encoding Options Configuration**:
    - Configure options like bitrate, channels, and frame rate.
    - Adjust audio quality and performance with configurable settings.

3. **Audio Streaming**:
    - Create an audio stream (`AudioStream`) to process and send audio data.
    - Efficiently read from the `ffmpeg` output buffer and send audio data to the stream.

4. **Opus Encoding**:
    - Encode PCM data into Opus packets using `gopus`.
    - Send Opus packets to a Discord voice connection (`discordgo.VoiceConnection`).

## Installation

To use this library, you need to have `ffmpeg` installed and accessible in your system PATH. You also need to install the necessary Go dependencies.

1. Install `ffmpeg`:
    - On macOS: `brew install ffmpeg`
    - On Linux: `sudo apt-get install ffmpeg`
    - On Windows: Download and install from [ffmpeg.org](https://ffmpeg.org/)

2. Install Go dependencies:
    ```bash
    go get github.com/bwmarrin/discordgo
    go get layeh.com/gopus
    ```

## Basic Usage

Here is a basic example of how to use the library:

```go
import (
    "log"
    "github.com/bwmarrin/discordgo"
    dgoAudio "github.com/niko-labs/discordgo-play"
)

func main() {
    audioPath := "/path/to/your/audiofile.mp3"
    voiceConn := &discordgo.VoiceConnection{} // Assume the voice connection is already initialized

    err := loadSound(audioPath, voiceConn)
    if err != nil {
        log.Fatal("Error loading sound:", err)
    }
}

func loadSound(audioPath string, vc *discordgo.VoiceConnection) error {
    defer vc.Speaking(false)
    defer vc.Disconnect()

    options := dgoAudio.NewDefaultEncodeOptions()

    audioOut, err := dgoAudio.New(audioPath)
    if err != nil {
        log.Panicln(err)
    }

    audioOut.SetEncodeOptions(options)
    go audioOut.GeneratePCM(vc)

    err = audioOut.Load()
    if err != nil {
        log.Panicln(err)
    }

    return nil
}
```

## Configuration

You can configure the encoding options using the `EncodeOptions` struct. This allows you to adjust settings such as bitrate, channels, frame rate, and more to optimize audio quality and performance.

## Error Handling

The library includes detailed error logging to assist with debugging and issue resolution. Ensure you handle errors appropriately in your application to maintain stability and performance.

## Compatibility

This library is tested with `Go` version `1.22` and `ffmpeg` version `7.0`. Ensure your development environment meets these requirements for optimal performance.

## Conclusion

Version 1.0.0 marks the first stable release of the audio library, providing a solid foundation for audio manipulation and streaming in Discord applications. We hope this library proves useful for developers seeking a robust and efficient audio integration solution.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.