package main

import (
	"flag"
	"regexp"
	"errors"
	"strings"
	"fmt"
	"os"
	"path/filepath"
)

var sortedQualities []string = []string{
	QUALITY_HIGHRES,
	QUALITY_HD1080,
	QUALITY_HD720,
	QUALITY_LARGE,
	QUALITY_MEDIUM,
	QUALITY_SMALL,
	QUALITY_UNKNOWN,
}

var formatsTrigger map[string]string = map[string]string{
	FORMAT_MP4: "video/mp4",
	FORMAT_FLV: "video/x-flv",
	FORMAT_WEBM: "video/webm",
	FORMAT_3GP: "video/3gpp",
}

var sortedFormats []string = []string{
	FORMAT_MP4,
	FORMAT_FLV,
	FORMAT_WEBM,
	FORMAT_3GP,
	FORMAT_UNKNOWN,
}

// comma delimited parameters

type commaStringList struct {
	values []string
	allowed map[string]struct{}
}

func (sl *commaStringList) String() string {
	return strings.Join(sl.values, ",")
}

func (sl *commaStringList) Set(value string) error {
	sl.values = sl.values[0:0]
	var exists bool
	for _, s := range strings.Split(value, ",") {
		_, exists = sl.allowed[s]
		if len(sl.allowed) > 0 && !exists {
			return errors.New(fmt.Sprintf("non allowed value '%s'", s))
		}
		sl.values = append(sl.values, s)
	}
	return nil
}

func CreateCommaStringList(values []string, allowed []string) *commaStringList {
	sl := &commaStringList{[]string{}, map[string]struct{}{}}
	for _, value := range values {
		sl.values = append(sl.values, value)
	}
	for _, value := range allowed {
		sl.allowed[value] = struct{}{}
	}
	return sl
}

// our config struct

type Config struct {
	verbose bool
	output string	// path
	overwrite bool
	quality *commaStringList
	format *commaStringList
	videoId string
	toMp3 bool
}

var cfg *Config = &Config{
	false,
	DEFAULT_DESTINATION,
	false,
	CreateCommaStringList(
		[]string{QUALITY_HD720, QUALITY_MAX},
		append([]string{QUALITY_MAX, QUALITY_MIN}, sortedQualities...),
	),
	CreateCommaStringList(
		[]string{FORMAT_MP4, FORMAT_FLV, FORMAT_WEBM, FORMAT_3GP},
		sortedFormats,
	),
	"",
	false,
}

// reads the videoId property and try to find what we need inside
func (cfg *Config) findVideoId() (videoId string, err error) {
	videoId = cfg.videoId
	if strings.Contains(videoId, "youtu") || strings.ContainsAny(videoId, "\"?&/<%=") {
		log("Provided video id seems to be an url, trying to detect")
		re_list := []*regexp.Regexp{
			regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`([^"&?/=%]{11})`),
		}
		for _, re := range re_list {
			if is_match := re.MatchString(videoId); is_match {
				subs := re.FindStringSubmatch(videoId)
				videoId = subs[1]
			}
		}
	}
	log("Found video id: '%s'", videoId)
	if strings.ContainsAny(videoId, "?&/<%=") {
		return videoId, errors.New("invalid characters in video id")
	}
	if len(videoId) < 10 {
		return videoId, errors.New("the video id must be at least 10 characters long")
	}
	return videoId, nil
}

func (cfg *Config) isMp3() bool {
	return cfg.toMp3
}

func (cfg *Config) OutputPath(stream stream) (path string) {
	path = strings.Replace(cfg.output, "%format%", stream.Format(), -1)
	path = strings.Replace(path, "%title%", stream["title"], -1)
	path = strings.Replace(path, "%author%", stream["author"], -1)
	return
}

func (cfg *Config) selectStream(streams streamList) (stream stream, err error) {
	if len(streams) < 1 {
		return nil, errors.New("no streams found")
	}
	valid_streams := streamList{}
	for _, format := range cfg.format.values {
		for _, s := range streams {
			if s.Format() == format {
				valid_streams = append(valid_streams, s)
			}
		}
		if len(valid_streams) >= 1 {
			log("Found format '%s', with %d streams", format, len(valid_streams))
			break
		}
	}
	if len(valid_streams) < 1 {
		return nil, errors.New("no streams match the requested formats")
	}
	streams = valid_streams
	valid_streams = streamList{}
	for _, quality := range cfg.quality.values {
		for _, s := range streams {
			if s.Quality() == quality {
				valid_streams = append(valid_streams, s)
			}
		}
		if len(valid_streams) >= 1 {
			log("Found quality '%s', with %d streams", quality, len(valid_streams))
			break
		}
	}
	if len(valid_streams) < 1 {
		return nil, errors.New("no streams match the requested qualities")
	}
	return valid_streams[0], nil
}

