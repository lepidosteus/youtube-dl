youtube-dl
=================

Go tool to download videos from youtube, with quality and format selection

**NOTE**: YouTube has changed their code so that url's signature need to be deciphered, and they rotate the ciphering code quite often. As such, this tool won't be able to download the video files anymore.

Since I'm not actively using it I won't bother working around it (the ciphering code is in youtube js player file).
If you need to download youtube videos, I recommend https://rg3.github.io/youtube-dl/

Install
-------

    (sudo) go get github.com/lepidosteus/youtube-dl

Your can then start using it using youtube-dl

    youtube-dl http://www.youtube.com/watch?v=xjfsmKmK9qc


Automatic naming
----------------

If used in the output path, "%author%" and "%title%" will be replaced by their value (respectively the uploader's name and the video's title)

    youtube-dl -output "%title%.mp3" http://www.youtube.com/watch?v=xjfsmKmK9qc

Usage
-----

    youtube-dl [-verbose -mp3 -overwrite -output /p/a/t/h -quality list -format list] videoId|url

Examples:

    youtube-dl http://www.youtube.com/watch?v=xjfsmKmK9qc

    youtube-dl youtu.be/xjfsmKmK9qc

    youtube-dl xjfsmKmK9qc

This video has highres and 1080p version, but we can ask for 720p only:

    youtube-dl -quality hd720 9dgSa4wmMzk

This one will fail, as there is no 720p version available for this video

    youtube-dl -output learning_go.mp4 -quality hd720 -format mp4 xjfsmKmK9qc

This, however, will work and download the video in the next best quality (medium)

    youtube-dl -output learning_go.mp4 -quality hd720,max -format mp4 xjfsmKmK9qc

Video ID detection
------------------

The parsing tries to be smart about finding out what the video's ID is, you can give it an url, an id, an embed code fragment, ... And it will do its best.

MP3 convertion
--------------

If ffmpeg is installed, it is possible to extract the audio stream to an mp3 file on the fly. Either pass the -mp3 parameter, or give an output path ending in .mp3

    youtube-dl -mp3 http://www.youtube.com/watch?v=xjfsmKmK9qc

    youtube-dl -output "audio.mp3" http://www.youtube.com/watch?v=xjfsmKmK9qc

Parameters
----------

<table>
  <tr>
    <th>Parameter</th><th>Default value</th><th>Allowed values</th><th>Example</th>
  </tr>
  <tr>
    <td>-overwrite<br>-overwrite=VALUE</td><td>false</td><td>true<br>false</td><td>-overwrite</td>
  </tr>
  <tr>
    <td colspan="4">if true, the destination file will be overwritten if it already exists</td>
  </tr>
  <tr>
    <td>-verbose<br>-verbose=VALUE</td><td>false</td><td>true<br>false</td><td>-verbose</td>
  </tr>
  <tr>
    <td colspan="4">if true, various status messages will be shown</td>
  </tr>
  <tr>
    <td>-mp3<br>-mp3=VALUE</td><td>false</td><td>true<br>false</td><td>-mp3</td>
  </tr>
  <tr>
    <td colspan="4">if true, the file's audio stream will be converted to an mp3 file</td>
  </tr>
  <tr>
    <td>-audio-bitrate</td><td>-audio-bitrate 0<td>0<br />any positive number</td><td>-audio-bitrate 128</td>
  </tr>
  <tr>
    <td colspan="4">The bitrate to use for audio files when converting to mp3. If set to 0 (which is the default) the bitrate will be set automatically depending on the quality of the downloaded video file</td>
  </tr>
  <tr>
    <td>-output VALUE</td><td>./video.%format%</td><td>a valid path<br>Tokens:<br>%format%<br>%author%<br>%title%</td><td>-output "$HOME/funny_video.%format%"
  </tr>
  <tr>
    <td colspan="4">path where to write the downloaded file<br>Use %format% for dynamic extension depending on format selected (eg: 'video.%format%' would be written as 'video.mp4' if the mp4 format is selected).<br>%author% and %title% will be replaced by the uploader's name and the video's title, respectively.<br>Use the .mp3 extension to convert the video to an mp3 file on the fly (eg: -ouput 'audio.mp3')</td>
  </tr>
  <tr>
    <td>-quality VALUE[,VALUE...]</td><td>hd720,max</td><td>highres<br>hd1080<br>hd720<br>large<br>medium<br>small<br>min<br>max</td><td>-quality small,min</td>
  </tr>
  <tr>
    <td colspan="4">comma separated list of desired video quality, in decreasing priority. Use 'max' (or 'min') to automatically select the best (or worst) possible quality available for this video; exemple: '-quality hd720,max': select hd720 quality, if not available then select the best quality available</td>
  </tr>
  <tr>
    <td>-format VALUE[,VALUE...]</td><td>mp4,flv,webm,3gp</td><td>mp4<br>flv<br>webm<br>3gp</td><td>-format mp4,flv</td>
  </tr>
  <tr>
    <td colspan="4">comma separated list of desired video format, in decreasing priority</td>
  </tr>
</table>
