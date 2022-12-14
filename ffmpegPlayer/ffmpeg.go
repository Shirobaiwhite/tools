package main

import (
	"io"
	"os/exec"
)

type FfmpegStreamer struct {
	subproc *exec.Cmd
	stdin   io.WriteCloser
}

func NewFfmpegStreamer(ffmpegPath, frameRes, fps, rtmpUrl string) *FfmpegStreamer {
	return NewFfmpegStreamer2(ffmpegPath,
		"-stream_loop", "-1",
		"-re",
		"-i", "./replayvid.mp4",
		"-f", "flv",
		"-vcodec", "libx264",
		"rtmp://localhost/live/livestream")
}

func NewFfmpegStreamer2(ffmpegPath string, ffmpegArgs ...string) *FfmpegStreamer {
	subproc := exec.Command(ffmpegPath, ffmpegArgs...)
	if subproc == nil {
		return nil
	}
	stdin, err := subproc.StdinPipe()
	if err != nil {
		return nil
	}
	if err := subproc.Start(); err != nil {
		return nil
	}

	return &FfmpegStreamer{subproc: subproc, stdin: stdin}
}
