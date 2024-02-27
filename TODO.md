## drawtext mode @feature
We could, as it was done in the prototype, rely on ffmpeg's drawtext mode
instead of LaTeX+magick.

## floating-point printing precision @feature
Two decimals should be enough for our present purposes, but we
may want to allow the precision to be configurable.

## automatic (text) duration computation @feature
The fact that the text input is LaTeX may add some
extra difficulties, but a basic implementation should
be easy.

## specify text position @feature @text-position
Find ways to move text around easily.

## allow videos and not just text @feature

## sound effect @feature
E.g. play a small sound everytime the text change.

## complete man page

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
