package main

import (
	"testing"
	"strings"
	"path/filepath"
	"text/template"
	"os"
)

// XXX assumed to exists AND to be absolute, which is
// otherwise guaranteed by 'mkshort.go:/^func doInit\('
var D = filepath.Join(os.Getenv("PWD"), ".cache")

// The default template './mkshort.go:/^var textTmpl' has change, and will
// change in the future. This should be sufficient for tests.
var T = `\documentclass[preview,convert={density=600,outext=.png,command=\unexpanded{ {\convertexe\space -density \density\space\infile\space \ifx\size\empty\else -resize \size\fi\space -quality 90 -trim +repage -background "rgba(50,50,50,0.5)" -bordercolor "rgba(50,50,50,0.5)" -border 25 -flatten \outfile} } }]{standalone}
% Requires lualatex
\usepackage{emoji}

\usepackage{xcolor}
\usepackage{amsmath}
\begin{document}

\begin{center}
\textcolor{white}{ {{ .text }} }
\end{center}
\end{document}
`

func TestCompile(t *testing.T) {
	doTests(t, []test{
		// Generated command might be invalid in some cases
		{
			"~empty input, ~no config",
			compile,
			[]any{
				[]string{},
				[]string{},
				[]string{},
				[]string{},
				"",
				&State{
					faststart   : false,
					overwrite   : false,
					width       : 0,
					height      : 0,
					framerate   : 0,
					defaultWait : 0.8,
					indent      : "\t",
					input       : nil,
					output      : "",
					pixFmt      : "",
					cacheDir    : "",
					textTmpl    : T,
					tmpl        : nil,
					latexCmd    : "",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",

				},
			}, []any{
`ffmpeg \
	-filter_complex "
		concat=n=0:v=1:a=0:unsafe=1 [v]
	" \
	-map "[v]" ""`,
			},
		},
		{
			"~empty input, overwrite",
			compile,
			[]any{
				[]string{},
				[]string{},
				[]string{},
				[]string{},
				"",
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 0,
					height      : 0,
					framerate   : 0,
					defaultWait : 0.8,
					indent      : "\t",
					input       : nil,
					output      : "",
					pixFmt      : "",
					cacheDir    : "",
					textTmpl    : T,
					tmpl        : nil,
					latexCmd    : "",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",

				},
			}, []any{
`ffmpeg -y \
	-filter_complex "
		concat=n=0:v=1:a=0:unsafe=1 [v]
	" \
	-map "[v]" ""`,
			},
		},
		{
			"~empty input, overwrite/pix_fmt/framerate/output",
			compile,
			[]any{
				[]string{},
				[]string{},
				[]string{},
				[]string{},
				"",
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 0,
					height      : 0,
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : nil,
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : "",
					textTmpl    : T,
					tmpl        : nil,
					latexCmd    : "",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			}, []any{
`ffmpeg -y \
	-filter_complex "
		concat=n=0:v=1:a=0:unsafe=1 [v]
	" \
	-pix_fmt yuv420p -r 30 -map "[v]" "reel.mp4"`,
			},
		},
		{
			"One file, one text, but it's just a pause",
			compile,
			[]any{
				[]string{
					`-r 1 -t 2.00 -loop 1 -i "/tmp/foo.jpg"`,
				},
				[]string{
					"[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0]",
				},
				[]string{},
				[]string{"[in0]"},
				"",
				&State{
					faststart   : true,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : nil,
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : "",
					textTmpl    : T,
					tmpl        : nil,
					latexCmd    : "",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			}, []any{
`ffmpeg -y \
	-r 1 -t 2.00 -loop 1 -i "/tmp/foo.jpg" \
	-filter_complex "
		[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0];

		[in0] concat=n=1:v=1:a=0:unsafe=1 [v]
	" \
	-pix_fmt yuv420p -r 30 -movflags faststart -map "[v]" "reel.mp4"`,
			},
		},
		{
			"One file, one text (not a pause)",
			compile,
			[]any{
				[]string{
					`-r 1 -t 2.00 -loop 1 -i "/tmp/foo.jpg"`,
					`-i "`+filepath.Join(
						D,
						"a4795a4ae1128364caeffcaa73115cfcea550ca6da13adbbdd8e94e5d0f95b27.png",
					)+`"`,
				},
				[]string{
					"[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0]",
				},
				[]string{
					"[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.00,2.00)' [in0,1]",
				},
				[]string{"[in0,1]"},
				"",
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : nil,
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : D,
					textTmpl    : S.textTmpl, // './mkshort.go:/^var textTmpl'
					tmpl        : nil,
					latexCmd    : "lualatex",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			// XXX filepath.Join()
			}, []any{
`ffmpeg -y \
	-r 1 -t 2.00 -loop 1 -i "/tmp/foo.jpg" \
	-i "`+D+"/"+`a4795a4ae1128364caeffcaa73115cfcea550ca6da13adbbdd8e94e5d0f95b27.png" \
	-filter_complex "
		[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0];

		[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.00,2.00)' [in0,1];

		[in0,1] concat=n=1:v=1:a=0:unsafe=1 [v]
	" \
	-pix_fmt yuv420p -r 30 -map "[v]" "reel.mp4"`,
			},
		},
		{
			"One file, one text (not a pause), audio track",
			compile,
			[]any{
				[]string{
					`-r 1 -t 2.00 -loop 1 -i "/tmp/foo.jpg"`,
					`-i "`+filepath.Join(
						D,
						"a4795a4ae1128364caeffcaa73115cfcea550ca6da13adbbdd8e94e5d0f95b27.png",
					)+`"`,
					`-i "foo.mp3"`,
				},
				[]string{
					"[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0]",
				},
				[]string{
					"[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.00,2.00)' [in0,1]",
				},
				[]string{"[in0,1]"},
				"[2:a] afade=type=in:start_time=0:duration=4.00, afade=type=out:start_time=-1.00:duration=5.00 [a]",
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : nil,
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : D,
					textTmpl    : S.textTmpl, // './mkshort.go:/^var textTmpl'
					tmpl        : nil,
					latexCmd    : "lualatex",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			// XXX filepath.Join()
			}, []any{
`ffmpeg -y \
	-r 1 -t 2.00 -loop 1 -i "/tmp/foo.jpg" \
	-i "`+D+"/"+`a4795a4ae1128364caeffcaa73115cfcea550ca6da13adbbdd8e94e5d0f95b27.png" \
	-i "foo.mp3" \
	-filter_complex "
		[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0];

		[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.00,2.00)' [in0,1];

		[in0,1] concat=n=1:v=1:a=0:unsafe=1 [v];

		[2:a] afade=type=in:start_time=0:duration=4.00, afade=type=out:start_time=-1.00:duration=5.00 [a]
	" \
	-pix_fmt yuv420p -r 30 -map "[v]" -map "[a]" -c:a aac -shortest "reel.mp4"`,
			},
		},
	})
}

