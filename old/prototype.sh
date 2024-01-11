#!/bin/sh

# EDIT: Included just in case; won't run as-is.
exit 1

# We could add options to move the overlay around.

# XXX Can't have simple quote because of the escaping mess...
# XXX Simple text-simley ":-)" isn't shown.

# Not using an image with the correct size seems to misplace
# the text overlays (can't see them anywhere).
#
# convert Leonardo*.jpg -resize 1080x1920 -size 1080x1920 xc:black +swap -gravity center -composite virgin-of-the-rocks-1080x1920.jpg
# convert virgin-of-the-rocks-london.jpg -resize 1080x1920 -size 1080x1920 xc:black +swap -gravity center -composite virgin-of-the-rocks-london-1080x1920.jpg

# The following font can be used for b&w smileys.
# /usr/share/blender/4.0/datafiles/fonts/NotoEmoji-VariableFont_wght.woff2

ffile="fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff"
fcolor="fontcolor=white"
fsize="fontsize=120"
fsize="fontsize=h/20"
fbox="box=1:boxcolor=black@0.5:boxborderw=10"
fpos="x=(w-text_w)/2:y=(h-text_h+1)/2+t"
ftiming="enable='between(t,0,1.5)'"
falign="text_align=M:text_align=C"

mkshort() {
awk '
	BEGIN {
		start    = 0
		duration = 0
		prevend  = 0
		text     = ""
#		texts    =
		ntexts   = 0
		font     = ""
		print "ffmpeg -y \\"
		print "-loop 1 \\"
		print "-i \"'$img'\" \\"
		print "-framerate 30 \\"
		print "-vf \"scale=1080x1920\" \\"
		print "-movflags faststart \\"
		print "-vf \"[in]"
	}

	function flush() {
		end = start + duration
		ftiming = sprintf("drawtext=enable='"'"'between(t,%.1f,%.1f)'"'"'", start, end)
		if (font == "")
			font = "'"$ffile"'"
		texts[ntexts++] = sprintf("%s:%s'":$fcolor:$fsize:$fbox:$fpos:$falign"':text='"'"'%s'"'"'", ftiming, font, text)
		font = ""
	}

	function maybeflush() {
		if (text != "") flush(); text = ""
	}

	! /^[\t ]/ {
		maybeflush()

		text     = ""
		start    = $1
		duration = $2

		if (substr(start, 1, 1) == "+")
			start = prevend + substr(start, 2)

		prevend = start + duration

		if ($3 != "") {
			font = "fontfile=" $3
		}

		next
	}
	/^[\t ]/ {
		if (text != "")
			text = text "\n"
		text = text "" substr($0, 2)
	}
	END {
		maybeflush()
		for (i = 0; i < ntexts-1; i++) {
			printf("%s,\n", texts[i])
		}
		printf("%s\n", texts[i])

		print "[out]\" \\"
		print "-t " prevend" \\"
		print "'"$out"'"
	}
'
}

set -e

cat > /tmp/blbl <<EOF
file /home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-1080x1920.mp4
file /home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-london-1080x1920.mp4
EOF

rm -f /home/mb/gits/posts/virgin-of-the-rocks/reel.mp4
ffmpeg -f concat -safe 0 -i /tmp/blbl -c copy /home/mb/gits/posts/virgin-of-the-rocks/reel.mp4

# ffmpeg -y -i reel.mp4 -pix_fmt yuv420p real-yuv420p.mp4

exit 0

img="/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-1080x1920.jpg"
out="/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-1080x1920.mp4"
cmd="/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks.cmd"
short="/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks.short"
mkshort < "$short" | tee "$cmd"

rm -f "$out"
sh "$cmd"

img="/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-london-1080x1920.jpg"
out="/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-london-1080x1920.mp4"
cmd="/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-london.cmd"
short="/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-london.short"
mkshort < "$short" | tee "$cmd"

rm -f "$out"
sh "$cmd"

# ffmpeg -y -i reel.mp4 -pix_fmt yuv420p real-yuv420p.mp4

mplayer /home/mb/gits/posts/me-profile-ink/reel.mp4