// display usage and quit
func error_usage() {
	fmt.Println("usage: youtube-dl [-verbose -mp3 -overwrite -output /p/a/t/h -quality list -format list] videoId|url")
	flag.PrintDefaults()
	os.Exit(1)
}

// load config from cli
func init() {
	flag.BoolVar(&cfg.verbose, "verbose", false, "if true, various status messages will be shown.")

	flag.BoolVar(&cfg.toMp3, "mp3", false, "if true, the file's audio stream will be converted to an mp3 file.")

	flag.BoolVar(&cfg.overwrite, "overwrite", false, "if true, the destination file will be overwritten if it already exists.")

	flag.StringVar(&cfg.output, "output", "", "path where to write the downloaded file. Use %format% for dynamic extension depending on format selected (eg: -output 'video.%format%' would be written as 'video.mp4' if the mp4 format is selected). %author% and %title% will be replaced by the uploader's name and the video's title, respectively. Use the .mp3 extension to convert the video to an mp3 file on the fly (eg: -ouput 'audio.mp3').")

	flag.Var(cfg.quality, "quality", "comma separated list of desired video quality, in decreasing priority. Use 'max' (or 'min') to automatically select the best (or worst) possible quality available for this video. Allowed values: " + strings.Join(sortedQualities, ", ") + ". Exemple: '-quality hd720,max': select hd720 quality, if not available then select the best quality available.")

	flag.Var(cfg.format, "format", "comma separated list of desired video format, in decreasing priority. Allowed values: " + strings.Join(sortedFormats, ", ") + ".")

	flag.Parse()

	if len(cfg.output) <= 0 {
		// if no path given, guess one

		if cfg.isMp3() {
			cfg.output = DEFAULT_DESTINATION_MP3
		} else {
			cfg.output = DEFAULT_DESTINATION
		}
		log("No output specified, defaulting to '%s'", cfg.output)
	} else if (filepath.Ext(cfg.output) == ".mp3") {
		// if a path is given and its for a mp3, make sure to convert

		cfg.toMp3 = true
		log("Converting video to mp3 due to output parameter")
	} else if (filepath.Ext(cfg.output) == ".%format%" && cfg.toMp3) {
		// if a path is given, and its a dynamic extension, and we asked for mp3, change path now
		cfg.output = cfg.output[:len(cfg.output)-8] + "mp3"
		log("Replacing output with '%s' due to parameters", cfg.output)
	} else if (filepath.Ext(cfg.output) != ".mp3" && cfg.toMp3) {
		// if we ask for mp3 but the output we gave doesn't have it, append it

		cfg.output = cfg.output + ".mp3"
		log("Replacing output with '%s' due to parameters", cfg.output)
	}

	log("Configuration:")

	log("\tVerbose: %t", cfg.verbose)
	log("\tOverwrite: %t", cfg.overwrite)
	log("\tQuality: %s", strings.Join(cfg.quality.values, ","))
	log("\tFormat: %s", strings.Join(cfg.format.values, ","))
	log("\tOutput: %s", cfg.output)

	// replace min/max quality by their actual values
	for idx := len(cfg.quality.values) - 1; idx >= 0; idx = idx - 1 {
		quality := cfg.quality.values[idx]
		if quality == QUALITY_MAX || quality == QUALITY_MIN {
			plug := append([]string{}, sortedQualities...)
			if quality == QUALITY_MIN {
				// reverse the order
				for i, j := 0, len(plug) - 1; i < j; i, j = i + 1, j - 1 {
					plug[i], plug[j] = plug[j], plug[i]
				}
			}
			cfg.quality.values = append(
				cfg.quality.values[:idx],
				append(
					plug,
					cfg.quality.values[idx + 1:]...,
				)...,
			)
		}
	}

	log("\tExtended quality: %s", strings.Join(cfg.quality.values, ","))

	if flag.NArg() != 1 {
		fmt.Println("ERROR: no videoId or url given")
		error_usage()
	}

	cfg.videoId = flag.Arg(0)

	log("\tVideo: %s", cfg.videoId)
}










