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
