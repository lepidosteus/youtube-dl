package main

import (
	"fmt"
	"errors"
	"net/url"
	"net/http"
	"io/ioutil"
	"strings"
)

func getVideoInfo(videoId string) (string, error) {
	url := "http://youtube.com/get_video_info?video_id=" + videoId
	log("Requesting url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New(fmt.Sprintf("An error occured while requesting the video information: '%s'", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("An error occured while requesting the video information: non 200 status code received: '%s'", err))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("An error occured while reading the video information: '%s'", err))
	}
	log("Got %d bytes answer", len(body))
	return string(body), nil
}

func decodeVideoInfo(response string) (streams streamList, err error) {
	// decode

	answer, err := url.ParseQuery(response)
	if err != nil {
		err = errors.New(fmt.Sprintf("parsing the server's answer: '%s'", err))
		return
	}

	// check the status

	status, ok := answer["status"]
	if !ok {
		err = errors.New(fmt.Sprint("no response status found in the server's answer"))
		return
	}
	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			err = errors.New(fmt.Sprintf("'fail' response status found in the server's answer, reason: '%s'", reason[0]))
		} else {
			err = errors.New(fmt.Sprint("'fail' response status found in the server's answer, no reason given"))
		}
		return
	}
	if status[0] != "ok" {
		err = errors.New(fmt.Sprintf("non-success response status found in the server's answer (status: '%s')", status))
		return
	}

	log("Server answered with a success code")

	// read the streams map

	stream_map, ok := answer["url_encoded_fmt_stream_map"]
	if !ok {
		err = errors.New(fmt.Sprint("no stream map found in the server's answer"))
		return
	}

	// read each stream

	streams_list := strings.Split(stream_map[0], ",")

	log("Found %d streams in answer", len(streams_list))

	for stream_pos, stream_raw := range streams_list {
		stream_qry, err := url.ParseQuery(stream_raw)
		if err != nil {
			log(fmt.Sprintf("An error occured while decoding one of the video's stream's information: stream %d: %s\n", stream_pos, err))
			continue
		}
		stream := stream{
			"quality": stream_qry["quality"][0],
			"type": stream_qry["type"][0],
			"url": stream_qry["url"][0],
			"sig": stream_qry["sig"][0],
		}
		streams = append(streams, stream)

		log("Stream found: quality '%s', format '%s'", stream.getQuality(), stream.getFormat())
	}

	log("Successfully decoded %d streams", len(streams))

	return
}
