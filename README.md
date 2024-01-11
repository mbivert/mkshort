## Introduction
``mkshort`` is a small go tool to create short videos containing
a diaporama of images with overlayed text, via [ffmpeg(1)][ffmpeg-doc].
The input is a text file like this one:

	:virgin-of-the-rocks-paris.jpg
	0.7 3.5
		Here is what seems \\
		to be \\
		a little-known fact \\
	+0.7 3
		About this famous \\
		Leonardo painting, \\

	...

	:virgin-of-the-rocks-london.jpg
	+0.7 4
		This also holds in the \\
		"London" version of the \\
		painting \\
	+0.7 2

![Sample image ](https://github.com/mbivert/mkshort/blob/master/virgin-of-the-sample.jpg)

An sample output is available
[here](https://www.ganjingworld.com/shorts/1gemffsenl01fJ8I4DxARv8qi1k81c)
or [here](https://youtube.com/shorts/d4iW-_-ETb4);
the [audio](http://www.baroquemusic.org/19Web.html) has been
manually added with:
```
ffmpeg -y \
	-i virgin-of-the-rocks.mp4 \
	-i $MUSIC/baroquemusic.com/BMC19T1VivaldiSeasonsSpring.mp3 \
	-map 0:v -map 1:a -c:v copy -shortest virgin-of-the-rocks-audio.mp4
```

While [ffmpeg(1)][ffmpeg-doc] provides a text overlay filter, it
can be limited. Instead, we're using LaTeX (LuaLaTeX) to provide
a more complete features set (colored [emojis (.pdf)][emoji-pkg],
mathematical typesetting, etc.).

A quick man page is available: [mkshort(1)][mkshort-1]. An old, semi-automatic
[awk(1)][awk-1]/[sh(1)][sh-1]-based prototype is available in
[``old/prototype.sh``][old-prototype-sh]. A "Virgin of the rocks" example
is provided. Quick help:

```
$ make help
Available targets:
	mkshort    : build mkshort
	all        : build mkshort; run tests
	clean      : removed compiled files
	tests      : run automated tests
	example    : display the ffmpeg(1) command,
	             build (a bit slow) and play the example.
	install    : install to /bin/ andqq /usr/share/man/man1/
	uninstall  : remove installed files
$ mkshort -help
Usage of mkshort:
  -b string
    	rescale padding color (default "black")
  -c string
    	LaTeX command, splitted on spaces (default "lualatex")
  -d string
    	cache directory (default "/home/mb/.mkshort")
  -e string
    	extension for the compiled text images (default ".png")
  -f string
    	output pixel format (default "yuv420p")
  -h int
    	output height (default 1920)
  -i string
    	text indentation (default "\t")
  -l string
    	shell to run the compiled ffmpeg(1) command (default "/bin/sh")
  -m	Enable/disable faststart (default true)
  -p string
    	LaTeX template (path)
  -r int
    	output framerate (default 30)
  -s float
    	default waiting time between text (default 0.8)
  -t string
    	LaTeX template (string) (default "\\documentclass[preview,convert={density=600,outext=.png,command=\\unexpanded{ {\\convertexe\\space -density \\density\\space\\infile\\space \\ifx\\size\\empty\\else -resize \\size\\fi\\space -quality 90 -trim +repage -background \"rgba(50,50,50,0.5)\" -bordercolor \"rgba(50,50,50,0.5)\" -border 25 -flatten \\outfile} } }]{standalone}\n% Requires lualatex\n\\usepackage{emoji}\n\n\\usepackage{xcolor}\n\\usepackage{amsmath}\n\\begin{document}\n\n\\begin{center}\n\\textcolor{white}{ {{ .text }} }\n\\end{center}\n\\end{document}\n")
  -w int
    	output width (default 1080)
  -x	Compiled command printed, not ran
  -y	Allow automatic output overwrite (default true)
```

[ffmpeg-doc]: https://ffmpeg.org/ffmpeg.html
[mkshort-1]: https://github.com/mbivert/mkshort/blob/master/mkshort.1
[awk-1]: https://man.openbsd.org/awk.1
[sh-1]: https://man.openbsd.org/sh.1
[old-prototype-sh]: https://github.com/mbivert/mkshort/blob/master/old/prototype.sh
[emoji-pkg]: https://ctan.math.washington.edu/tex-archive/macros/luatex/latex/emoji/emoji-doc.pdf
