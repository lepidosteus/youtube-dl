youtube-dl
=================

Go tool to download videos from youtube, with quality and format selection

Install
-------

    (sudo) go get github.com/lepidosteus/youtube-dl

Your can then start using it using youtube-dl

    youtube-dl http://www.youtube.com/watch?v=xjfsmKmK9qc

Usage
-----

    youtube-dl [-verbose -overwrite -output /p/a/t/h -quality list -format list] videoId|url

Examples:

    youtube-dl http://www.youtube.com/watch?v=xjfsmKmK9qc

    youtube-dl youtu.be/xjfsmKmK9qc

    youtube-dl xjfsmKmK9qc

    youtube-dl -output funny_video.mp4 -quality hd720 -format mp4 xjfsmKmK9qc

Parameters
----------

<table>
  <tr>
    <th>Parameter</th><th>Default value</th><th>Allowed values</th><th>Description</th><th>Example</th>
  </tr>
  <tr>
    <td>-overwrite|-overwrite=VALUE</td><td>false</td><td>true,false</td><td>if true, the destination file will be overwritten if it already exists</td><td>-overwrite</td>
  </tr>
  <tr>
    <td>-verbose|-verbose=VALUE</td><td>false</td><td>true,false</td><td>if true, various status messages will be shown</td><td>-verbose</td>
  </tr>
  <tr>
    <td>-output VALUE</td><td>./video.%format%</td><td>a valid path</td><td>path where to write the downloaded file, use %format% for dynamic extension depending on format selected (eg: 'video.%format%' would be written as 'video.mp4' if the mp4 format is selected)</td><td>-output "$HOME/funny_video.%format%"
  </tr>
  <tr>
    <td>-quality VALUE[,VALUE...]</td><td>hd720,max</td><td>hd720,large,medium,small,min,max</td><td>comma separated list of desired video quality, in decreasing priority. Use 'max' (or 'min') to automatically select the best (or worst) possible quality available for this video; exemple: '-quality hd720,max': select hd720 quality, if not available then select the best quality available</td><td>-quality small,min</td>
  </tr>
  <tr>
    <td>-format VALUE[,VALUE...]</td><td>mp4,flv,webm,3gp</td><td>mp4,flv,webm,3gp</td><td>comma separated list of desired video format, in decreasing priority</td><td>-format mp4,flv</td>
  </tr>
</table>