func TestParse(t *testing.T) {
	doTests(t, []test{
		{
			"no input",
			parse,
			[]any{
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 0,
					height      : 0,
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(``),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : "",
					textTmpl    : T,
					tmpl        : nil,
					latexCmd    : "",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			[]any{
				[]string{},
				[]string{},
				[]string{},
				[]string{},
				"",
				nil,
			},
		},
		{
			"just comments",
			parse,
			[]any{
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 0,
					height      : 0,
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(
						"# hello, world\n#foo bar\n",
					),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : "",
					textTmpl    : "",
					tmpl        : nil,
					latexCmd    : "",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			[]any{
				[]string{},
				[]string{},
				[]string{},
				[]string{},
				"",
				nil,
			},
		},
		{
			"One image, no texts",
			parse,
			[]any{
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(`:/tmp/foo.jpg`),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : "",
					textTmpl    : T,
					tmpl        : nil,
					latexCmd    : "",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			[]any{
				[]string{
					`-r 30 -t 0.00 -loop 1 -i "/tmp/foo.jpg"`,
				},
				[]string{
					"[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0]",
				},
				[]string{},
				[]string{"[in0]"},
				"",
				nil,
			},
		},
		{
			"Two images, no texts",
			parse,
			[]any{
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(
						":/tmp/foo.jpg\n:/tmp/bar.jpg",
					),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : "",
					textTmpl    : T,
					tmpl        : nil,
					latexCmd    : "",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			[]any{
				[]string{
					`-r 30 -t 0.00 -loop 1 -i "/tmp/foo.jpg"`,
					`-r 30 -t 0.00 -loop 1 -i "/tmp/bar.jpg"`,
				},
				[]string{
					"[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0]",
					"[1:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in1]",
				},
				[]string{},
				[]string{
					"[in0]",
					"[in1]",
				},
				"",
				nil,
			},
		},
		{
			"One image, empty text",
			parse,
			[]any{
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(
						":/tmp/foo.jpg\n0 2\n",
					),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : "",
					textTmpl    : T,
					tmpl        : nil,
					latexCmd    : "",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			[]any{
				[]string{
					`-r 30 -t 2.00 -loop 1 -i "/tmp/foo.jpg"`,
				},
				[]string{
					"[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0]",
				},
				[]string{},
				[]string{"[in0]"},
				"",
				nil,
			},
		},
		{
			"Two images, no texts but pauses",
			parse,
			[]any{
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(
						":/tmp/foo.jpg\n0 2\n:/tmp/bar.jpg\n+ 3.5",
					),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : "",
					textTmpl    : T,
					tmpl        : nil,
					latexCmd    : "",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			[]any{
				[]string{
					`-r 30 -t 2.00 -loop 1 -i "/tmp/foo.jpg"`,
					`-r 30 -t 3.50 -loop 1 -i "/tmp/bar.jpg"`,
				},
				[]string{
					"[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0]",
					"[1:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in1]",
				},
				[]string{},
				[]string{
					"[in0]",
					"[in1]",
				},
				"",
				nil,
			},
		},
		{
			"One image, non-empty text",
			parse,
			[]any{
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(
						":/tmp/foo.jpg\n0 2\n\thello, world!",
					),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : D,
					textTmpl    : T,
					tmpl        : template.Must(template.New("").Parse(T)),
					latexCmd    : "lualatex",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			[]any{
				[]string{
					`-r 30 -t 2.00 -loop 1 -i "/tmp/foo.jpg"`,
					`-i "`+
						filepath.Join(
							D,
							"a4795a4ae1128364caeffcaa73115cfcea550ca6da13adbbdd8e94e5d0f95b27.png",
						)+`"`,
				},
				[]string{
					"[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0]",
				},
				[]string{
					"[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.00,2.00)' [in0,1]",
				},
				[]string{"[in0,1]"},
				"",
				nil,
			},
		},
		{
			"Complete test: Virgin of the rocks",
			parse,
			[]any{
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(`
# This is a comment; will be ignored
:virgin-of-the-rocks-paris.jpg
0.7 3.5
	Here is what seems \\
	to be \\
	a little-known fact \\
# Another comments
+0.7 3
	About this famous \\
# Comments are okay anywhere as long as the line starts with a '#'
	Leonardo painting, \\
+0.7 3
	« Virgin of the Rocks »: \\
+0.7 3
	Of the four characters \\
+0.7 3
	Jesus, the \\
	bottom right child \\
+0.7 5
	Is the only one \\
	looking at the \\
	\textit{Source of Light} \\
	\emoji{face-in-clouds} \\
+0.7 3
	Pay attention to his \\
	chest and eyes \\
+0.7 2
+0.7 4
	The source of light can \\
	be deduced from the \\
	highlights and shadow \\
	patterns. \\
+0.7 1
:virgin-of-the-rocks-london.jpg
+0.7 4
	This also holds in the \\
	"London" version of the \\
	painting \\
+0.7 2
+0.7 2.5
	I am not sure how \\
	well-known this is \\
+0.7 3
	But I have not found it \\
	mentioned anywhere, \\
	so far \\
+0.7 3
`),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : D,
					textTmpl    : T,
					tmpl        : template.Must(template.New("").Parse(T)),
					latexCmd    : "lualatex",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			[]any{
				[]string{
					`-r 30 -t 36.80 -loop 1 -i "virgin-of-the-rocks-paris.jpg"`,
					`-i "`+filepath.Join(
						D,
						"90ed2c0bdd3b03acdd5ac3a66e78983c197933eac4d771a64cdd353d7f12542d.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"0874052b84eaea473405843666c658223881605489935376624fed4601c3ab50.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"24292a73f7cf7157b6b62f4799931e9c8fca518208eeda4782c8bb75ef36a2dc.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"2b96a193886a1a0a329a24bcd3468f98f13cb337b26b4045e5843b5563252b47.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"23ac3136c67210652b5582c7847a47787f68f50e2288d4381b92d1607489e42d.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"b954eb86f84bbd64ddb38f90c98b7f7a3afb7829daa81f957b2737b447fc9ec0.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"242215143f65b91164cfe12d9ba126b7fe52bf35eab7df8de227f1916e10e8b7.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"c9ed7cf3d1331df8d11fccf2d964fe221b1e037980e2346dd579604b94d9c53e.png",
					)+`"`,
					`-r 30 -t 17.30 -loop 1 -i "virgin-of-the-rocks-london.jpg"`,
					`-i "`+filepath.Join(
						D,
						"3d426114dc910999e9f5c9583b24444111dffda5ad5acf4e711d7079c83b4008.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"87f8eab5e96fcfde670fb93dfaa0649ba964286be63314c8fec524d1eafc184b.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"7f55fb922a6b776d224d061a9fcf1c725833680970f7d652965c735fbdde4342.png",
					)+`"`,
				},
				[]string{
					"[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0]",
					"[9:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in1]",
				},
				[]string{
					"[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.20)' [in0,1]",
					"[in0,1] [2:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,4.90,7.90)' [in0,2]",
					"[in0,2] [3:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.60,11.60)' [in0,3]",
					"[in0,3] [4:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,12.30,15.30)' [in0,4]",
					"[in0,4] [5:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,16.00,19.00)' [in0,5]",
					"[in0,5] [6:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,19.70,24.70)' [in0,6]",
					"[in0,6] [7:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,25.40,28.40)' [in0,7]",
					"[in0,7] [8:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,31.80,35.80)' [in0,9]",

					"[in1] [10:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.70)' [in1,1]",
					"[in1,1] [11:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.10,10.60)' [in1,3]",
					"[in1,3] [12:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,11.30,14.30)' [in1,4]",
				},
				[]string{"[in0,9]","[in1,4]"},
				"",
				nil,
			},
		},
		{
			"Complete test with audio: Virgin of the rocks",
			parse,
			[]any{
				&State{
					faststart   : false,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(`
# This will be overidden by the next call to @
@4 42 BMC19T1VivaldiSeasonsSummer.mp3

# Setup audio track; fade-in/out of 4sec both
@4 4 BMC19T1VivaldiSeasonsSpring.mp3

# This is a comment; will be ignored
:virgin-of-the-rocks-paris.jpg
0.7 3.5
	Here is what seems \\
	to be \\
	a little-known fact \\
# Another comments
+0.7 3
	About this famous \\
# Comments are okay anywhere as long as the line starts with a '#'
	Leonardo painting, \\
+0.7 3
	« Virgin of the Rocks »: \\
+0.7 3
	Of the four characters \\
+0.7 3
	Jesus, the \\
	bottom right child \\
+0.7 5
	Is the only one \\
	looking at the \\
	\textit{Source of Light} \\
	\emoji{face-in-clouds} \\
+0.7 3
	Pay attention to his \\
	chest and eyes \\
+0.7 2
+0.7 4
	The source of light can \\
	be deduced from the \\
	highlights and shadow \\
	patterns. \\
+0.7 1
:virgin-of-the-rocks-london.jpg
+0.7 4
	This also holds in the \\
	"London" version of the \\
	painting \\
+0.7 2
+0.7 2.5
	I am not sure how \\
	well-known this is \\
+0.7 3
	But I have not found it \\
	mentioned anywhere, \\
	so far \\
+0.7 3
`),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : D,
					textTmpl    : T,
					tmpl        : template.Must(template.New("").Parse(T)),
					latexCmd    : "lualatex",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			[]any{
				[]string{
					`-r 30 -t 36.80 -loop 1 -i "virgin-of-the-rocks-paris.jpg"`,
					`-i "`+filepath.Join(
						D,
						"90ed2c0bdd3b03acdd5ac3a66e78983c197933eac4d771a64cdd353d7f12542d.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"0874052b84eaea473405843666c658223881605489935376624fed4601c3ab50.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"24292a73f7cf7157b6b62f4799931e9c8fca518208eeda4782c8bb75ef36a2dc.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"2b96a193886a1a0a329a24bcd3468f98f13cb337b26b4045e5843b5563252b47.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"23ac3136c67210652b5582c7847a47787f68f50e2288d4381b92d1607489e42d.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"b954eb86f84bbd64ddb38f90c98b7f7a3afb7829daa81f957b2737b447fc9ec0.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"242215143f65b91164cfe12d9ba126b7fe52bf35eab7df8de227f1916e10e8b7.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"c9ed7cf3d1331df8d11fccf2d964fe221b1e037980e2346dd579604b94d9c53e.png",
					)+`"`,
					`-r 30 -t 17.30 -loop 1 -i "virgin-of-the-rocks-london.jpg"`,
					`-i "`+filepath.Join(
						D,
						"3d426114dc910999e9f5c9583b24444111dffda5ad5acf4e711d7079c83b4008.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"87f8eab5e96fcfde670fb93dfaa0649ba964286be63314c8fec524d1eafc184b.png",
					)+`"`,
					`-i "`+filepath.Join(
						D,
						"7f55fb922a6b776d224d061a9fcf1c725833680970f7d652965c735fbdde4342.png",
					)+`"`,
					`-i "BMC19T1VivaldiSeasonsSpring.mp3"`,
				},
				[]string{
					"[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0]",
					"[9:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in1]",
				},
				[]string{
					"[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.20)' [in0,1]",
					"[in0,1] [2:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,4.90,7.90)' [in0,2]",
					"[in0,2] [3:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.60,11.60)' [in0,3]",
					"[in0,3] [4:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,12.30,15.30)' [in0,4]",
					"[in0,4] [5:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,16.00,19.00)' [in0,5]",
					"[in0,5] [6:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,19.70,24.70)' [in0,6]",
					"[in0,6] [7:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,25.40,28.40)' [in0,7]",
					"[in0,7] [8:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,31.80,35.80)' [in0,9]",

					"[in1] [10:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.70)' [in1,1]",
					"[in1,1] [11:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.10,10.60)' [in1,3]",
					"[in1,3] [12:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,11.30,14.30)' [in1,4]",
				},
				[]string{"[in0,9]","[in1,4]"},
				"[13:a] afade=type=in:start_time=0:duration=4.00, afade=type=out:start_time=50.10:duration=4.00 [a]",
				nil,
			},
		},
	})
}

func TestParseAndCompile(t *testing.T) {
	doTests(t, []test{
		{
			"Virgin of the rocks v2",
			parseAndCompile,
			[]any{
				&State{
					faststart   : true,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(`
:virgin-of-the-rocks-paris.jpg
0.7 3.5
	Here is what seems \\
	to be \\
	a little-known fact \\
+0.7 3
	About this famous \\
	Leonardo painting, \\
+0.7 3
	« Virgin of the Rocks »: \\
+0.7 3
	Of the four characters \\
+0.7 3
	Jesus, the \\
	bottom right child \\
+0.7 5
	Is the only one \\
	looking at the \\
	\textit{Source of Light} \\
	\emoji{face-in-clouds} \\
+0.7 3
	Pay attention to his \\
	chest and eyes \\
+0.7 2
+0.7 4
	The source of light can \\
	be deduced from the \\
	highlights and shadow \\
	patterns. \\
+0.7 1
:virgin-of-the-rocks-london.jpg
+0.7 4
	This also holds in the \\
	"London" version of the \\
	painting \\
+0.7 2
+0.7 2.5
	I am not sure how \\
	well-known this is \\
+0.7 3
	But I have not found it \\
	mentioned anywhere, \\
	so far \\
+0.7 3
`),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : D,
					textTmpl    : T,
					tmpl        : template.Must(template.New("").Parse(T)),
					latexCmd    : "lualatex",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			// XXX filepath.Join() would be better to build the output
			// paths (lazy sed//)
			[]any{
`ffmpeg -y \
	-r 30 -t 36.80 -loop 1 -i "virgin-of-the-rocks-paris.jpg" \
	-i "`+D+"/"+`90ed2c0bdd3b03acdd5ac3a66e78983c197933eac4d771a64cdd353d7f12542d.png" \
	-i "`+D+"/"+`0874052b84eaea473405843666c658223881605489935376624fed4601c3ab50.png" \
	-i "`+D+"/"+`24292a73f7cf7157b6b62f4799931e9c8fca518208eeda4782c8bb75ef36a2dc.png" \
	-i "`+D+"/"+`2b96a193886a1a0a329a24bcd3468f98f13cb337b26b4045e5843b5563252b47.png" \
	-i "`+D+"/"+`23ac3136c67210652b5582c7847a47787f68f50e2288d4381b92d1607489e42d.png" \
	-i "`+D+"/"+`b954eb86f84bbd64ddb38f90c98b7f7a3afb7829daa81f957b2737b447fc9ec0.png" \
	-i "`+D+"/"+`242215143f65b91164cfe12d9ba126b7fe52bf35eab7df8de227f1916e10e8b7.png" \
	-i "`+D+"/"+`c9ed7cf3d1331df8d11fccf2d964fe221b1e037980e2346dd579604b94d9c53e.png" \
	-r 30 -t 17.30 -loop 1 -i "virgin-of-the-rocks-london.jpg" \
	-i "`+D+"/"+`3d426114dc910999e9f5c9583b24444111dffda5ad5acf4e711d7079c83b4008.png" \
	-i "`+D+"/"+`87f8eab5e96fcfde670fb93dfaa0649ba964286be63314c8fec524d1eafc184b.png" \
	-i "`+D+"/"+`7f55fb922a6b776d224d061a9fcf1c725833680970f7d652965c735fbdde4342.png" \
	-filter_complex "
		[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0];
		[9:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in1];

		[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.20)' [in0,1];
		[in0,1] [2:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,4.90,7.90)' [in0,2];
		[in0,2] [3:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.60,11.60)' [in0,3];
		[in0,3] [4:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,12.30,15.30)' [in0,4];
		[in0,4] [5:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,16.00,19.00)' [in0,5];
		[in0,5] [6:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,19.70,24.70)' [in0,6];
		[in0,6] [7:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,25.40,28.40)' [in0,7];
		[in0,7] [8:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,31.80,35.80)' [in0,9];
		[in1] [10:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.70)' [in1,1];
		[in1,1] [11:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.10,10.60)' [in1,3];
		[in1,3] [12:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,11.30,14.30)' [in1,4];

		[in0,9] [in1,4] concat=n=2:v=1:a=0:unsafe=1 [v]
	" \
	-pix_fmt yuv420p -r 30 -movflags faststart -map "[v]" "reel.mp4"`,
				nil,
			},
		},
		{
			"Virgin of the rocks v2, with audio",
			parseAndCompile,
			[]any{
				&State{
					faststart   : true,
					overwrite   : true,
					width       : 1080,
					height      : 1920,
					padColor    : "black",
					framerate   : 30,
					defaultWait : 0.8,
					indent      : "\t",
					input       : strings.NewReader(`
@4 4 BMC19T1VivaldiSeasonsSpring.mp3

:virgin-of-the-rocks-paris.jpg
0.7 3.5
	Here is what seems \\
	to be \\
	a little-known fact \\
+0.7 3
	About this famous \\
	Leonardo painting, \\
+0.7 3
	« Virgin of the Rocks »: \\
+0.7 3
	Of the four characters \\
+0.7 3
	Jesus, the \\
	bottom right child \\
+0.7 5
	Is the only one \\
	looking at the \\
	\textit{Source of Light} \\
	\emoji{face-in-clouds} \\
+0.7 3
	Pay attention to his \\
	chest and eyes \\
+0.7 2
+0.7 4
	The source of light can \\
	be deduced from the \\
	highlights and shadow \\
	patterns. \\
+0.7 1
:virgin-of-the-rocks-london.jpg
+0.7 4
	This also holds in the \\
	"London" version of the \\
	painting \\
+0.7 2
+0.7 2.5
	I am not sure how \\
	well-known this is \\
+0.7 3
	But I have not found it \\
	mentioned anywhere, \\
	so far \\
+0.7 3
`),
					output      : "reel.mp4",
					pixFmt      : "yuv420p",
					cacheDir    : D,
					textTmpl    : T,
					tmpl        : template.Must(template.New("").Parse(T)),
					latexCmd    : "lualatex",
					dryRun      : false,
					binsh       : "",
					textImgExt  : ".png",
					imgPrefix   : ":",
					audioPrefix : "@",
					headerSep   : " ",
				},
			},
			// XXX filepath.Join() would be better to build the output
			// paths (lazy sed//)
			[]any{
`ffmpeg -y \
	-r 30 -t 36.80 -loop 1 -i "virgin-of-the-rocks-paris.jpg" \
	-i "`+D+"/"+`90ed2c0bdd3b03acdd5ac3a66e78983c197933eac4d771a64cdd353d7f12542d.png" \
	-i "`+D+"/"+`0874052b84eaea473405843666c658223881605489935376624fed4601c3ab50.png" \
	-i "`+D+"/"+`24292a73f7cf7157b6b62f4799931e9c8fca518208eeda4782c8bb75ef36a2dc.png" \
	-i "`+D+"/"+`2b96a193886a1a0a329a24bcd3468f98f13cb337b26b4045e5843b5563252b47.png" \
	-i "`+D+"/"+`23ac3136c67210652b5582c7847a47787f68f50e2288d4381b92d1607489e42d.png" \
	-i "`+D+"/"+`b954eb86f84bbd64ddb38f90c98b7f7a3afb7829daa81f957b2737b447fc9ec0.png" \
	-i "`+D+"/"+`242215143f65b91164cfe12d9ba126b7fe52bf35eab7df8de227f1916e10e8b7.png" \
	-i "`+D+"/"+`c9ed7cf3d1331df8d11fccf2d964fe221b1e037980e2346dd579604b94d9c53e.png" \
	-r 30 -t 17.30 -loop 1 -i "virgin-of-the-rocks-london.jpg" \
	-i "`+D+"/"+`3d426114dc910999e9f5c9583b24444111dffda5ad5acf4e711d7079c83b4008.png" \
	-i "`+D+"/"+`87f8eab5e96fcfde670fb93dfaa0649ba964286be63314c8fec524d1eafc184b.png" \
	-i "`+D+"/"+`7f55fb922a6b776d224d061a9fcf1c725833680970f7d652965c735fbdde4342.png" \
	-i "BMC19T1VivaldiSeasonsSpring.mp3" \
	-filter_complex "
		[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0];
		[9:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in1];

		[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.20)' [in0,1];
		[in0,1] [2:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,4.90,7.90)' [in0,2];
		[in0,2] [3:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.60,11.60)' [in0,3];
		[in0,3] [4:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,12.30,15.30)' [in0,4];
		[in0,4] [5:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,16.00,19.00)' [in0,5];
		[in0,5] [6:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,19.70,24.70)' [in0,6];
		[in0,6] [7:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,25.40,28.40)' [in0,7];
		[in0,7] [8:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,31.80,35.80)' [in0,9];
		[in1] [10:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.70)' [in1,1];
		[in1,1] [11:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.10,10.60)' [in1,3];
		[in1,3] [12:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,11.30,14.30)' [in1,4];

		[in0,9] [in1,4] concat=n=2:v=1:a=0:unsafe=1 [v];

		[13:a] afade=type=in:start_time=0:duration=4.00, afade=type=out:start_time=50.10:duration=4.00 [a]
	" \
	-pix_fmt yuv420p -r 30 -movflags faststart -map "[v]" -map "[a]" -c:a aac -shortest "reel.mp4"`,
				nil,
			},
		},
	})
}

/*

	-movflags faststart
*/