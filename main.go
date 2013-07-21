package main

import (
	"fmt"
)

const (
	QUALITY_HD720 = "hd720"
	QUALITY_LARGE = "large"
	QUALITY_MEDIUM = "medium"
	QUALITY_SMALL = "small"
	QUALITY_MIN = "min"
	QUALITY_MAX = "max"
	QUALITY_UNKNOWN = "unknown"

	FORMAT_MP4 = "mp4"
	FORMAT_WEBM = "webm"
	FORMAT_FLV = "flv"
	FORMAT_3GP = "3ggp"
	FORMAT_UNKNOWN = "unknown"

	DEFAULT_DESTINATION = "./video.%format%"
)

func log(format string, params ...interface{}) {
	if cfg.verbose {
		fmt.Printf(format + "\n", params...)
	}
}

func main() {
	videoId, err := cfg.findVideoId()
	if err != nil {
		fmt.Printf("ERROR: unable to detect the video id: %s\n", err)
		return
	}

	response, err := getVideoInfo(videoId)
	if err != nil {
		fmt.Printf("ERROR: unable to request the video information: %s\n", err)
		return
	}

	streams, err := decodeVideoInfo(response)
	if err != nil {
		fmt.Printf("ERROR: unable to decode the server's answer: %s\n", err)
		return
	}

	stream, err := cfg.selectStream(streams)
	if err != nil {
		fmt.Printf("ERROR: unable to select a stream: %s\n", err)
		return
	}

	err = stream.download(cfg.getOutputPath(stream), cfg.overwrite)
	if err != nil {
		fmt.Printf("ERROR: unable to download the stream: %s\n", err)
		return
	}

	return
}
