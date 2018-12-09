package main

import (
	"io"
	"os"
	"os/exec"
	"fmt"
)

// custom io.WriteCloser to handle on the fly mp3 convertion

const FFMPEG = "ffmpeg"

type FFMpegWriter struct {
	cmd *exec.Cmd
	stdin io.WriteCloser
}

func (w *FFMpegWriter) Write(p []byte) (n int, err error) {
	return w.stdin.Write(p)
}

func (w *FFMpegWriter) Close() error {
	w.stdin.Close()
	return w.cmd.Wait()
}

func getFFmpegWriter(path string, audio_bitrate uint) (w *FFMpegWriter, err error) {
	_, err = exec.LookPath(FFMPEG)
	if err != nil {
		return nil, fmt.Errorf("you need to install ffmpeg to convert to mp3: %s", err)
	}

	w = &FFMpegWriter{
		exec.Command(FFMPEG, "-i", "-", "-ab", fmt.Sprintf("%dk", audio_bitrate), path),
		nil,
	}
	w.stdin, err = w.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	w.cmd.Start()
	return w, nil
}

func getWriter(cfg *Config, stream stream) (out io.WriteCloser, err error) {
	path := cfg.OutputPath(stream)

	if _, err = os.Stat(path); err == nil && cfg.overwrite == false {
		return nil, fmt.Errorf("the destination file '%s' already exists and overwrite set to false", path)
	}

	if cfg.isMp3() {
		fmt.Printf("Converting video to mp3 file at '%s' ...\n", path)
		out, err = getFFmpegWriter(path, cfg.AudioBitrate(stream))
	} else {
		fmt.Printf("Downloading video to disk at '%s' ...\n", path)
		out, err = os.Create(path)
	}

	if err != nil {
		return nil, fmt.Errorf("opening destination file: %s", err)
	}

	log("Destination opened at '%s'", path)

	return out, nil
}










