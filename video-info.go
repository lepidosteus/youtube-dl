package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func getVideoInfo(videoId string) (string, error) {
	url := "http://youtube.com/get_video_info?video_id=" + videoId
	log("Requesting url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("An error occured while requesting the video information: '%s'", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("An error occured while requesting the video information: non 200 status code received: '%s'", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("An error occured while reading the video information: '%s'", err)
	}
	log("Got %d bytes answer", len(body))
	return string(body), nil
}

func ensureFields(source url.Values, fields []string) (err error) {
	for _, field := range fields {
		if _, exists := source[field]; !exists {
			return fmt.Errorf("Field '%s' is missing in url.Values source", field)
		}
	}
	return nil
}

func getAudioData(response string) (resp *http.Response, err error) {
	answer, err := url.ParseQuery(response)
	if err != nil {
		err = fmt.Errorf("aaww snap, the server's answer: '%v'", err)
		return
	}
	err = ensureFields(answer, []string{"status", "player_response"})
	if err != nil {
		err = fmt.Errorf("Missing fields in the server's answer: '%s'", err)
		return
	}

	status := answer["status"]
	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			err = fmt.Errorf("'fail' response status found in the server's answer, reason: '%s'", reason[0])
		} else {
			err = errors.New(fmt.Sprint("'fail' response status found in the server's answer, no reason given"))
		}
		return
	}
	if status[0] != "ok" {
		err = fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
		return
	}
	log("Server answered with a success code")

	var URLm4a string

	// Getting the m4a url using regexp
	playerResponse := answer["player_response"]
	PResponseList := strings.Split(playerResponse[0], `},{`)
	for _, ListItem := range PResponseList {
		matched, _ := regexp.MatchString(`"itag":140`, ListItem)
		if matched {
			re := regexp.MustCompile(`https:\/\/[^"]+`)
			urls := re.FindStringSubmatch(ListItem)
			URLm4a = strings.Replace(urls[0], `\u0026`, "&", -1)
			break
		}
	}

	resp, err = http.Get(URLm4a)
	if err != nil {
		return nil, err
	}

	return

}

func decodeVideoInfo(response string) (streams streamList, err error) {
	// decode

	answer, err := url.ParseQuery(response)

	if err != nil {
		err = fmt.Errorf("parsing the server's answer: '%s'", err)
		return
	}

	err = ensureFields(answer, []string{"status", "url_encoded_fmt_stream_map", "title", "author"})
	if err != nil {
		err = fmt.Errorf("Missing fields in the server's answer: '%s'", err)
		return
	}

	status := answer["status"]
	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			err = fmt.Errorf("'fail' response status found in the server's answer, reason: '%s'", reason[0])
		} else {
			err = errors.New(fmt.Sprint("'fail' response status found in the server's answer, no reason given"))
		}
		return
	}
	if status[0] != "ok" {
		err = fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
		return
	}

	log("Server answered with a success code")

	/*
		for k, v := range answer {
			log("%s: %#v", k, v)
		}
	*/

	// read the streams map

	stream_map := answer["url_encoded_fmt_stream_map"]
	// read each stream

	streams_list := strings.Split(stream_map[0], ",")
	// Stream list is []string and inside got url=    &quality=  &type=
	log("Found %d streams in answer", len(streams_list))
	log("stream list:", streams_list)

	for stream_pos, stream_raw := range streams_list {
		stream_qry, err := url.ParseQuery(stream_raw)
		if err != nil {
			log(fmt.Sprintf("An error occured while decoding one of the video's stream's information: stream %d: %s\n", stream_pos, err))
			continue
		}
		err = ensureFields(stream_qry, []string{"quality", "type", "url"})
		if err != nil {
			log(fmt.Sprintf("Missing fields in one of the video's stream's information: stream %d: %s\n", stream_pos, err))
			continue
		}
		/* dumps the raw streams
		log(fmt.Sprintf("%v\n", stream_qry))
		*/
		stream := stream{
			"quality": stream_qry["quality"][0],
			"type":    stream_qry["type"][0],
			"url":     stream_qry["url"][0],
			"sig":     "",
			"title":   answer["title"][0],
			"author":  answer["author"][0],
		}
		if sig, exists := stream_qry["sig"]; exists { // old one
			stream["sig"] = sig[0]
		}

		if sig, exists := stream_qry["s"]; exists { // now they use this
			stream["sig"] = sig[0]
		}

		streams = append(streams, stream)

		quality := stream.Quality()
		if quality == QUALITY_UNKNOWN {
			log("Found unknown quality '%s'", stream["quality"])
		}

		format := stream.Format()
		if format == FORMAT_UNKNOWN {
			log("Found unknown format '%s'", stream["type"])
		}

		log("Stream found: quality '%s', format '%s'", quality, format)
	}

	log("Successfully decoded %d streams", len(streams))

	return
}
