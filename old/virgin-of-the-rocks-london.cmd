ffmpeg -y \
-loop 1 \
-i "/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-london-1080x1920.jpg" \
-framerate 30 \
-vf "scale=1080x1920" \
-movflags faststart \
-vf "[in]
drawtext=enable='between(t,0.7,4.7)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='This also holds in the
"London" version of the
painting',
drawtext=enable='between(t,8.1,10.6)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='I am not sure how
well-known this is',
drawtext=enable='between(t,11.3,14.3)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='But I have not found it
mentioned anywhere,
so far'
[out]" \
-t 18 \
/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-london-1080x1920.mp4
