package main

/*
 * The .short files are parsed via ':/^func parse\(',
 * to a few arrays. Those arrays can then be compiled
 * to a ffmpeg(1) command (string) via ':/func compile\('
 *
 * The parsing is implemented with a basic state machine;
 * the input .short file is processed line per line.
 *
 * The State type ':/type State struct {' contains all
 * the globals which can affect the processing, including
 * flags default values.
 */

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// All those are essentially global variables, stored
// in a struct to help with tests mainly.
type State struct {
	faststart   bool
	overwrite   bool
	width       int
	height      int
	padColor    string
	framerate   int
	defaultWait float64
	indent      string
	output      string
	pixFmt      string
	cacheDir    string
	textTmpl    string
	tmpl        *template.Template
	latexCmd    string
	input       io.Reader
	dryRun      bool
	binsh       string
	textImgExt  string

	// "Internal" stuff (.short file format)
	imgPrefix   string
	headerSep   string
}

var textTmpl = `\documentclass[preview,convert={density=600,outext=.png,command=\unexpanded{ {\convertexe\space -density \density\space\infile\space \ifx\size\empty\else -resize \size\fi\space -quality 90 -trim +repage -background "rgba(50,50,50,0.5)" -bordercolor "rgba(50,50,50,0.5)" -border 25 -flatten \outfile} } }]{standalone}
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

// TODO
var autoDuration = false

var S State = State{
	faststart   :     true,
	overwrite   :     true,
	width       :     1080,
	height      :     1920,
	padColor    :     "black",
	framerate   :     30,
	defaultWait :     0.8,
	indent      :     "\t",
	input       :     nil, // computed in ':/^func init\('
	output      :     "reel.mp4",
	pixFmt      :     "yuv420p",
	cacheDir    :     filepath.Join(os.Getenv("HOME"), ".mkshort"),
	textTmpl    :     textTmpl,
	tmpl        :     nil, // computed in ':/^func init\('
	latexCmd    :     "lualatex",
	dryRun      :     false,
	binsh       :     "/bin/sh",
	textImgExt  :    ".png",

	imgPrefix   :     ":",
	headerSep   :     " ",
}

type Text struct {
	// Eventually automatically computed (both of them)
	Start    float64
	Duration float64

	// Automatically computed ("internal")
	Path     string  // path to compiled image
}

func getCompiledTextBaseFn(s string, S *State) string {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
	return filepath.Join(S.cacheDir, hash)
}

func getCompiledTextImgFn(s string, S *State) string {
	return getCompiledTextBaseFn(s, S)+S.textImgExt
}

// TODO: ".tex" should be in conf if we want to run
// something other than LaTeX
func doCompileText(s string, S *State) (string, error) {
	tex := getCompiledTextBaseFn(s, S)+".tex"
	png := getCompiledTextImgFn(s, S)

	if err := os.WriteFile(tex, []byte(s), 0644); err != nil {
		return "", err
	}

	// TODO: quirk (-shell-escape; see that when we flag.Parse())
	cmd := exec.Command(S.latexCmd, "-shell-escape", tex)
	cmd.Dir = S.cacheDir

	// TODO: do something with output
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Println(tex, S.cacheDir)
		return "", fmt.Errorf("'%s': %s -- %s\n", S.latexCmd, err, out)
	}

	return png, nil
}

func getCached(s string, S *State) string {
	png := getCompiledTextImgFn(s, S)
	if _, err := os.Stat(png); errors.Is(err, os.ErrNotExist) {
		return ""
	}

	return png
}

// compile text or use cached data if any
func compileText(raw string, t *Text, S *State) (*Text, error) {
	var s strings.Builder
	var err error

	// Just a pause (TODO: dangerous semantic)
	if raw == "" {
		return t, nil
	}

	if err = S.tmpl.Execute(&s, map[string]any{ "text" : raw }); err != nil {
		return nil, err
	}
	t.Path = getCached(s.String(), S)
	if t.Path == "" {
		t.Path, err = doCompileText(s.String(), S)
	}
	return t, err
}

// The parsing is basically a state machine, mimicking
// the behavior of an early awk(1) prototype.
//
// We read S.input line per line, progressively fill a few
// string arrays which can later be compiled to a ffmpeg(1)
// command by ':/^func compile\('.
func parse(S *State) (
		[]string, []string, []string, []string, error) {

	// All input files to be fed to ffmpeg via -i.
	// They will be stored in the following order:
	// (indentation is for clarity only)
	//	img0
	//		overlay-img0-0
	//		overlay-img0-1
	//		...
	//	img1
	//		overlay-img1-0
	//		...
	//	...
	ins   := []string{}

	// Current image (path); it'll be added to ins
	// once either we reach EOF or the next image.
	img := ""

	// Collection of all overlay texts for the current
	// image which have been fully read already.
	texts := []*Text{}

	// What has been read so far of the current overlay.
	// Once the overlay text has been fully read, it'll
	// be compiled to a *Text, and append to the previous
	// texts variable.
	raw := ""

	// Total number of images so far. Those are the
	// "real" images, not the text overlays compiled
	// as images.
	//
	// It's mostly used to compute unique temporary stream
	// names for the final ffmpeg(1) command.
	nimg := 0

	// Starting time (sec) of the current text;
	// it's reset to 0. at the beginning of each
	// image: it's relative to the temporary
	// stream corresponding to a single image, not
	// of the final output stream.
	start := 0.

	// Duration of the current raw text
	duration := 0.

	// End of the previous text (time). This allows
	// to concatenate text, which is the only tested
	// behavior so far. Having parallel streams of
	// text at once could require some special care
	// here (but likely, will require e.g. on @text-position)
	prevend := 0.

	// Prepared scale filters. There's one for each
	// "real", input image, which are to be systematically
	// rescaled to S.width x S.height before we start
	// overlaying the texts. There's no rescaling
	// for the texts compiled from LaTeX (at least
	// for now).
	//
	// NOTE: thus, nimg == len(scales)
	scales := []string{}

	// Prepared overlay filters (chains).
	overs := []string{}

	// There's one overlay chain for each image: the
	// final stream created by each chain is the image
	// with text overlapped.
	//
	// The name of the final stream of each chain is
	// stored here, so that we can concatenate them.
	//
	// NOTE: we could be clever and have some extra
	// code to pick them from overs, but it's simpler this way.
	concats := []string{}

	// Compute total duration for which we'll want to
	// display the current image.
	calcDuration := func() float64 {
		d := 0.
		if len(texts) > 0 {
			d = texts[len(texts)-1].Start+
				texts[len(texts)-1].Duration-
				texts[0].Start
		}
		return d
	}

	addInput := func(r int, d float64, p string) int {
		ins = append(ins,
			fmt.Sprintf(`-r %d -t %.2f -loop 1 -i "%s"`, r, d, p))
		return len(ins)-1
	}

	addImg := func() { addInput(S.framerate, calcDuration(), img) }

	// Creates and registers a scale filter to scales (it's
	// actually a scale filter followed by a pad filter).
	addScale := func() {
		scales = append(scales, fmt.Sprintf(
			"[%d:v] "+
			"scale=%d:%d:force_original_aspect_ratio=decrease,"+
			"pad=%d:%d:(ow-iw)/2:(oh-ih)/2:color=%s"+
			" [in%d]",
			len(ins)-1, S.width, S.height, S.width, S.height, S.padColor, nimg,
		))
	}

	// register a single overlay
	addOverlay := func(from, ts, to string, start, duration float64) string {
		overs = append(overs, fmt.Sprintf(
			"%s %s " +
			"overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,%.2f,%.2f)'"+
			" %s", from, ts, start, start+duration, to,
		))

		return to
	}

	// register the current texts as an overlay chain; returns
	// the final stream name
	addOverlays := func() string {

		// Add the overlay chain for this image
		from := fmt.Sprintf("[in%d]", nimg)
		for i, t := range texts {
			// That's just a pause
			if len(t.Path) == 0 {
				continue
			}

			// New input file for this overlay; grab the
			// corresponding entry stream number: it's
			// this one we'll want to overlay.
			n := addInput(1, t.Duration, t.Path)

			// Add overlay & chain;
			from = addOverlay(from,
				fmt.Sprintf("[%d:v]", n),
				fmt.Sprintf("[in%d,%d]", nimg, i+1),
				t.Start, t.Duration,
			)
		}
		return from
	}

	// register a final overlay chain stream to be concatenated
	addConcat := func(n string) { concats = append(concats, n) }

	flushImg := func() {
		// Register everything. Order matters.
		addImg()
		addScale()
		addConcat(addOverlays())

		// Ready for the next image
		nimg++; img = ""

		// Remember: those are relative to the per-image
		// streams, and not to the final stream
		prevend = 0; start = 0

		texts = []*Text{}
	}

	flushRaw := func() error {
		t, err := compileText(raw, &Text{start, duration, ""}, S)
		if err != nil {
			return err
		}

		texts = append(texts, t)
		prevend = start+duration
		raw = ""
		duration = 0
		return nil
	}

	maybeFlushImg := func() {
		if img != "" {
			flushImg()
		}
	}

	maybeFlushRaw := func() error {
		// raw can be empty: no text is overlayed,
		// but we wait for the given duration nevertheless
		if duration > 0. {
			return flushRaw()
		}
		return nil
	}

	s := bufio.NewScanner(S.input)

	// main loop
	for s.Scan() {
		// read a line
		x := s.Text()

		// Empty: skipped
		if len(strings.TrimSpace(x)) == 0 {
			continue
		}

		// Comments
		// TODO: tests
		if strings.HasPrefix(x, "#") {
			continue
		}

		// Keep slurping raw LaTeX
		// TODO: indented text met before a start/duration
		// indication will be silently discarded (document at
		// least)
		if strings.HasPrefix(x, S.indent) {
			// TODO: eol cli args?
			raw = raw + "\n" + strings.TrimPrefix(x, S.indent)
			continue
		}

		if err := maybeFlushRaw(); err != nil {
			return []string{}, []string{}, []string{}, []string{}, err
		}

		// New image
		if strings.HasPrefix(x, S.imgPrefix) {
			maybeFlushImg()
			img = strings.TrimPrefix(x, S.imgPrefix)
			continue
		}

		// New text overlay header
		xs := strings.SplitN(x, S.headerSep, 2)

		// can't be zero (empty lines are skipped) nor can
		// it be strictly more than 2 because of the SplitN.
		if len(xs) == 1 {
			xs = append(xs, "")
		}

		// assert(len(xs) == 2)

		s := 0.
		if strings.HasPrefix(xs[0], "+") {
			s = prevend
			xs[0] = strings.TrimPrefix(xs[0], "+")
		}
		// xs[0] was a single "+"
		if xs[0] == "" {
			start = prevend + S.defaultWait // == s + S.defaultWait

		// Wait time is actually a number.
		} else {
			a, err := strconv.ParseFloat(xs[0], 64)
			if err != nil {
				return []string{}, []string{}, []string{}, []string{}, err
			}
			start = s + a
		}

		if xs[1] == "" {
			// Special value: we'll automatically
			// guess the duration once we've read all
			// the raw text (TODO)
			panic("Automatic duration not implemented yet")
			duration = -1.
		} else {
			a, err := strconv.ParseFloat(xs[1], 64)
			if err != nil {
				return []string{}, []string{}, []string{}, []string{}, err
			}
			duration = a
		}
	}

	if err := maybeFlushRaw(); err != nil {
		return []string{}, []string{}, []string{}, []string{}, err
	}
	maybeFlushImg()

	return ins, scales, overs, concats, s.Err()
}

// Compile a parsed .short file to a ffmpeg(1) command.
// The compiled .short is represented by the ins/scales/overs/concats
// arrays.
func compile(ins, scales, overs, concats []string, S *State) string {
	cmd := "ffmpeg "
	if S.overwrite {
		cmd = cmd + "-y "
	}
	cmd = cmd + "\\\n"
	for _, x := range ins {
		cmd += "\t" + x + " \\\n"
	}
	cmd += "\t-filter_complex \"\n"

	{
		for _, x := range scales {
			cmd += "\t\t" + x + ";\n"
		}
		if len(scales) > 0 {
			cmd += "\n"
		}
		for _, x := range overs {
			cmd += "\t\t" + x + ";\n"
		}
		if len(overs) > 0 {
			cmd += "\n"
		}
		cmd += "\t\t"
		for _, x := range concats {
			cmd += x + " "
		}
		cmd += fmt.Sprintf("concat=n=%d:v=1:a=0:unsafe=1 [v]", len(concats))
	}
	cmd += "\n\t\" \\\n\t"

	if len(S.pixFmt) > 0 {
		cmd += fmt.Sprintf("-pix_fmt %s ", S.pixFmt)
	}
	if S.framerate > 0 {
		cmd += fmt.Sprintf("-r %d ", S.framerate)
	}
	if S.faststart {
		cmd += "-movflags faststart "
	}

	// XXX what if no output? (can we stdout?)
	cmd += fmt.Sprintf("-map \"[v]\" \"%s\"", S.output)

	return cmd
}

func parseAndCompile(S *State) (string, error) {
	ins, scales, overs, concats, err := parse(S)
	if err != nil {
		return "", err
	}
	return compile(ins, scales, overs, concats, S), nil
}

// Essentially the entry point.
func parseCompileAndMaybeRun(S *State) error {
	scmd, err := parseAndCompile(S)
	if err != nil {
		return err
	}

	if S.dryRun {
		fmt.Println(scmd)
		return nil
	}

	// Write scmd to the stdin of a /bin/sh.
	cmd := exec.Command(S.binsh)

	// Forward child's stdout/stderr to parent
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Grab child's stdin
    stdin, err := cmd.StdinPipe()
    if err != nil {
		return err
    }

    go func() {
        defer stdin.Close()
        io.WriteString(stdin, scmd)
    }()

	return cmd.Run()
}

// For the most part, allows the global State variable S
// to be updated from its default values via CLI flags.
//
// NOTE: not named init() so that it's not ran with tests
func doInit() {
	// ffmpeg -movflags faststart
	flag.BoolVar(&S.faststart, "m", S.faststart, "Enable/disable faststart")
	// ffmpeg -y
	flag.BoolVar(&S.overwrite, "y", S.overwrite, "Allow automatic output overwrite")

	// Input files will be rescaled accordingly
	flag.IntVar(&S.width,  "w", S.width,  "output width")
	flag.IntVar(&S.height, "h", S.height, "output height")

	// https://ffmpeg.org/ffmpeg-utils.html#color-syntax
	// ("b" for background)
	flag.StringVar(&S.padColor, "b", S.padColor, "rescale padding color")

	// Output framerate (-r is the ffmpeg flag name)
	flag.IntVar(&S.framerate, "r", S.framerate, "output framerate")

	// s for "sleep"
	flag.Float64Var(&S.defaultWait, "s", S.defaultWait,
		"default waiting time between text")

	flag.StringVar(&S.indent, "i", S.indent, "text indentation")

	flag.StringVar(&S.pixFmt, "f", S.pixFmt, "output pixel format")
	flag.StringVar(&S.cacheDir, "d", S.cacheDir, "cache directory")

	// XXX clumsy
	flag.StringVar(&S.latexCmd, "c", S.latexCmd, "LaTeX command, splitted on spaces")

	// NOTE: file opening panic() when filename is too long, which happens
	// with the default textTmpl: we can't use a single CLI argument which
	// can be either a path or a template.
	flag.StringVar(&S.textTmpl, "t", S.textTmpl, "LaTeX template (string)")
	var tmplFn = flag.String("p", "", "LaTeX template (path)")

	flag.BoolVar(&S.dryRun, "x", S.dryRun, "Compiled command printed, not ran")

	flag.StringVar(&S.binsh, "l", S.binsh, "shell to run the compiled ffmpeg(1) command")
	flag.StringVar(&S.textImgExt, "e", S.textImgExt, "extension for the compiled text images")

	// TODO? imgPrefix, headerSep

	flag.Parse()

	// No template file => use the string template
	if _, err := os.Stat(*tmplFn); errors.Is(err, os.ErrNotExist) {
		S.tmpl = template.Must(template.New("").Parse(S.textTmpl))
	} else {
		S.tmpl = template.Must(template.New("").ParseFiles(*tmplFn))
	}

	// Expected to exists later
	if err := os.MkdirAll(S.cacheDir, 0750); err != nil {
		log.Fatal(err)
	}

	// XXX/NOTE: Not quite sturdy, but ensuring cacheDir to be absolute
	// makes things simpler when running ffmpeg(1)/latex(1), considering
	// .short image paths are typically relative.
	if !strings.HasPrefix(S.cacheDir, "/") {
		S.cacheDir = filepath.Join(os.Getenv("PWD"), S.cacheDir)
	}

	// Set input/output. Input is the .short file; it
	// needs to be opened. Output OTOH is a mere string
	// provided to ffmpeg
	var input string
	switch len(flag.Args()) {
	case 0:
		S.output, input = "/dev/stdout", "-"
	case 1:
		S.output, input = flag.Args()[0], "-"
	default:
		S.output, input = flag.Args()[0], flag.Args()[1]
	}

	if input == "-" {
		input = "/dev/stdin"
	}

	var err error
	if S.input, err = os.Open(input); err != nil {
		log.Fatal("Opening", input, ": ", err)
	}
}

func main() {
	doInit()

	if err := parseCompileAndMaybeRun(&S); err != nil {
		log.Fatal(err)
	}
}
