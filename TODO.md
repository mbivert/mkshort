## "1 frame" bug @bug
There's one (small) remaining issue for which I have no (clean) solution yet:
consider the following overlay chain:

	[0:v] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.20)' [in0,1];
	[in0,1] [2:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,4.90,7.90)' [in0,2];
	[in0,2] [3:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.60,11.60)' [in0,3];
	...

Then, for about a frame, `[1:v]` is displayed at the beginning of the stream,
before being hidden for `0.7` and shown until `4.20` as it should. I haven't
found (clean) ways to get rid of this one-frame display; surprisingly, weird
manipulations like this one don't help:

	[0:v] split, overlay=enable='between(t,0,0.7)' [in0];
	[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.20)' [in0,1];
	[in0,1] [2:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,4.90,7.90)' [in0,2];
	...

## drawtext mode @feature
We could, as it was done in the prototype, rely on ffmpeg's drawtext mode
instead of LaTeX+magick.

## floating-point printing precision @feature
Two decimals should be enough for our present purposes, but we
may want to allow the precision to be configurable.

## automatic duration computation @feature
The fact that the text input is LaTeX may add some
extra difficulties, but a basic implementation should
be easy.

## specify text position @feature @text-position
Find ways to move text around easily.

## allow videos and not just text @feature

## sound effect @feature
E.g. play a small sound everytime the text change.

## audio track @feature
For now, a global audio track can be added manually
to the final stream, but we could allow it in here.
We may want to be able to change the soundtrack
as we go, and to generate both a silent and non-silent
output streams at once.

## complete man page

## clean code

## bad file format @bugs
test and fix e.g.
```
1 2
	foo
:hello.jpg
``
Or
```
	foo
:hello.jpg
```
