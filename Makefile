# Default installation directory.
#	make install dir=$HOME/bin
dir ?= /bin/
mandir ?= /usr/share/man/man1/
root ?= root
group ?= root


.PHONY: all
all: mkshort tests

.PHONY: help
help:
	@echo Available targets:
	@echo "	mkshort    : build mkshort"
	@echo "	all        : build mkshort; run tests"
	@echo "	clean      : removed compiled files"
	@echo "	tests      : run automated tests"
	@echo "	example    : display the ffmpeg(1) command,"
	@echo "	             build (a bit slow) and play the example."
	@echo "	install    : install to ${dir} and ${mandir}"
	@echo "	uninstall  : remove installed files"

.PHONY: tests
tests:
	@echo Running tests...
	@go test -v mkshort_test.go ftests.go mkshort.go

mkshort: mkshort.go
	@echo Building $@...
	@go build $^

# Audio track used for the example
BMC19T1VivaldiSeasonsSpring.mp3:
	@echo Fetching $@...
	@wget 'http://www.baroquemusic.org/DLower/BMC19T1VivaldiSeasonsSpring.mp3' -O "$@"

.PHONY: clean
clean:
	@echo Remove compiled binaries...
	@rm -f mkshort virgin-of-the-rocks.cmd virgin-of-the-rocks.mp4

virgin-of-the-rocks.cmd: virgin-of-the-rocks.short
	@echo Building $@...
	@go run mkshort.go -d .cache -x < virgin-of-the-rocks.short  > $@

virgin-of-the-rocks.mp4: virgin-of-the-rocks.short BMC19T1VivaldiSeasonsSpring.mp3
	@echo Building $@...
	@go run mkshort.go -d .cache virgin-of-the-rocks.mp4 virgin-of-the-rocks.short

virgin-of-the-rocks.gif: virgin-of-the-rocks.mp4
	@echo Building $@...
	@ffmpeg -y -i $^ $@

.PHONY: example
example:                                \
		mkshort                         \
		BMC19T1VivaldiSeasonsSpring.mp3 \
		virgin-of-the-rocks-paris.jpg   \
		virgin-of-the-rocks-london.jpg  \
		virgin-of-the-rocks.cmd         \
		virgin-of-the-rocks.mp4         \
		virgin-of-the-rocks.gif
	@echo Generated command:
	@cat virgin-of-the-rocks.cmd
	@echo Playing virgin-of-the-rocks.mp4...
	@mplayer virgin-of-the-rocks.mp4

.PHONY: install
install: mkshort
	@echo Installing mkshort to ${dir}/mkshort...
	@install -o ${root} -g ${group} -m 755 mkshort ${dir}/mkshort
	@echo Installing mkshort.1 to ${mandir}/mkshort.1...
	@install -o ${root} -g ${group} -m 644 mkshort.1 ${mandir}/mkshort.1

.PHONY: uninstall
uninstall:
	@echo Removing ${dir}/mkshort...
	@rm -f ${dir}/mkshort
	@echo Removing ${mandir}/mkshort.1...
	@rm -f ${mandir}/mkshort.1
